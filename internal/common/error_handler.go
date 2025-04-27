package common

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func HandleValidationError(err error) []ValidationError {
	var validationErrors []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			var message string

			switch e.Tag() {
			case "required":
				message = fmt.Sprintf("%s là bắt buộc", e.Field())
			case "email":
				message = fmt.Sprintf("%s không phải là email hợp lệ", e.Field())
			case "min":
				message = fmt.Sprintf("%s phải có ít nhất %s ký tự", e.Field(), e.Param())
			case "max":
				message = fmt.Sprintf("%s không được vượt quá %s ký tự", e.Field(), e.Param())
			case "len":
				message = fmt.Sprintf("%s phải có chính xác %s ký tự", e.Field(), e.Param())
			case "numeric":
				message = fmt.Sprintf("%s phải là số", e.Field())
			case "uuid4":
				message = fmt.Sprintf("%s phải là UUID phiên bản 4 hợp lệ", e.Field())
			case "oneof":
				message = fmt.Sprintf("%s phải có giá trị là: %s", e.Field(), e.Param())
			default:
				message = fmt.Sprintf("%s không hợp lệ", e.Field())
			}

			validationErrors = append(validationErrors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: message,
			})
		}
	}

	return validationErrors
}
