package error

import (
	"errors"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

var (
	UserIsNotFound      = errors.New("user is not found")
	PasswordIsIncorrect = errors.New("password is incorrect")
	SessionNotCreated   = errors.New("session was not created")
)

var ValidationLoginErrorsMap = map[string]string{
	"Login.required":    "Отсутствует логин или пароль",
	"Login.min":         "Длина логина не может быть меньше 3 символов",
	"Password.required": "Отсутствует логин или пароль",
	"Password.min":      "Длина пароля не может быть меньше 5 символов",
}

var (
	LoginRequestNotFound = errors.New("loginRequest is not found")
)

func MapUserErrors(err error) error {
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

	default:
		logger.CreateErrLog(err.Error(), "Ошибка работы с данными о пользователе")
		return err
	}
}
