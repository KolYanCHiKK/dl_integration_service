package error

import (
	"errors"

	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/ory/fosite"
	log "github.com/sirupsen/logrus"
)

var (
	AccessTokenNotFound  = errors.New("access token is not found")
	RefreshTokenNotFound = errors.New("refresh token is not found")
)

func MapTokenError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, AccessTokenNotFound):
		logger2.CreateInfoLog(log.Fields{"error": err.Error()}, "Access Token не был найден")
		return fosite.ErrNotFound
	case errors.Is(err, RefreshTokenNotFound):
		logger2.CreateInfoLog(log.Fields{"error": err.Error()}, "Refresh Token не был найден")
		return fosite.ErrNotFound
	default:
		logger2.CreateErrLog(err.Error(), "Ошибка работы с Access или Refresh токеном")
		return fosite.ErrServerError
	}
}
