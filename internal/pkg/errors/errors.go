package errors

import (
	"fmt"
	"net/http"

	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/status"
)

type Error struct {
	httpStatusCode int
	status         string
	message        string
}

func New(httpStatusCode int, status string, message string) *Error {
	return &Error{
		httpStatusCode: httpStatusCode,
		status:         status,
		message:        message,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s: %s", e.httpStatusCode, e.status, e.message)
}

func (e Error) HTTPStatusCode() int {
	return e.httpStatusCode
}

func (e Error) Status() string {
	return e.status
}

func (e Error) Message() string {
	return e.message
}

func Destruct(err error) *Error {
	if err == nil {
		return nil
	}
	ae, ok := err.(*Error)
	if !ok {
		return New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, "an error occured while attempting request")
	}

	return ae
}

func MatchStatus(err error, s string) bool {
	if err == nil {
		return false
	}
	ae := Destruct(err)

	return ae.Status() == s
}
