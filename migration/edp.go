package migration

import (
	connection "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/models"
	user "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/models"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"gorm.io/gorm"
)

const schema = "edp"

func CreateEdpTables(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`CREATE SCHEMA IF NOT EXISTS ` + schema).Error; err != nil {
			return err
		}

		if err := tx.Exec(`SET search_path TO ` + schema).Error; err != nil {
			return err
		}

		if err := tx.AutoMigrate(
			&connection.Connection{},
			&connection.ConnectionAccess{},
			&user.User{},
		); err != nil {
			return err
		}

		logger.CreateOperationLogs("migration", schema)

		return nil
	})
}

func DropUsers(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`SET LOCAL search_path TO ` + schema).Error; err != nil {
			return err
		}

		if err := tx.Migrator().DropTable(
			&connection.Connection{},
			&user.User{},
		); err != nil {
			return err
		}

		logger.CreateOperationLogs("migration rollback", schema)

		return nil
	})
}
