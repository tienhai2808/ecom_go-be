package common

import (
	"encoding/json"
	"errors"
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

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
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
		return validationErrors
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   strings.ToLower(unmarshalTypeError.Field),
			Message: fmt.Sprintf("%s phải là kiểu %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()),
		})
		return validationErrors
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "",
			Message: fmt.Sprintf("JSON không hợp lệ tại byte %d", syntaxError.Offset),
		})
		return validationErrors
	}

	if err != nil {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "",
			Message: err.Error(),
		})
	}

	return validationErrors
}

