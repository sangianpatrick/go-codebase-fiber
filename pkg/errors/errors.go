package errors

import "fmt"

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
