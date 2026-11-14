package postgres

import (
	"context"
	"errors"
	"time"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/error"
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
			Status:      oauth.SigningKeyStatusRetired,
			RotatedAt:   &rotatedAt,
			RetireAfter: &retireAfter,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CheckActiveJWKSigningKey(ctx context.Context) error {
	var signingKeyActive models.SigningKey

	err := p.WithContext(ctx).
		Where("status = ? AND not_before <= ? AND is_deleted = false", oauth.SigningKeyStatusActive, time.Now().UTC()).
		First(&signingKeyActive).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return error2.JwkActiveIsNotFound
	}
	if err != nil {
		return err
	}

	return nil
}
