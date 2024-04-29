package use_case_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/entity"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/request"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/use_case"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/errors"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/validator"
	mock_customer "github.com/sangianpatrick/go-codebase-fiber/mock/module/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUpSuccess(t *testing.T) {
	logger := applogger.GetZap()
	vld := validator.Get()
	customerRepo := new(mock_customer.CustomerRepository)
	customerRepo.On("BeginTx", mock.Anything).Return(new(sql.Tx), nil)
	customerRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*sql.Tx")).Return(entity.Customer{}, errors.NotFound)
	customerRepo.On("Save", mock.Anything, mock.AnythingOfType("entity.Customer"), mock.AnythingOfType("*sql.Tx")).Return(nil)
	customerRepo.On("CommitTx", mock.Anything, mock.AnythingOfType("*sql.Tx")).Return(nil)

	props := use_case.CustomerUseCaseProperty{
		Logger:     logger,
		Timeout:    time.Second * 10,
		Secret:     uuid.NewString(),
		Validator:  vld,
		Repository: customerRepo,
	}

	u := props.Create()

	req := request.SignUpRequest{
		Firstname: "patrick",
		Lastname:  "sangian",
		Email:     "patrick@mail.com",
		Password:  "12345678",
	}

	resp, err := u.SignUp(context.TODO(), req)

	assert.NoError(t, err)
	assert.Equal(t, req.Email, resp.Email)
}
