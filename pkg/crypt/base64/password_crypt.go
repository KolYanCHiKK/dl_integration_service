package base64

import (
	"errors"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var bcryptCostValue = 12

var ErrorPasswordIsTooLong = errors.New("password is long")

func CreateHashedPassword(password string) (string, error) {
	bytePassword := []byte(password)
	if len(bytePassword) > 72 {
		return "", ErrorPasswordIsTooLong
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCostValue)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashedPassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	switch {
	case err == nil:
		return true, nil

	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		logger.CreateErrLog(err.Error(), "Переданный хэш не совпадает с паролем")
		return false, bcrypt.ErrMismatchedHashAndPassword

	default:
		return false, err
	}
}
