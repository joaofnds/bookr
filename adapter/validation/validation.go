package validation

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var Module = fx.Provide(func() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
})

func ErrorMessages(err validator.ValidationErrors) []string {
	var errorMessages []string
	for _, err := range err {
		errorMessages = append(errorMessages, err.Error())
	}
	return errorMessages
}
