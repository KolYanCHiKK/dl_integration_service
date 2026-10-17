package models

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"gorm.io/datatypes"
)

type Request struct {
	ID              oauth.AuthRequestId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FositeRequestID string              `gorm:"size:64;not null;index"`                  // ID запроса от fosite: общий для code/access/refresh одной выдачи. НЕ уникален
	ClientID        string              `gorm:"type:uuid;not null;index"`                // Client_id строкой (fosite работает со строками)
	Subject         string              `gorm:"size:64;not null;index"`                  // Users.id строкой
	RequestedAt     time.Time           `gorm:"type:timestamptz;not null;default:now()"` // Время получения исходного запроса
	Form            datatypes.JSON      `gorm:"type:jsonb;not null"`                     // Url.Values исходного запроса
	Session         []byte              `gorm:"type:bytea;not null"`                     // Сохранённое внутреннее состояние fosite для конкретного OAuth-запроса
	IsDeleted       bool                `gorm:"not null;type:boolean;default:false"`
	CreatedAt       time.Time           `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt       time.Time           `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (Request) TableName() string { return "edp.oauth_request" }

type RequestScope struct {
	ID        oauth.RequestScopeId   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequestID oauth.AuthRequestId    `gorm:"type:uuid;uniqueIndex:ux_request_scope_kind"`
	Kind      oauth.RequestScopeKind `gorm:"size:16;uniqueIndex:ux_request_scope_kind;check:chk_request_scope_kind,kind IN ('requested','granted')"`
	Scope     string                 `gorm:"size:128;uniqueIndex:ux_request_scope_kind"` // Scope прямо с запроса и ответа
	IsDeleted bool                   `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time              `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time              `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Request *Request `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE"`
}

func (RequestScope) TableName() string { return "edp.request_scope" }

type RequestResource struct {
	ID        oauth.RequestResourceId   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequestID oauth.AuthRequestId       `gorm:"type:uuid;uniqueIndex:ux_request_resourcse_kind"`
	Kind      oauth.RequestResourceKind `gorm:"size:16;uniqueIndex:ux_request_resourcse_kind;check:chk_request_resource_kind,kind IN ('requested','granted')"`
	Resource  string                    `gorm:"size:255;uniqueIndex:ux_request_resourcse_kind"` // Ресурсы прямо с запроса и ответа
	IsDeleted bool                      `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time                 `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time                 `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Request *Request `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE"`
}

func (RequestResource) TableName() string {
	return "edp.request_resource"
}

type TokenRecord struct {
	ID        oauth.TokenId       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Signature string              `gorm:"size:255;uniqueIndex"`
	RequestID oauth.AuthRequestId `gorm:"type:uuid;not null;uniqueIndex"`
	ExpiresAt time.Time           `gorm:"type:timestamptz;not null;index"`
	Active    bool                `gorm:"not null;default:true"`

	Request *Request `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE"`
}

type AuthorizeCode struct {
	TokenRecord `gorm:"embedded"`
	CreatedAt   time.Time `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (AuthorizeCode) TableName() string { return "edp.authorize_code" }

type AccessToken struct {
	TokenRecord `gorm:"embedded"`
	CreatedAt   time.Time `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (AccessToken) TableName() string { return "edp.access_token" }

type RefreshToken struct {
	TokenRecord          `gorm:"embedded"`
	AccessTokenSignature string     `gorm:"size:255;index"`         // access token, выданный вместе с refresh (нужен при rotation)
	FirstUsedAt          *time.Time `gorm:"type:timestamptz;index"` // для grace period при rotation (параллельные запросы клиента при замене refresh токена)
	CreatedAt            time.Time  `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt            time.Time  `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (RefreshToken) TableName() string { return "edp.refresh_token" }

type PKCERequest struct {
	ID              oauth.PKCERequestId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequestID       oauth.AuthRequestId `gorm:"type:uuid;not null;uniqueIndex"`
	CodeChallenge   string              `gorm:"size:255;not null"`
	ChallengeMethod string              `gorm:"size:16;not null;default:S256"`
	ExpiresAt       time.Time           `gorm:"not null;type:timestamptz;index"`
	IsDeleted       bool                `gorm:"not null;type:boolean;default:false"`
	CreatedAt       time.Time           `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt       time.Time           `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Request *Request `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE"`
}

func (PKCERequest) TableName() string {
	return "edp.pkce_request"
}

type OIDCSession struct {
	ID        oauth.OIDCSessionId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RequestID oauth.AuthRequestId `gorm:"type:uuid;not null;uniqueIndex"`
	Signature string              `gorm:"size:255;uniqueIndex"`
	Subject   string              `gorm:"size:255;not null;index"`
	Nonce     string              `gorm:"size:255"`
	AuthTime  *time.Time          `gorm:"type:timestamptz"`
	IsDeleted bool                `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time           `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time           `gorm:"not null;type:timestamptz;autoUpdateTime"`

	Request *Request `gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE"`
}

func (OIDCSession) TableName() string {
	return "edp.oidc_session"
}

type JTIBlacklist struct {
	ID        oauth.JTIBlacklistId `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	JTI       string               `gorm:"size:255;not null;uniqueIndex"`
	ExpiresAt time.Time            `gorm:"not null;type:timestamptz;index"`
	IsDeleted bool                 `gorm:"not null;type:boolean;default:false"`
	CreatedAt time.Time            `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt time.Time            `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (JTIBlacklist) TableName() string {
	return "edp.jti_blacklist"
}

type SigningKey struct {
	ID            oauth.SigningKeyId     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	KID           string                 `gorm:"size:64;not null;uniqueIndex"`                                                                                 // Идентификатор ключа, попадает в заголовок JWT. Не переиспользуется, поэтому индекс не частичный
	Algorithm     string                 `gorm:"size:16;not null;check:chk_signing_key_algorithm,algorithm IN ('RS256','PS256','ES256')"`                      // Алгоритм подписи
	PublicJWK     datatypes.JSON         `gorm:"type:jsonb;not null"`                                                                                          // Публичная часть в формате JWK, отдаётся в JWKS endpoint
	PrivateKeyEnc []byte                 `gorm:"type:bytea;not null"`                                                                                          // Приватный ключ, зашифрованный мастер-ключом (KMS/HSM). В открытом виде не хранить
	Status        oauth.SigningKeyStatus `gorm:"size:16;not null;default:active;index;check:chk_signing_key_status,status IN ('active','retiring','retired')"` // active подписывает, retiring только проверяется, retired убран из JWKS
	NotBefore     time.Time              `gorm:"type:timestamptz;not null"`                                                                                    // С какого момента ключ можно использовать
	RotatedAt     *time.Time             `gorm:"type:timestamptz"`                                                                                             // Когда ключ перестал быть активным
	RetireAfter   *time.Time             `gorm:"type:timestamptz;index"`                                                                                       // Когда перевести в retired (после истечения последних токенов)
	IsDeleted     bool                   `gorm:"not null;type:boolean;default:false"`
	CreatedAt     time.Time              `gorm:"not null;type:timestamptz;autoCreateTime"`
	UpdatedAt     time.Time              `gorm:"not null;type:timestamptz;autoUpdateTime"`
}

func (SigningKey) TableName() string {
	return "edp.signing_key"
}
