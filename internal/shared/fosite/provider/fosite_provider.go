package provider

import (
	"context"
	"errors"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/commands"
	commands2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/token/commands"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
)

type FositeStore struct {
	*commands.AuthorizeCommand
	*commands2.TokenCommand
}

type SigningKeyProvider interface {
	GetSigningKey(ctx context.Context) (any, error)
}

func NewProvider(store *FositeStore, fositeConfig *configs.FositeAuthConfig, signingProvider SigningKeyProvider) (fosite.OAuth2Provider, error) {
	if len(fositeConfig.FositeGlobalSecret) < 32 {
		err := errors.New("global secret must be at least 32 bytes")
		logger2.CreateErrLog(err.Error(), "Ошибка создания провайдера fosite")
		return nil, err
	}

	cfg := &fosite.Config{
		GlobalSecret:                []byte(fositeConfig.FositeGlobalSecret),
		AccessTokenIssuer:           fositeConfig.AccessTokenIssuer,
		TokenURL:                    fositeConfig.TokenURL,
		ScopeStrategy:               fositeConfig.ScopeStrategy,
		AuthorizeCodeLifespan:       fositeConfig.AuthorizeCodeLifespan,
		AccessTokenLifespan:         fositeConfig.AccessTokenLifespan,
		RefreshTokenLifespan:        fositeConfig.RefreshTokenLifespan,
		EnforcePKCE:                 fositeConfig.EnforcePKCE,
		EnforcePKCEForPublicClients: fositeConfig.EnforcePKCEForPublicClients,
	}

	hmacStrategy := compose.NewOAuth2HMACStrategy(cfg)

	strategy := &compose.CommonStrategy{
		CoreStrategy: compose.NewOAuth2JWTStrategy(signingProvider.GetSigningKey, hmacStrategy, cfg),
	}

	return compose.Compose(
		cfg, store, strategy,
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2PKCEFactory,
		compose.OAuth2RefreshTokenGrantFactory,
	), nil
}
