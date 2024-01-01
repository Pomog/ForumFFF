package types

import (
	"database/sql"

	"forum-authentication/config"
)

type Categories struct {
	Id        int
	Name      string
	Name_slug string
}

type PostCategories struct {
	Id         int
	PostId     int
	CategoryId int
}

func (c *Categories) GetCategories() ([]Categories, error) {
	var categories []Categories
	stmt := `SELECT * FROM categories`

	res, err := config.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var category Categories
		if err := res.Scan(&category.Id, &category.Name, &category.Name_slug); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *Categories) GetCategoryBySlug(slug string) ([]Categories, error) {
	var categories []Categories
	stmt := `SELECT * FROM categories WHERE name_slug = ?`

	err := config.DB.QueryRow(stmt, slug).Scan(&c.Id, &c.Name, &c.Name_slug)

	categories = append(categories, *c)

	if err != nil {
		if err == sql.ErrNoRows {
			return categories, err
		}
		return categories, err
	}
	return categories, nil
}

func (c *Categories) GetCurrentCategory(cat string) (Categories, error) {
	if cat != "" {
		categories, err := c.GetCategoryBySlug(cat)
		if err != nil || len(categories) == 0 {
			return *c, err
		}
		return categories[0], nil
	}
	return *c, nil
}

func (c *Categories) CreatePostCategory(postCategories *PostCategories) (int64, error) {
	insertStmt := `INSERT INTO posts_category (post_id, category_id) VALUES (?, ?)`

	stmt, err := config.DB.Prepare(insertStmt)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(postCategories.PostId, postCategories.CategoryId)
	if err != nil {
		return 0, err
	}

	postCategoryID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postCategoryID, nil
}
