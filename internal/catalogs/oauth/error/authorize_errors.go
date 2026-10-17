package error

import (
	"context"
	"errors"

	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/ory/fosite"
)

var (
	ClientNotFound        = errors.New("client is not found")
	InvalidClientID       = errors.New("invalid clientId")
	GrantTypeNotFound     = errors.New("grantType is not found")
	ResponseTypeNotFound  = errors.New("responseType is not found")
	RedirectURINotFound   = errors.New("redirectURI is not found")
	ScopeNotFound         = errors.New("scope is not found")
	ResourceNotFound      = errors.New("resource is not found")
	AuthorizeCodeNotFound = errors.New("authorizeCode is not found")
	PKCERequestNotFound   = errors.New("PKCE request is not found")
	OauthRequestNotFound  = errors.New("oauthRequest is not found")
)

func MapClientLoadError(err error) error {
	switch {
	// Клиента нет или идентификатор не uuid
	case errors.Is(err, ClientNotFound), errors.Is(err, InvalidClientID):
		logger2.CreateErrLog(err.Error(), "[GetClient] Ошибка загрузки клиента")
		return fosite.ErrNotFound

	// Клиент есть, но настроен неполно: ошибка конфигурации в БД
	case
		errors.Is(err, GrantTypeNotFound),
		errors.Is(err, ResponseTypeNotFound),
		errors.Is(err, RedirectURINotFound),
		errors.Is(err, ScopeNotFound),
		errors.Is(err, ResourceNotFound):
		logger2.CreateErrLog(err.Error(), "[GetClient] Клиент настроен некорректно")
		return fosite.ErrNotFound

	case errors.Is(err, AuthorizeCodeNotFound):
		logger2.CreateErrLog(err.Error(), "Код авторизации не был найден")
		return fosite.ErrNotFound

	case errors.Is(err, PKCERequestNotFound):
		logger2.CreateErrLog(err.Error(), "Данные PKCE запроса не были найдены")
		return fosite.ErrNotFound

	case errors.Is(err, OauthRequestNotFound):
		logger2.CreateErrLog(err.Error(), "Данные Request запроса не были найдены")
		return fosite.ErrNotFound

	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err

	default:
		logger2.CreateErrLog(err.Error(), "[GetClient] Ошибка загрузки клиента")
		return err
	}
}

func MapSessionErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, AuthorizeCodeNotFound), errors.Is(err, PKCERequestNotFound), errors.Is(err, OauthRequestNotFound):
		logger2.CreateErrLog(err.Error(), "Ошбика поиска данных при создании fosite сессии")
		return fosite.ErrNotFound
	default:
		logger2.CreateErrLog(err.Error(), "Ошбика создания fosite сессии")
		return fosite.ErrServerError
	}
}
