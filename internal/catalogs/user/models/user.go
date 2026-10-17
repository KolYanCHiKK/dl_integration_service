package models

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
	"gorm.io/datatypes"
)

type User struct {
	ID             user.UserId    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Surname        string         `gorm:"type:varchar(100)"`
	Name           string         `gorm:"type:varchar(100);not null"`
	LastName       string         `gorm:"type:varchar(100)"`
	Phone          string         `gorm:"type:varchar(65)"`
	BirthDate      datatypes.Date `gorm:"default:null"`
	Login          string         `gorm:"not null;type:varchar(120);uniqueIndex"`
	Email          string         `gorm:"not null;type:varchar(120)"`
	HashedPassword string         `gorm:"not null;type:varchar(1000)"`
	IsDeleted      bool           `gorm:"not null;type:boolean;default:false"`
	CreatedAt      time.Time      `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (User) TableName() string {
	return "edp.user"
}

type Role struct {
	ID          user.RoleId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string      `gorm:"size:64;not null;uniqueIndex"`
	Description string      `gorm:"size:255"`
	IsDeleted   bool        `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time   `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (Role) TableName() string {
	return "edp.role"
}

type UserRole struct {
	ID        user.UserRoleId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoleID    user.RoleId     `gorm:"type:uuid;not null;index:user_role_idx,unique"`
	UserID    user.UserId     `gorm:"type:uuid;not null;index:user_role_idx,unique"`
	IsDeleted bool            `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time       `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Role Role `gorm:"foreignKey:RoleID;references:ID"`
	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (UserRole) TableName() string {
	return "edp.user_role"
}

type LoginSession struct {
	ID        user.LoginSessionId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    user.UserId         `gorm:"type:uuid;not null;index"`
	AuthTime  time.Time           `gorm:"not null;type:timestamptz"`
	IP        string              `gorm:"size:45"`
	UserAgent string              `gorm:"size:512"`
	ExpiresAt time.Time           `gorm:"not null;type:timestamptz;index"`
	RevokedAt *time.Time          `gorm:"type:timestamptz"`
	IsDeleted bool                `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time           `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time           `gorm:"not null;type:timestamptz;autoUpdateTime"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (LoginSession) TableName() string {
	return "edp.login_session"
}
