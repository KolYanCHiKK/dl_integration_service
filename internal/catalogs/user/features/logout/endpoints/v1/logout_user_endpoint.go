package v1

import (
	"errors"
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/response/rest"
	"github.com/google/uuid"
)

type LogoutUserEndpoint struct {
	*LogoutUserHandler
}

func NewLogoutUserEndpoint(params *LogoutUserHandler) interfaces.Endpoint {
	return &LogoutUserEndpoint{params}
}

func (ep *LogoutUserEndpoint) MapEndpoint() {
	ep.Mux.Post("/users/logout", ep.handle())
}

func (ep *LogoutUserEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

		clientId, err := uuid.Parse(r.PostFormValue("client_id"))
		if err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка получения clientId в запросе")
			rest.GenerateBadRequestResponse(w, nil, "Неверный идентификатор клиента", err.Error())
			return
		}

		sessionCookie, err := r.Cookie("session")
		if err != nil {
			ep.createErrSessionBody(w, err)
			return
		}

		sessionId, err := uuid.Parse(sessionCookie.Value)
		if err != nil {
			ep.createErrSessionBody(w, err)
			return
		}

		redirectUri, err := ep.Dependents.Command.LogoutUser(r.Context(), &mappings.LogoutUserParams{
			SessionId:   sessionId,
			ClientId:    clientId,
			RedirectUri: r.PostFormValue("redirect_uri"),
		})

		if errors.Is(err, error2.ClientLogoutRedirectURINotFound) {
			rest.GenerateBadRequestResponse(w, nil, "Невозможно выполнить редирект по данному URI", err.Error())
			return
		}
		if err != nil {
			rest.GenerateBadRequestResponse(w, nil, "Ошибка деактивации сессии пользователя", err.Error())
			return
		}

		if redirectUri == "" {
			rest.BuildRestResponse(w, nil, http.StatusNoContent, nil)
			return
		}

		http.Redirect(w, r, redirectUri, http.StatusSeeOther)
	}
}

func (ep *LogoutUserEndpoint) createErrSessionBody(w http.ResponseWriter, err error) {
	logger.CreateErrLog(err.Error(), "Ошибка получения sessionId в cookie")
	rest.GenerateBadRequestResponse(w, nil, "Неверный идентифкатор сессии для деактивации", err.Error())
}
