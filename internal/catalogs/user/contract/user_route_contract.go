package contract

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/commands"
	loginv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/validation"
	"github.com/go-chi/chi/v5"
)

func InitUserHandlers(mux *chi.Mux, database *postgres2.Db) *logger.DefaultErrorLogsParams {
	// Валидаторы
	validate, err := validation.NewRestRequesterValidator([]string{"phone"})
	if err != nil {
		return err
	}

	// Репозиторий
	userPostgresRepository := postgres.NewPostgres(database)

	// Слой комманд
	loginCommand := commands.NewLoginUserCommand(userPostgresRepository)

	// Хендлеры
	loginv1.NewCreateUserHandler(mux, &loginv1.LoginUserHandlerDeps{
		Command:   loginCommand,
		Validator: validate,
	})

	return nil
}
