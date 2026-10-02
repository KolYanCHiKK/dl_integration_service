package connection

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/connection"
)

const (
	Postgres connection.DBType = "postgres"
	MySQL    connection.DBType = "mysql"
	MongoDB  connection.DBType = "mongodb"
)

const (
	RoleOwner  connection.AccessRole = "owner"
	RoleEditor connection.AccessRole = "editor"
	RoleViewer connection.AccessRole = "viewer"
)
