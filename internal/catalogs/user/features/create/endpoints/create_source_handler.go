package endpoints

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/features/create/commands"
	"github.com/go-chi/chi/v5"
)

type CreateSourceHandler struct {
	Mux        *chi.Mux
	Dependents *CreateSourceHandlerDeps
}

type CreateSourceHandlerDeps struct {
	Command *commands.CreateSourceCommand
}

func NewCreateSourceHandler(mux *chi.Mux, deps *CreateSourceHandlerDeps) {
	handler := &CreateSourceHandler{
		Mux:        mux,
		Dependents: deps,
	}

	CreateAddSourceEndpoint(handler).MapEndpoint()
}
