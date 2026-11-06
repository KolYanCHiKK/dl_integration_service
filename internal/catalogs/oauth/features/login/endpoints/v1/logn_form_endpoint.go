package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/notice/auth"
	"github.com/google/uuid"
)

type LoginEndpoint struct {
	*LoginHandler
}

func NewLoginEndpoint(params *LoginHandler) interfaces.Endpoint {
	return &LoginEndpoint{params}
}

func (ep *LoginEndpoint) MapEndpoint() {
	ep.Mux.Get("/oauth2/login", ep.handle())
}

func (ep *LoginEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		clientID := q.Get("client_id")
		loginRequestID := q.Get("request")

		if _, err := uuid.Parse(loginRequestID); err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка обработки loginRequestID")
			ep.Dependents.Command.RenderLogin(w, mappings.LoginPage{}, []auth.Notice{{
				Kind:  auth.NoticeError,
				Title: "Ошибка авторизации",
				Text:  "Ссылка для входа некорректна. Вернитесь в приложение и начните вход заново",
			}}, http.StatusBadRequest)
			return
		}

		if clientID == "" {
			ep.Dependents.Command.RenderLogin(w, mappings.LoginPage{}, []auth.Notice{{
				Kind:  auth.NoticeError,
				Title: "Ошибка авторизации",
				Text:  "Ссылка для входа некорректна. Вернитесь в приложение и начните вход заново",
			}}, http.StatusBadRequest)
			return
		}

		ep.Dependents.Command.RenderLogin(w, mappings.LoginPage{LoginRequestID: loginRequestID}, nil, http.StatusOK)
	}
}
