package dtos

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/connection"
)

type CreateSourceRequest struct {
	Name     string            `json:"name"`
	Type     connection.DBType `json:"type"`
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	Database string            `json:"database"`
}
