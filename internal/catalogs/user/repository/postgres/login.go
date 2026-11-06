package postgres

import (
	"context"
	"errors"
	"time"

	error3 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	models2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/models"
	oauth2 "github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *Postgres) GetUserLoginParamsByLogin(ctx context.Context, login string) (*mappings.UserLoginParams, error) {
	var userLoginParams mappings.UserLoginParams
	err := p.WithContext(ctx).
		Select("id, login, hashed_password").Model(&models.User{}).
		Where("login = ? AND is_deleted = false", login).First(&userLoginParams).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error2.UserIsNotFound
	}
	if err != nil {
		return nil, err
	}

	return &userLoginParams, nil
}

func (p *Postgres) LoginUser(tx *gorm.DB, ctx context.Context, userId user.UserId, sessionParams *mappings.LoginRequestParams) (*user.LoginSessionId, string, error) {
	// Проставляем признак использования запроса и получаем редирект URL
	var loginRequest models.LoginRequest

	res := tx.WithContext(ctx).Model(&loginRequest).Clauses(clause.Returning{}).
		Where("id = ? AND used_at is NULL AND expires_at > ?", sessionParams.LoginRequest, time.Now()).
		Update("used_at", time.Now())

	if res.Error != nil {
		return nil, "", res.Error
	}
	if res.RowsAffected == 0 {
		return nil, "", error2.LoginRequestNotFound
	}

	// Создаем сессию
	session := models.LoginSession{
		UserID:    userId,
		AuthTime:  sessionParams.AuthTime,
		IP:        sessionParams.IP,
		UserAgent: sessionParams.UserAgent,
		ExpiresAt: sessionParams.ExpiresAt,
		RevokedAt: nil,
		IsDeleted: false,
	}

	err := tx.WithContext(ctx).Create(&session).Error
	if err != nil {
		return nil, "", err
	}

	return &session.ID, loginRequest.RedirectURL, nil
}

func (p *Postgres) IsValidClientID(ctx context.Context, clientID oauth.ClientId) error {
	var client models2.Client
	err := p.WithContext(ctx).Where("id = ? AND status = ? AND is_deleted = false", clientID, oauth2.ClientStatusActive).
		First(&client).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return error3.ClientNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetClientLogoutRedirectURI(ctx context.Context, clientID oauth.ClientId) ([]string, error) {
	var redirectUris []string
	err := p.WithContext(ctx).Model(&models2.ClientPostLogoutRedirectURI{}).
		Where("client_id = ? AND is_deleted = false", clientID).
		Pluck("uri", &redirectUris).Error

	if err != nil {
		return nil, err
	}

	return redirectUris, nil
}

func (p *Postgres) LogoutUser(ctx context.Context, sessionID user.LoginSessionId) error {
	result := p.WithContext(ctx).Model(&models.LoginSession{}).
		Where("id = ? AND revoked_at IS NULL AND is_deleted = false", sessionID).
		Update("revoked_at", time.Now().UTC())

	if result.Error != nil {
		return result.Error
	}

	return nil
}
