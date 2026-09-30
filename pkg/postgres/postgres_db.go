package postgres

import (
	"context"
	"fmt"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db struct {
	*gorm.DB
}

func NewDb(ctx context.Context, conf *configs.ApplicationConfig) (*Db, error) {
	gdb, err := gorm.Open(postgres.Open(conf.Postgres.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(conf.Postgres.MaxOpenConn)
	sqlDB.SetMaxIdleConns(conf.Postgres.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(conf.Postgres.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(conf.Postgres.ConnMaxIdleTime)

	// Проверяем соединение
	pingCtx, cancel := context.WithTimeout(ctx, conf.Postgres.ConnectTimeout)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		defer func() {
			if err := sqlDB.Close(); err != nil {
				logger2.CreateErrLog(err.Error(), "Ошибка закрытия соединения Postgres при старте приложения")
			}
		}()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Db{gdb}, nil
}

func (d *Db) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
