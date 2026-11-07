package commands

import (
	"context"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
)

type JWKCommand struct {
	Postgres *postgres.Postgres
}

func NewJWKCommand(postgres *postgres.Postgres) *JWKCommand {
	return &JWKCommand{
		Postgres: postgres,
	}
}

func (c *JWKCommand) GetJWKSigningKey(ctx context.Context) (*mappings.GetJWKParameters, error) {
	jsonJWK, err := c.Postgres.GetJWKSigningKey(ctx)
	if err != nil {
		return nil, error2.MapWellKnowError(err)
	}

	return &mappings.GetJWKParameters{JsonJWK: mappings.MapJSONToSliceRawMessages(jsonJWK)}, nil
}
