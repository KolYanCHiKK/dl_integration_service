package commands

import (
	"context"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"
)

type ProfileCommand struct {
	Postgres *postgres.Postgres
}

func NewProfileCommand(postgres *postgres.Postgres) *ProfileCommand {
	return &ProfileCommand{Postgres: postgres}
}

func (c *ProfileCommand) GetUserProfile(ctx context.Context, params *mappings.CreateUserParams) (*mappings.CreateUserResponseParams, error) {
	return nil, nil
}
