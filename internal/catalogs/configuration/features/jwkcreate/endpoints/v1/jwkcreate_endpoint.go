package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/response/rest"
)

type JWKCreateEndpoint struct {
	*JWKCreateHandler
}

func NewJWKCreateEndpointEndpoint(params *JWKCreateHandler) interfaces.Endpoint {
	return &JWKCreateEndpoint{params}
}

func (ep *JWKCreateEndpoint) MapEndpoint() {
	ep.Mux.Post("/oauth2/configuration/rotate-jwk", ep.handle())
}

func (ep *JWKCreateEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signingKeyParams, err := ep.Dependents.Command.RotateJWKSigningKey(r.Context(), ep.Dependents.AuthConfig.JWKRotationInterval)
		if err != nil {
			rest.GenerateInternalServerErrorResponse(w, nil, err.Error())
			return
		}

		rest.BuildRestResponse(w, mappings.MapSigningKeyDTOParamsToDTO(signingKeyParams), 201, nil)
	}
}
