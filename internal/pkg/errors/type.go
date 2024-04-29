package errors

import (
	"net/http"

	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/status"
)

var (
	NotFound = New(http.StatusNotFound, status.NOT_FOUND, ``)
)
