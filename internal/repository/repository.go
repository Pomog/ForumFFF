package repasitory

import "github.com/Pomog/ForumFFF/internal/models"

type Database interface {
	UserPresent(userName, email string) (bool, error)
	InsertUser(r models.User) error
}
