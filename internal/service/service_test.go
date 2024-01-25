package service_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hulla-hoop/restapi/internal/config"
	"github.com/hulla-hoop/restapi/internal/modeldb"
	"github.com/hulla-hoop/restapi/internal/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var s *service.Service

func TestMain(m *testing.M) {
	os.Setenv("AGEAPI", "https://api.agify.io/?name=%s")
	os.Setenv("NATIONAPI", "https://api.nationalize.io/?name=%s")
	os.Setenv("GENDERAPI", "https://api.genderize.io/?name=%s")
	e := logrus.New()
	cfgT := config.NewCfgApi()
	s = service.New(e, cfgT)
	exitVal := m.Run()
	os.Exit(exitVal)

}

func TestEncrement(t *testing.T) {

	userTest := modeldb.User{
		Name:       "Dmitriy",
		Surname:    "Olegovich",
		Patronymic: "Kirirlov",
	}

	userExpected := modeldb.User{
		Name:        "Dmitriy",
		Surname:     "Olegovich",
		Patronymic:  "Kirirlov",
		Age:         43,
		Gender:      "male",
		Nationality: "UA",
	}

	comparer := cmp.Comparer(func(x, y modeldb.User) bool {
		return x.Name == y.Name && x.Surname == y.Surname && x.Patronymic == y.Patronymic && x.Age == y.Age && x.Gender == y.Gender && x.Nationality == y.Nationality
	})

	err := s.Encriment(&userTest)

	if diff := cmp.Diff(userExpected, userTest, comparer); diff != "" {
		t.Errorf(diff, err)
	}

}

func TestCheckErr(t *testing.T) {

	data := []struct {
		name string
		user modeldb.User
		want error
	}{
		{
			name: "u1",
			user: modeldb.User{Name: "Shamil",
				Surname: "Suleimanov"},
			want: nil,
		},
		{
			name: "u2",
			user: modeldb.User{Name: "Sha111mil",
				Surname: "Suleimanov"},
			want: errors.Errorf("Неверный формат"),
		},
		{
			name: "u3",
			user: modeldb.User{Name: "Shamil",
				Surname: "Suleima1111nov"},
			want: errors.Errorf("Неверный формат"),
		},
		{
			name: "u4",
			user: modeldb.User{Name: "",
				Surname: "Suleima1111nov"},
			want: errors.Errorf("Нет обязательного поля"),
		}, {
			name: "u5",
			user: modeldb.User{Name: "Shamil",
				Surname: ""},
			want: errors.Errorf("Нет обязательного поля"),
		},
		{
			name: "u6",
			user: modeldb.User{Name: "Sham%il",
				Surname: ""},
			want: errors.Errorf("Неверный формат"),
		},
		{
			name: "u7",
			user: modeldb.User{
				Surname: ""},
			want: errors.Errorf("Нет обязательного поля"),
		},
		{
			name: "u8",
			user: modeldb.User{Name: "Shamil"},
			want: errors.Errorf("Нет обязательного поля"),
		}, {
			name: "u9",
			user: modeldb.User{Name: "Sh am l", Surname: "Suleimanov"},
			want: errors.Errorf("Неверный формат поля имя"),
		}, {
			name: "u10",
			user: modeldb.User{Name: "Shamil",
				Surname:    "Suleimanov",
				Patronymic: ""},
			want: nil,
		}, {
			name: "u11",
			user: modeldb.User{Name: "Shamil",
				Surname:    "Suleimanov",
				Patronymic: "Alievich444"},
			want: errors.Errorf("Неверный формат поля отчество"),
		},
		{
			name: "u12",
			user: modeldb.User{Name: " Shamil",
				Surname: "Suleimanov"},
			want: errors.Errorf("Неверный формат поля имя"),
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := s.CheckErr(&d.user)
			if result != nil {
				t.Errorf("Expected  %v , got  %v", d.want, result)
			}
		})
	}

}
