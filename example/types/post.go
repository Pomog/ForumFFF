package types

import (
	"database/sql"
	"time"

	"forum-authentication/config"
)

type Post struct {
	Id       int
	Title    string
	Content  string
	Created  time.Time
	UserId   int
	Username string
}

type PostRating struct {
	Id     int
	PostId int
	UserId int
	Rating int
}

type PostReply struct {
	Id       int
	PostId   int
	UserId   int
	Username string
	Content  string
	Created  time.Time
	Dislikes int
	Likes    int
}

type ReplyRating struct {
	Id          int
	PostReplyId int
	UserId      int
	Rating      int
}

func (p *Post) CreatePost(post Post) (int64, error) {

	insertStmt := `INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)`

	stmt, err := config.DB.Prepare(insertStmt)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(post.Title, post.Content, post.UserId)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (p *Post) GetCategoryPosts(category Categories) ([]Post, error) {

	if category.Id == 0 {
		return nil, nil
	}

	stmt := `
	SELECT posts.*
	FROM posts
	JOIN posts_category ON posts.id = posts_category.post_id
	WHERE posts_category.category_id = ?
	`

	var posts []Post
	res, err := config.DB.Query(stmt, category.Id)
	if err != nil {
		panic(err)
	}

	defer res.Close()

	for res.Next() {
		var post Post

		err = res.Scan(&post.Id, &post.Title, &post.Content, &post.Created, &post.UserId)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Post) GetCategoryLikedPosts(category Categories, user_id int) ([]Post, error) {
	if category.Id == 0 {
		return nil, nil
	}

	stmt := `
	SELECT posts.*
	FROM posts
	JOIN posts_category ON posts.id = posts_category.post_id
	JOIN posts_rating ON posts.id = posts_rating.post_id
	WHERE posts_category.category_id = ? AND posts_rating.user_id = ? AND posts_rating.rating = 1
	`

	var posts []Post
	res, err := config.DB.Query(stmt, category.Id, user_id)
	if err != nil {
		panic(err)
	}

	defer res.Close()

	for res.Next() {
		var post Post
		err = res.Scan(&post.Id, &post.Title, &post.Content, &post.Created, &post.UserId)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Post) GetCategoryCreatedPosts(category Categories, user_id int) ([]Post, error) {
	if category.Id == 0 {
		return nil, nil
	}

	stmt := `
	SELECT posts.*
	FROM posts
	JOIN posts_category ON posts.id = posts_category.post_id
	WHERE posts_category.category_id = ? AND posts.user_id = ?
	`

	var posts []Post
	res, err := config.DB.Query(stmt, category.Id, user_id)
	if err != nil {
		panic(err)
	}

	defer res.Close()

	for res.Next() {
		var post Post
		err = res.Scan(&post.Id, &post.Title, &post.Content, &post.Created, &post.UserId)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *Post) GetPostById(id string) (Post, error) {
	var post Post
	stmt := `
		SELECT posts.*, u.Username
		FROM posts
		JOIN Users u ON posts.user_id = u.id
		WHERE posts.id = ?
	`

	err := config.DB.QueryRow(stmt, id).Scan(&post.Id, &post.Title, &post.Content, &post.Created, &post.UserId, &post.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return post, err
		}
		return post, err
	}
	return post, nil
}

func (p *PostRating) HandlePostRating(id, user_id int, rating string) {
	stmt := `SELECT * FROM posts_rating WHERE post_id = ? AND user_id = ?`
	err := config.DB.QueryRow(stmt, id, user_id).Scan(&p.Id, &p.PostId, &p.UserId, &p.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			p.CreatePostRating(id, user_id, rating)
		}
	}

	p.UpdatePostRating(id, user_id, rating)
}

func (p *ReplyRating) HandleReplyRating(id, user_id int, rating string) {
	stmt := `SELECT * FROM posts_replies_rating WHERE post_reply_id = ? AND user_id = ?`
	err := config.DB.QueryRow(stmt, id, user_id).Scan(&p.Id, &p.PostReplyId, &p.UserId, &p.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			p.CreateReplyRating(id, user_id, rating)
		}
	}

	p.UpdateReplyRating(id, user_id, rating)
}

func (p *ReplyRating) CreateReplyRating(id int, user_id int, rating string) (int64, error) {
	insertStmt := `INSERT INTO posts_replies_rating (post_reply_id, user_id, rating) VALUES (?, ?, ?)`
	stmt, err := config.DB.Prepare(insertStmt)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(id, user_id, rating)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (p *ReplyRating) UpdateReplyRating(id, user_id int, rating string) (int64, error) {
	updateStmt := `UPDATE posts_replies_rating SET rating = ? WHERE post_reply_id = ? AND user_id = ?`

	stmt, err := config.DB.Prepare(updateStmt)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(rating, id, user_id)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (p *PostReply) GetReplyRatings(id int) (int, int, error) {
	/*
		Iterates through all the ratings for a post and returns the number of likes and dislikes
	*/
	stmt := `SELECT * FROM posts_replies_rating WHERE post_reply_id = ?`

	dislikes := 0
	likes := 0

	res, err := config.DB.Query(stmt, id)
	if err != nil {
		panic(err)
	}

	defer res.Close()

	for res.Next() {
		var replyRating ReplyRating
		err = res.Scan(&replyRating.Id, &replyRating.PostReplyId, &replyRating.UserId, &replyRating.Rating)

		if err != nil {
			panic(err)
		}

		if replyRating.Rating == 0 {
			dislikes++
		} else {
			likes++
		}
	}

	return dislikes, likes, err
}

func (p *PostRating) CreatePostRating(id int, user_id int, rating string) (int64, error) {
	insertStmt := `INSERT INTO posts_rating (post_id, user_id, rating) VALUES (?, ?, ?)`

	stmt, err := config.DB.Prepare(insertStmt)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(id, user_id, rating)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (p *PostRating) UpdatePostRating(id, user_id int, rating string) (int64, error) {

	updateStmt := `UPDATE posts_rating SET rating = ? WHERE post_id = ? AND user_id = ?`

	stmt, err := config.DB.Prepare(updateStmt)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(rating, id, user_id)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (p *PostRating) GetPostRatings(id string) (int, int, error) {
	/*
		Iterates through all the ratings for a post and returns the number of likes and dislikes
	*/
	stmt := `SELECT * FROM posts_rating WHERE post_id = ?`

	dislikes := 0
	likes := 0

	res, err := config.DB.Query(stmt, id)
	if err != nil {
		panic(err)
	}

	defer res.Close()

	for res.Next() {
		var postRating PostRating
		err = res.Scan(&postRating.Id, &postRating.PostId, &postRating.UserId, &postRating.Rating)

		if err != nil {
			panic(err)
		}

		if postRating.Rating == 0 {
			dislikes++
		} else {
			likes++
		}
	}

	return dislikes, likes, err
}

func (p *PostReply) CreatePostReply(id int, user_id int, content string) (int64, error) {
	insertStmt := `INSERT INTO posts_replies (post_id, user_id, content, created) VALUES (?, ?, ?, datetime('now', 'localtime'))`

	stmt, err := config.DB.Prepare(insertStmt)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(id, user_id, content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (p *PostReply) GetPostReplies(id string) ([]PostReply, error) {
	/*
		Iterates through all the ratings for a post and returns the number of likes and dislikes
	*/
	stmt := `SELECT pr.*, u.username FROM posts_replies pr JOIN Users u ON pr.user_id = u.id WHERE pr.post_id = ?`

	var postReplies []PostReply
	res, err := config.DB.Query(stmt, id)
	if err != nil {
		panic(err)
	}

	defer res.Close()

	for res.Next() {
		var postReply PostReply

		if err != nil {
			panic(err)
		}

		err = res.Scan(&postReply.Id, &postReply.PostId, &postReply.UserId, &postReply.Content, &postReply.Created, &postReply.Username)
		dislikes, likes, err := p.GetReplyRatings(postReply.Id)
		postReply.Dislikes = dislikes
		postReply.Likes = likes

		if err != nil {
			panic(err)
		}
		postReplies = append(postReplies, postReply)
	}
	return postReplies, err
}
