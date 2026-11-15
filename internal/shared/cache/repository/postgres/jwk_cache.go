package postgres

import (
	"context"
	"errors"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"gorm.io/gorm"
)

func (p *Postgres) GetActiveSigningKey(ctx context.Context) (*models.SigningKey, error) {
	var key models.SigningKey
	err := p.WithContext(ctx).
		Where("status = ? AND not_before <= now() AND is_deleted = false", oauth.SigningKeyStatusActive).
		Order("not_before DESC").
		First(&key).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error2.SigningKeysNotFound
	}
	if err != nil {
		return nil, err
	}

	return &key, nil
}
