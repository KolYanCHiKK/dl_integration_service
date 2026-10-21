package provider

import (
	"errors"
	"time"

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

func NewProvider(store *FositeStore, accessIssuer string, secret []byte, tokenUrl string, scopeStrategy func([]string, string) bool,
	authCodeLifespan time.Duration, accessTokenLifespan time.Duration, refreshTokenLifespan time.Duration, enforcePKCE bool) (fosite.OAuth2Provider, error) {
	if len(secret) < 32 {
		err := errors.New("global secret must be at least 32 bytes")
		logger2.CreateErrLog(err.Error(), "Ошибка создания провайдера fosite")
		return nil, err
	}

	cfg := &fosite.Config{
		GlobalSecret:          secret,
		AccessTokenIssuer:     accessIssuer,
		TokenURL:              tokenUrl,
		ScopeStrategy:         scopeStrategy,
		AuthorizeCodeLifespan: authCodeLifespan,
		AccessTokenLifespan:   accessTokenLifespan,
		RefreshTokenLifespan:  refreshTokenLifespan,
		EnforcePKCE:           enforcePKCE,
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
