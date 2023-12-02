package dbrepo

import (
	"context"
	"time"

	"github.com/Pomog/ForumFFF/internal/models"
)

func (m *SqliteBDRepo) UserPresent(userName, email string) (bool, error) {
func (m *SqliteBDRepo) UserPresent(userName, email string) (bool, error) {
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

func (m *SqliteBDRepo) CreateUser(r models.User) error {
func (m *SqliteBDRepo) GetUserByID(ID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * 
	from users
	where id = $1
	`
	var user models.User

	row := m.DB.QueryRowContext(ctx, query, ID)

	err := row.Scan(&user.ID, &user.UserName, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.Created, &user.Picture, &user.LastActivity)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *SqliteBDRepo) CreateUser(r models.User) error {
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

func (m *SqliteBDRepo) CreateThread(userID int, thread models.Thread) error {
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

// CreatePost insert post into SQLite DB
func (m *SqliteBDRepo) CreatePost(post models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into post
	(subject, content, threadID, userID)
	values ($1, $2, $3, $4
	)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		post.Subject,
		post.Content,
		post.ThreadId,
		post.UserID,
	)

	if err != nil {
		return err
	}
	return nil
}

// isThreadExist returns true if Thread with same Subject exist
func (m *SqliteBDRepo) IsThreadExist(thread models.Thread) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) 
	from thread
	where subject = $1
	`
	var numRows int
	row := m.DB.QueryRowContext(ctx, query, thread.Subject)

	err := row.Scan(&numRows)
	if err != nil {
		return false, nil
	}

	if numRows == 0 {
		return false, nil
	}

	return true, nil
}

// GetAllPostsFromThread returns all slice of all Posts with given ThreadID, nil if there are no Posts
func (m *SqliteBDRepo) GetAllPostsFromThread(threadID int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * 
	from post
	where threadID = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Subject, &post.Content, &post.Created, &post.ThreadId, &post.UserID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetAllThreads returns all Threads, nil if there are no threads in DB
func (m *SqliteBDRepo) GetAllThreads() ([]models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * 
	from thread
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread

	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(&thread.ID, &thread.Subject, &thread.Created, &thread.UserID)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}
