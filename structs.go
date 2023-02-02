package main

import "database/sql"

type Memberlist struct {
	ID       int
	Username string
	Password string
	Email    string
}
type Forumdata struct {
	Id          int
	Author      string
	Title       string
	Content     string
	Date_posted string
	Comments    []Commentdata
	Likes       int
	Dislikes    int
	Loggedin    bool
}
type Commentdata struct {
	Id       int
	Content  string
	UserId   int
	PostId   int
	Likes    int
	Dislikes int
}

type PostLike struct {
	Id     int
	Name   string
	UserId int
	PostId int
}

type Forumstuff struct {
	Loggedin    bool
	LoggedUser  Memberlist
	CurrentPost Forumdata
	Forum_data  []Forumdata
	User_data   []Memberlist
	ErrorMsg    string
}

var DataBase *sql.DB
var Web Forumstuff
