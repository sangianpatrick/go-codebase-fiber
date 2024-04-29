package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/handler"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/request"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/response"
	rpkg "github.com/sangianpatrick/go-codebase-fiber/internal/pkg/response"
	mock_customer "github.com/sangianpatrick/go-codebase-fiber/mock/module/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUpSuccess(t *testing.T) {
	req := request.SignUpRequest{
		Firstname: "patrick",
		Lastname:  "sangian",
		Email:     "patrick@mail.com",
		Password:  "12345678",
	}
	reqBuff, _ := json.Marshal(req)

	u := new(mock_customer.CustomerUseCase)
	u.On("SignUp", mock.Anything, mock.AnythingOfType("request.SignUpRequest")).Return(response.SignUpResponse{}, nil)

	h := handler.HTTPHandler{
		UseCase: u,
	}

	timeoutMs := 10
	router := fiber.New(fiber.Config{
		ErrorHandler: rpkg.FiberErrorHandler,
	})

	router.Post("/test", h.SignUp)

	httpReq := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(reqBuff))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, _ := router.Test(httpReq, timeoutMs)
	defer resp.Body.Close()

	respBuf, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Fatal(string(respBuf))
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
