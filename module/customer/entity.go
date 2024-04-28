package customer

import (
	"time"

	"github.com/google/uuid"
)

const (
	VerificationStatus_Unverified = "UNVERIFIED"
	VerificationStatus_Verified   = "VERIFIED"
	MemberStatus_Active           = "ACTIVE"
	MemberStatus_Inactive         = "INACTIVE"
)

type Customer struct {
	ID                 uuid.UUID
	Email              string
	Firstname          string
	Lastname           string
	VerificationStatus string
	MemberStatus       string
	Password           string
	PasswordSalt       string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
