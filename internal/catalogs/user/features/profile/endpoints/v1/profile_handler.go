package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/create/commands"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/validation"
	"github.com/go-chi/chi/v5"
)

type ProfileHandler struct {
	Mux        *chi.Mux
	Dependents *ProfileHandlerDeps
}

type ProfileHandlerDeps struct {
	Command   *commands.CreateUserCommand
	Validator *validation.RestRequesterValidator
}

func NewProfileHandlerDeps(mux *chi.Mux, deps *ProfileHandlerDeps) {
	handler := &ProfileHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewProfileE(handler).MapEndpoint()
}
