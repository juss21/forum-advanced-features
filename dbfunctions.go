package main

import (
	"database/sql"
	"fmt"
	"time"
)

func sendPost(database *sql.DB, originalposter string, header string, content string) (bool, string) {
	canSendContent := false
	for i := 0; i < len(forum_data); i++ {
		if forum_data[i].Post_title == header {
			return false, "Post with the title '" + header + "' already exists!"
		}
	}
	for c := 0; c < len(content); c++ {
		if rune(content[c]) >= 32 && rune(content[c]) <= 126 {
			canSendContent = true
		}
	}
	if !canSendContent {
		return false, "Content cannot be empty!"
	}
	if len(content) < 10 {
		return false, "Too few characters in content, spam detected!"
	}
	currentTime := time.Now().Format("02.01.2006 15:04")

	likes := 0
	dislikes := 0

	forum_data = append(forum_data, forumfamily{Originalposter: originalposter, Post_title: header, Post_content: content, Date_posted: currentTime, Post_likes: 0, Post_disLikes: 0})

	statement, _ := database.Prepare("INSERT INTO forum (originalposter, post_header, post_content, likes, dislikes) VALUES (?,?,?,?,?)")
	statement.Exec(originalposter, header, content, likes, dislikes) // exec first name, last name

	return true, ""
}

func sendComment(database *sql.DB, commenter string, forum_Commentbox string, forum_header string) {
	fmt.Println(forum_header, ">", commenter, "lisas kommentaari:", forum_Commentbox)

	statement, _ := database.Prepare("INSERT INTO commentdb (commentor, forum_comments, post_header, likes, dislikes, date) VALUES (?,?,?,?,?,?)")
	likes := 0
	dislikes := 0
	currentTime := time.Now().Format("02.01.2006 15:04")

	if len(forum_Commentbox) == 0 {
		return
	}

	for i := 0; i < len(forum_data); i++ {
		if forum_data[i].Post_title == forum_header {
			forum_data[i].Commentor_data = append(forum_data[i].Commentor_data, commentpandemic{Commentor: commenter, Forum_comment: forum_Commentbox, Post_header: forum_header, Comment_likes: 0, Comment_disLikes: 0, Date: currentTime})
		}
	}
	statement.Exec(commenter, forum_Commentbox, forum_header, likes, dislikes, currentTime) // exec first name, last name

}

func sendRegister(database *sql.DB, username string, password string, email string) {
	statement, _ := database.Prepare("INSERT INTO userdata (username, password, email) VALUES (?,?,?)")
	statement.Exec(username, password, email) // exec first name, last name
	userlist = append(userlist, memberlist{ID: len(userlist) + 1, Username: username, Password: password, Email: email})
	//fmt.Println("Server:", username, "has successfully registered!", " <", email, ">")

	// kasutajate printimine konsooli
	// for i := 0; i < len(userlist); i++ {
	// 	fmt.Println(userlist[i])
	// }
}

func saveAllPosts(database *sql.DB) {
	rows, _ := database.Query("SELECT originalposter, post_header, post_content, likes, dislikes FROM forum")
	var op string
	var header string
	var content string
	var likes int
	var disLikes int
	for rows.Next() {
		rows.Scan(&op, &header, &content, &likes, &disLikes)
		forum_data = append(forum_data, forumfamily{Originalposter: op, Post_title: header, Post_content: content, Post_likes: likes, Post_disLikes: disLikes})
	}
}

func saveAllUsers(database *sql.DB) {
	rows, _ := database.Query("SELECT id, username, password, email FROM userdata")
	var id int
	var username string
	var password string
	var email string

	for rows.Next() {
		rows.Scan(&id, &username, &password, &email)
		userlist = append(userlist, memberlist{ID: id, Username: username, Password: password, Email: email})
	}
}

func saveAllComments(database *sql.DB) {
	rows, _ := database.Query("SELECT commentor, forum_comments, post_header, likes, dislikes, date FROM commentdb")
	var commentor string
	var comments string
	var header string
	var likes int
	var disLikes int
	var date string
	for rows.Next() {
		rows.Scan(&commentor, &comments, &header, &likes, &disLikes, &date)
		for i := 0; i < len(forum_data); i++ {
			if forum_data[i].Post_title == header {
				forum_data[i].Commentor_data = append(forum_data[i].Commentor_data, commentpandemic{Commentor: commentor, Forum_comment: comments, Post_header: header, Comment_likes: likes, Comment_disLikes: disLikes, Date: date})
			}
		}
	}
}
