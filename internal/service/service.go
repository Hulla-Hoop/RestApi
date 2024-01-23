package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/hulla-hoop/restapi/internal/config"
	"github.com/hulla-hoop/restapi/internal/modeldb"
	"github.com/pkg/errors"
)

type Service struct {
	errLogger *log.Logger
	cfg       *config.ConfigApi
}

func New(errLogger *log.Logger, cfg *config.ConfigApi) *Service {
	return &Service{
		errLogger: errLogger,
		cfg:       cfg,
	}
}

type Age struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func (s *Service) EncrimentAge(uName *string) (int, error) {
	userAge := new(Age)
	url := (fmt.Sprintf(s.cfg.AGEAPI, *uName))
	r, err := http.Get(url)
	if err != nil {
		s.errLogger.Println("Server is not available. Check connection", err)
		time.Sleep(5 * time.Second)
		age, err := s.EncrimentAge(uName)
		if err != nil {
			return 0, err
		}
		return age, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.errLogger.Println(err)
		return 0, err
	}

	err = json.Unmarshal(body, &userAge)
	if err != nil {
		s.errLogger.Println(err)
		return 0, err
	}

	return userAge.Age, nil
}

type Gender struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (s *Service) EncrimentGender(uName string) (string, error) {
	userGender := Gender{}
	url := (fmt.Sprintf(s.cfg.GENDERAPI, uName))
	r, err := http.Get(url)
	if err != nil {
		s.errLogger.Println("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		name, err := s.EncrimentGender(uName)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.Wrap(err, "Internal server error")
	}

	err = json.Unmarshal(body, &userGender)
	if err != nil {
		return "", errors.Wrap(err, "Internal server error")
	}

	return userGender.Gender, nil
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

func (s *Service) EncrimentCountry(uName string) (string, error) {
	userNati := Natonality{}
	url := (fmt.Sprintf(s.cfg.NATIONAPI, uName))
	r, err := http.Get(url)

	if err != nil {

		s.errLogger.Println("Server is not available. Check connection")
		time.Sleep(5 * time.Second)
		name, err := s.EncrimentCountry(uName)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.errLogger.Println(err)
		return "", errors.Wrap(err, "Internal server error")
	}

	err = json.Unmarshal(body, &userNati)
	if err != nil {
		s.errLogger.Println(err)
		return "", errors.Wrap(err, "Internal server error")
	}

	return userNati.Country[0].CountryId, nil
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

func (s *Service) Encriment(u *modeldb.User) (*modeldb.User, error) {
	err := s.CheckErr(u)
	if err != nil {
		return nil, err
	}
	u.Age, err = s.EncrimentAge(&u.Name)
	if err != nil {
		return nil, err
	}

	u.Gender, err = s.EncrimentGender(u.Name)
	if err != nil {
		return nil, err
	}

	u.Nationality, err = s.EncrimentCountry(u.Name)
	if err != nil {
		return nil, err
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return u, nil
}
