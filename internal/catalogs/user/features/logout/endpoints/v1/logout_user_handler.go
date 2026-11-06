package v1

import (
	commands2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/logout/commands"
	"github.com/go-chi/chi/v5"
)

type LogoutUserHandler struct {
	Mux        *chi.Mux
	Dependents *LogoutUserHandlerDeps
}

type LogoutUserHandlerDeps struct {
	Command *commands2.LogoutUserCommand
}

func NewLogoutUserHandler(mux *chi.Mux, deps *LogoutUserHandlerDeps) {
	handler := &LogoutUserHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewLogoutUserEndpoint(handler).MapEndpoint()
}
