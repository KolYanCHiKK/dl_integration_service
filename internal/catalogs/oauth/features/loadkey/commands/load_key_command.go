package commands

import (
	"context"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/loadkey/cache"
)

type LoadKeyCommand struct {
	JWKCache   *cache.SigningKeyProviderCache
	AuthConfig *configs.AuthJwtConfig
}

func PrepareSigningKeyProviderCache(jwkDeps *cache.JWKCacheDependents, authConfig *configs.AuthJwtConfig) (*LoadKeyCommand, error) {
	cacheReloadTTL := authConfig.JWKReloadCacheTTL
	signingKeyCache := cache.NewSigningKeyProviderCache(jwkDeps, cacheReloadTTL, nil, time.Now().UTC())

	loadCommand := &LoadKeyCommand{
		JWKCache: signingKeyCache,
	}

	_, err := loadCommand.JWKCache.ReloadCache(context.Background())
	if err != nil {
		return nil, err
	}

	return loadCommand, nil
}

// Пока работаем в одном инстансе, поэтому и кэш только в воркере. При масшабировании поменяем

func (c *LoadKeyCommand) GetSigningKey(ctx context.Context) (any, error) {
	key, loadedAt := c.JWKCache.GetCurrentAndLoadedAt()

	if key == nil || loadedAt.Add(c.JWKCache.Ttl).Before(time.Now().UTC()) {
		newKey, err := c.JWKCache.ReloadCache(ctx)
		if err != nil {
			return nil, err
		}
		return newKey, nil
	}

	return key, nil
}
