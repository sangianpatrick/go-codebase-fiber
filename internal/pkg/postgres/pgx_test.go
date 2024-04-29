package postgres_test

import (
	"testing"

	"github.com/sangianpatrick/go-codebase-fiber/internal/pkg/postgres"
	"github.com/stretchr/testify/assert"
)

func TestGetDatabase(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USERNAME", "username")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_DB", "mydb")
	t.Setenv("POSTGRES_SCHEMA", "public")
	t.Setenv("POSTGRES_SSLMODE", "disable")
	t.Setenv("POSTGRES_POOL_MAX_CONNS", "1")

	db := postgres.GetDatabase()
	db2 := postgres.GetDatabase()
	assert.NotNil(t, db)
	assert.NotNil(t, db2)
	assert.Equal(t, db, db2, "should have the same memory address")
}
