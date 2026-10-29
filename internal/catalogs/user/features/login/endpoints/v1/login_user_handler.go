package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/commands"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/validation"
	"github.com/go-chi/chi/v5"
)

type LoginUserHandler struct {
	Mux        *chi.Mux
	Dependents *LoginUserHandlerDeps
}

type LoginUserHandlerDeps struct {
	Command    *commands.LoginUserCommand
	Validator  *validation.RestRequesterValidator
	AuthConfig *configs.AuthJwtConfig
}

func NewCreateUserHandler(mux *chi.Mux, deps *LoginUserHandlerDeps) {
	handler := &LoginUserHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewLoginUserEndpoint(handler).MapEndpoint()
}
