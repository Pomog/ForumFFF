package repository

import "github.com/Pomog/ForumFFF/internal/models"

type DatabaseInt interface {
	UserPresent(userName, email string) (bool, error)
	UserPresentLogin(email, password string) (int, error)
	CreateUser(r models.User) error
	CreateThread(thread models.Thread) error
	CreatePost(post models.Post) error
	IsThreadExist(thread models.Thread) (bool, error)
	GetAllPostsFromThread(threadID int) ([]models.Post, error)
	GetUserByID(ID int) (models.User, error)
	GetAllThreads() ([]models.Thread, error)
	GetThreadByID(ID int) (models.Thread, error)
	GetSessionIDForUserID(userID int) (string, error)
	GetUserIDForSessionID(sessionID string) (int, error) 
	InsertSessionintoDB(sessionID string, userID int) error
}
