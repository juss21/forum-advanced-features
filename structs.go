package main

import (
	"database/sql"
)

type Memberlist struct {
	ID          int
	Username    string
	Password    string
	Email       string
	DateCreated string
	Session     string
}

type Category struct {
	Id   int
	Name string
}
type Forumdata struct {
	Id          int
	UserId      int
	Author      string
	Title       string
	Content     string
	Date_posted string
	Category    string
	Comments    []Commentdata
	Likes       int
	Dislikes    int
	Loggedin    bool
	Image       string
}
type Commentdata struct {
	Id             int
	Content        string
	UserId         int
	PostId         int
	Date_commented string
	Likes          int
	Dislikes       int
	Username       string
}

type PostLike struct {
	Id     int
	Name   string
	UserId int
	PostId int
}
type Createdstuff struct {
	UserID    int
	PostID    int
	PostTopic string
}
type Likedstuff struct {
	User      string
	PostID    int
	CommentId int
	Title     string
}
type Forumstuff struct {
	Loggedin       bool
	LoggedUser     Memberlist
	CreatedPosts   []Createdstuff
	LikedStuff     []PostLike
	LikedComments  []Likedstuff
	CurrentPost    Forumdata
	Forum_data     []Forumdata
	User_data      []Memberlist
	ErrorMsg       string
	Categories     []Category
	SelectedFilter string
}
type MyError struct{}

var (
	DataBase *sql.DB
	Web      Forumstuff
)
