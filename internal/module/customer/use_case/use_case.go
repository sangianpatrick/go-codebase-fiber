package use_case

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/entity"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/repository"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/request"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/response"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/errors"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/status"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/util"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/validator"
)

type CustomerUseCase interface {
	SignUp(ctx context.Context, req request.SignUpRequest) (response.SignUpResponse, error)
}

type useCase struct {
	logger     applogger.AppLogger
	timeout    time.Duration
	secret     string
	validator  *validator.Validator
	repository repository.CustomerRepository
}

type CustomerUseCaseProperty struct {
	Logger     applogger.AppLogger
	Timeout    time.Duration
	Secret     string
	Validator  *validator.Validator
	Repository repository.CustomerRepository
}

func (p CustomerUseCaseProperty) Create() CustomerUseCase {
	return &useCase{
		logger:     p.Logger,
		timeout:    p.Timeout,
		secret:     p.Secret,
		validator:  p.Validator,
		repository: p.Repository,
	}
}

// SignUp implements UseCase.
func (u *useCase) SignUp(ctx context.Context, req request.SignUpRequest) (response.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if err := u.validator.ValidateStruct(ctx, req); err != nil {
		return response.SignUpResponse{}, err
	}

	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return response.SignUpResponse{}, err
	}

	_, err = u.repository.FindByEmail(ctx, req.Email, tx)
	if err == nil {
		u.repository.RollbackTx(ctx, tx)
		return response.SignUpResponse{}, errors.New(http.StatusConflict, status.ALREADY_EXIST, "customer is already registered")
	}
	if !errors.MatchStatus(err, status.NOT_FOUND) {
		u.repository.RollbackTx(ctx, tx)
		return response.SignUpResponse{}, err
	}

	now := time.Now()
	ID := uuid.New()
	passwordSalt := util.GenerateRandomHEX(32)
	passwordHash := util.GenerateSecret(fmt.Sprintf("%s:%s", u.secret, req.Password), passwordSalt, 256)

	c := entity.Customer{
		ID:                 ID,
		Email:              req.Email,
		Firstname:          req.Firstname,
		Lastname:           req.Lastname,
		VerificationStatus: entity.VerificationStatus_Verified,
		MemberStatus:       entity.MemberStatus_Active,
		Password:           passwordHash,
		PasswordSalt:       passwordSalt,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := u.repository.Save(ctx, c, tx); err != nil {
		u.repository.RollbackTx(ctx, tx)
		return response.SignUpResponse{}, err
	}

	u.repository.CommitTx(ctx, tx)

	resp := response.SignUpResponse{
		ID:                 c.ID,
		Email:              c.Email,
		Firstname:          c.Firstname,
		Lastname:           c.Lastname,
		VerificationStatus: c.VerificationStatus,
		CreatedAt:          c.CreatedAt,
	}

	return resp, nil
}
