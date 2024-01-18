package models

import (
	"time"
)

type User struct {
	ID           int
	UserName     string
	Password     string
	FirstName    string
	LastName     string
	Email        string
	Created      time.Time
	Picture      string
	LastActivity time.Time
	Type         string
}

type Thread struct {
	ID             int
	Subject        string
	Created        time.Time
	UserID         int
	Image          string
	Category       string
	Classification string
}

type ThreadDataForMainPage struct {
	Subject                       string
	Created                       string
	UserNameWhoCreatedThread      string
	UserNameWhoCreatedLastPost    string
	PictureUserWhoCreatedThread   string
	PictureUserWhoCreatedLastPost string
	Posts                         []Post
	ThreadID                      int
	Image                         string
	Category                      string
	Classification                string
	UserID                        int
}

type PostDataForThemePage struct {
	ID                        int
	Subject                   string
	Content                   string
	Image                     string
	Created                   string
	UserNameWhoCreatedPost    string
	UserIDWhoCreatedPost      int
	PictureUserWhoCreatedPost string
	UserRegistrationDate      string
	UserPostsAmmount          int
	Likes                     int
	Dislikes                  int
	Classification            string
}

type Post struct {
	ID             int
	Subject        string
	Content        string
	Created        time.Time
	ThreadId       int
	UserID         int
	Image          string
	Classification string
}

type Votes struct {
	ID        int
	UpCount   int
	DownCount int
	PostId    int
}

// CommentType represents the type of a comment
type TextClassification string

var Classifications = []TextClassification{"irrelevant",
	"obscene",
	"illegal",
	"insulting",
	"unsorted",
	"approved",
}

type PM struct {
	ID             int
	Content        string
	Created        time.Time
	SenderUserID   int
	ReceiverUserID int
}
