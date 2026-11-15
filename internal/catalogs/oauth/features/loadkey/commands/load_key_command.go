package commands

import (
	"context"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/loadkey/cache"
)

type LoadKeyCommand struct {
	JWKCache cache.SigningKeyLoadKeyCacheImplementation
}

func PrepareSigningKeyProviderCache(signingKeyCache cache.SigningKeyLoadKeyCacheImplementation) (*LoadKeyCommand, error) {
	loadCommand := &LoadKeyCommand{
		JWKCache: signingKeyCache,
	}

	_, err := loadCommand.JWKCache.ReloadCache(context.Background())
	if err != nil {
		return nil, err
	}

	return loadCommand, nil
}

// Пока работаем в одном инстансе, поэтому и кэш только в воркере. При масшабировании поменяем на Redis

func (c *LoadKeyCommand) GetSigningKey(ctx context.Context) (any, error) {
	key, loadedAt, ttl := c.JWKCache.Get()

	if key == nil || loadedAt.Add(ttl).Before(time.Now().UTC()) {
		newKey, err := c.JWKCache.ReloadCache(ctx)
		if err != nil {
			return nil, err
		}
		return newKey, nil
	}

	return key, nil
}
