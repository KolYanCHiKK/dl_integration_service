package contracts

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/features/create/commands"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/features/create/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/repository/postgres"
	"github.com/go-chi/chi/v5"
)

func InitSourceHandlers(mux *chi.Mux, config *configs.ApplicationConfig) {
	sourcePostgresRepository := postgres.NewPostgres(config)

	v1.NewCreateSourceHandler(mux, &v1.CreateSourceHandlerDeps{
		Command: &commands.CreateSourceCommand{
			Postgres: sourcePostgresRepository,
		},
	})
}
