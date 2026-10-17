package auth

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	contract2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/contract"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func StartOauthServer(mux *chi.Mux, database *postgres.Db, fositeConfig *configs.FositeAuthConfig) error {
	err, errInfo := contract2.InitOauthHandlers(mux, database, fositeConfig)
	if err != nil {
		logger2.CreateErrLog(err.Error(), errInfo)
		return err
	}

	return nil
}
