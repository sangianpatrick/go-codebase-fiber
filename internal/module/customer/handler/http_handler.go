package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/request"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/use_case"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/errors"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/response"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/status"
)

type HTTPHandler struct {
	UseCase use_case.CustomerUseCase
}

func InitHTTPHandler(router *fiber.App, useCase use_case.CustomerUseCase) {
	handler := HTTPHandler{
		UseCase: useCase,
	}
	router.Post("/user/v1/customers/sign-up", handler.SignUp)
}

func (handler HTTPHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req request.SignUpRequest

	if err := c.BodyParser(&req); err != nil {
		return errors.New(http.StatusUnprocessableEntity, status.UNPROCESSABLE_ENTITY, err.Error())
	}

	resp, err := handler.UseCase.SignUp(ctx, req)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(response.RESTEnvelope{
		Status:  status.CREATED,
		Message: "customer has been successfully signed up",
		Data:    resp,
	})
}
