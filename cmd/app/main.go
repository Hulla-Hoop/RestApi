package main

import (
	"github.com/hulla-hoop/restapi/internal/DB/psql"
	"github.com/hulla-hoop/restapi/internal/pkg/app"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	//Инициализируем логеры
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("Загружаем переменные окружения")
	//Загружаем .env содержащий все конфиги
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Не загружается .env файл", err)
	}
	logger.Info("Подключаемся к БД")
	//Инициализируем базу данных с стандартной библиотекой
	db, err := psql.InitDb(logger)
	if err != nil {
		logger.Fatal("Проблемы иниициализации БД", err)
	}
	logger.Info("Запускаем приложение")
	//Инициализируем echo роутер и запскаем его
	a := app.New(db, logger)

	err = a.Start()
	if err != nil {
		logger.Fatal("Проблемы запуска приложения", err)
	}

}
