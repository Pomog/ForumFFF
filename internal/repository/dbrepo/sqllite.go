package dbrepo

import (
	"context"
	"time"

	"github.com/Pomog/ForumFFF/internal/models"
)

func (m *sqliteBDRepo) UserPresent(userName, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) 
	from users
	where username = $1 and
	email = $2
	`
	var numRows int
	row := m.DB.QueryRowContext(ctx, query, userName, email)

	err := row.Scan(&numRows)
	if err != nil {
		return false, nil
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

func (m *sqliteBDRepo) InsertUser(r models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into users
	(username, password, first_name, last_name, email, picture)
	values ($1, $2, $3, $4, $5, $6)
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		r.UserName,
		r.Password,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Picture,
	)
	if err != nil {
		return err
	}
	return nil
}

/*
var userTable = `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username varchar(100) DEFAULT "",
    password varchar(100) DEFAULT "",
    first_name varchar(100) DEFAULT "",
    last_name varchar(100) DEFAULT "",
    email varchar(254) DEFAULT "",
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    picture TEXT DEFAULT "static/ava/pomog_ava.png",
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`
*/
