package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	commands2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/token/commands"
	"github.com/go-chi/chi/v5"
	"github.com/ory/fosite"
)

type RevokeHandler struct {
	Mux        *chi.Mux
	Dependents *RevokeHandlerDeps
}

type RevokeHandlerDeps struct {
	Command     *commands2.TokenCommand
	OAuth       fosite.OAuth2Provider
	OauthConfig *configs.FositeAuthConfig
}

func NewRevokeHandler(mux *chi.Mux, deps *RevokeHandlerDeps) {
	handler := &RevokeHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewRevokeEndpoint(handler).MapEndpoint()
}
