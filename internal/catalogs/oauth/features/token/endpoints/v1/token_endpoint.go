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

// Выдача Access Token + замена на новый через Refresh
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

		// Проверяем сессию авторизации пользователя
		isLoginSessionValid, err := ep.Dependents.Command.CheckLoginSessionActiveByRequest(r.Context(), tokenReq.GetID())
		if isLoginSessionValid != true || err != nil {
			ep.createMismatchingTypeResponse(w, r, tokenReq, err)
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

		// При Refresh просходи клонирование, поэтому записываем именно в объект их fosite запроса
		sessionInRequest := tokenReq.GetSession().(*custom.Session)

		// Добавляем роли и пермишены в токен
		sessionInRequest.Roles, err = ep.Dependents.Command.GetUserRoles(r.Context(), userId, clientId)
		if err != nil {
			ep.createErrorGeneratingResponse(w, r, tokenReq, err)
			return
		}
		sessionInRequest.Permissions, err = ep.Dependents.Command.GetUserPermission(r.Context(), userId, clientId)
		if err != nil {
			ep.createErrorGeneratingResponse(w, r, tokenReq, err)
			return
		}

		sessionInRequest.ClientId = clientId

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
	ep.Dependents.OAuth.WriteAccessError(r.Context(), w, tokenReq,
		fosite.ErrInvalidGrant.
			WithDescription("The user data of this grant is invalid.").WithWrap(err),
	)
	return
}

func (ep *TokenEndpoint) createErrorGeneratingResponse(w http.ResponseWriter, r *http.Request, tokenReq fosite.AccessRequester, err error) {
	logger.CreateErrLog(err.Error(), "Ошибка создания /token ответа в fosite")
	ep.Dependents.OAuth.WriteAccessError(r.Context(), w, tokenReq, err)
	return
}
