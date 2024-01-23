package modeldb

import (
	"encoding/json"
	"time"
)

type User struct {
	Id          uint `gorm:"autoIncrement" gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string
	Nationality string
}

func (i User) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

type Users []User

func (i Users) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
