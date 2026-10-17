package models

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"gorm.io/datatypes"
)

type Client struct {
	ID                     oauth.ClientId     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name                   string             `gorm:"size:255;not null"`
	Status                 oauth.ClientStatus `gorm:"size:16;not null;default:active;index"`
	Owner                  string             `gorm:"size:255"`
	SecretHash             string             `gorm:"size:255"`
	Public                 bool               `gorm:"not null;default:false"`
	AuthMethod             string             `gorm:"size:32;not null;default:client_secret_basic"`
	AuthSigningAlg         string             `gorm:"size:16"`
	JWKS                   datatypes.JSON     `gorm:"type:jsonb"`
	JWKSURI                string             `gorm:"size:1024;default:null"`
	AccessTokenTTLSeconds  *int               `gorm:"column:access_token_ttl_seconds"`
	RefreshTokenTTLSeconds *int               `gorm:"column:refresh_token_ttl_seconds"`
	IsDeleted              bool               `gorm:"not null;type:boolean;default:false"`
	CreatedAt              time.Time          `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt              time.Time          `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (Client) TableName() string {
	return "edp.client"
}

type ClientRedirectURI struct {
	ID        oauth.ClientRedirectUriId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID  oauth.ClientId            `gorm:"type:uuid;not null;uniqueIndex:ux_client_redirect_uri"`
	URI       string                    `gorm:"size:1024;not null;uniqueIndex:ux_client_redirect_uri"`
	IsDeleted bool                      `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time                 `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time                 `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Client *Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

func (ClientRedirectURI) TableName() string {
	return "edp.client_redirect_uri"
}

type ClientGrantType struct {
	ID        oauth.ClientGrantTypeId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID  oauth.ClientId          `gorm:"type:uuid;index:ux_client_grant_type,unique"`
	GrantType oauth.GrantType         `gorm:"size:64;index:ux_client_grant_type,unique;check:chk_client_grant_type,grant_type IN ('authorization_code','refresh_token','client_credentials')"`
	CreatedAt time.Time               `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time               `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Client *Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

func (ClientGrantType) TableName() string {
	return "edp.client_grant_type"
}

type ClientResponseType struct {
	ID           oauth.ResponseTypeId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID     oauth.ClientId       `gorm:"type:uuid;index:ux_client_response_type,unique"`
	ResponseType string               `gorm:"size:64;index:ux_client_response_type,unique;check:chk_client_response_type,response_type IN ('code','id_token','token')"`
	IsDeleted    bool                 `gorm:"not null;type:boolean;default:false"`
	CreatedAt    time.Time            `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt    time.Time            `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Client *Client `gorm:"foreignKey:ClientID;references:ID;constraint:OnDelete:CASCADE"`
}

func (ClientResponseType) TableName() string {
	return "edp.client_response_type"
}

type Scope struct {
	ID          oauth.ScopeId   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        oauth.ScopeName `gorm:"size:255;not null;uniqueIndex"`
	Description string          `gorm:"size:1024"`
	IsDeleted   bool            `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time       `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (Scope) TableName() string { return "edp.oauth_scope" }

type ClientScope struct {
	ID        oauth.ClientScopeId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID  oauth.ClientId      `gorm:"type:uuid;index:ux_client_scope,unique"`
	ScopeID   oauth.ScopeId       `gorm:"type:uuid;index:ux_client_scope,unique"`
	CreatedAt time.Time           `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time           `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Scope  *Scope  `gorm:"foreignKey:ScopeID;references:ID;constraint:OnDelete:CASCADE"`
	Client *Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

func (ClientScope) TableName() string { return "edp.client_scope" }

type Resource struct {
	ID          oauth.ResourceId   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        oauth.ResourceName `gorm:"size:255;not null;uniqueIndex"`
	Description string             `gorm:"size:1024"`
	IsDeleted   bool               `gorm:"not null;type:boolean;default:false"`
	CreatedAt   time.Time          `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time          `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (Resource) TableName() string { return "edp.oauth_resource" }

type ClientResource struct {
	ID         oauth.ClientResourceId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID   oauth.ClientId         `gorm:"type:uuid;index:ux_client_resource,unique"`
	ResourceID oauth.ResourceId       `gorm:"type:uuid;index:ux_client_resource,unique"`
	IsDeleted  bool                   `gorm:"not null;type:boolean;default:false"`
	CreatedAt  time.Time              `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt  time.Time              `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Resource *Resource `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:CASCADE"`
	Client   *Client   `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

func (ClientResource) TableName() string {
	return "edp.client_resource"
}

type ClientSecretRotation struct {
	ID         oauth.ClientSecretRotationId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ClientID   oauth.ClientId               `gorm:"type:uuid"`
	SecretHash string                       `gorm:"size:255;not null"`
	ExpiresAt  time.Time                    `gorm:"type:timestamptz;not null;index"`
	IsDeleted  bool                         `gorm:"not null;type:boolean;default:false"`
	CreatedAt  time.Time                    `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt  time.Time                    `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Client *Client `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

func (ClientSecretRotation) TableName() string {
	return "edp.client_secret_rotation"
}
