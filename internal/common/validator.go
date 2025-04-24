package common

import (
	"fmt"

	"github.com/go-playground/locales/vi"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vitrans "github.com/go-playground/validator/v10/translations/vi"
)

var (
	Validator  *validator.Validate
	Translator ut.Translator
)

func InitValidator() error {
	vn := vi.New()
	uni := ut.New(vn, vn)
	trans, found := uni.GetTranslator("vi")
	if !found {
		return fmt.Errorf("translator not found")
	}
	v := validator.New()

	if err := vitrans.RegisterDefaultTranslations(v, trans); err != nil {
		return err
	}

	Validator = v
	Translator = trans
	return nil
}
