package v1

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
)

type ProfileEndpoint struct {
	*ProfileHandler
}

func NewProfileEndpoint(params *ProfileHandler) interfaces.Endpoint {
	return &ProfileEndpoint{params}
}

func (ep *ProfileEndpoint) MapEndpoint() {
	ep.Mux.Get("/users/profile", ep.handle())
}

func (ep *ProfileEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
