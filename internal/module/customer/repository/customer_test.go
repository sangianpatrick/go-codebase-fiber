package repository_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sangianpatrick/go-codebase-fiber/internal/module/customer/entity"
)

const (
	email = "patrick@mail.com"
)

var (
	insertQuery      = "INSERT INTO customer"
	findByEmailQuery = "SELECT id, email, firstname, lastname, verification_status, member_status, password, password_salt, created_at, updated_at FROM customer WHERE email = \\$1 LIMIT 1"
)

func NewCustomer(now time.Time) entity.Customer {
	return entity.Customer{
		ID:                 uuid.New(),
		Email:              "patrick@mail.com",
		Firstname:          "patrick",
		Lastname:           "sangian",
		VerificationStatus: entity.VerificationStatus_Verified,
		MemberStatus:       entity.MemberStatus_Active,
		Password:           "password",
		PasswordSalt:       "password_salt",
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func SetEnv(t *testing.T) {
	t.Setenv("SERVICE_NAME", "test-service")
}
