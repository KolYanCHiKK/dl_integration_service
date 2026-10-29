package validation

import (
	"math"
	"regexp"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type ValidationFuncMapType map[string]func(f1 validator.FieldLevel) bool

var validationFuncMap = ValidationFuncMapType{
	"scale": scaleValidate,
	"phone": phoneNumberValidate,
}

var AllValidationTags = []string{
	"scale", "phone",
}

func scaleValidate(f1 validator.FieldLevel) bool {
	price := f1.Field().Float()
	scaled := price * 100
	if math.Abs(scaled-math.Round(scaled)) > 0.00001 {
		return false
	}
	return true
}

func phoneNumberValidate(f1 validator.FieldLevel) bool {
	phoneNumberRegular, err := regexp.Compile(
		`^(8|\+7)[0-9]{7,10}$`,
	)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка валидации регулярного выражения")
		return false
	}

	isValid := phoneNumberRegular.MatchString(f1.Field().String())
	if !isValid {
		return false
	}

	return true
}
