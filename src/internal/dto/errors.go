package dto

type IValidationError interface {
}

type ValidationError struct {
	Msg              string
	InputData        interface{}
	ValidationErrors IValidationError
}
