package custom

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/utils"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/jwt"
)

type Session struct {
	*oauth2.JWTSession
	Roles       []string       `json:"-"`
	Permissions []string       `json:"-"`
	ClientId    oauth.ClientId `json:"-"`
}

func NewEmptySession() *Session {
	return &Session{
		JWTSession: &oauth2.JWTSession{
			JWTClaims: &jwt.JWTClaims{},
			JWTHeader: &jwt.Headers{},
		},
	}
}

func (s *Session) GetJWTClaims() jwt.JWTClaimsContainer {
	claims, ok := s.JWTSession.GetJWTClaims().(*jwt.JWTClaims)
	if !ok {
		return s.JWTSession.GetJWTClaims()
	}

	if claims.Extra == nil {
		claims.Extra = map[string]interface{}{}
	}
	claims.Extra["roles"] = utils.NonNil(s.Roles)
	claims.Extra["permissions"] = utils.NonNil(s.Permissions)
	claims.Extra["client_id"] = s.ClientId

	return claims
}

func (s *Session) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	var inner *oauth2.JWTSession
	if s.JWTSession != nil {
		inner, _ = s.JWTSession.Clone().(*oauth2.JWTSession)
	}

	return &Session{
		JWTSession:  inner,
		Roles:       append([]string(nil), s.Roles...),
		Permissions: append([]string(nil), s.Permissions...),
	}
}
