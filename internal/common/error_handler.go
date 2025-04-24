package common

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

func HandleValidationError(err error) []string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var translated []string
		for _, fe := range ve {
			translated = append(translated, fe.Translate(Translator))
		}
		return translated
	}
	return []string{"Dữ liệu không hợp lệ"}
}
