package migration

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
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

		// Проводим автомиграцию
		if err := tx.AutoMigrate(
			&user.User{},
			&user.Role{},
			&user.UserRole{},
			&user.LoginSession{},

			&models.Client{},
			&models.ClientRedirectURI{},
			&models.ClientGrantType{},
			&models.ClientResponseType{},
			&models.ClientScope{},
			&models.Scope{},
			&models.Resource{},
			&models.ClientResource{},
			&models.ClientSecretRotation{},

			&models.Request{},
			&models.RequestScope{},
			&models.RequestResource{},
			&models.AuthorizeCode{},
			&models.AccessToken{},
			&models.RefreshToken{},
			&models.PKCERequest{},
			&models.OIDCSession{},
			&models.JTIBlacklist{},
			&models.SigningKey{},
		); err != nil {
			return err
		}

		logger.CreateOperationLogs("migration", "schema", schema)

		return nil
	})
}

func DropEdpTables(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`SET LOCAL search_path TO ` + schema).Error; err != nil {
			return err
		}

		if err := tx.Migrator().DropTable(
			&user.User{},
			&user.Role{},
			&user.UserRole{},
			&user.LoginSession{},

			&models.Client{},
			&models.ClientRedirectURI{},
			&models.ClientGrantType{},
			&models.ClientResponseType{},
			&models.ClientScope{},
			&models.Scope{},
			&models.Resource{},
			&models.ClientResource{},
			&models.ClientSecretRotation{},

			&models.Request{},
			&models.RequestScope{},
			&models.RequestResource{},
			&models.AuthorizeCode{},
			&models.AccessToken{},
			&models.RefreshToken{},
			&models.PKCERequest{},
			&models.OIDCSession{},
			&models.JTIBlacklist{},
			&models.SigningKey{},
		); err != nil {
			return err
		}

		logger.CreateOperationLogs("migration rollback", schema)

		return nil
	})
}
