package v1

import (
	"errors"
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/response/rest"
)

type JWKEndpoint struct {
	*JWKHandler
}

func NewAuthorizeEndpoint(params *JWKHandler) interfaces.Endpoint {
	return &JWKEndpoint{params}
}

func (ep *JWKEndpoint) MapEndpoint() {
	ep.Mux.Get("/oauth2/well-known/jwk.json", ep.handle())
}

func (ep *JWKEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwkParams, err := ep.Dependents.Command.GetJWKSigningKey(r.Context())
		if errors.Is(err, error2.SigningKeysNotFound) {
			rest.GenerateBadRequestResponse(w, nil, "Не найдено активных конфигураций JWK", err.Error())
			return
		}
		if err != nil {
			rest.GenerateInternalServerErrorResponse(w, nil, err.Error())
			return
		}

		rest.BuildRestResponse(w, mappings.MapGetJWKParametersToDTO(jwkParams), 200, nil)
	}
}
