package contract

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/cache/list/signingkey"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/cache/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
)

type ApplicationCacheElementInMemory struct {
	SigningKey *signingkey.SigningKeyProviderCache
}

func PrepareAppCache(database *postgres2.Db, authConfig *configs.AuthJwtConfig, cipher *ecdsa.CipherEcdsa) (ApplicationCacheElementInMemory, error) {
	// Инициализируем репозитории
	postgresRepository := postgres.NewPostgres(database)

	// Создаем кэш структуры
	signingKeyCache := signingkey.NewSigningKeyProviderCache(&signingkey.JWKCacheDependents{
		Postgres:    postgresRepository,
		CipherEcdsa: cipher,
	}, authConfig.JWKReloadCacheTTL, nil, time.Now().UTC())

	// Собираем тип для работы с кэшем на сервере

	serverInMemoryCache := ApplicationCacheElementInMemory{
		SigningKey: signingKeyCache,
	}

	return serverInMemoryCache, nil
}
