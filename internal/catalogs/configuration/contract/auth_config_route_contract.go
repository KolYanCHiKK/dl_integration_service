package contract

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/features/jwkcreate/commands"
	createjwksv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/features/jwkcreate/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func InitConfigurationHandlers(mux *chi.Mux, database *postgres2.Db, authConfig *configs.AuthJwtConfig, cipher *ecdsa.CipherEcdsa) *logger.DefaultErrorLogsParams {
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
