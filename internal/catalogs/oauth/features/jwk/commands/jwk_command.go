package commands

import (
	"context"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	"github.com/ory/fosite"
)

type JWKCommand struct {
	Postgres *postgres.Postgres
	Cipher   *ecdsa.CipherEcdsa
}

func NewJWKCommand(postgres *postgres.Postgres, cipher *ecdsa.CipherEcdsa) *JWKCommand {
	return &JWKCommand{
		Postgres: postgres,
		Cipher:   cipher,
	}
}

func (c *JWKCommand) GetJWK(ctx context.Context, id string) (fosite.Client, error) {
	return nil, nil
}
