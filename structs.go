package main

import (
	"database/sql"
	"fmt"
	"os"
)

type memberlist struct {
	ID          int
	Username    string
	Password    string
	DateCreated string
	Email       string
}
type forumfamily struct {
	ID             int
	Originalposter string
	Post_title     string
	Post_content   string
	Commentor_data []commentpandemic
	Date_posted    string
	Post_likes     int
	Post_disLikes  int
	Loggedin       bool
	Currentuser    string
}
type commentpandemic struct {
	ID               int
	Date             string
	Commentor        string
	Forum_comment    string
	Post_header      string
	Comment_likes    int
	Comment_disLikes int
}

type commentlikes struct {
	CommentID int
	UserID    int
	Status    int
}
type commentdislikes struct {
	CommentID int
	UserID    int
	Status    int
}

type forumlikes struct {
	TopicID int
	UserID  int
	Status  int
}
type forumdislikes struct {
	TopicID int
	UserID  int
	Status  int
}
type webstuff struct {
	Loggedin        bool
	Currentuser     string
	Currentpage     string
	Sqlbase         *sql.DB
	TopicDisLikes   []forumdislikes
	TopicLikes      []forumlikes
	CommentLikes    []commentlikes
	CommentDisLikes []commentdislikes
	Userlist        []memberlist
	Forum_data      []forumfamily
	allposts        int
	allcomments     int
	tempint         int
	tempcomments    []commentpandemic
	ErrorMsg        string
}

var (
	Web   webstuff
	debug bool = true
)

// function for errorchecking
func errorCheck(err error, exit bool) {
	if err != nil {
		fmt.Println(err)
		if exit {
			os.Exit(0)
		}
	}
}

// function for log printing, can be turned off by setting debug boolean to false
func printLog(a ...any) (n int, err error) {
	if debug {
		return fmt.Fprintln(os.Stdout, a...)
	}
	return 0, nil
}
