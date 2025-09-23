package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func HandleValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s là bắt buộc", strings.ToLower(e.Field()))
			case "email":
				return fmt.Sprintf("%s không phải là email hợp lệ", strings.ToLower(e.Field()))
			case "min":
				return fmt.Sprintf("%s phải có ít nhất %s ký tự", strings.ToLower(e.Field()), e.Param())
			case "max":
				return fmt.Sprintf("%s không được vượt quá %s ký tự", strings.ToLower(e.Field()), e.Param())
			case "len":
				return fmt.Sprintf("%s phải có chính xác %s ký tự", strings.ToLower(e.Field()), e.Param())
			case "numeric":
				return fmt.Sprintf("%s phải là số", strings.ToLower(e.Field()))
			case "uuid4":
				return fmt.Sprintf("%s phải là UUID phiên bản 4 hợp lệ", strings.ToLower(e.Field()))
			case "oneof":
				return fmt.Sprintf("%s phải có giá trị là: %s", strings.ToLower(e.Field()), e.Param())
			default:
				return fmt.Sprintf("%s không hợp lệ", strings.ToLower(e.Field()))
			}
		}
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return fmt.Sprintf("%s phải là kiểu %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String())
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return fmt.Sprintf("JSON không hợp lệ tại byte %d", syntaxError.Offset)
	}

	if err != nil {
		return err.Error()
	}

	return "Dữ liệu không hợp lệ"
}
