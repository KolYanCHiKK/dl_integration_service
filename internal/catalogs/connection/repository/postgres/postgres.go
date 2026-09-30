package postgres

import "github.com/KolYanCHiKK/dl_integration_service/configs"

type Postgres struct {
	Db *configs.ApplicationConfig
}

func NewPostgres(db *configs.ApplicationConfig) *Postgres {
	return &Postgres{Db: db}
}
