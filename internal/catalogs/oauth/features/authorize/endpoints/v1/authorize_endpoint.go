package v1

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/app"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/google/uuid"
	"github.com/ory/fosite"
)

type AuthorizeEndpoint struct {
	*AuthorizeHandler
}

func NewAuthorizeEndpoint(params *AuthorizeHandler) interfaces.Endpoint {
	return &AuthorizeEndpoint{params}
}

func (ep *AuthorizeEndpoint) MapEndpoint() {
	ep.Mux.Get("/oauth2/authorize", ep.handle())
}

func (ep *AuthorizeEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authReq, err := ep.Dependents.OAuth.NewAuthorizeRequest(r.Context(), r)
		if err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка валидации тела /authorize запроса в fosite")
			ep.Dependents.OAuth.WriteAuthorizeError(r.Context(), w, authReq, err)
			return
		}

		// Проверяем передали ли скоупы и ресурcы
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

		// Если не прислали в cookie идентификатор сессии, то редиректим на окно авторизации
		sessionCookie, err := r.Cookie("session")
		if err != nil || sessionCookie.Value == "" {
			ep.redirectToLoginForm(w, r, authReq)
			return
		}

		// Если cookie передали, то получаем параметры сессии по sessionId
		sessionId, err := uuid.Parse(sessionCookie.Value)
		if err != nil {
			ep.redirectToLoginForm(w, r, authReq) //
			return
		}

		sessionParams, err := ep.Dependents.Command.GetDefaultSessionParameters(r.Context(), sessionId)
		if errors.Is(err, error2.LoginSessionNotFound) {
			ep.redirectToLoginForm(w, r, authReq)
			return
		}
		if err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка получении параметров сессии авторизации пользователя")
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

		// Добавим параметры в сессию
		session := fosite.DefaultSession{
			Username: sessionParams.Username,
			Subject:  sessionParams.Subject.String(),
			Extra:    map[string]any{app.LoginSessionExtraKey: sessionParams.SessionID},
		}

		authResp, err := ep.Dependents.OAuth.NewAuthorizeResponse(r.Context(), authReq, &session)
		if err != nil {
			ep.Dependents.OAuth.WriteAuthorizeError(r.Context(), w, authReq, err)
			return
		}

		ep.Dependents.OAuth.WriteAuthorizeResponse(r.Context(), w, authReq, authResp)
	}
}

func (ep *AuthorizeEndpoint) redirectToLoginForm(w http.ResponseWriter, r *http.Request, authReq fosite.AuthorizeRequester) {
	// Создадим запрос на авторизацию
	loginRequestId, err := ep.Dependents.Command.CreateLoginRequest(r.Context(), r.URL.RequestURI(), ep.Dependents.AuthConfig.LoginRequestTTL)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка сохрания запроса на авторизацию пользователя (loginRequestId)")
		ep.Dependents.OAuth.WriteAuthorizeError(r.Context(), w, authReq, err)
		return
	}

	// Удалим cookie с браузера
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Соберем редирект
	query := url.Values{}
	query.Set("client_id", authReq.GetClient().GetID())
	query.Set("request", loginRequestId.String())

	w.Header().Set("Cache-Control", "no-store")
	http.Redirect(w, r, "/oauth2/login?"+query.Encode(), http.StatusSeeOther)
}
