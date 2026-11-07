package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
)

type JWKEndpoint struct {
	*JWKHandler
}

func NewAuthorizeEndpoint(params *JWKHandler) interfaces.Endpoint {
	return &JWKEndpoint{params}
}

func (ep *JWKEndpoint) MapEndpoint() {
	ep.Mux.Get("/oauth2/well-known/jwk", ep.handle())
}

func (ep *JWKEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
