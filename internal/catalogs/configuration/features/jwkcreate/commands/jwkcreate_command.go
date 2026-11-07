package commands

import (
	"context"
	"encoding/json"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
)

type JWKCreateCommand struct {
	Postgres *postgres.Postgres
	Cipher   *ecdsa.CipherEcdsa
}

func NewJWKCreateCommand(postgres *postgres.Postgres, cipher *ecdsa.CipherEcdsa) *JWKCreateCommand {
	return &JWKCreateCommand{
		Postgres: postgres,
		Cipher:   cipher,
	}
}

func (c *JWKCreateCommand) CreateJWKSigningKey(ctx context.Context) (*mappings.SigningKeyDTOParams, error) {
	generatedKey, err := ecdsa.NewGenerateSigningKey(c.Cipher)
	if err != nil {
		return nil, error2.MapJWKErrors(error2.JwkIsNotCripted)
	}

	signingKey, err := c.Postgres.CreateJWKSigningKeys(ctx, &mappings.SigningKeyParams{GeneratedSigningKey: generatedKey})
	if err != nil {
		return nil, error2.MapJWKErrors(err)
	}

	return &mappings.SigningKeyDTOParams{
		PublicJWK: json.RawMessage(signingKey.PublicJWK),
		Status:    signingKey.Status,
		NotBefore: signingKey.NotBefore,
		CreatedAt: signingKey.CreatedAt,
	}, nil
}
