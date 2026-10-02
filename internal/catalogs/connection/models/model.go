package models

import (
	"time"

	userModels "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/connection"
	userTypes "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
)

type Connection struct {
	Id                connection.ConnectionId   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name              string                    `gorm:"size:255;not null"`
	Type              connection.DBType         `gorm:"size:100;not null;index"`
	Host              string                    `gorm:"size:255;not null"`
	Port              int                       `gorm:"not null"`
	User              string                    `gorm:"size:255;not null"`
	PasswordEncrypted string                    `gorm:"column:password_encrypted;size:512;not null"`
	Database          string                    `gorm:"size:255;not null"`
	Settings          connection.JSONStructType `gorm:"default: null"`
	IsDeleted         bool                      `gorm:"not null;type:boolean;default:false"`
	CreatedAt         time.Time                 `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt         time.Time                 `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (Connection) TableName() string {
	return "edp.connection"
}

type ConnectionAccess struct {
	Id           connection.ConnectionAccessId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ConnectionID connection.ConnectionId       `gorm:"type:uuid;not null;index:conn_user_role_idx,unique"`
	UserID       userTypes.UserId              `gorm:"type:uuid;not null;index:conn_user_role_idx,unique"`
	Role         connection.AccessRole         `gorm:"size:100;not null;default:'viewer';index:conn_user_role_idx,unique"`
	IsDeleted    bool                          `gorm:"not null;type:boolean;default:false"`
	CreatedAt    time.Time                     `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt    time.Time                     `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Connection Connection      `gorm:"foreignKey:ConnectionID;references:Id"`
	User       userModels.User `gorm:"foreignKey:UserID;references:ID"`
}

func (ConnectionAccess) TableName() string {
	return "edp.connection_access"
}
