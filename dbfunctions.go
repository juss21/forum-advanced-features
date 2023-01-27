package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func sendPost(database *sql.DB, originalposter string, header string, content string) (bool, string) {
	canSendContent := false
	for i := 0; i < len(Web.Forum_data); i++ {
		if Web.Forum_data[i].Post_title == header {
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

	Web.Forum_data = append(Web.Forum_data, forumfamily{Originalposter: originalposter, Post_title: header, Post_content: content, Date_posted: currentTime, Post_likes: 0, Post_disLikes: 0})

	statement, _ := database.Prepare("INSERT INTO forum (originalposter, post_header, post_content, likes, dislikes, date) VALUES (?,?,?,?,?,?)")
	statement.Exec(originalposter, header, content, likes, dislikes, currentTime) // exec first name, last name

	return true, ""
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

func sendLike(database *sql.DB, user string, topic string, comment bool) {
	if !comment {
		fmt.Println(topic, ">", user, "pani meeldivaks.")
		userid := getUserID(user)
		topicid := getTopicID(topic)
		//kui ei ole märgitud ei meeldivaks ega ebameeldivaks
		if !strings.Contains(Web.Userlist[userid].Likedcontent, topic) && !strings.Contains(Web.Userlist[userid].Dislikedcontent, topic) {

			Web.Userlist[userid].Likedcontent += topic
			Web.Forum_data[topicid].Post_likes += 1

		} else if !strings.Contains(Web.Userlist[userid].Likedcontent, topic) && strings.Contains(Web.Userlist[userid].Dislikedcontent, topic) {
			//kui ei ole meeldiv vaid hoopis on ebameeldiv!
			rebuild := ""
			reg := strings.Split(Web.Userlist[userid].Dislikedcontent, topic)
			for regloop := 0; regloop > len(reg); regloop++ {
				rebuild += reg[regloop]
			}
			Web.Userlist[userid].Dislikedcontent = rebuild
			Web.Forum_data[topicid].Post_likes += 1
			Web.Forum_data[topicid].Post_disLikes -= 1
		} else if strings.Contains(Web.Userlist[userid].Likedcontent, topic) && !strings.Contains(Web.Userlist[userid].Dislikedcontent, topic) {
			// kui on meeldiv ega olegi ebameeldiv!
			rebuild := ""
			reg := strings.Split(Web.Userlist[userid].Likedcontent, topic)
			for regloop := 0; regloop > len(reg); regloop++ {
				rebuild += reg[regloop]
			}
			Web.Userlist[userid].Likedcontent = rebuild
			Web.Forum_data[topicid].Post_likes -= 1
		}
		Web.Sqlbase.Exec("UPDATE forum SET likes = Web.Forum_data[topicid].Post_likes WHERE post_header = title")
		Web.Sqlbase.Exec("UPDATE forum SET dislikes = Web.Forum_data[topicid].Post_dislikes WHERE post_header = title")
		Web.Sqlbase.Exec("UPDATE userdata SET likedcontent = Web.Userlist[userid].Likedcontent WHERE username = user")
		Web.Sqlbase.Exec("UPDATE userdata SET dislikedcontent = Web.Userlist[userid].Dislikedcontent WHERE username = user")
	}
}

func sendDisLike(database *sql.DB, user string, topic string, comment bool) {
	userid := getUserID(user)
	topicid := getTopicID(topic)

	if !comment {
		fmt.Println(topic, ">", user, "pani ebameeldivaks.")
		//kui ei ole märgitud ei meeldivaks ega ebameeldivaks
		if !strings.Contains(Web.Userlist[userid].Likedcontent, topic) && !strings.Contains(Web.Userlist[userid].Dislikedcontent, topic) {

			Web.Userlist[userid].Dislikedcontent += topic
			Web.Forum_data[topicid].Post_disLikes += 1

		} else if !strings.Contains(Web.Userlist[userid].Dislikedcontent, topic) && strings.Contains(Web.Userlist[userid].Likedcontent, topic) {
			//kui on meeldiv aga mitte ebameeldiv
			rebuild := ""
			reg := strings.Split(Web.Userlist[userid].Likedcontent, topic)
			for regloop := 0; regloop > len(reg); regloop++ {
				rebuild += reg[regloop]
			}
			Web.Userlist[userid].Likedcontent = rebuild
			Web.Userlist[userid].Dislikedcontent += topic
			Web.Forum_data[topicid].Post_disLikes += 1
			Web.Forum_data[topicid].Post_likes -= 1

		} else if strings.Contains(Web.Userlist[userid].Dislikedcontent, topic) && !strings.Contains(Web.Userlist[userid].Likedcontent, topic) {
			//kui on märgitud ebameeldivaks aga mitte meeldivaks
			rebuild := ""
			reg := strings.Split(Web.Userlist[userid].Dislikedcontent, topic)
			for regloop := 0; regloop > len(reg); regloop++ {
				rebuild += reg[regloop]
			}
			Web.Userlist[userid].Dislikedcontent = rebuild
			Web.Forum_data[topicid].Post_disLikes -= 1
		}

		Web.Sqlbase.Exec("UPDATE forum SET likes = Web.Forum_data[topicid].Post_likes WHERE post_header = title")
		Web.Sqlbase.Exec("UPDATE forum SET dislikes = Web.Forum_data[topicid].Post_dislikes WHERE post_header = title")
		Web.Sqlbase.Exec("UPDATE userdata SET likedcontent = Web.Userlist[userid].Likedcontent WHERE username = user")
		Web.Sqlbase.Exec("UPDATE userdata SET dislikedcontent = Web.Userlist[userid].Dislikedcontent WHERE username = user")
	}
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

	for i := 0; i < len(Web.Forum_data); i++ {
		if Web.Forum_data[i].Post_title == forum_header {
			Web.Forum_data[i].Commentor_data = append(Web.Forum_data[i].Commentor_data, commentpandemic{Commentor: commenter, Forum_comment: forum_Commentbox, Post_header: forum_header, Comment_likes: 0, Comment_disLikes: 0, Date: currentTime})
		}
	}
	statement.Exec(commenter, forum_Commentbox, forum_header, likes, dislikes, currentTime) // exec first name, last name

}

func sendRegister(database *sql.DB, username string, password string, email string) {
	statement, _ := database.Prepare("INSERT INTO userdata (username, password, email, likedcontent, dislikedcontent) VALUES (?,?,?,?,?)")
	statement.Exec(username, password, email, "", "") // exec first name, last name
	Web.Userlist = append(Web.Userlist, memberlist{ID: len(Web.Userlist) + 1, Username: username, Password: password, Email: email, Likedcontent: "", Dislikedcontent: ""})
	//fmt.Println("Server:", username, "has successfully registered!", " <", email, ">")

	// kasutajate printimine konsooli
	// for i := 0; i < len(userlist); i++ {
	// 	fmt.Println(userlist[i])
	// }
}

func saveAllPosts(database *sql.DB) {
	rows, _ := database.Query("SELECT originalposter, post_header, post_content, likes, dislikes, date FROM forum")
	var op string
	var header string
	var content string
	var likes int
	var disLikes int
	var date string
	for rows.Next() {
		rows.Scan(&op, &header, &content, &likes, &disLikes, &date)
		Web.Forum_data = append(Web.Forum_data, forumfamily{Originalposter: op, Post_title: header, Post_content: content, Post_likes: likes, Post_disLikes: disLikes, Date_posted: date})
	}
}

func saveAllUsers(database *sql.DB) {
	rows, _ := database.Query("SELECT id, username, password, email, likedcontent, dislikedcontent FROM userdata")
	var id int
	var username string
	var password string
	var email string
	var likedcontent string
	var dislikedcontent string
	for rows.Next() {
		rows.Scan(&id, &username, &password, &email, &likedcontent, &dislikedcontent)
		Web.Userlist = append(Web.Userlist, memberlist{ID: id, Username: username, Password: password, Email: email, Likedcontent: likedcontent, Dislikedcontent: dislikedcontent})
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
		for i := 0; i < len(Web.Forum_data); i++ {
			if Web.Forum_data[i].Post_title == header {
				Web.Forum_data[i].Commentor_data = append(Web.Forum_data[i].Commentor_data, commentpandemic{Commentor: commentor, Forum_comment: comments, Post_header: header, Comment_likes: likes, Comment_disLikes: disLikes, Date: date})
			}
		}
	}
}
