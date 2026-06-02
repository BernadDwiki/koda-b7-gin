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
	// fallback generic message for non-validator errors
	return "Data tidak valid"
}

func parseValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()

	// specific field messages (Indonesian)
	switch field {
	case "Amount":
		switch tag {
		case "required":
			return "Nominal top up wajib diisi"
		case "gt":
			return "Nominal top up harus lebih besar dari 0"
		case "gte":
			return "Minimal top up adalah Rp" + err.Param()
		case "min":
			return "Nominal top up harus minimal " + err.Param()
		default:
			return "Nominal top up tidak valid"
		}

	case "PaymentMethodID":
		return "Metode pembayaran wajib dipilih"
	}

	// generic tag-based messages (Indonesian)
	switch tag {
	case "required":
		return field + " wajib diisi"
	case "email":
		return field + " harus berupa email yang valid"
	case "min":
		return field + " harus minimal " + err.Param() + " karakter"
	case "max":
		return field + " harus maksimal " + err.Param() + " karakter"
	case "len":
		return field + " harus tepat " + err.Param() + " karakter"
	case "numeric":
		return field + " harus berupa angka"
	case "eqfield":
		return field + " harus sama dengan " + err.Param()
	case "unique":
		return field + " sudah ada"
	default:
		return "Data tidak valid"
	}
}
