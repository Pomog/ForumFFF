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
}

type Thread struct {
	ID      int
	Subject string
	Created time.Time
	UserID  int
}

type ThreadDataForMainPage struct {
	Subject  string
	Created  string
	UserName string
	Picture  string
	Posts    []Post
}

type Post struct {
	ID       int
	Subject  string
	Content  string
	Created  time.Time
	ThreadId int
	UserID   int
}

type Votes struct {
	ID        int
	UpCount   int
	DownCount int
	PostId    int
}
