package auth

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	contract3 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/contract"
	contract2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/contract"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/contract"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func StartOauthServer(mux *chi.Mux, database *postgres.Db, applicationConfig *configs.ApplicationConfig) *logger2.DefaultErrorLogsParams {
	err, errInfo := contract2.InitOauthHandlers(mux, database, applicationConfig.OauthFosite, applicationConfig.JWTAuth)
	if err != nil {
		return &logger2.DefaultErrorLogsParams{Error: err, Info: errInfo}
	}

	errLoggerParams := contract.InitUserHandlers(mux, database, applicationConfig.JWTAuth, applicationConfig.Network)
	if errLoggerParams != nil {
		return errLoggerParams
	}

	errLoggerParams = contract3.InitConfigurationHandlers(mux, database, applicationConfig.JWTAuth)
	if errLoggerParams != nil {
		return errLoggerParams
	}

	return nil
}
