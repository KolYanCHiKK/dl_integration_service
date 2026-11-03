package postgres

import (
	"context"
	"errors"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) IsLoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	err := p.WithContext(ctx).Model(&models.User{}).
		Select("count(*) > 0").
		Where("login = ?", login).
		Find(&exists).Error

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *Postgres) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := p.WithContext(ctx).Model(&models.User{}).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).Error

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *Postgres) CheckHigherId(ctx context.Context, higherId *user.UserId) (bool, error) {
	if higherId == nil {
		return true, nil
	}

	var exists bool
	err := p.WithContext(ctx).Model(&models.User{}).
		Select("count(*) > 0").Where("id = ? AND is_deleted = false", higherId).
		Find(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (p *Postgres) CheckDepartmentId(ctx context.Context, departmentId *user.DepartmentId) (bool, error) {
	if departmentId == nil {
		return true, nil
	}

	var exists bool
	err := p.WithContext(ctx).Model(&models.Department{}).
		Select("count(*) > 0").Where("id = ? AND is_deleted = false", departmentId).
		Find(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (p *Postgres) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := p.WithContext(ctx).Create(user).Error

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return nil, error2.LoginAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
