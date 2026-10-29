package dtos

import "github.com/google/uuid"

type LoginUserRequest struct {
	Login          string    `json:"login"`
	Password       string    `json:"password"`
	LoginRequestId uuid.UUID `json:"loginRequestId"`
}
