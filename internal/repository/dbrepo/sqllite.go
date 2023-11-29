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
		return false, nil
	}

	return true, nil
}

func (m *sqliteBDRepo) CreatetUser(r models.User) error {
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

func (m *sqliteBDRepo) CreatetThread(userID int, thread models.Thread) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into thread
	(subject, userID)
	values ($1, $2, )
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		thread.Subject,
		userID,
	)

	if err != nil {
		return err
	}
	return nil
}
