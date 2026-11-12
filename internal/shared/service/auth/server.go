package auth

import (
	"encoding/base64"
	"strings"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	contract3 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/contract"
	contract2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/contract"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/contract"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func StartOauthServer(mux *chi.Mux, database *postgres.Db, applicationConfig *configs.ApplicationConfig) *logger2.DefaultErrorLogsParams {
	// Создаем объект шифрования для хранения сессии
	sessionKey, err := base64.StdEncoding.DecodeString(strings.TrimSpace(applicationConfig.OauthFosite.SessionEncKey))
	if err != nil {
		return &logger2.DefaultErrorLogsParams{
			Error: err,
			Info:  "SESSION_ENC_KEY не является корректным base64",
		}
	}

	// Создаем объекты шифрования
	sessionCipher, err := aead.NewSessionCipher(sessionKey)
	if err != nil {
		return &logger2.DefaultErrorLogsParams{
			Error: err,
			Info:  "Ошибка создания объекта шифрования",
		}
	}

	// Создаем объект шифрования для хранения приватного ключа JWK
	jwkKey, err := base64.StdEncoding.DecodeString(strings.TrimSpace(applicationConfig.JWTAuth.JWKEncKey))
	if err != nil {
		return &logger2.DefaultErrorLogsParams{
			Error: err,
			Info:  "JWK_ENC_KEY не является корректным base64",
		}
	}

	jwkCipher, err := ecdsa.NewJWKCipher(jwkKey)
	if err != nil {
		return &logger2.DefaultErrorLogsParams{
			Error: err,
			Info:  "Ошибка создания объекта шифрования",
		}
	}

	userComponent, errLoggerParams := contract.InitUserHandlers(mux, database, applicationConfig.JWTAuth, applicationConfig.Network)
	if errLoggerParams != nil {
		return errLoggerParams
	}

	err, errInfo := contract2.InitOauthHandlers(mux, database, applicationConfig.OauthFosite, applicationConfig.JWTAuth, sessionCipher, jwkCipher, userComponent.PostgresRepository)
	if err != nil {
		return &logger2.DefaultErrorLogsParams{Error: err, Info: errInfo}
	}

	errLoggerParams = contract3.InitConfigurationHandlers(mux, database, applicationConfig.JWTAuth, jwkCipher)
	if errLoggerParams != nil {
		return errLoggerParams
	}

	return nil
}
