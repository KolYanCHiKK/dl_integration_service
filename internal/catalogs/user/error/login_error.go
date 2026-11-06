package error

import (
	"errors"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

var (
	UserIsNotFound                  = errors.New("user is not found")
	PasswordIsIncorrect             = errors.New("password is incorrect")
	SessionNotCreated               = errors.New("session was not created")
	ClientNoFound                   = error2.ClientNotFound
	ClientLogoutRedirectURINotFound = errors.New("logout redirect uri is not found")
)

var ValidationLoginErrorsMap = map[string]string{
	"Login.required":          "Отсутствует логин или пароль",
	"Login.min":               "Длина логина не может быть меньше 3 символов",
	"Password.required":       "Отсутствует логин или пароль",
	"Password.min":            "Длина пароля не может быть меньше 5 символов",
	"LoginRequestId.required": "Отсутствует идентификатор запроса на авторизацию",
	"LoginRequestId.uuid":     "Наверный тип данных идентификатора сессии",
	"Credentials.required":    "Отсутствуют параметры авторизации в запросе",
}

var (
	LoginRequestNotFound = errors.New("loginRequest is not found")
)

func MapUserLoginErrors(err error) error {
	switch {

	case errors.Is(err, UserIsNotFound):
		logger.CreateErrLog(err.Error(), "Пользователь не найден")
		return UserIsNotFound

	case errors.Is(err, PasswordIsIncorrect):
		logger.CreateErrLog(err.Error(), "Пароль неверный")
		return PasswordIsIncorrect

	case errors.Is(err, SessionNotCreated):
		logger.CreateErrLog(err.Error(), "Ошибка создания сессии авторизации пользователя")
		return SessionNotCreated

	case errors.Is(err, LoginRequestNotFound):
		logger.CreateErrLog(err.Error(), "Запрос на создание сессии не найден")
		return LoginRequestNotFound

	case errors.Is(err, ClientNoFound):
		logger.CreateErrLog(err.Error(), "Клиент не был найден")
		return ClientNoFound

	case errors.Is(err, ClientLogoutRedirectURINotFound):
		logger.CreateErrLog(err.Error(), "Ссылка для редиректа не была найдена")
		return ClientLogoutRedirectURINotFound

	default:
		logger.CreateErrLog(err.Error(), "Ошибка работы с данными о пользователе")
		return err
	}
}
