package commands

import (
	"context"

	error3 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"
)

type LogoutUserCommand struct {
	Postgres *postgres.Postgres
}

func NewLogoutUserCommand(postgres *postgres.Postgres) *LogoutUserCommand {
	return &LogoutUserCommand{Postgres: postgres}
}

func (c *LogoutUserCommand) LogoutUser(ctx context.Context, params *mappings.LogoutUserParams) (string, error) {
	err := c.Postgres.IsValidClientID(ctx, params.ClientId)
	if err != nil {
		return "", error3.MapUserLoginErrors(err)
	}

	err = c.Postgres.LogoutUser(ctx, params.SessionId)
	if err != nil {
		return "", error3.MapUserLoginErrors(err)
	}

	if params.RedirectUri == "" {
		return "", nil
	}

	redirectUris, err := c.Postgres.GetClientLogoutRedirectURI(ctx, params.ClientId)
	if err != nil {
		return "", error3.MapUserLoginErrors(err)
	}

	for _, value := range redirectUris {
		if value == params.RedirectUri {
			return value, nil
		}
	}

	return "", error3.MapUserLoginErrors(error3.ClientLogoutRedirectURINotFound)
}
