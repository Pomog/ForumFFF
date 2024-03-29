package repository

import "github.com/Pomog/ForumFFF/internal/models"

type DatabaseInt interface {
	UserPresent(userName, email string) (bool, error)
	UserPresentLogin(email, password string) (int, error)
	CreateUser(r models.User) error
	CreateThread(thread models.Thread) (int64, error)
	CreatePost(post models.Post) error
	IsThreadExist(thread models.Thread) (bool, error)
	GetAllPostsFromThread(threadID int) ([]models.Post, error)
	GetUserByID(ID int) (models.User, error)
	GetAllThreads() ([]models.Thread, error)
	GetThreadByID(ID int) (models.Thread, error)
	GetSessionIDForUserID(userID int) (string, error)
	GetUserIDForSessionID(sessionID string) (int, error)
	InsertSessionintoDB(sessionID string, userID int) error
	GetTotalPostsAmmountByUserID(userID int) (int, error)
	LikePostByUserIdAndPostId(userID, postID int) error
	DislikePostByUserIdAndPostId(userID, postID int) error
	CountLikesAndDislikesForPostByPostID(postID int) (likes, dislikes int, err error)
	GetGuestID() (int, error)
	GetSearchedThreads(search string) ([]models.Thread, error)
	GetPostByID(ID int) (models.Post, error)
	EditPost(post models.Post) error
	EditTopic(topic models.Thread) error
	DeletePost(post models.Post) error
	GetSearchedThreadsByCategory(search string) ([]models.Thread, error)
	GetAllThreadsByUserID(userID int) ([]models.Thread, error)
	GetAllPostsByUserID(userID int) ([]models.Post, error)
	GetAllLikedPostsByUserID(userID int) ([]models.Post, error)
	EditUserType(user models.User) error
	DelSessionByUserID(userID int) error
	EditPostClassification(post models.Post, classification models.TextClassification) error
	EditTopicClassification(topic models.Thread, classification models.TextClassification) error
	GetAllPostsByClassification(classification models.TextClassification) ([]models.Post, error)
	GetAllThreadsByClassification(classification models.TextClassification) ([]models.Thread, error)
	CreatePM(pm models.PM) error
	DeletePM(pm models.PM) error
	GetPMbyReceiverUserID(userID int) ([]models.PM, error)
	GetPMbysenderUserID(userID int) ([]models.PM, error)
}
