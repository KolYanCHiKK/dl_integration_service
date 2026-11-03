package postgres

import (
	"context"
	"errors"
	"time"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/models"
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
