package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/middleware"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	timeoutMs := 10
	router := fiber.New(fiber.Config{
		ErrorHandler: response.FiberErrorHandler,
	})

	logger := applogger.GetZap()
	router.Get("/test", middleware.RecoveryMiddleware(logger), func(c *fiber.Ctx) error {
		panic("panic test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	resp, _ := router.Test(req, timeoutMs)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
