package models

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/settings/types"
	"gorm.io/datatypes"
)

type User struct {
	UserId         types.UserId   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Surname        string         `gorm:"type:varchar(100)"`
	Name           string         `gorm:"type:varchar(100);not null"`
	LastName       string         `gorm:"type:varchar(100)"`
	Phone          string         `gorm:"type:varchar(65)"`
	BirthDate      datatypes.Date `gorm:"default:null"`
	Email          string         `gorm:"not null;type:varchar(120)"`
	HashedPassword string         `gorm:"not null;type:varchar(1000)"`
	IsDeleted      bool           `gorm:"not null;type:boolean;default:true"`
	CreatedAt      time.Time      `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (User) TableName() string {
	return "edp.user"
}
