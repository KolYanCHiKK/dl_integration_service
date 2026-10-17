package uow

import (
	"context"

	"gorm.io/gorm"
)

type PostgresRepository interface {
	StartTransaction() *gorm.DB
}

func ExecuteTransaction(ctx context.Context, repo PostgresRepository, fn func(tx *gorm.DB) error) error {
	err := repo.StartTransaction().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
	if err != nil {
		return err
	}

	return nil
}
