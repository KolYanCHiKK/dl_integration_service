package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/library/custom"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/google/uuid"
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
		session := custom.NewEmptySession()

		tokenReq, err := ep.Dependents.OAuth.NewAccessRequest(r.Context(), r, session)
		if err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка валидации тела /token запроса в fosite")
			ep.Dependents.OAuth.WriteAccessError(r.Context(), w, tokenReq, err)
			return
		}

		userId, err := uuid.Parse(tokenReq.GetSession().GetSubject())
		if err != nil {
			ep.createMismatchingTypeResponse(w, r, tokenReq, err)
			return
		}
		clientId, err := uuid.Parse(tokenReq.GetClient().GetID())
		if err != nil {
			ep.createMismatchingTypeResponse(w, r, tokenReq, err)
			return
		}

		// Добавляем роли и пермишены в токен
		session.Roles, err = ep.Dependents.Command.GetUserRoles(r.Context(), userId, clientId)
		if err != nil {
			ep.createErrorGeneratingResponse(w, r, tokenReq, err)
			return
		}
		session.Permissions, err = ep.Dependents.Command.GetUserPermission(r.Context(), userId, clientId)
		if err != nil {
			ep.createErrorGeneratingResponse(w, r, tokenReq, err)
			return
		}

		session.ClientId = clientId

		// Формируем ответ
		tokenResp, err := ep.Dependents.OAuth.NewAccessResponse(r.Context(), tokenReq)
		if err != nil {
			ep.createErrorGeneratingResponse(w, r, tokenReq, err)
			return
		}

		ep.Dependents.OAuth.WriteAccessResponse(r.Context(), w, tokenReq, tokenResp)
	}
}

func (ep *TokenEndpoint) createMismatchingTypeResponse(w http.ResponseWriter, r *http.Request, tokenReq fosite.AccessRequester, err error) {
	logger.CreateErrLog(err.Error(), "Ошибка валидации пользовательских данных в сессии запроса")
	ep.Dependents.OAuth.WriteAccessError(r.Context(), w, tokenReq, err)
	return
}

func (ep *TokenEndpoint) createErrorGeneratingResponse(w http.ResponseWriter, r *http.Request, tokenReq fosite.AccessRequester, err error) {
	logger.CreateErrLog(err.Error(), "Ошибка создания /token ответа в fosite")
	ep.Dependents.OAuth.WriteAccessError(r.Context(), w, tokenReq, err)
	return
}
