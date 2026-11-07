package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/features/jwkcreate/commands"
	"github.com/go-chi/chi/v5"
)

type JWKCreateHandler struct {
	Mux        *chi.Mux
	Dependents *JWKCreateHandlerDeps
}

type JWKCreateHandlerDeps struct {
	Command    *commands.JWKCreateCommand
	AuthConfig *configs.AuthJwtConfig
}

func NewJWKCreateHandler(mux *chi.Mux, deps *JWKCreateHandlerDeps) {
	handler := &JWKCreateHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewJWKCreateEndpointEndpoint(handler).MapEndpoint()
}
