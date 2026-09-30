package commands

import "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/repository/postgres"

type CreateSourceCommand struct {
	Postgres *postgres.Postgres
}
