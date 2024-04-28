package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sangianpatrick/go-codebase-fiber/config"
	"github.com/sangianpatrick/go-codebase-fiber/module/customer"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/monitoring"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/postgres"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/response"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/status"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/validator"
	"go.uber.org/zap"
)

func main() {
	c := config.Get()

	logger := applogger.GetZapLogger()

	mon := monitoring.NewOpenTelemetry(
		c.Service.Name,
		c.Service.Environment,
		c.GCP.ProjectID,
	)

	mon.Start(context.Background())

	vld := validator.Get()
	db := postgres.GetDatabase()
	if err := db.Ping(); err != nil {
		logger.Error(context.Background(), err.Error(), zap.Error(err))
	}

	router := fiber.New(fiber.Config{
		AppName:      c.Service.Name,
		ServerHeader: c.Service.Name,
	})
	router.Use(
		otelfiber.Middleware(
			otelfiber.WithServerName(c.Service.Name),
		),
	)
	router.Hooks().OnShutdown(OnShutdown(logger))

	router.Get("/user", HealthCheck)

	customerRepository := customer.NewRepository(logger, db)
	customerUseCase := customer.UseCaseProperty{
		Logger:     logger,
		Timeout:    c.Service.Timeout,
		Secret:     c.Service.Secret,
		Validator:  vld,
		Repository: customerRepository,
	}.Create()
	customer.InitHTTPHandler(router, customerUseCase)

	go func() {
		router.Listen(fmt.Sprintf(":%d", c.Service.Port))
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	router.ShutdownWithTimeout(c.Service.Timeout)
	db.Close()
	mon.Stop(context.Background())
}

func OnShutdown(l *applogger.ZapLogger) func() error {
	return func() error {
		l.Info(context.Background(), "server shutdown")
		return nil
	}
}

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(response.RESTEnvelope{
		Status:  status.OK,
		Message: "service is running properly",
	})
}
