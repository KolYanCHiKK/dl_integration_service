package validation

import (
	"errors"
	"fmt"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/go-playground/validator/v10"
)

var ValidateError = errors.New("validate is failed")

type RestRequesterValidator struct {
	validate *validator.Validate
}

func NewRestRequesterValidator(functionTags []string) (*RestRequesterValidator, *logger.DefaultErrorLogsParams) {
	validate := validator.New()

	for _, tag := range functionTags {
		funcByTag, ok := validationFuncMap[tag]
		if !ok {
			return nil, &logger.DefaultErrorLogsParams{
				Error: ValidateError,
				Info:  fmt.Sprintf("Не найдена функция валидации для %s", tag),
			}
		}

		err := validate.RegisterValidation(tag, funcByTag)
		if err != nil {
			return nil, &logger.DefaultErrorLogsParams{
				Error: err,
				Info:  "Ошибка добавления функции валидации",
			}
		}
	}

	restValidator := &RestRequesterValidator{validate: validate}

	return restValidator, nil
}

func (v *RestRequesterValidator) ValidateBody(requestBody any, errorValidateMap map[string]string) ([]string, error) {
	var validateErr []string
	err := v.validate.Struct(requestBody)
	if err == nil {
		return nil, nil
	}

	var errSlice validator.ValidationErrors
	// Проверяем ошибку самого валидатора
	if !errors.As(err, &errSlice) {
		return nil, err
	}

	for _, value := range errSlice {
		key := value.Field() + "." + value.Tag()
		errMess, ok := errorValidateMap[key]
		if !ok {
			errMess = fmt.Sprintf("Поле %s не прошло проверку %s", value.Field(), value.Tag())
		}
		validateErr = append(validateErr, errMess)
	}
	return validateErr, nil
}
