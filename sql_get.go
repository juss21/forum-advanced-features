package main

import (
	"database/sql"
)

func saveAllPosts(database *sql.DB) {
	rows, _ := database.Query("SELECT originalposter, post_header, post_content, likes, dislikes, date FROM forum")
	var op, header, content, date string
	var likes, disLikes int

	for rows.Next() {
		rows.Scan(&op, &header, &content, &likes, &disLikes, &date)
		Web.Forum_data = append(Web.Forum_data, forumfamily{Originalposter: op, Post_title: header, Post_content: content, Post_likes: likes, Post_disLikes: disLikes, Date_posted: date})
	}
}

func saveAllUsers(database *sql.DB) {
	rows, _ := database.Query("SELECT id, username, password, email, likedcontent, dislikedcontent,datecreated FROM userdata")
	var id int
	var username, password, created, email, likedcontent, dislikedcontent string
	for rows.Next() {
		rows.Scan(&id, &username, &password, &email, &likedcontent, &dislikedcontent, &created)
		Web.Userlist = append(Web.Userlist, memberlist{ID: id, Username: username, Password: password, Email: email, Likedcontent: likedcontent, Dislikedcontent: dislikedcontent, DateCreated: created})
	}
}

func saveAllComments(database *sql.DB) {
	rows, _ := database.Query("SELECT commentid,commentor, forum_comments, post_header, likes, dislikes, date, likedbyusers, dislikedbyusers FROM commentdb")
	var commentor, comments, header, date, likedBy, disLikedBy string
	var commentid, likes, disLikes int
	for rows.Next() {
		rows.Scan(&commentid, &commentor, &comments, &header, &likes, &disLikes, &date, &likedBy, &disLikedBy)
		for i := 0; i < len(Web.Forum_data); i++ {
			if Web.Forum_data[i].Post_title == header {
				Web.allcomments += 1
				Web.Forum_data[i].Commentor_data = append(Web.Forum_data[i].Commentor_data, commentpandemic{ID: commentid, Commentor: commentor, Forum_comment: comments, Post_header: header, Comment_likes: likes, Comment_disLikes: disLikes, Date: date, Likedby: likedBy, Dislikedby: disLikedBy})
			}
		}
	}
}

func getUserID(username string) int {
	for i := 0; i < len(Web.Userlist); i++ {
		if Web.Userlist[i].Username == username {
			return i
		}
	}
	return 0
}

func getTopicID(topic string) int {
	for i := 0; i < len(Web.Userlist); i++ {
		if Web.Forum_data[i].Post_title == topic {
			return i
		}
	}
	return 0
}
