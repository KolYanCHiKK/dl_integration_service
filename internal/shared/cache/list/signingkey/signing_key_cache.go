package signingkey

import (
	"context"
	"sync"
	"time"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/cache/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/go-jose/go-jose/v3"
)

type SigningKeyProviderCache struct {
	Ttl time.Duration // Частота перечитываия

	mu       sync.RWMutex
	current  *jose.JSONWebKey
	loadedAt time.Time

	JWKCacheDependents *JWKCacheDependents
}

type JWKCacheDependents struct {
	Postgres    *postgres.Postgres
	CipherEcdsa *ecdsa.CipherEcdsa
}

func NewSigningKeyProviderCache(deps *JWKCacheDependents, cacheTTL time.Duration, current *jose.JSONWebKey, loadedAt time.Time) *SigningKeyProviderCache {
	return &SigningKeyProviderCache{
		Ttl:                cacheTTL,
		mu:                 sync.RWMutex{},
		current:            current,
		loadedAt:           loadedAt,
		JWKCacheDependents: deps,
	}
}

func (c *SigningKeyProviderCache) ReloadCache(ctx context.Context) (*jose.JSONWebKey, error) {
	signingKeyModel, err := c.JWKCacheDependents.Postgres.GetActiveSigningKey(ctx)
	if err != nil {
		return nil, error2.MapWellKnowError(err)
	}

	priv, err := ecdsa.DecryptSigningKey(c.JWKCacheDependents.CipherEcdsa, signingKeyModel.PrivateKeyEnc, signingKeyModel.KID)
	if err != nil {
		return nil, error2.MapWellKnowError(err)
	}

	key := &jose.JSONWebKey{
		Key:       priv,
		KeyID:     signingKeyModel.KID,
		Algorithm: signingKeyModel.Algorithm,
		Use:       "sig",
	}

	c.SetCurrentAndLoadedAt(key)
	c.AddLog()

	return key, nil
}

func (c *SigningKeyProviderCache) Get() (*jose.JSONWebKey, time.Time, time.Duration) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current, c.loadedAt, c.Ttl
}

func (c *SigningKeyProviderCache) SetCurrentAndLoadedAt(k *jose.JSONWebKey) {
	c.mu.Lock()
	c.current, c.loadedAt = k, time.Now().UTC()
	c.mu.Unlock()
}

func (c *SigningKeyProviderCache) AddLog() {
	logger.CreateInMemoryCacheOperationInfoLog("SigningKeyProvider")
}
