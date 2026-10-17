package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/commands"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	"github.com/go-chi/chi/v5"
	"github.com/ory/fosite"
)

type AuthorizeHandler struct {
	Mux        *chi.Mux
	Dependents *AuthorizeHandlerDeps
}

type AuthorizeHandlerDeps struct {
	Command     *commands.AuthorizeCommand
	OAuth       fosite.OAuth2Provider
	OauthConfig *configs.FositeAuthConfig
	Cipher      *aead.Cipher
}

func NewCAuthorizeHandler(mux *chi.Mux, deps *AuthorizeHandlerDeps) {
	handler := &AuthorizeHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewAuthorizeEndpoint(handler).MapEndpoint()
}
