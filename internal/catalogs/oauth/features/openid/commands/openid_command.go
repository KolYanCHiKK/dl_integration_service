package commands

import (
	"context"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
)

type OpenIDCommand struct {
	Postgres *postgres.Postgres
}

func NewOpenIDCommand(postgres *postgres.Postgres) *OpenIDCommand {
	return &OpenIDCommand{
		Postgres: postgres,
	}
}

func (c *OpenIDCommand) GetOpenIDConfigurator(ctx context.Context) (*mappings.GetJWKParameters, error) {
	return nil, nil
}
