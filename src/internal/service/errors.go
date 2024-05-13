package service

import (
	"errors"

	"github.com/VincentSchmid/htmx-templ-app-template/internal/dto"
)

type ValidationError struct {
	Msg              string
	InputData        interface{}
	ValidationErrors dto.IValidationError
}

var _ dto.IValidationError = (*ValidationError)(nil)

func (e *ValidationError) Error() string {
	return e.Msg
}

type InternalServerError error

var (
	ErrInternal InternalServerError = errors.New("internal server error")
)
