package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	commands2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/token/commands"
	"github.com/go-chi/chi/v5"
	"github.com/ory/fosite"
)

type TokenHandler struct {
	Mux        *chi.Mux
	Dependents *TokenHandlerDeps
}

type TokenHandlerDeps struct {
	Command     *commands2.TokenCommand
	OAuth       fosite.OAuth2Provider
	OauthConfig *configs.FositeAuthConfig
}

func NewTokenHandler(mux *chi.Mux, deps *TokenHandlerDeps) {
	handler := &TokenHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewTokenEndpoint(handler).MapEndpoint()
}
