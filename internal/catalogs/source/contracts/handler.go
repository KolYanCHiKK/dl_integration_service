package contracts

import (
	chi "github.com/go-chi/chi/v5"
)

type Handler struct {
	Dependents *HandlerDeps
}

type HandlerDeps struct {
}

func New() {
	chi.NewRouter()
}
