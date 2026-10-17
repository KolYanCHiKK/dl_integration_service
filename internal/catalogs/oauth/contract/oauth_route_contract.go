package contract

import (
	"encoding/base64"
	"strings"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/commands"
	v1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/fosite/provider"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func InitOauthHandlers(mux *chi.Mux, database *postgres2.Db, fositeConfig *configs.FositeAuthConfig) (error, string) {
	// Хэндер Authorize

	// Создаем объект шифрования для хранения сессии
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(fositeConfig.SessionEncKey))
	if err != nil {
		return err, "SESSION_ENC_KEY не является корректным base64"
	}

	cipher, err := aead.NewSessionCipher(key)
	if err != nil {
		return err, "Ошибка создания объекта шифрования"
	}

	// Создаем репозиторий и команду
	oauthPostgresRepository := postgres.NewPostgres(database)
	oauthCommand := commands.NewAuthorizeCommand(oauthPostgresRepository, cipher)

	// Инициализация fosite
	providerFosite, err := provider.NewProvider(
		&provider.FositeStore{AuthorizeCommand: oauthCommand},
		fositeConfig.AccessTokenIssuer,
		[]byte(fositeConfig.FositeGlobalSecret),
		fositeConfig.TokenURL,
		fositeConfig.ScopeStrategy,
		fositeConfig.AuthorizeCodeLifespan,
		fositeConfig.EnforcePKCE,
	)

	if err != nil {
		info := "Ошибка создания провайдера fosite"
		logger2.CreateErrLog(err.Error(), info)
		return err, info
	}

	// Инициализация хэндлера
	v1.NewCAuthorizeHandler(mux, &v1.AuthorizeHandlerDeps{
		Command: oauthCommand,
		OAuth:   providerFosite,
		Cipher:  cipher,
	})

	return nil, ""
}
