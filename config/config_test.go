package config_test

import (
	"testing"

	"github.com/sangianpatrick/go-codebase-fiber/config"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	c1 := config.Get()
	c2 := config.Get()

	assert.NotNil(t, c1)
	assert.NotNil(t, c2)
	assert.Equal(t, c1, c2)
}
