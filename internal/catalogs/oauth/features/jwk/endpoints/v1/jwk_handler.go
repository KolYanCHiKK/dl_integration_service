package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/authorize/commands"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	"github.com/go-chi/chi/v5"
	"github.com/ory/fosite"
)

type JWKHandler struct {
	Mux        *chi.Mux
	Dependents *JWKHandlerDeps
}

type JWKHandlerDeps struct {
	Command     *commands.AuthorizeCommand
	OAuth       fosite.OAuth2Provider
	OauthConfig *configs.FositeAuthConfig
	AuthConfig  *configs.AuthJwtConfig
	Cipher      *aead.Cipher
}

func NewAuthorizeHandler(mux *chi.Mux, deps *JWKHandlerDeps) {
	handler := &JWKHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewAuthorizeEndpoint(handler).MapEndpoint()
}
