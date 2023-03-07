package app

import (
	"database/sql"
	"html/template"
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
	Id   string
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
	Edit        bool
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
	Edit           bool
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
	CurrentComment Commentdata
	Forum_data     []Forumdata
	User_data      []Memberlist
	ErrorMsg       string
	ErrorPage      *template.Template
	Categories     []Category
	SelectedFilter string
	Notifications  []Notifications
}
type Notifications struct {
	UserID  int
	PostID  int
	User    string
	Content string
}

type MyError struct{}

var (
	DataBase           *sql.DB
	Web                Forumstuff
	userRequestAmounts = make(map[string]int)
)
