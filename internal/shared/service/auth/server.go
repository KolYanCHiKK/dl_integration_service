package auth

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/contract"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func StartIdentityServer(mux *chi.Mux, database *postgres.Db) error {
	contract.InitUserHandlers(mux, database)
	return nil
}
