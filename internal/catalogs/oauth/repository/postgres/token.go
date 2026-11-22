package postgres

import (
	"context"
	"errors"
	"time"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"gorm.io/gorm"
)

func (p *Postgres) CreateAccessToken(tx *gorm.DB, ctx context.Context, requestId oauth.AuthRequestId, signature string, expiresAt time.Time) error {
	var accessToken = models.AccessToken{
		Signature: signature,
		ExpiresAt: expiresAt,
		Active:    true,
		RequestID: requestId,
	}
	err := tx.WithContext(ctx).Create(&accessToken).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetAccessToken(ctx context.Context, signature string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	err := p.WithContext(ctx).Where("signature = ?", signature).First(&accessToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error2.AccessTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (p *Postgres) DeactivateAccessToken(tx *gorm.DB, ctx context.Context, signature string) error {
	res := tx.WithContext(ctx).Model(&models.AccessToken{}).
		Where("signature = ?", signature).Update("active", false)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return error2.AccessTokenNotFound
	}

	return nil
}

func (p *Postgres) RevokeAccessTokensByRequest(tx *gorm.DB, ctx context.Context, fositeRequestID string) error {
	requestIDsSub := tx.WithContext(ctx).Model(&models.Request{}).Select("id").
		Where("fosite_request_id = ?", fositeRequestID)

	res := p.WithContext(ctx).
		Model(&models.AccessToken{}).Where("request_id IN (?) AND active = true", requestIDsSub).
		Update("active", false)

	// Не проверяем RowAffected = 0, так как отзыв идемпотентный и если его уже отозвали, то все ок
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (p *Postgres) CreateRefreshToken(tx *gorm.DB, ctx context.Context, requestId oauth.AuthRequestId, signature, accessSignature string, expiresAt time.Time) error {
	var accessToken models.AccessToken
	err := tx.WithContext(ctx).Model(&models.AccessToken{}).
		Where("signature = ? AND active = true", accessSignature).
		Select("id").Take(&accessToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return error2.AccessTokenNotFound
	}
	if err != nil {
		return err
	}

	token := models.RefreshToken{
		Signature:     signature,
		ExpiresAt:     expiresAt,
		Active:        true,
		RequestID:     requestId,
		AccessTokenId: &accessToken.ID,
		UsedAt:        nil,
	}

	err = tx.WithContext(ctx).Create(&token).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetRefreshToken(ctx context.Context, signature string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := p.WithContext(ctx).Where("signature = ?", signature).
		Take(&refreshToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error2.RefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (p *Postgres) DeactivateRefreshToken(tx *gorm.DB, ctx context.Context, signature string) error {
	res := tx.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("signature = ? AND active = true", signature).Update("active", false)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return error2.RefreshTokenNotFound
	}
	return nil
}

func (p *Postgres) RevokeRefreshTokensByRequest(tx *gorm.DB, ctx context.Context, fositeRequestID string) error {
	requestIDs := tx.WithContext(ctx).Model(&models.Request{}).
		Select("id").Where("fosite_request_id = ?", fositeRequestID)

	res := tx.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("request_id IN (?) AND active = true", requestIDs).Update("active", false)

	// Не проверяем RowAffected = 0, так как отзыв идемпотентный и если его уже отозвали, то все ок
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// Замена токенов с помощью refresh
func (p *Postgres) RotateRefreshToken(tx *gorm.DB, ctx context.Context, fositeRequestID, signature string) error {
	requestIDs := tx.WithContext(ctx).Model(&models.Request{}).
		Select("id").Where("fosite_request_id = ?", fositeRequestID)

	// Отзываем access-токены запроса
	accessTokenRes := tx.WithContext(ctx).Model(&models.AccessToken{}).
		Where("request_id IN (?) AND active = true", requestIDs).Update("active", false)

	// Не проверяем RowAffected = 0, так как отзыв идемпотентный и если его уже отозвали, то все ок
	if accessTokenRes.Error != nil {
		return accessTokenRes.Error
	}

	// Помечаем refresh-токен использованным (used_at ставится один раз)
	res := tx.WithContext(ctx).
		Model(&models.RefreshToken{}).Where("signature = ? AND active = true", signature).
		Updates(map[string]any{
			"active":  false,
			"used_at": gorm.Expr("COALESCE(used_at, NOW())"),
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return error2.RefreshTokenNotFound
	}

	return nil
}

func (p *Postgres) IsLoginSessionActiveByRequest(ctx context.Context, fositeRequestId string) (bool, error) {
	var countLoginSessionRaws int64

	err := p.WithContext(ctx).Table("edp.oauth_request AS o").
		Joins("JOIN edp.login_session AS ls ON ls.id = o.login_session_id").
		Joins(`JOIN edp."user" AS u ON u.id = ls.user_id`).
		Where("o.fosite_request_id = ?", fositeRequestId).
		Where("ls.expires_at > now() AND ls.revoked_at IS NULL AND ls.is_deleted = false").
		Where("u.is_deleted = false").Count(&countLoginSessionRaws).Error

	if err != nil {
		return false, err
	}
	if countLoginSessionRaws == 0 {
		return false, error2.LoginSessionNotFound
	}

	return true, nil
}
