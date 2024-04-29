package validator

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	vld      *Validator
	syncOnce sync.Once
)

type Validator struct {
	v *validator.Validate
}

func construct() *Validator {
	validate := validator.New()

	return &Validator{
		v: validate,
	}
}

func Get() *Validator {
	syncOnce.Do(func() {
		vld = construct()
	})

	return vld
}

func (v Validator) ValidateStruct(ctx context.Context, s interface{}) error {
	err := v.v.StructCtx(ctx, s)
	if err == nil {
		return nil
	}

	errorFields := err.(validator.ValidationErrors)

	errMessages := make([]string, len(errorFields))

	for k, errorField := range errorFields {
		errMessages[k] = fmt.Sprintf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value())
	}

	errorMessage := strings.Join(errMessages, ", ")

	return fmt.Errorf(errorMessage)
}
