package echoendpoint

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hulla-hoop/restapi/internal/DB"
	"github.com/hulla-hoop/restapi/internal/modeldb"
	"github.com/hulla-hoop/restapi/internal/service"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	Db      DB.DB
	logger  *logrus.Logger
	service *service.Service
}

func New(db DB.DB, logger *logrus.Logger, service *service.Service) *Endpoint {
	return &Endpoint{Db: db,
		logger:  logger,
		service: service}
}

func (e *Endpoint) Insert(c echo.Context) error {
	u := new(modeldb.User)
	err := c.Bind(u)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = e.service.Encriment(u)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := e.Db.Create(u)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	u.Id = uint(*id)
	e.logger.Info("Пользователь добавлен")
	return c.JSON(http.StatusCreated, u)
}

func (e *Endpoint) Delete(c echo.Context) error {
	id := c.Param("id")
	idi, err := strconv.Atoi(id)
	if err != nil {
		e.logger.Error(err)
	}
	err = e.Db.Delete(idi)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	e.logger.Info("Пользователь удален")
	return c.JSON(http.StatusAccepted, "Пользователь удален")
}

func (e *Endpoint) Update(c echo.Context) error {
	u := new(modeldb.User)
	err := c.Bind(u)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = e.service.CheckErr(u)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	id, _ := strconv.Atoi(c.Param("id"))
	err = e.Db.Update(u, id)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	e.logger.Info("Пользователь изменен")
	return c.JSON(http.StatusOK, "Пользователь изменен")
}

func (e *Endpoint) Sort(c echo.Context) error {
	users := []modeldb.User{}

	valueStr, err := c.FormParams()
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	field := valueStr["sort"]
	e.logger.Debug(field[0])
	users, err = e.Db.Sort(field[0])
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

//Пагинация через Limit без Offset

func (e *Endpoint) UserPagination(c echo.Context) error {
	valueStr, err := c.FormParams()
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pageStr := valueStr["page"]
	limitStr := valueStr["limit"]

	e.logger.Debug(pageStr, limitStr)

	page, err := strconv.Atoi(pageStr[0])
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	limit, err := strconv.Atoi(limitStr[0])
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	u, err := e.Db.InsertPage(uint(page)-1, limit)
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, u)

}

func (e *Endpoint) UserFilter(c echo.Context) error {

	valueStr, err := c.FormParams()
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var validParam string

	for p := range valueStr {
		if p != "age" && p != "name" && p != "surname" && p != "patronumic" && p != "gender" && p != "nationality" && p != "id" {
			continue
		} else {
			validParam = p
		}

	}
	param := strings.Split(valueStr[validParam][0], " ")
	u, err := e.Db.Filter(validParam, param[0], param[1])
	if err != nil {
		e.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, u)
}
