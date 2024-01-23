package app

import (
	"log"

	"github.com/hulla-hoop/restapi/internal/DB"
	"github.com/hulla-hoop/restapi/internal/config"
	"github.com/hulla-hoop/restapi/internal/echoendpoint"
	"github.com/hulla-hoop/restapi/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	e         *echoendpoint.Endpoint
	echo      *echo.Echo
	DB        DB.DB
	inflogger *log.Logger
	errLogger *log.Logger
}

func New(db DB.DB, inflogger *log.Logger, errLogger *log.Logger) *App {
	a := App{}
	cfg := config.NewCfgApi()
	service := service.New(errLogger, cfg)
	a.DB = db
	a.errLogger = errLogger
	a.inflogger = inflogger
	a.e = echoendpoint.New(a.DB, a.inflogger, a.errLogger, service)
	a.echo = echo.New()

	a.echo.Use(middleware.Logger())
	a.echo.Use(middleware.Recover())

	a.echo.POST("/user", a.e.Insert)

	// /user/sort?sort=value
	a.echo.GET("/user/sort", a.e.Sort)
	a.echo.DELETE("/user/:id", a.e.Delete)
	a.echo.PUT("/user/:id", a.e.Update)

	// /user/filter?needField=lt|le|gt|ge|eq|ne value
	// в данной примере фильтрация возможна по одному полю и без сортировки результата
	a.echo.GET("/user/filter", a.e.UserFilter)

	// /user/?page=value&limit=value
	a.echo.GET("/user/", a.e.UserPagination)

	return &a

}

func (a *App) Start() {
	a.inflogger.Println("Запуск сервера на localhost:1234")
	a.echo.Start(":1234")
}
