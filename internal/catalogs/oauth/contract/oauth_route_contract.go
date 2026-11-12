package contract

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/commands"
	authv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/endpoints/v1"
	commands4 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/jwk/commands"
	jwkv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/jwk/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/loadkey/cache"
	commands6 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/loadkey/commands"
	commands3 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/login/commands"
	loginv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/login/endpoints/v1"
	commands5 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/openid/commands"
	openidv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/openid/endpoints/v1"
	commands2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/token/commands"
	tokenv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/token/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/fosite/provider"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func InitOauthHandlers(mux *chi.Mux, database *postgres2.Db, fositeConfig *configs.FositeAuthConfig, authConfig *configs.AuthJwtConfig,
	sessionCipher *aead.Cipher, jwkCipher *ecdsa.CipherEcdsa, rbacRepo commands2.UserRBACRepository) (error, string) {

	// Создаем репозиторий
	oauthPostgresRepository := postgres.NewPostgres(database)

	// Создаем комманды
	oauthCommand := commands.NewAuthorizeCommand(oauthPostgresRepository, sessionCipher)
	tokenCommand := commands2.NewTokenCommand(oauthPostgresRepository, oauthCommand, rbacRepo)
	loginCommand := commands3.LoginCommand{}
	jwkCommand := commands4.NewJWKCommand(oauthPostgresRepository)
	openIDCommand := commands5.NewOpenIDCommand(oauthPostgresRepository)

	// Подготоавливаем кеш для SigningKeyProvider
	loadKeyCommand, err := commands6.PrepareSigningKeyProviderCache(
		&cache.JWKCacheDependents{
			Postgres:    oauthPostgresRepository,
			CipherEcdsa: jwkCipher,
		},
		authConfig,
	)

	if err != nil {
		return err, "Ошибка загрузка данных в кеш для SigningKeyProvider"
	}

	// Инициализация fosite
	providerFosite, err := provider.NewProvider(
		&provider.FositeStore{
			AuthorizeCommand: oauthCommand,
			TokenCommand:     tokenCommand,
		},
		fositeConfig,
		loadKeyCommand,
	)

	if err != nil {
		info := "Ошибка создания провайдера fosite"
		logger2.CreateErrLog(err.Error(), info)
		return err, info
	}

	// Инициализация хендлеров
	authv1.NewAuthorizeHandler(mux, &authv1.AuthorizeHandlerDeps{
		Command:    oauthCommand,
		OAuth:      providerFosite,
		AuthConfig: authConfig,
		Cipher:     sessionCipher,
	})
	tokenv1.NewTokenHandler(mux, &tokenv1.TokenHandlerDeps{
		Command:     tokenCommand,
		OAuth:       providerFosite,
		OauthConfig: fositeConfig,
	})

	loginv1.NewLoginHandler(mux, &loginv1.LoginHandlerDeps{
		Command: loginCommand,
	})

	jwkv1.NewJWKHandlerHandler(mux,
		&jwkv1.JWKHandlerDeps{Command: jwkCommand},
	)

	openidv1.NewJOpenIDHandler(mux, &openidv1.OpenIDHandlerDeps{
		Command:      openIDCommand,
		FositeConfig: fositeConfig,
	})

	return nil, ""
}
