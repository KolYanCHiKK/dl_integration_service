package oauth

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
)

const (
	RequestScopeKindRequested oauth.RequestScopeKind = "requested"
	RequestScopeKindGranted   oauth.RequestScopeKind = "granted"
)

const (
	RequestResourceKindRequested oauth.RequestResourceKind = "requested"
	RequestResourceKindGranted   oauth.RequestResourceKind = "granted"
)

const (
	SigningKeyStatusActive   oauth.SigningKeyStatus = "active"
	SigningKeyStatusRetiring oauth.SigningKeyStatus = "retiring"
	SigningKeyStatusRetired  oauth.SigningKeyStatus = "retired"
)
