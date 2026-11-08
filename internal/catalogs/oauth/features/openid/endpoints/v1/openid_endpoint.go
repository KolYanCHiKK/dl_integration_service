package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/response/rest"
)

type OpenIDEndpoint struct {
	*OpenIDHandler
}

func NewOpenIDEndpoint(params *OpenIDHandler) interfaces.Endpoint {
	return &OpenIDEndpoint{params}
}

func (ep *OpenIDEndpoint) MapEndpoint() {
	ep.Mux.Get("/oauth2/well-known/openid-configuration", ep.handle())
}

func (ep *OpenIDEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		openIdParams, err := ep.Dependents.Command.GetOpenIDConfigurator(r.Context(), ep.Dependents.FositeConfig.AccessTokenIssuer)
		if err != nil {
			rest.GenerateInternalServerErrorResponse(w, nil, err.Error())
			return
		}

		rest.BuildRestResponse(w, mappings.MapOpenIDResponseParamsToDTO(openIdParams), 200, nil)
	}
}
