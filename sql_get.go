package main

import (
	"database/sql"
)

func buildTopicLikesStruct(allposts int, allusers int) {
	// looping till lenght of commentlikes is equal to allcomments*allusers
	for len(Web.TopicLikes) < Web.allposts*allusers {
		//forum loop
		for forum := 0; forum < len(Web.Forum_data); forum++ {
			//looping all users
			for u := 0; u < len(Web.Userlist); u++ {
				Web.TopicLikes = append(Web.TopicLikes, forumlikes{TopicID: Web.Forum_data[forum].ID, UserID: Web.Userlist[u].ID, Status: 0})
			}
		}
	}
}

func buildTopicDisLikesStruct(allposts int, allusers int) {
	// looping till lenght of commentlikes is equal to allcomments*allusers
	for len(Web.TopicDisLikes) < Web.allposts*allusers {
		//forum loop
		for forum := 0; forum < len(Web.Forum_data); forum++ {
			//looping all users
			for u := 0; u < len(Web.Userlist); u++ {
				Web.TopicDisLikes = append(Web.TopicDisLikes, forumdislikes{TopicID: Web.Forum_data[forum].ID, UserID: Web.Userlist[u].ID, Status: 0})
			}
		}
	}
}
func buildLikesStruct(allcomments int, allusers int) {
	// looping till lenght of commentlikes is equal to allcomments*allusers
	for len(Web.CommentLikes) < Web.allcomments*allusers {
		//forum loopg
		for forum := 0; forum < len(Web.Forum_data); forum++ {
			//looping all comments
			for ac := 0; ac < len(Web.Forum_data[forum].Commentor_data); ac++ {
				//looping all users
				for u := 0; u < len(Web.Userlist); u++ {
					//if Web.Forum_data[forum].Commentor_data[ac].ID == Web.Userlist[u].ID {
					Web.CommentLikes = append(Web.CommentLikes, commentlikes{CommentID: Web.Forum_data[forum].Commentor_data[ac].ID, UserID: Web.Userlist[u].ID, Status: 0})
					//	fmt.Println(Web.CommentLikes)
					//}
				}
			}
		}
	}
}
func buildDisLikesStruct(allcomments int, allusers int) {
	// looping till lenght of commentlikes is equal to allcomments*allusers
	for len(Web.CommentDisLikes) < Web.allcomments*allusers {
		//forum loop
		for forum := 0; forum < len(Web.Forum_data); forum++ {
			//looping all comments
			for ac := 0; ac < len(Web.Forum_data[forum].Commentor_data); ac++ {
				//looping all users
				for u := 0; u < len(Web.Userlist); u++ {
					//if Web.Forum_data[forum].Commentor_data[ac].ID == Web.Userlist[u].ID {
					Web.CommentDisLikes = append(Web.CommentDisLikes, commentdislikes{CommentID: Web.Forum_data[forum].Commentor_data[ac].ID, UserID: Web.Userlist[u].ID, Status: 0})
					//	fmt.Println(Web.CommentLikes)
					//}
				}
			}
		}
	}
}

func saveAllPosts(database *sql.DB) {
	rows, _ := database.Query("SELECT forumid, originalposter, post_header, post_content, likes, dislikes, date FROM forum")
	var op, header, content, date string
	var forumid, likes, disLikes int

	for rows.Next() {
		rows.Scan(&forumid, &op, &header, &content, &likes, &disLikes, &date)
		Web.allposts += 1
		Web.Forum_data = append(Web.Forum_data, forumfamily{ID: forumid, Originalposter: op, Post_title: header, Post_content: content, Post_likes: likes, Post_disLikes: disLikes, Date_posted: date})
	}
}

func saveAllUsers(database *sql.DB) {
	rows, _ := database.Query("SELECT id, username, password, email, datecreated FROM userdata")
	var id int
	var username, password, created, email string
	for rows.Next() {
		rows.Scan(&id, &username, &password, &email, &created)
		Web.Userlist = append(Web.Userlist, memberlist{ID: id, Username: username, Password: password, Email: email, DateCreated: created})
	}
}

func saveAllComments(database *sql.DB) {
	rows, _ := database.Query("SELECT commentid, commentor, comments, post_header, likes, dislikes, date FROM commentdb")
	var commentor, comments, header, date string
	var commentid, likes, disLikes int
	for rows.Next() {
		rows.Scan(&commentid, &commentor, &comments, &header, &likes, &disLikes, &date)
		for i := 0; i < len(Web.Forum_data); i++ {
			if Web.Forum_data[i].Post_title == header {
				Web.allcomments += 1
				Web.Forum_data[i].Commentor_data = append(Web.Forum_data[i].Commentor_data, commentpandemic{ID: commentid, Commentor: commentor, Forum_comment: comments, Post_header: header, Comment_likes: likes, Comment_disLikes: disLikes, Date: date})
			}
		}
	}
}

func saveAllLikes(database *sql.DB) {
	rows, _ := database.Query("SELECT comment_id, user_id, status FROM commentlikes")
	var comment_id, user_id int
	var status int
	for rows.Next() {
		rows.Scan(&comment_id, &user_id, &status)
		Web.CommentLikes = append(Web.CommentLikes, commentlikes{CommentID: comment_id, UserID: user_id, Status: status})
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
