package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/login/dtos"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/body/rest/request"
	rest2 "github.com/KolYanCHiKK/dl_integration_service/pkg/response/rest"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/utils"
)

type LoginUserEndpoint struct {
	*LoginUserHandler
}

func NewLoginUserEndpoint(params *LoginUserHandler) interfaces.Endpoint {
	return &LoginUserEndpoint{params}
}

func (ep *LoginUserEndpoint) MapEndpoint() {
	ep.Mux.Post("/users/login", ep.handle())
}

func (ep *LoginUserEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.DecodeRestRequest[dtos.LoginUserRequest](w, r)
		if err != nil {
			rest2.GenerateBadRequestResponse(w, nil, err.Error())
			return
		}

		validationErrs, err := ep.Dependents.Validator.ValidateBody(body, error2.ValidationLoginErrorsMap)
		if err != nil {
			rest2.GenerateInternalServerErrorResponse(w, nil, err.Error())
			return
		}
		if validationErrs != nil {
			rest2.GenerateBadRequestResponse(w, nil, "", validationErrs...)
			return
		}

		sessionId, redirectURL, err := ep.Dependents.Command.LoginUser(
			r.Context(),
			&mappings.LoginRequestParams{
				LoginRequest: body.LoginRequestId,
				IP:           utils.GetClientIP(r, ep.Dependents.NetworkConfig.TrustedProxies),
				UserAgent:    r.Header.Get("User-Agent"),
				ExpiresAt:    time.Now().Add(ep.Dependents.AuthConfig.LoginSessionTTl),
			},
			body.Credentials.Login,
			body.Credentials.Password,
		)

		if errors.Is(err, error2.UserIsNotFound) || errors.Is(err, error2.PasswordIsIncorrect) {
			rest2.GenerateBadRequestResponse(w, nil, "Неверный логин или пароль", err.Error())
			return
		}
		if errors.Is(err, error2.LoginRequestNotFound) {
			rest2.GenerateNotFoundErrorResponse(w, nil, "Запрос на создание сессии не найден", err.Error())
			return
		}
		if err != nil {
			rest2.GenerateInternalServerErrorResponse(w, nil, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionId.String(),
			Path:     "/",
			MaxAge:   3600, // секунд
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
