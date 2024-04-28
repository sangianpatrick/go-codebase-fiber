package customer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/applogger"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/errors"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/status"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/util"
	"github.com/sangianpatrick/go-codebase-fiber/pkg/validator"
)

type UseCase interface {
	SignUp(ctx context.Context, req SignUpRequest) (SignUpResponse, error)
}

type useCase struct {
	logger     *applogger.ZapLogger
	timeout    time.Duration
	secret     string
	validator  *validator.Validator
	repository Repository
}

type UseCaseProperty struct {
	Logger     *applogger.ZapLogger
	Timeout    time.Duration
	Secret     string
	Validator  *validator.Validator
	Repository Repository
}

func (p UseCaseProperty) Create() UseCase {
	return &useCase{
		logger:     p.Logger,
		timeout:    p.Timeout,
		secret:     p.Secret,
		validator:  p.Validator,
		repository: p.Repository,
	}
}

// SignUp implements UseCase.
func (u *useCase) SignUp(ctx context.Context, req SignUpRequest) (SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if err := u.validator.ValidateStruct(ctx, req); err != nil {
		return SignUpResponse{}, err
	}

	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return SignUpResponse{}, err
	}

	_, err = u.repository.FindByEmail(ctx, req.Email, tx)
	if err == nil {
		u.repository.RollbackTx(ctx, tx)
		return SignUpResponse{}, errors.New(http.StatusConflict, status.ALREADY_EXIST, "customer is already registered")
	}
	if !errors.MatchStatus(err, status.NOT_FOUND) {
		u.repository.RollbackTx(ctx, tx)
		return SignUpResponse{}, err
	}

	now := time.Now()
	expiresAt := now.Add(time.Hour * 3)
	ID := uuid.New()
	passwordSalt := util.GenerateRandomHEX(32)
	passwordHash := util.GenerateSecret(fmt.Sprintf("%s:%s", u.secret, req.Password), passwordSalt, 256)

	c := Customer{
		ID:                 ID,
		Email:              req.Email,
		Firstname:          req.Firstname,
		Lastname:           req.Lastname,
		VerificationStatus: VerificationStatus_Unverified,
		MemberStatus:       MemberStatus_Active,
		Password:           passwordHash,
		PasswordSalt:       passwordSalt,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := u.repository.Save(ctx, c, tx); err != nil {
		u.repository.RollbackTx(ctx, tx)
		return SignUpResponse{}, err
	}

	u.repository.CommitTx(ctx, tx)

	resp := SignUpResponse{
		Customer: CustomerResponse{
			ID:                 c.ID,
			Email:              c.Email,
			Firstname:          c.Firstname,
			Lastname:           c.Lastname,
			VerificationStatus: c.VerificationStatus,
			CreatedAt:          c.CreatedAt,
		},
		VerificationExpiresAt: expiresAt,
	}

	return resp, nil
}
