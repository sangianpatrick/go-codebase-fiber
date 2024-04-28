package customer

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/errors"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/response"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/status"
)

type HTTPHandler struct {
	UseCase UseCase
}

func InitHTTPHandler(router *fiber.App, useCase UseCase) {
	handler := HTTPHandler{
		UseCase: useCase,
	}
	router.Post("/user/v1/customers/sign-up", handler.SignUp)
}

func (handler HTTPHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var req SignUpRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(response.RESTEnvelope{
			Status:  status.UNPROCESSABLE_ENTITY,
			Message: err.Error(),
		})
	}

	resp, err := handler.UseCase.SignUp(ctx, req)
	if err != nil {
		ae := errors.Destruct(err)
		return c.Status(ae.HTTPStatusCode()).JSON(response.RESTEnvelope{
			Status:  ae.Status(),
			Message: ae.Message(),
			Data:    nil,
			Meta:    nil,
		})
	}

	return c.Status(http.StatusCreated).JSON(response.RESTEnvelope{
		Status:  status.CREATED,
		Message: "customer has been successfully signed up",
		Data:    resp,
		Meta:    nil,
	})
}
