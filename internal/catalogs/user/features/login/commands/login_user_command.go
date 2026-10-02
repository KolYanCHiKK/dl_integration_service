package commands

import "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"

type LoginUserCommand struct {
	Postgres *postgres.Postgres
}
