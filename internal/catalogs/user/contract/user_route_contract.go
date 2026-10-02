package contract

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/commands"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func InitUserHandlers(mux *chi.Mux, database *postgres2.Db) {
	userPostgresRepository := postgres.NewPostgres(database)

	v1.NewCreateUserHandler(mux, &v1.LoginUserHandlerDeps{
		Command: &commands.LoginUserCommand{
			Postgres: userPostgresRepository,
		},
	})
}
