package oauth

import "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"

const (
	ClientStatusActive  oauth.ClientStatus = "active"
	ClientStatusBlocked oauth.ClientStatus = "blocked"
	ClientStatusDeleted oauth.ClientStatus = "deleted"
)

const (
	AuthMethodSecretBasic   oauth.AuthMethod = "client_secret_basic"
	AuthMethodSecretPost    oauth.AuthMethod = "client_secret_post"
	AuthMethodPrivateKeyJWT oauth.AuthMethod = "private_key_jwt"
	AuthMethodNone          oauth.AuthMethod = "none"
)

const (
	GrantTypeAuthorizationCode oauth.GrantType = "authorization_code"
	GrantTypeRefreshToken      oauth.GrantType = "refresh_token"
	GrantTypeClientCredentials oauth.GrantType = "client_credentials"
)

const (
	ResponseTypeCode    oauth.ResponseType = "code"
	ResponseTypeIDToken oauth.ResponseType = "id_token"
	ResponseTypeToken   oauth.ResponseType = "token"
)

const (
	ScopeOpenID  oauth.ScopeName = "openid"
	ScopeProfile oauth.ScopeName = "profile"
	ScopeEmail   oauth.ScopeName = "email"
	ScopeAddress oauth.ScopeName = "address"
	ScopePhone   oauth.ScopeName = "phone"
	ScopeAPI     oauth.ScopeName = "api"
)

const (
	ResourceIntegrationUI oauth.ResourceName = "web-app"
)
