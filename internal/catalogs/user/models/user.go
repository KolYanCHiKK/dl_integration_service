package models

import (
	"time"

	oauthmodel "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	oauthtypes "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
	"gorm.io/datatypes"
)

type Department struct {
	ID          user.DepartmentId  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string             `gorm:"size:100;not null;uniqueIndex"`
	Description string             `gorm:"size:255"`
	Higher      *user.DepartmentId `gorm:"type:uuid"`
	IsDeleted   bool               `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time          `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time          `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`
}

func (Department) TableName() string {
	return "edp.department"
}

type User struct {
	ID             user.UserId        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Surname        string             `gorm:"type:varchar(100)"`
	Name           string             `gorm:"type:varchar(100);not null"`
	LastName       string             `gorm:"type:varchar(100)"`
	Phone          string             `gorm:"type:varchar(65)"`
	BirthDate      *datatypes.Date    `gorm:"default:null"`
	Login          string             `gorm:"not null;type:varchar(120);uniqueIndex"`
	Email          string             `gorm:"not null;type:varchar(120)"`
	DepartmentID   *user.DepartmentId `gorm:"type:uuid"`
	Higher         *user.UserId       `gorm:"type:uuid"`
	HashedPassword string             `gorm:"not null;type:varchar(1000)"`
	IsDeleted      bool               `gorm:"not null;type:boolean;default:false"`
	CreatedAt      time.Time          `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt      time.Time          `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`

	Department *Department `gorm:"foreignKey:DepartmentID;references:ID;constraint:OnDelete:SET NULL"`
}

func (User) TableName() string {
	return "edp.user"
}

type Role struct {
	ID          user.RoleId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string      `gorm:"size:100;not null;uniqueIndex"`
	Description string      `gorm:"size:255"`
	IsSystem    bool        `gorm:"not null;type:boolean;default:false"`
	IsDeleted   bool        `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time   `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`
}

func (Role) TableName() string {
	return "edp.role"
}

type ClientRole struct {
	ID        user.ClientRoleId   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoleID    user.RoleId         `gorm:"type:uuid;not null;index:ux_role_id_client_id,unique"`
	ClientID  oauthtypes.ClientId `gorm:"type:uuid;not null;index:ux_role_id_client_id,unique"`
	IsDeleted bool                `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time           `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time           `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`

	Client oauthmodel.Client `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE"`
	Role   Role              `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
}

func (ClientRole) TableName() string {
	return "edp.client_role"
}

type UserRole struct {
	ID        user.UserRoleId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    user.UserId     `gorm:"type:uuid;not null;index:user_role_active_idx,unique,where:revoked_at IS NULL"`
	RoleID    user.RoleId     `gorm:"type:uuid;not null;index:user_role_active_idx,unique,where:revoked_at IS NULL"`
	RevokedAt *time.Time      `gorm:"type:timestamptz"`
	CreatedAt time.Time       `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time       `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`

	Role Role `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserRole) TableName() string {
	return "edp.user_role"
}

type PermissionTarget struct {
	ID          user.PermissionTargetId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code        string                  `gorm:"size:255;not null;index:ux_name_client_id,unique"`
	ClientID    oauthtypes.ClientId     `gorm:"type:uuid;not null;index:ux_name_client_id,unique"`
	Description string                  `gorm:"size:255"`
	IsDeleted   bool                    `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time               `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time               `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`

	Client oauthmodel.Client `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE"`
}

func (PermissionTarget) TableName() string {
	return "edp.permission_target"
}

type Permission struct {
	ID          user.PermissionId       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code        string                  `gorm:"size:32;not null;index:permission_target_name_idx,unique"`
	Description string                  `gorm:"size:255"`
	TargetID    user.PermissionTargetId `gorm:"type:uuid;not null;index:permission_target_name_idx,unique"`
	IsDeleted   bool                    `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time               `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time               `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`

	PermissionTarget PermissionTarget `gorm:"foreignKey:TargetID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Permission) TableName() string {
	return "edp.permission"
}

type RolePermission struct {
	ID           user.RolePermissionId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoleID       user.RoleId           `gorm:"type:uuid;not null;index:role_permission_idx,unique"`
	PermissionID user.PermissionId     `gorm:"type:uuid;not null;index:role_permission_idx,unique"`
	CreatedAt    time.Time             `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`

	Role       Role       `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
	Permission Permission `gorm:"foreignKey:PermissionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (RolePermission) TableName() string {
	return "edp.role_permission"
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
	CreatedAt time.Time           `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time           `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (LoginSession) TableName() string {
	return "edp.login_session"
}

type LoginRequest struct {
	ID          user.LoginRequestId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RedirectURL string              `gorm:"not null;type:text"`
	UsedAt      *time.Time          `gorm:"type:timestamptz;index"`
	ExpiresAt   time.Time           `gorm:"not null;type:timestamptz;index"`
	CreatedAt   time.Time           `gorm:"not null;default:now();type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time           `gorm:"not null;default:now();type:timestamptz;autoUpdateTime"`
}

func (LoginRequest) TableName() string {
	return "edp.login_request"
}
