package uow

import (
	"context"

	"gorm.io/gorm"
)

type PostgresRepository interface {
	StartTransaction() *gorm.DB
}

type txKey struct{}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return nil, false
	}

	return tx, ok
}

func ExecuteTransaction(ctx context.Context, repo PostgresRepository, fn func(tx *gorm.DB) error) error {
	// Проверяем наличие внешней транзакции
	if tx, ok := TxFromContext(ctx); ok {
		return fn(tx)
	}

	// Если не нашли, то создаем новую
	if err := repo.StartTransaction().WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			return fn(tx)
		},
	); err != nil {
		return err
	}

	return nil
}
