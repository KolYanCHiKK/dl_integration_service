package oauth

import "github.com/google/uuid"

type ClientStatus = string

type ClientId = uuid.UUID
type AuthMethod = string

type ClientRedirectUriId = uuid.UUID

type ClientGrantTypeId = uuid.UUID
type GrantType = string

type ResponseTypeId = uuid.UUID
type ResponseType = string

type ScopeId = uuid.UUID
type ScopeName = string

type ClientScopeId = uuid.UUID

type ResourceId = uuid.UUID
type ResourceName = string

type ClientResourceId = uuid.UUID

type ClientSecretRotationId = uuid.UUID
