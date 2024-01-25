package app

import (
	"github.com/hulla-hoop/restapi/internal/DB"
	"github.com/hulla-hoop/restapi/internal/config"
	"github.com/hulla-hoop/restapi/internal/echoendpoint"
	"github.com/hulla-hoop/restapi/internal/service"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	e      *echoendpoint.Endpoint
	echo   *echo.Echo
	DB     DB.DB
	logger *logrus.Logger
}

func New(db DB.DB, logger *logrus.Logger) *App {
	a := App{}
	cfg := config.NewCfgApi()
	service := service.New(logger, cfg)
	a.DB = db
	a.logger = logger
	a.e = echoendpoint.New(a.DB, logger, service)
	a.echo = echo.New()
	a.echo.Use(middleware.Recover())

	user := a.echo.Group("/user")

	//  Создание записи
	//  Пример:
	//  /user
	//  body json
	//  {
	//      "name":"данные без цифр и пробелов",
	//      "surname":"данные без цифр и пробелов",
	//      "patronymic":"не обязательное поле|данные без цифр и пробелов"
	//  }

	user.POST("", a.e.Insert)

	//  Удаление записи
	//  Пример:
	//  /user/:id
	// id-пользователя

	user.DELETE("/:id", a.e.Delete)

	//  Обновление записи
	//  Пример:
	//  /user/:id
	//  id-пользователя
	//  body json
	//  {
	//      "name": "данные без цифр и пробелов", - обязательные поля
	//      "surname": "данные без цифр и пробелов", - обязательные поля
	//      "patronymic": "данные без цифр и пробелов", - не обязательное поле
	//      "age":  "данные без букв и пробелов", - не обязательное поле
	//      "Gender":  "данные без цифр и пробелов", - не обязательное поле
	//      "Nationality":  "данные без цифр и пробелов", - не обязательное поле
	//  }

	user.PUT("/:id", a.e.Update)

	//  Сортировка по полю от меньшего к большему
	//  Пример:
	//  /user/sort?sort=age

	user.GET("/sort", a.e.Sort)

	//  Более гибкая фильтрация данных
	//  Пример:
	//  /user/filter?name=eq 'Your value'
	//  Доступные фильтры: eq-равно|ne-не равно|lt-меньше|le-меньше или равно|gt-больше|ge-больше или равно|

	user.GET("/filter", a.e.UserFilter)

	//  Пагинация данных
	//  Пример:
	//  /user/?page=1&limit=3
	//  page-отображаемая страница|limit-количество данных на странице

	user.GET("/", a.e.UserPagination)

	return &a

}

func (a *App) Start() error {
	a.logger.Info("Запуск сервера на localhost:1234")
	err := a.echo.Start(":1234")
	if err != nil {
		return err
	}
	return nil
}
