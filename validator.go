package http

import (
	"errors"
	"regexp"
	"sync"

	_validator "github.com/go-playground/validator/v10"
)

// Validator
/**
 * @description: validator instance
 */
var (
	validator *_validator.Validate
	once      sync.Once
)

// getValidator returns the global singleton validator instance
func getValidator() *_validator.Validate {
	once.Do(func() {
		// 启用 RequiredStructEnabled 以自动验证嵌套结构体
		validator = _validator.New(_validator.WithRequiredStructEnabled())
		err := validator.RegisterValidation("phone", validatePhone)
		if err != nil {
			panic("Unable to register validator, error: " + err.Error())
		}
	})
	return validator
}

// validatePhone
func validatePhone(f _validator.FieldLevel) bool {
	phone := f.Field().String()
	regx := `^1[3-9]\d{9}$`
	return regexp.MustCompile(regx).MatchString(phone)
}

// ValidateStruct
/**
 * @description: validate struct
 */
func ValidateStruct(s interface{}) (err error) {
	err = getValidator().Struct(s)
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
