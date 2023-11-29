package repository

import "github.com/Pomog/ForumFFF/internal/models"

type Database interface {
	UserPresent(userName, email string) (bool, error)
	CreatetUser(r models.User) error
	CreatetThread(userID int, thread models.Thread) error
}
