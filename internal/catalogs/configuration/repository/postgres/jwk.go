package postgres

import (
	"context"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
)

func (p *Postgres) CreateJWKSigningKeys(ctx context.Context, jwkCreateParams *mappings.SigningKeyParams) (*models.SigningKey, error) {
	var signingKey = mappings.MapGeneratedSigningKeyToSigningModel(jwkCreateParams)

	err := p.WithContext(ctx).Create(&signingKey).Error
	if err != nil {
		return nil, err
	}

	return &signingKey, nil

}
