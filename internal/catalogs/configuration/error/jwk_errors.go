package error

import (
	"errors"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

var (
	JwkIsNotCripted = errors.New("jwk crypted is failed")
)

func MapJWKErrors(err error) error {
	switch {

	case errors.Is(err, JwkIsNotCripted):
		logger.CreateErrLog(err.Error(), "Ошибка создания криптограции JWK")
		return JwkIsNotCripted

	default:
		logger.CreateErrLog(err.Error(), "Ошибка работы с данными JWK")
		return err
	}
}
