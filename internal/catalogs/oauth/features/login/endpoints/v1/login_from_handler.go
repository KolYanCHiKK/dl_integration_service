package v1

import (
	commands2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/login/commands"
	"github.com/go-chi/chi/v5"
)

type LoginHandler struct {
	Mux        *chi.Mux
	Dependents *LoginHandlerDeps
}

type LoginHandlerDeps struct {
	Command commands2.LoginCommand
}

func NewLoginHandler(mux *chi.Mux, deps *LoginHandlerDeps) {
	handler := &LoginHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewLoginEndpoint(handler).MapEndpoint()
}
