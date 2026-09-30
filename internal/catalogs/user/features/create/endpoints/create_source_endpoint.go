package endpoints

import (
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
)

type AddSourceEndpoint struct {
	*CreateSourceHandler
}

func CreateAddSourceEndpoint(params *CreateSourceHandler) interfaces.Endpoint {
	return &AddSourceEndpoint{params}
}

func (ep *AddSourceEndpoint) MapEndpoint() {
	ep.Mux.Post("/sourcses", ep.handle())
}

func (ep *AddSourceEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
