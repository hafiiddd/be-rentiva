package validatorpkg

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func New() *CustomValidator {
	v := validator.New()
	return &CustomValidator{Validator: v}
}
