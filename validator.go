package http

import (
	"errors"
	"regexp"
	"sync"

	_validator "github.com/go-playground/validator/v10"
)

var (
	validator  *_validator.Validate
	once       sync.Once
	phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
)

func getValidator() *_validator.Validate {
	once.Do(func() {
		validator = _validator.New(_validator.WithRequiredStructEnabled())
		err := validator.RegisterValidation("phone", validatePhone)
		if err != nil {
			panic("Unable to register validator, error: " + err.Error())
		}
	})
	return validator
}

func validatePhone(f _validator.FieldLevel) bool {
	return phoneRegex.MatchString(f.Field().String())
}

func ValidateStruct(s interface{}) error {
	err := getValidator().Struct(s)
	if err != nil {
		var validationErrors _validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				return e
			}
		}
	}
	return err
}
