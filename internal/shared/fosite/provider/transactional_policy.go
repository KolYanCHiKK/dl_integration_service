package provider

import (
	"context"
	"errors"

	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/uow"
)

type Transactional struct {
	*postgres2.Db
}

func (t *Transactional) BeginTX(ctx context.Context) (context.Context, error) {
	tx := t.Db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}
	return uow.WithTx(ctx, tx), nil // транзакция добавляется в контекст
}

func (t *Transactional) Commit(ctx context.Context) error {
	tx, ok := uow.TxFromContext(ctx)
	if !ok {
		return errors.New("commit: no transaction in context")
	}
	return tx.Commit().Error
}

func (t *Transactional) Rollback(ctx context.Context) error {
	tx, ok := uow.TxFromContext(ctx)
	if !ok {
		return errors.New("rollback: no transaction in context")
	}
	return tx.Rollback().Error
}
