package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/applogger"
)

func RecoveryMiddleware(logger applogger.AppLogger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			r := recover()
			if r != nil {
				err = fmt.Errorf(fmt.Sprint(r))
				stack := debug.Stack()
				logger.Error(c.UserContext(), err.Error(), applogger.Error(err), applogger.String("stack_trace", string(stack)))
			}
		}()

		err = c.Next()

		return
	}
}
