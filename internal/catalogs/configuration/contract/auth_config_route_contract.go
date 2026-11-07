package contract

import (
	"encoding/base64"
	"strings"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/features/jwkcreate/commands"
	createjwksv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/features/jwkcreate/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func InitConfigurationHandlers(mux *chi.Mux, database *postgres2.Db, authConfig *configs.AuthJwtConfig) *logger.DefaultErrorLogsParams {
	// Создаем объект шифрования для хранения сессии
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(authConfig.JWKEncKey))
	if err != nil {
		return &logger.DefaultErrorLogsParams{
			Error: err,
			Info:  "JWK_ENC_KEY не является корректным base64",
		}
	}

	cipher, err := ecdsa.NewJWKCipher(key)
	if err != nil {
		return &logger.DefaultErrorLogsParams{
			Error: err,
			Info:  "Ошибка создания объекта шифрования",
		}
	}

	// Создаем репозиторий
	configurationRepository := postgres.NewPostgres(database)

	// Создаем комманды
	createJWKCommand := commands.NewJWKCreateCommand(configurationRepository, cipher)

	// Инициализация хендлеров
	createjwksv1.NewJWKCreateHandler(mux, &createjwksv1.JWKCreateHandlerDeps{
		Command:    createJWKCommand,
		AuthConfig: authConfig,
	})

	return nil
}
