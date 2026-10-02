package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
)

type LoginUserEndpoint struct {
	*LoginUserHandler
}

func NewLoginUserEndpoint(params *LoginUserHandler) interfaces.Endpoint {
	return &LoginUserEndpoint{params}
}

func (ep *LoginUserEndpoint) MapEndpoint() {
	ep.Mux.Post("/login", ep.handle())
}

func (ep *LoginUserEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
