package main

import (
	"context"
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/configs"
	"github.com/KolYanCHiKK/dl_integration_service/migration"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	postgres2 "github.com/KolYanCHiKK/dl_integration_service/pkg/postgres"
)

func main() {
	// Инициализация форматера логгера
	logger.InitLogger()

	// Добавление конфигов
	config, err := configs.LoadConfig()
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка добавления конфига")
		return
	}

	// Создание БД
	postgres, err := postgres2.NewDb(context.Background(), config)
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

	// Запуск сервера
	err = http.ListenAndServe(":8080", http.DefaultServeMux)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка запуска сервера")
		return
	}

}
