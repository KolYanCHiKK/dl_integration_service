package dtos

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/settings/types"
)

type CreateSourceRequest struct {
	Name     string       `json:"name"`
	Type     types.DBType `json:"type"`
	Host     string       `json:"host"`
	Port     int          `json:"port"`
	User     string       `json:"user"`
	Password string       `json:"password"`
	Database string       `json:"database"`
}
