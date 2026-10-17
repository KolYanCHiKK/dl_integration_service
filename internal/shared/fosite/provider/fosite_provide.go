package provider

import (
	"context"
	"errors"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/commands"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
)

type FositeStore struct {
	*commands.AuthorizeCommand
}

func NewProvider(store *FositeStore, accessIssuer string, secret []byte, tokenUrl string, scopeStrategy func([]string, string) bool,
	authCodeLifespan time.Duration, enforcePKCE bool) (fosite.OAuth2Provider, error) {
	if len(secret) < 32 {
		err := errors.New("global secret must be at least 32 bytes")
		logger2.CreateErrLog(err.Error(), "Ошибка создания провайдера fosite")
		return nil, err
	}

	cfg := &fosite.Config{
		GlobalSecret:        secret,
		AccessTokenIssuer:   accessIssuer,
		TokenURL:            tokenUrl,
		ScopeStrategy:       scopeStrategy,
		AccessTokenLifespan: authCodeLifespan,
		EnforcePKCE:         enforcePKCE,
	}

	strategy := &compose.CommonStrategy{
		CoreStrategy: compose.NewOAuth2HMACStrategy(cfg), // Токены пока непрозрачные
	}

	return compose.Compose(
		cfg, store, strategy,
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2PKCEFactory,
	), nil
}

// Просто заглушки для запуска. Уберутся при создании метода /token

func (s *FositeStore) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {
	return fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return nil, fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	return fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) CreateRefreshTokenSession(ctx context.Context, signature string, accessSignature string, req fosite.Requester) error {
	return fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return nil, fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	return fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) error {
	return fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) RevokeRefreshToken(ctx context.Context, requestID string) error {
	return fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}

func (s *FositeStore) RevokeAccessToken(ctx context.Context, requestID string) error {
	return fosite.ErrServerError.WithDebug("access token storage is not implemented yet")
}
