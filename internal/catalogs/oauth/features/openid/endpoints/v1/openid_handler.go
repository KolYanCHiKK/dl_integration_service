package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/openid/commands"
	"github.com/go-chi/chi/v5"
)

type OpenIDHandler struct {
	Mux        *chi.Mux
	Dependents *OpenIDHandlerDeps
}

type OpenIDHandlerDeps struct {
	Command      *commands.OpenIDCommand
	FositeConfig *configs.FositeAuthConfig
}

func NewJOpenIDHandler(mux *chi.Mux, deps *OpenIDHandlerDeps) {
	handler := &OpenIDHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewOpenIDEndpoint(handler).MapEndpoint()
}
