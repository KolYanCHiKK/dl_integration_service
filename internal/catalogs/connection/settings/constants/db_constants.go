package constants

import "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/settings/types"

const (
	Postgres types.DBType = "postgres"
	MySQL    types.DBType = "mysql"
	MongoDB  types.DBType = "mongodb"
)

const (
	RoleOwner  types.AccessRole = "owner"
	RoleEditor types.AccessRole = "editor"
	RoleViewer types.AccessRole = "viewer"
)
