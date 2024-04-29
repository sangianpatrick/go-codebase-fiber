package response

import (
	"time"

	"github.com/google/uuid"
)

type CustomerResponse struct {
	ID                 uuid.UUID `json:"id"`
	Email              string    `json:"email"`
	Firstname          string    `json:"firstname"`
	Lastname           string    `json:"lastname"`
	VerificationStatus string    `json:"verification_status"`
	CreatedAt          time.Time `json:"created_at"`
}

type SignUpResponse CustomerResponse
