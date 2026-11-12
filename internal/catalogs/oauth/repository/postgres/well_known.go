package postgres

import (
	"context"
	"errors"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	oauth2 "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

func (p *Postgres) GetAllScopes(ctx context.Context) ([]oauth2.ScopeName, error) {
	var scopeName []oauth2.ScopeName
	err := p.WithContext(ctx).Where("is_deleted = false").
		Model(&models.Scope{}).Order("name").Pluck("name", &scopeName).Error

	if err != nil {
		return nil, err
	}

	return scopeName, nil
}

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
