package postgres

import (
	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
)

type Postgres struct {
	*postgres.Db
}

func NewPostgres(db *postgres.Db) *Postgres {
	return &Postgres{Db: db}
}
