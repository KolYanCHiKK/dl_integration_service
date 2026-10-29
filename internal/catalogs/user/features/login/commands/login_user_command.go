package commands

import (
	"context"
	"errors"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/base64"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/uow"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginUserCommand struct {
	Postgres *postgres.Postgres
}

func NewLoginUserCommand(postgres *postgres.Postgres) *LoginUserCommand {
	return &LoginUserCommand{Postgres: postgres}
}

func (c *LoginUserCommand) LoginUser(ctx context.Context, loginParams *mappings.LoginRequestParams, login string, password string) (*user.LoginSessionId, string, error) {
	userParams, err := c.Postgres.GetUserLoginParamsByLogin(ctx, login)
	if err != nil {
		return nil, "", error2.MapUserErrors(err)
	}

	_, err = base64.CompareHashedPassword(userParams.UserPassword, password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, "", error2.MapUserErrors(error2.PasswordIsIncorrect)
	}
	if err != nil {
		return nil, "", error2.MapUserErrors(err)
	}

	var sessionParams mappings.SessionParams
	err = uow.ExecuteTransaction(ctx, c.Postgres, func(tx *gorm.DB) error {
		sessionId, redirectUrl, err := c.Postgres.LoginUser(tx, ctx, userParams.UserId, loginParams)
		if err != nil {
			return err
		}
		sessionParams.SessionId = *sessionId
		sessionParams.RedirectURL = redirectUrl

		return nil
	})

	if err != nil {
		return nil, "", err
	}

	return &sessionParams.SessionId, sessionParams.RedirectURL, nil
}
