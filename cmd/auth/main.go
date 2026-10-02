package main

import (
	"context"
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/service/auth"
	"github.com/KolYanCHiKK/dl_integration_service/migration"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Инициализация форматера логгера
	logger.InitLogger()

	// Добавление конфигов
	appConfig, err := configs.LoadConfig()
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка добавления конфига")
		return
	}

	// Создание БД
	postgres, err := postgres2.NewDb(context.Background(), appConfig)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка подключения к Postgres")
		return
	}

	defer func() {
		if err := postgres.Close(); err != nil {
			logger.CreateErrLog(err.Error(), "Ошибка закрытия соединений Postgres")
			return
		}
	}()

	// Миграция таблиц
	err = migration.CreateEdpTables(postgres.DB)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка миграции Postgres")
		return
	}

	// Инициализация и создание роутера
	mux := chi.NewMux()
	logger.CreateOperationLogs("starting_server", "port", "8080")

	// Инициализация и подянтие сервера
	err = auth.StartIdentityServer(mux, postgres)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка подготовки запуска сервера")
		return
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка запуска сервера")
		return
	}

}
