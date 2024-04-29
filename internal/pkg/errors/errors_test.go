package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/errors"
	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/status"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	ae := errors.New(http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, "ise")

	assert.Equal(t, http.StatusInternalServerError, ae.HTTPStatusCode())
	assert.Equal(t, status.INTERNAL_SERVER_ERROR, ae.Status())
	assert.Equal(t, "ise", ae.Message())
	assert.Equal(t, fmt.Sprintf("%d %s: %s", http.StatusInternalServerError, status.INTERNAL_SERVER_ERROR, "ise"), ae.Error())

	t.Run("should destruct error", func(t *testing.T) {
		destructed := errors.Destruct(ae)
		assert.Equal(t, ae.HTTPStatusCode(), destructed.HTTPStatusCode())
	})

	t.Run("should return nil error", func(t *testing.T) {
		destructed := errors.Destruct(nil)
		assert.Nil(t, destructed)
	})

	t.Run("should return Error type with status internal server", func(t *testing.T) {
		err := fmt.Errorf("abcd")

		destructed := errors.Destruct(err)
		assert.Equal(t, http.StatusInternalServerError, destructed.HTTPStatusCode())
	})

	t.Run("should match", func(t *testing.T) {
		match := errors.MatchStatus(errors.NotFound, status.NOT_FOUND)

		assert.True(t, match)
	})

	t.Run("should not match", func(t *testing.T) {
		match := errors.MatchStatus(errors.NotFound, status.INTERNAL_SERVER_ERROR)

		assert.False(t, match)

		match = errors.MatchStatus(nil, status.NOT_FOUND)

		assert.False(t, match)

		match = errors.MatchStatus(fmt.Errorf("asdf"), status.NOT_FOUND)

		assert.False(t, match)
	})
}
