package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/ory/fosite"
)

type AuthorizeEndpoint struct {
	*AuthorizeHandler
}

func NewAuthorizeEndpoint(params *AuthorizeHandler) interfaces.Endpoint {
	return &AuthorizeEndpoint{params}
}

func (ep *AuthorizeEndpoint) MapEndpoint() {
	ep.Mux.Get("/authorize", ep.handle())
}

func (ep *AuthorizeEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authReq, err := ep.Dependents.OAuth.NewAuthorizeRequest(r.Context(), r)
		if err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка валидации тела /authorize запроса в fosite")
			ep.Dependents.OAuth.WriteAuthorizeError(r.Context(), w, authReq, err)
			return
		}

		// Проверяем передали ли скоупы и ресуры
		if len(authReq.GetRequestedScopes()) == 0 {
			err := fosite.ErrInvalidScope.WithHint("scope is required")
			logger.CreateErrLog(err.Error(), "Отсутсвуют scopes в запрос /authorize")
			ep.Dependents.OAuth.WriteAuthorizeError(r.Context(), w, authReq, err)
			return
		}
		if len(authReq.GetRequestedAudience()) == 0 {
			err := fosite.ErrInvalidRequest.WithHint("audience is required")
			logger.CreateErrLog(err.Error(), "Отсутсвуют audience в запрос /authorize")
			ep.Dependents.OAuth.WriteAuthorizeError(r.Context(), w, authReq, err)
			return
		}

		/* Добавляем все скоупы и ресурсы из запроса в выданные, так как если один неверный, то вернется
		   ошибка на этапе валидации */
		for _, value := range authReq.GetRequestedScopes() {
			authReq.GrantScope(value)
		}
		for _, value := range authReq.GetRequestedAudience() {
			authReq.GrantAudience(value)
		}

		// Временная заглушка с пользователем
		session := fosite.DefaultSession{
			Username: "TestUser",
			Subject:  "test_user",
		}

		authResp, err := ep.Dependents.OAuth.NewAuthorizeResponse(r.Context(), authReq, &session)
		if err != nil {
			ep.Dependents.OAuth.WriteAuthorizeError(r.Context(), w, authReq, err)
			return
		}

		ep.Dependents.OAuth.WriteAuthorizeResponse(r.Context(), w, authReq, authResp)
	}
}
