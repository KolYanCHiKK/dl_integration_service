package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/jwk/commands"
	"github.com/go-chi/chi/v5"
)

type JWKHandler struct {
	Mux        *chi.Mux
	Dependents *JWKHandlerDeps
}

type JWKHandlerDeps struct {
	Command *commands.JWKCommand
}

func NewJWKHandlerHandler(mux *chi.Mux, deps *JWKHandlerDeps) {
	handler := &JWKHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewJWKEndpoint(handler).MapEndpoint()
}
