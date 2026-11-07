package commands

import (
	"context"
	"encoding/json"
	"time"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/uow"
	"gorm.io/gorm"
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

func (c *JWKCreateCommand) RotateJWKSigningKey(ctx context.Context, rotatedInterval time.Duration) (*mappings.SigningKeyDTOParams, error) {
	generatedKey, err := ecdsa.NewGenerateSigningKey(c.Cipher)
	if err != nil {
		return nil, error2.MapJWKErrors(error2.JwkIsNotCripted)
	}

	var signingParams *mappings.SigningKeyDTOParams

	err = uow.ExecuteTransaction(ctx, c.Postgres, func(tx *gorm.DB) error {
		err := c.Postgres.RotateSigningKey(ctx, tx, time.Now().Add(rotatedInterval))
		if err != nil {
			return err
		}

		signingKey, err := c.Postgres.CreateJWKSigningKeys(ctx, tx, &mappings.SigningKeyParams{GeneratedSigningKey: generatedKey})
		if err != nil {
			return err
		}

		signingParams = &mappings.SigningKeyDTOParams{
			PublicJWK: json.RawMessage(signingKey.PublicJWK),
			Status:    signingKey.Status,
			NotBefore: signingKey.NotBefore,
			CreatedAt: signingKey.CreatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, error2.MapJWKErrors(err)
	}

	return signingParams, nil
}
