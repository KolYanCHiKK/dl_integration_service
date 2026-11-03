package v1

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/create/commands"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/validation"
	"github.com/go-chi/chi/v5"
)

type CreateUserHandler struct {
	Mux        *chi.Mux
	Dependents *CreateUserHandlerDeps
}

type CreateUserHandlerDeps struct {
	Command   *commands.CreateUserCommand
	Validator *validation.RestRequesterValidator
}

func NewCreateUserHandler(mux *chi.Mux, deps *CreateUserHandlerDeps) {
	handler := &CreateUserHandler{
		Mux:        mux,
		Dependents: deps,
	}

	NewCreateUserEndpoint(handler).MapEndpoint()
}
