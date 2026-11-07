package postgres

import (
	"context"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"gorm.io/gorm"
)

func (p *Postgres) CreateJWKSigningKeys(ctx context.Context, tx *gorm.DB, jwkCreateParams *mappings.SigningKeyParams) (*models.SigningKey, error) {
	var signingKey = mappings.MapGeneratedSigningKeyToSigningModel(jwkCreateParams)

	err := tx.WithContext(ctx).Create(&signingKey).Error
	if err != nil {
		return nil, err
	}

	return &signingKey, nil

}

func (p *Postgres) RotateSigningKey(ctx context.Context, tx *gorm.DB, retireAfter time.Time) error {
	rotatedAt := time.Now().UTC()
	err := tx.WithContext(ctx).Where("status = ? AND is_deleted = false", oauth.SigningKeyStatusActive).
		Updates(models.SigningKey{
			Status:      oauth.SigningKeyStatusRetiring,
			RotatedAt:   &rotatedAt,
			RetireAfter: &retireAfter,
		}).Error

	if err != nil {
		return err
	}

	return nil
}
