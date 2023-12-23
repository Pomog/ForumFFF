package dbrepo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Pomog/ForumFFF/internal/models"
)

func (m *SqliteBDRepo) UserPresent(userName, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) 
	from users
	where username = $1 OR
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

func (m *SqliteBDRepo) UserPresentLogin(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id 
	from users
	where email = $1 and
	password = $2
	`

	row := m.DB.QueryRowContext(ctx, query, email, password)

	userID := 0
	err := row.Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

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

	if user.ID == 0 || user.UserName == "" || user.Email == "" {
		return user, errors.New("wrong User Data")
	}

	if user.ID == 0 || user.UserName == "" || user.Email == "" {
		return user, errors.New("wrong User Data")
	}

	return user, nil
}

func (m *SqliteBDRepo) GetThreadByID(ID int) (models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * 
	from thread
	where id = $1
	`
	var thread models.Thread

	row := m.DB.QueryRowContext(ctx, query, ID)

	err := row.Scan(&thread.ID, &thread.Subject, &thread.Created, &thread.UserID, &thread.Image)
	if err != nil {
		return thread, err
	}
	return thread, nil
}

func (m *SqliteBDRepo) GetPostByID(ID int) (models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * 
	from post
	where id = $1
	`
	var post models.Post

	row := m.DB.QueryRowContext(ctx, query, ID)

	err := row.Scan(&post.ID, &post.Subject, &post.Content, &post.Created, &post.ThreadId, &post.UserID, &post.Image)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (m *SqliteBDRepo) CreateUser(r models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if r.UserName == "" || r.Email == "" || r.FirstName == "" {
		return errors.New("wrong User Data")
	}

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

func (m *SqliteBDRepo) CreateThread(thread models.Thread) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if strings.TrimSpace(thread.Subject) == "" {
		return 0, errors.New("empty thread can not be created")
	}

	if len(thread.Subject) > 1500 {
		return 0, errors.New("the text is to long, 1500 syblos allowed")
	}

	user, err := m.GetUserByID(thread.UserID)
	if err != nil {
		return 0, errors.New("guest can not create a thread")
	}
	userName := user.UserName

	if userName == "guest" || strings.TrimSpace(userName) == "" {
		return 0, errors.New("guest can not create a thread")
	}

	stmt := `insert into thread
	(subject, userID, threadImage)
	values ($1, $2, $3)
	`

	sqlRes, err := m.DB.ExecContext(ctx, stmt,
		thread.Subject,
		thread.UserID,
		thread.Image,
	)

	if err != nil {
		return 0, err
	}
	id, err := sqlRes.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreatePost insert post into SQLite DB
func (m *SqliteBDRepo) CreatePost(post models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if strings.TrimSpace(post.Content) == "" {
		return errors.New("empty post can not be created")
	}

	if len(post.Content) > 500 {
		return errors.New("the post is to long, 500 syblos allowed")
	}

	stmt := `insert into post
	(subject, content, threadID, userID, postImage)
	values ($1, $2, $3, $4, $5
	)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		post.Subject,
		post.Content,
		post.ThreadId,
		post.UserID,
		post.Image,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *SqliteBDRepo) EditPost(post models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if strings.TrimSpace(post.Content) == "" {
		return errors.New("empty post can not be created")
	}

	if len(post.Content) > 2500 {
		return errors.New("the post is to long, 2500 syblos allowed")
	}

	stmt := `UPDATE post
	SET subject = $1, content = $2, threadID = $3, userID = $4
	WHERE id = $5;
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		post.Subject,
		post.Content,
		post.ThreadId,
		post.UserID,
		post.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *SqliteBDRepo) DeletePost(post models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM post
	WHERE id = $1;
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		post.ID,
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
		err := rows.Scan(&post.ID, &post.Subject, &post.Content, &post.Created, &post.ThreadId, &post.UserID, &post.Image)
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
		err := rows.Scan(&thread.ID, &thread.Subject, &thread.Created, &thread.UserID, &thread.Image)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}

func (m *SqliteBDRepo) GetSessionIDForUserID(userID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select sessionID 
	from sessionId WHERE
	userID = $1
	`
	sessionID := ""

	row := m.DB.QueryRowContext(ctx, query, userID)

	err := row.Scan(&sessionID)
	if err != nil {
		return sessionID, err
	}

	return sessionID, nil
}

func (m *SqliteBDRepo) GetUserIDForSessionID(sessionID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select userID 
	from sessionId WHERE
	sessionID = $1
	`
	var userID int

	row := m.DB.QueryRowContext(ctx, query, sessionID)

	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (m *SqliteBDRepo) InsertSessionintoDB(sessionID string, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into sessionId
	(sessionID, userID)
	values ($1, $2
	)
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		sessionID,
		userID,
	)

	if err != nil {
		return err
	}
	return nil

}

func (m *SqliteBDRepo) GetTotalPostsAmmountByUserID(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(*) from post where userID = $1;	`
	var numberOfPosts int

	row := m.DB.QueryRowContext(ctx, query, userID)

	err := row.Scan(&numberOfPosts)
	if err != nil {
		return 0, err
	}

	return numberOfPosts, nil
}

func (m *SqliteBDRepo) LikePostByUserIdAndPostId(userID, postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//reverse posts like flag
	stmt := `INSERT OR REPLACE INTO votes (id, like, dislike, postID, userID)
	VALUES (
		COALESCE((SELECT id FROM votes WHERE userID = $1 AND postID = $2), NULL),
		1, -- Setting like to true
		CASE WHEN NOT COALESCE((SELECT like FROM votes WHERE userID = $1 AND postID = $2), 0)
		THEN 0 ELSE COALESCE((SELECT dislike FROM votes WHERE userID = $1 AND postID = $2), 0) END,
		$2,	$1
	);
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		userID,
		postID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *SqliteBDRepo) DislikePostByUserIdAndPostId(userID, postID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//reverse posts dislike flag
	stmt := `INSERT OR REPLACE INTO votes (id, dislike, like, postID, userID)
	VALUES (
		COALESCE((SELECT id FROM votes WHERE userID = $1 AND postID = $2), NULL),
		1, -- Setting dislike to true
		CASE WHEN NOT COALESCE((SELECT dislike FROM votes WHERE userID = $1 AND postID = $2), 0)
		THEN 0 
		ELSE COALESCE((SELECT like FROM votes WHERE userID = $1 AND postID = $2), 0)
		END,
		$2,	$1
	);
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		userID,
		postID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *SqliteBDRepo) CountLikesAndDislikesForPostByPostID(postID int) (likes, dislikes int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT 
    SUM(CASE WHEN postID = $1 AND like = TRUE THEN 1 ELSE 0 END) AS like_count,
    SUM(CASE WHEN postID = $1 AND dislike = TRUE THEN 1 ELSE 0 END) AS dislike_count
	FROM votes;
	`

	row := m.DB.QueryRowContext(ctx, query, postID)

	likes = 0
	dislikes = 0

	row.Scan(&likes, &dislikes)

	err = nil

	return
}

func (m *SqliteBDRepo) GetGuestID() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id 
	from users
	where email = 'guest@gmail.com'
	`
	var guestID int

	row := m.DB.QueryRowContext(ctx, query)

	err := row.Scan(&guestID)
	if err != nil {
		return guestID, err
	}
	return guestID, nil
}

func (m *SqliteBDRepo) GetSearchedThreads(search string) ([]models.Thread, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * 
	from thread
	WHERE subject LIKE '%' || $1 || '%';
	`
	rows, err := m.DB.QueryContext(ctx, query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []models.Thread

	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(&thread.ID, &thread.Subject, &thread.Created, &thread.UserID, &thread.Image)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}
