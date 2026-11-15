package cache

import (
	"context"
	"time"

	"github.com/go-jose/go-jose/v3"
)

type SigningKeyJWKCreateCacheImplementation interface {
	ReloadCache(ctx context.Context) (*jose.JSONWebKey, error)
	Get() (*jose.JSONWebKey, time.Time, time.Duration)
	SetCurrentAndLoadedAt(k *jose.JSONWebKey)
}
