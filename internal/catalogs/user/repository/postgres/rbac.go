package postgres

import (
	"context"

	oauthTypes "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	userTypes "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
)

func (p *Postgres) GetUserRoles(ctx context.Context, userId userTypes.UserId, clientId oauthTypes.ClientId) ([]string, error) {
	var roles []string
	err := p.WithContext(ctx).
		Table(`edp."user" u`).Distinct("r.code").
		Joins("JOIN edp.user_role ur ON ur.user_id = u.id").
		Joins(`JOIN edp."role" r ON r.id = ur.role_id AND r.is_deleted = false`).
		Joins("JOIN edp.client_role cr ON cr.role_id = r.id AND cr.is_deleted = false").
		Where("u.id = ? AND cr.client_id = ? AND u.is_deleted = false", userId, clientId).
		Order("r.code").Pluck("r.code", &roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (p *Postgres) GetUserPermissionsForClient(ctx context.Context, userId userTypes.UserId, clientId oauthTypes.ClientId) ([]string, error) {
	var permissions []string

	err := p.WithContext(ctx).
		Table(`edp."user" AS u`).Distinct("p.code").
		Joins("JOIN edp.user_role AS ur ON ur.user_id = u.id").
		Joins(`JOIN edp."role" AS r ON r.id = ur.role_id AND r.is_deleted = false`).
		Joins("JOIN edp.client_role AS cr ON cr.role_id = r.id AND cr.is_deleted = false").
		Joins("JOIN edp.permission_role AS pr ON pr.role_id = r.id").
		Joins(`JOIN edp."permission" AS p ON p.id = pr.permission_id AND p.is_deleted = false`).
		Where("u.id = ? AND cr.client_id = ? AND u.is_deleted = false", userId, clientId).
		Order("p.code").Pluck("p.code", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
