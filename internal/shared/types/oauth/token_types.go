package oauth

import "github.com/google/uuid"

type AuthRequestId = uuid.UUID

type RequestScopeId = uuid.UUID
type RequestScopeKind = string

type RequestResourceId = uuid.UUID
type RequestResourceKind = string

type TokenId = uuid.UUID

type PKCERequestId = uuid.UUID

type OIDCSessionId = uuid.UUID

type JTIBlacklistId = uuid.UUID

type SigningKeyId = uuid.UUID
type SigningKeyStatus = string
