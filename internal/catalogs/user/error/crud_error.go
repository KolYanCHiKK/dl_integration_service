package error

import (
	"errors"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

var ValidationCreateUserErrorsMap = map[string]string{
	"Surname.max":        "Фамилия не может быть длиннее 100 символов",
	"Name.required":      "Отсутствует имя пользователя",
	"Name.max":           "Имя не может быть длиннее 100 символов",
	"LastName.max":       "Отчество не может быть длиннее 100 символов",
	"Phone.max":          "Номер телефона не может быть длиннее 65 символов",
	"Phone.phone":        "Неверный формат номера телефона",
	"BirthDate.datetime": "Неверный формат даты рождения. Ожидается ГГГГ-ММ-ДД",
	"Login.required":     "Отсутствует логин",
	"Login.max":          "Логин не может быть длиннее 120 символов",
	"Email.required":     "Отсутствует адрес электронной почты",
	"Email.max":          "Адрес электронной почты не может быть длиннее 120 символов",
	"Email.email":        "Неверный формат адреса электронной почты",
	"Password.required":  "Отсутствует пароль",
	"Password.max":       "Пароль не может быть длиннее 72 символов",
	"DepartmentId.uuid":  "Неверный формат идентификатора подразделения",
	"HigherId.uuid":      "Неверный формат идентификатора руководителя",
}

var (
	LoginAlreadyExists    = errors.New("login already exists")
	EmailAlreadyExists    = errors.New("email already exists")
	PasswordIsTooLong     = errors.New("password is too long")
	HigherIsNotExists     = errors.New("higher is not exists")
	DepartmentIsNotExists = errors.New("department is not exists")
)

func MapUserCrudErrors(err error) error {
	switch {
	case errors.Is(err, LoginAlreadyExists):
		logger.CreateErrLog(err.Error(), "Пользователь с таким логином уже существует")
		return LoginAlreadyExists

	case errors.Is(err, EmailAlreadyExists):
		logger.CreateErrLog(err.Error(), "Пользователь с таким email уже существует")
		return EmailAlreadyExists

	case errors.Is(err, PasswordIsTooLong):
		logger.CreateErrLog(err.Error(), "Длина пароля превышает 72 байта")
		return PasswordIsTooLong

	case errors.Is(err, HigherIsNotExists):
		logger.CreateErrLog(err.Error(), "Не был найден руководитель в таблице пользователей")
		return HigherIsNotExists

	case errors.Is(err, DepartmentIsNotExists):
		logger.CreateErrLog(err.Error(), "Не был найден департамент/отдел в справочнике подразделений")
		return HigherIsNotExists

	default:
		logger.CreateErrLog(err.Error(), "Ошибка создания пользователя")
		return err
	}
}
