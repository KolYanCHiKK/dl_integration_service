package error

import (
	"errors"

	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

var (
	SigningKeysNotFound = errors.New("signingKey is not found")
)

func MapWellKnowError(err error) error {
	switch {

	case err == nil:
		return nil

	case errors.Is(err, SigningKeysNotFound):
		logger2.CreateErrLog(err.Error(), "JWK отсутсвуют или нет активных записей")
		return SigningKeysNotFound

	default:
		logger2.CreateErrLog(err.Error(), "Ошибка работы с конфигурациями сервера авторизации")
		return err
	}
}
