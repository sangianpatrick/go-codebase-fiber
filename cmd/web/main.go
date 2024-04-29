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
	customer_handler "github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/handler"
	customer_repository "github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/repository"
	customer_usecase "github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/use_case"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/middleware"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/monitoring"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/postgres"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/response"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/status"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/validator"
)

func main() {
	c := config.Get()

	logger := applogger.GetZap()

	mon := monitoring.NewOpenTelemetry(
		c.Service.Name,
		c.Service.Environment,
		c.GCP.ProjectID,
	)
	mon.Start(context.Background())

	vld := validator.Get()
	db := postgres.GetDatabase()
	if err := db.Ping(); err != nil {
		logger.Error(context.Background(), err.Error(), applogger.Error(err))
	}

	router := fiber.New(fiber.Config{
		AppName:      c.Service.Name,
		ServerHeader: c.Service.Name,
		ErrorHandler: response.FiberErrorHandler,
	})
	router.Use(
		otelfiber.Middleware(
			otelfiber.WithServerName(c.Service.Name),
		),
		middleware.RecoveryMiddleware(logger),
	)
	router.Hooks().OnShutdown(OnShutdown(logger))

	router.Get("/user", HealthCheck)

	customerRepository := customer_repository.NewRepository(logger, db)
	customerUseCase := customer_usecase.CustomerUseCaseProperty{
		Logger:     logger,
		Timeout:    c.Service.Timeout,
		Secret:     c.Service.Secret,
		Validator:  vld,
		Repository: customerRepository,
	}.Create()
	customer_handler.InitHTTPHandler(router, customerUseCase)

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

func OnShutdown(l applogger.AppLogger) func() error {
	return func() error {
		l.Info(context.Background(), "server shutdown")
		return nil
	}
}

func HealthCheck(c *fiber.Ctx) error {
	data := make([]int, 1)
	return c.Status(http.StatusOK).JSON(response.RESTEnvelope{
		Status:  status.OK,
		Message: "service is running properly",
		Data:    data[2],
	})
}
