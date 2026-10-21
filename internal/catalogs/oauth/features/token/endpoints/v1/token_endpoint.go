package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/ory/fosite"
)

type TokenEndpoint struct {
	*TokenHandler
}

func NewTokenEndpoint(params *TokenHandler) interfaces.Endpoint {
	return &TokenEndpoint{params}
}

func (ep *TokenEndpoint) MapEndpoint() {
	ep.Mux.Post("/oauth2/token", ep.handle())
}

func (ep *TokenEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := &fosite.DefaultSession{}
		tokenReq, err := ep.Dependents.OAuth.NewAccessRequest(r.Context(), r, session)
		if err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка валидации тела /token запроса в fosite")
			ep.Dependents.OAuth.WriteAccessError(r.Context(), w, tokenReq, err)
			return
		}

		tokenResp, err := ep.Dependents.OAuth.NewAccessResponse(r.Context(), tokenReq)
		if err != nil {
			ep.Dependents.OAuth.WriteAccessError(r.Context(), w, tokenReq, err)
			return
		}

		ep.Dependents.OAuth.WriteAccessResponse(r.Context(), w, tokenReq, tokenResp)
	}
}
