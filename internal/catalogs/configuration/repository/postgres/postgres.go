package postgres

import (
	"context"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/uow"
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

func (p *Postgres) connWithCtxAndTx(ctx context.Context) *gorm.DB {
	if tx, ok := uow.TxFromContext(ctx); ok {
		return tx
	}
	return p.WithContext(ctx)
}
