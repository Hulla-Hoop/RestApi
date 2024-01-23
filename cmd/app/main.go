package main

import (
	"log"
	"os"

	"github.com/hulla-hoop/restapi/internal/DB/psql"
	"github.com/hulla-hoop/restapi/internal/pkg/app"
	"github.com/joho/godotenv"
)

func main() {
	//Инициализируем логеры
	infLogger := log.New(os.Stdout, "\n\u001b[33m  INFO:  ", log.Ldate|log.Lshortfile)
	errLogger := log.New(os.Stdout, "\n\u001b[31m  ERROR:  ", log.Ldate|log.Lshortfile)

	//Загружаем .env содержащий все конфиги
	err := godotenv.Load()
	if err != nil {
		errLogger.Fatal("Не загружается .env файл")
	}
	//Инициализируем базу данных с стандартной библиотекой
	db, err := psql.InitDb()
	if err != nil {
		errLogger.Println("Проблемы иниициализации БД", err)
	}
	//

	//Инициализируем echo роутер и запскаем его
	a := app.New(db, infLogger, errLogger)

	a.Start()

}
