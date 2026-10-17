package postgres

import (
	"context"
	"errors"
	"time"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	oauth2 "github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Postgres struct {
	*postgres.Db
}

func NewPostgres(db *postgres.Db) *Postgres {
	return &Postgres{Db: db}
}

func (p *Postgres) StartTransaction() *gorm.DB {
	return p.DB
}

func (p *Postgres) GetClient(ctx context.Context, clientId oauth.ClientId) (*models.Client, error) {
	var client models.Client
	result := p.
		WithContext(ctx).
		Where("id = ? AND status <> ? AND is_deleted = false", clientId, oauth2.ClientStatusBlocked).
		First(&client)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, error2.ClientNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &client, nil
}

func (p *Postgres) GetClientGrantType(ctx context.Context, clientId oauth.ClientId) ([]string, error) {
	var clientGrantType []string
	err := p.
		WithContext(ctx).
		Model(&models.ClientGrantType{}).
		Where("client_id = ?", clientId).
		Pluck("grant_type", &clientGrantType).Error

	if err != nil {
		return nil, err
	}
	if len(clientGrantType) == 0 {
		return nil, error2.GrantTypeNotFound
	}

	return clientGrantType, nil
}

func (p *Postgres) GetClientResponseTypes(ctx context.Context, clientId oauth.ClientId) ([]string, error) {
	var out []string
	err := p.WithContext(ctx).
		Model(&models.ClientResponseType{}).
		Where("client_id = ? AND is_deleted = false", clientId).
		Pluck("response_type", &out).Error

	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, error2.ResponseTypeNotFound
	}
	return out, nil
}

func (p *Postgres) GetClientRedirectURIs(ctx context.Context, clientId oauth.ClientId) ([]string, error) {
	var out []string
	err := p.WithContext(ctx).
		Model(&models.ClientRedirectURI{}).
		Where("client_id = ? AND is_deleted = false", clientId).
		Pluck("uri", &out).Error
	return out, err
}

func (p *Postgres) GetClientScopes(ctx context.Context, clientId oauth.ClientId) ([]string, error) {
	var out []string
	err := p.WithContext(ctx).
		Table("edp.client_scope AS cs").
		Joins("JOIN edp.oauth_scope s ON s.id = cs.scope_id").
		Where("cs.client_id = ? AND s.is_deleted = false", clientId).
		Pluck("s.name", &out).Error

	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, error2.ScopeNotFound
	}

	return out, nil
}

func (p *Postgres) GetClientAudience(ctx context.Context, clientId oauth.ClientId) ([]string, error) {
	var out []string
	err := p.WithContext(ctx).
		Table("edp.client_resource AS cr").
		Joins("JOIN edp.oauth_resource r ON r.id = cr.resource_id").
		Where("cr.client_id = ? AND cr.is_deleted = false AND r.is_deleted = false", clientId).
		Pluck("r.name", &out).Error

	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, error2.ResourceNotFound
	}

	return out, nil
}

func (p *Postgres) GetClientRotatedSecrets(ctx context.Context, clientId oauth.ClientId) ([][]byte, error) {
	var hashes []string
	err := p.WithContext(ctx).
		Model(&models.ClientSecretRotation{}).
		Where("client_id = ? AND is_deleted = false AND expires_at > now()", clientId).
		Pluck("secret_hash", &hashes).Error
	if err != nil {
		return nil, err
	}
	out := make([][]byte, 0, len(hashes))
	for _, h := range hashes {
		out = append(out, []byte(h))
	}
	return out, nil
}

func (p *Postgres) GetOauthRequestIDByFositeID(tx *gorm.DB, ctx context.Context, fositeId string) (oauth.AuthRequestId, bool, error) {
	var req models.Request
	err := tx.WithContext(ctx).
		Select("id").
		Where("fosite_request_id = ?", fositeId).
		First(&req).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return oauth.AuthRequestId{}, false, nil // не нашли: это не ошибка
	}
	if err != nil {
		return oauth.AuthRequestId{}, false, err
	}
	return req.ID, true, nil
}

func (p *Postgres) CreateOauthRequest(tx *gorm.DB, ctx context.Context, fositeRequestId string, clientId string, subject string, requestedAt time.Time, form datatypes.JSON, session []byte) (*oauth.AuthRequestId, error) {
	request := models.Request{
		FositeRequestID: fositeRequestId,
		ClientID:        clientId,
		Subject:         subject,
		RequestedAt:     requestedAt,
		Form:            form,
		Session:         session,
	}
	err := tx.WithContext(ctx).Create(&request).Error
	if err != nil {
		return nil, err
	}

	return &request.ID, nil
}

func (p *Postgres) CreateRequestResources(tx *gorm.DB, ctx context.Context, requestId oauth.AuthRequestId, kind oauth.RequestResourceKind, resources []string) error {
	if len(resources) == 0 {
		return nil // нечего вставлять, пустой срез
	}

	rows := make([]models.RequestResource, 0, len(resources))
	for _, r := range resources {
		rows = append(rows, models.RequestResource{
			RequestID: requestId,
			Kind:      kind,
			Resource:  r,
		})
	}

	return tx.WithContext(ctx).Create(&rows).Error
}

func (p *Postgres) CreateRequestScopes(tx *gorm.DB, ctx context.Context, requestId oauth.AuthRequestId, kind oauth.RequestScopeKind, scopes []string) error {
	rows := make([]models.RequestScope, 0, len(scopes))
	for _, s := range scopes {
		rows = append(rows, models.RequestScope{RequestID: requestId, Kind: kind, Scope: s})
	}
	if len(rows) == 0 {
		return nil
	}
	err := tx.WithContext(ctx).Create(&rows).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreatePKCERequest(tx *gorm.DB, ctx context.Context, requestId oauth.AuthRequestId, codeChallenge string, challengeMethod string, expiresAt time.Time) error {
	pkceRequest := models.PKCERequest{
		RequestID:       requestId,
		CodeChallenge:   codeChallenge,
		ChallengeMethod: challengeMethod,
		ExpiresAt:       expiresAt,
	}

	err := tx.WithContext(ctx).Create(&pkceRequest).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateAuthorizeCode(tx *gorm.DB, ctx context.Context, requestId oauth.AuthRequestId, signature string, expiresAt time.Time) error {
	code := models.AuthorizeCode{
		Signature: signature,
		RequestID: requestId,
		ExpiresAt: expiresAt,
		Active:    true,
	}

	err := tx.WithContext(ctx).Create(&code).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetAuthorizeCode(ctx context.Context, signature string) (*models.AuthorizeCode, error) {
	var code models.AuthorizeCode

	err := p.WithContext(ctx).Where("signature = ?", signature).First(&code).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error2.AuthorizeCodeNotFound
	}
	if err != nil {
		return nil, err
	}

	return &code, nil
}

func (p *Postgres) GetPKCERequest(ctx context.Context, signature string) (*models.PKCERequest, error) {
	var pkce models.PKCERequest

	sub := p.WithContext(ctx).Where("signature = ?", signature).
		Model(&models.AuthorizeCode{}).Select("request_id")
	err := p.WithContext(ctx).Where("request_id = (?) AND is_deleted = false", sub).First(&pkce).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error2.PKCERequestNotFound
	}
	if err != nil {
		return nil, err
	}

	return &pkce, nil
}

func (p *Postgres) GetOauthRequest(ctx context.Context, id oauth.AuthRequestId) (*models.Request, error) {
	var request models.Request

	err := p.
		WithContext(ctx).Where("id = ? AND is_deleted = false", id).First(&request).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error2.OauthRequestNotFound
	}
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (p *Postgres) GetRequestScopes(ctx context.Context, requestId oauth.AuthRequestId, kind oauth.RequestScopeKind) ([]string, error) {
	var scopes []models.RequestScope

	err := p.WithContext(ctx).
		Where("request_id = ? AND kind = ? AND is_deleted = false", requestId, kind).
		Find(&scopes).Error

	if err != nil {
		return nil, err
	}

	var scopeNames []string
	for _, value := range scopes {
		scopeNames = append(scopeNames, value.Scope)
	}

	return scopeNames, nil // Отсуствие скоупов не ошибка по спецификации Oauth, но в хендлере проверяем
}

func (p *Postgres) GetRequestResources(ctx context.Context, requestId oauth.AuthRequestId, kind oauth.RequestResourceKind) ([]string, error) {
	var resources []models.RequestResource

	err := p.WithContext(ctx).
		Where("request_id = ? AND kind = ? AND is_deleted = false", requestId, kind).
		Find(&resources).
		Error

	if err != nil {
		return nil, err
	}

	var resourceNames []string
	for _, value := range resources {
		resourceNames = append(resourceNames, value.Resource)
	}

	return resourceNames, nil // Отсуствие ресурсов не ошибка по спецификации Oauth, но в хендлере проверяем
}

func (p *Postgres) DeactivateAuthorizeCode(ctx context.Context, signature string) error {
	res := p.
		WithContext(ctx).
		Model(&models.AuthorizeCode{}).
		Where("signature = ?", signature).
		Update("active", false)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return error2.AuthorizeCodeNotFound // такого кода нет
	}

	return nil
}

func (p *Postgres) DeletePKCERequest(ctx context.Context, signature string) error {
	sub := p.WithContext(ctx).Where("signature = ?", signature).
		Model(&models.AuthorizeCode{}).Select("request_id")
	err := p.WithContext(ctx).Where("request_id = (?)", sub).Delete(&models.PKCERequest{}).Error
	if err != nil {
		return err
	}

	return nil
}
