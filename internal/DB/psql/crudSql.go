package psql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hulla-hoop/restapi/internal/modeldb"
)

func (db *sqlPostgres) Create(user *modeldb.User) (*int, error) {

	var id int
	err := db.dB.QueryRow("INSERT INTO users(created_at,updated_at,name,surname,patronymic,age,gender,nationality) "+
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8) returning id",
		user.CreatedAt,
		user.UpdatedAt,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.Nationality).Scan(&id)
	if err != nil {

		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf("Пользователь добавлен но не удалось записать ID %s", err)
		default:
			return nil, fmt.Errorf("Ошибка при создании пользователя %s", err)
		}
	}

	fmt.Println("id созданого пользователя", id)
	return &id, nil

}

func (db *sqlPostgres) Update(user *modeldb.User, id int) error {
	var patronymic, gender, nationality, age string

	user.UpdatedAt = time.Now()
	if user.Patronymic == "" {
		patronymic = " "
	} else {
		patronymic = fmt.Sprintf("patronymic = '%s',", user.Patronymic)
	}
	if user.Gender == "" {
		gender = " "
	} else {
		gender = fmt.Sprintf("gender = '%s',", user.Gender)
	}
	if user.Nationality == "" {
		nationality = " "
	} else {
		nationality = fmt.Sprintf("nationality = '%s',", user.Nationality)
	}
	if user.Age == 0 {
		age = " "
	} else {
		age = fmt.Sprintf("age = '%d',", user.Age)
	}

	update := fmt.Sprintf("UPDATE users SET updated_at=$1,name=$2, %s %s %s %s surname=$3  WHERE id=$4 ", patronymic, age, gender, nationality)
	_, err := db.dB.Exec(
		update,
		user.UpdatedAt,
		user.Name,
		user.Surname, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *sqlPostgres) Delete(id int) error {
	result, err := db.dB.Exec(
		"DELETE "+
			"FROM users "+
			"WHERE id = $1 ",
		id)
	if err != nil {
		return err
	}
	fmt.Println(result.RowsAffected())
	return nil
}

func (db *sqlPostgres) InsertAll() ([]modeldb.User, error) {
	rows, err := db.dB.Query(
		"SELECT id,name,surname,patronymic,age,gender,nationality " +
			"FROM users ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := []modeldb.User{}

	for rows.Next() {
		u := modeldb.User{}

		err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality)
		if err != nil {
			fmt.Println(err)
			continue
		}
		user = append(user, u)
	}

	return user, nil
}

func (db *sqlPostgres) InsertPage(page uint, limit int) (modeldb.Users, error) {

	cashPage := page*uint(limit) - 1

	rows, err := db.dB.Query(
		"SELECT id,name,surname,patronymic,age,gender,nationality "+
			"FROM users "+
			"WHERE id > $1 "+
			"ORDER BY id ASC "+
			"LIMIT $2 ", cashPage, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := []modeldb.User{}

	for rows.Next() {
		u := modeldb.User{}

		err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cashPage = u.Id
		user = append(user, u)
	}

	return user, nil
}

func (db *sqlPostgres) Sort(field string) ([]modeldb.User, error) {

	query := "SELECT id,name,surname,patronymic,age,gender,nationality " +
		"FROM users " +
		"ORDER BY %s"

	queryR := fmt.Sprintf(query, field)

	rows, err := db.dB.Query(queryR)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := []modeldb.User{}

	for rows.Next() {
		u := modeldb.User{}

		err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality)
		if err != nil {
			fmt.Println(err)
			continue
		}
		user = append(user, u)
	}

	return user, nil
}

func (db *sqlPostgres) Filter(field string, operator string, value string) ([]modeldb.User, error) {

	query := "SELECT id,name,surname,patronymic,age,gender,nationality " +
		"FROM users " +
		"WHERE %s %s %s "

	operatorMap := make(map[string]string)

	operatorMap["eq"] = "="
	operatorMap["ne"] = "!="
	operatorMap["gt"] = ">"
	operatorMap["ge"] = ">="
	operatorMap["lt"] = "<"
	operatorMap["le"] = "<="

	q := fmt.Sprintf(query, field, operatorMap[operator], value)

	rows, err := db.dB.Query(q)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := []modeldb.User{}

	for rows.Next() {
		u := modeldb.User{}

		err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality)
		if err != nil {
			fmt.Println(err)
			continue
		}
		user = append(user, u)
	}

	return user, nil

}
