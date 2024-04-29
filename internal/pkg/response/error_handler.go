package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/errors"
)

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	e := errors.Destruct(err)

	return c.Status(e.HTTPStatusCode()).JSON(RESTEnvelope{
		Status:  e.Status(),
		Message: e.Message(),
	})
}
