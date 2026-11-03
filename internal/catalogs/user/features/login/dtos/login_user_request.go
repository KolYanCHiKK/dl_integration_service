package dtos

import "github.com/google/uuid"

type Credentials struct {
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=5"`
}

type LoginUserRequest struct {
	Credentials    *Credentials `json:"credentials"  validate:"required"`
	LoginRequestId uuid.UUID    `json:"loginRequestId"`
}
