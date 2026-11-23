package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

type RevokeEndpoint struct {
	*RevokeHandler
}

func NewRevokeEndpoint(params *RevokeHandler) interfaces.Endpoint {
	return &RevokeEndpoint{params}
}

func (ep *RevokeEndpoint) MapEndpoint() {
	ep.Mux.Post("/oauth2/revoke", ep.handle())
}

func (ep *RevokeEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := ep.Dependents.OAuth.NewRevocationRequest(r.Context(), r)
		if err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка отзыва токена в /revoke")
		}

		ep.Dependents.OAuth.WriteRevocationResponse(r.Context(), w, err)
	}
}
