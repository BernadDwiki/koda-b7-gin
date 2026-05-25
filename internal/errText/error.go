package errText

import (
	"github.com/go-playground/validator/v10"
)

func GetValidationErrorMessage(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		if len(validationErrors) > 0 {
			return parseValidationError(validationErrors[0])
		}
	}
	return err.Error()
}

func parseValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + err.Param() + " characters"
	case "max":
		return field + " must be at most " + err.Param() + " characters"
	case "len":
		return field + " must be exactly " + err.Param() + " characters"
	case "numeric":
		return field + " must contain only numbers"
	case "eqfield":
		return field + " must match " + err.Param()
	case "unique":
		return field + " already exists"
	default:
		return "invalid " + field
	}
}
