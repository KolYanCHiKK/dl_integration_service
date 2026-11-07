package postgres

import (
	"context"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"gorm.io/datatypes"
)

func (p *Postgres) GetJWKSigningKey(ctx context.Context) ([]datatypes.JSON, error) {
	var signingKeysJSON []datatypes.JSON

	err := p.WithContext(ctx).Model(&models.SigningKey{}).
		Where("(status = ? OR (status = ? AND retire_after > now())) AND is_deleted = false",
			oauth.SigningKeyStatusActive, oauth.SigningKeyStatusRetiring).
		Order("not_before DESC").Pluck("public_jwk", &signingKeysJSON).Error

	if err != nil {
		return nil, err
	}
	if len(signingKeysJSON) == 0 {
		return nil, error2.SigningKeysNotFound
	}

	return signingKeysJSON, nil
}
