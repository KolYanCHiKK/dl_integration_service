package contract

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	commands2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/create/commands"
	createv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/create/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/commands"
	loginv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/endpoints/v1"
	commands3 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/logout/commands"
	logoutv1 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/logout/endpoints/v1"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/validation"
	"github.com/go-chi/chi/v5"
)

func InitUserHandlers(mux *chi.Mux, database *postgres2.Db, authConfig *configs.AuthJwtConfig, networkConfig *configs.NetworkConfig) *logger.DefaultErrorLogsParams {
	// Валидаторы
	validate, err := validation.NewRestRequesterValidator([]string{"phone"})
	if err != nil {
		return err
	}

	// Репозиторий
	userPostgresRepository := postgres.NewPostgres(database)

	// Слой комманд
	loginCommand := commands.NewLoginUserCommand(userPostgresRepository)
	createUserCommand := commands2.NewCreateUserCommand(userPostgresRepository)
	logoutUserCommand := commands3.NewLogoutUserCommand(userPostgresRepository)

	// Хендлеры
	loginv1.NewLoginUserHandler(mux, &loginv1.LoginUserHandlerDeps{
		Command:       loginCommand,
		Validator:     validate,
		AuthConfig:    authConfig,
		NetworkConfig: networkConfig,
	})
	createv1.NewCreateUserHandler(mux, &createv1.CreateUserHandlerDeps{
		Command:   createUserCommand,
		Validator: validate,
	})

	logoutv1.NewLogoutUserHandler(mux, &logoutv1.LogoutUserHandlerDeps{
		Command: logoutUserCommand,
	})

	return nil
}
