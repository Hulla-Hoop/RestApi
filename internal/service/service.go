package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/hulla-hoop/restapi/internal/config"
	"github.com/hulla-hoop/restapi/internal/modeldb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	cfg    *config.ConfigApi
}

func New(errLogger *logrus.Logger, cfg *config.ConfigApi) *Service {
	return &Service{
		logger: errLogger,
		cfg:    cfg,
	}
}

type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (s *Service) EncrimentAge(uName *string) (*int, error) {
	userAge := new(Age)

	url := (fmt.Sprintf(s.cfg.AGEAPI, *uName))
	r, err := http.Get(url)
	if err != nil {
		s.logger.Error("Server is not available. Check connection", err)
		time.Sleep(5 * time.Second)
		age, err := s.EncrimentAge(uName)
		if err != nil {
			return nil, err
		}
		return age, nil
	}
	s.logger.Debug("Полученные данные", r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	err = json.Unmarshal(body, &userAge)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &userAge.Age, nil
}

type Gender struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (s *Service) EncrimentGender(uName *string) (*string, error) {
	userGender := Gender{}
	url := (fmt.Sprintf(s.cfg.GENDERAPI, *uName))
	r, err := http.Get(url)
	if err != nil {
		s.logger.Error("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		name, err := s.EncrimentGender(uName)
		if err != nil {
			return nil, err
		}
		return name, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Internal server error")
	}

	err = json.Unmarshal(body, &userGender)
	if err != nil {
		return nil, errors.Wrap(err, "Internal server error")
	}

	return &userGender.Gender, nil
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Natonality struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []Country
}

func (s *Service) EncrimentCountry(uName *string) (*string, error) {
	userNati := Natonality{}
	url := (fmt.Sprintf(s.cfg.NATIONAPI, *uName))
	r, err := http.Get(url)

	if err != nil {

		s.logger.Error("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		name, err := s.EncrimentCountry(uName)
		if err != nil {
			return nil, err
		}
		return name, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error(err)
		return nil, errors.Wrap(err, "Internal server error")
	}

	err = json.Unmarshal(body, &userNati)
	if err != nil {
		s.logger.Error(err)
		return nil, errors.Wrap(err, "Internal server error")
	}

	return &userNati.Country[0].CountryId, nil
}

func (s *Service) CheckErr(U *modeldb.User) error {
	if U.Name == "" || U.Surname == "" {

		return fmt.Errorf("Нет обязательного поля")
	}

	r, err := regexp.MatchString("^[a-zA-Z]+$", U.Name)
	if err != nil {
		return err
	}
	if r == false {
		return fmt.Errorf("Неверный формат поля имя")
	}

	r, err = regexp.MatchString("^[a-zA-Z]+$", U.Surname)
	if err != nil {
		return err
	}
	if r == false {
		return fmt.Errorf("Неверный формат поля фамилия")
	}
	if U.Patronymic == "" {
		return nil
	} else {
		r, err = regexp.MatchString("^[a-zA-Z]+$", U.Patronymic)
		if err != nil {
			return err
		}
		if r == false {
			return fmt.Errorf("Неверный формат поля отчество")
		}
	}

	return nil

}

func (s *Service) Encriment(u *modeldb.User) error {
	err := s.CheckErr(u)
	if err != nil {
		return err
	}
	age, err := s.EncrimentAge(&u.Name)
	if err != nil {
		return err
	}
	u.Age = *age

	gender, err := s.EncrimentGender(&u.Name)
	if err != nil {
		return err
	}
	u.Gender = *gender

	nationality, err := s.EncrimentCountry(&u.Name)
	if err != nil {
		return err
	}
	u.Nationality = *nationality
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
