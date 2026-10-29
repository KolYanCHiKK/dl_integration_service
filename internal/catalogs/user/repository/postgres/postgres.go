package postgres

import (
	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	*postgres.Db
}

func NewPostgres(db *postgres.Db) *Postgres {
	return &Postgres{Db: db}
}

func (p *Postgres) StartTransaction() *gorm.DB {
	return p.DB
}
