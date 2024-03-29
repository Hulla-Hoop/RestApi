package DB

import "github.com/hulla-hoop/restapi/internal/modeldb"

type DB interface {
	Create(user *modeldb.User) (*int, error)
	Delete(id int) error
	Update(user *modeldb.User, id int) error
	InsertPage(page uint, limit int) (modeldb.Users, error)
	Sort(field string) ([]modeldb.User, error)
	Filter(field string, operator string, value string) ([]modeldb.User, error)
}
