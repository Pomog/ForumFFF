package repository

import "github.com/Pomog/ForumFFF/internal/models"

type DatabaseInt interface {
	UserPresent(userName, email string) (bool, error)
	CreateUser(r models.User) error
	CreateThread(userID int, thread models.Thread) error
	CreatePost(post models.Post) error
	IsThreadExist(thread models.Thread) (bool, error)
	GetAllPostsFromThread(threadID int) ([]models.Post, error)
}
