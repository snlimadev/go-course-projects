package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("notblank", notBlank)
}

var notBlank validator.Func = func(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return len(strings.TrimSpace(field)) > 0
}
