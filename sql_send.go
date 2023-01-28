package main

import (
	"database/sql"
	"strings"
	"time"
)

func sendLike(database *sql.DB, user string, topic string, comment bool) {
	userid := getUserID(user)
	topicid := getTopicID(topic)
	if !comment {
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

		//see miskipärast ei uuenda!!! aga võiks(loe peaks)!
		Web.Sqlbase.Exec("UPDATE forum SET likes = $Web.Forum_data[topicid].Post_likes WHERE post_header = $title")
		Web.Sqlbase.Exec("UPDATE forum SET dislikes = $Web.Forum_data[topicid].Post_dislikes WHERE post_header = $title")
		Web.Sqlbase.Exec("UPDATE userdata SET likedcontent = $Web.Userlist[userid].Likedcontent WHERE username = $user")
		Web.Sqlbase.Exec("UPDATE userdata SET dislikedcontent = $Web.Userlist[userid].Dislikedcontent WHERE username = $user")
		printLog(topic, ">", user, "pani meeldivaks.")

	}
}

func sendDisLike(database *sql.DB, user string, topic string, comment bool) {
	userid := getUserID(user)
	topicid := getTopicID(topic)

	if !comment {
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

		//see miskipärast ei uuenda!!! aga võiks(loe peaks)!

		Web.Sqlbase.Exec("UPDATE forum SET likes = $Web.Forum_data[topicid].Post_likes WHERE post_header = $title")
		Web.Sqlbase.Exec("UPDATE forum SET dislikes = $Web.Forum_data[topicid].Post_dislikes WHERE post_header = $title")
		Web.Sqlbase.Exec("UPDATE userdata SET likedcontent = $Web.Userlist[userid].Likedcontent WHERE username = $user")
		Web.Sqlbase.Exec("UPDATE userdata SET dislikedcontent = $Web.Userlist[userid].Dislikedcontent WHERE username = $user")
		printLog(topic, ">", user, "pani ebameeldivaks.")
	}
}

func sendRegister(database *sql.DB, username string, password string, email string) {
	statement, _ := database.Prepare("INSERT INTO userdata (username, password, email, likedcontent, dislikedcontent, datecreated) VALUES (?,?,?,?,?,?)")
	currentTime := time.Now().Format("02.01 2006")
	Web.Userlist = append(Web.Userlist, memberlist{ID: len(Web.Userlist) + 1, Username: username, Password: password, DateCreated: currentTime, Email: email, Likedcontent: "", Dislikedcontent: ""})
	statement.Exec(username, password, email, "", "", currentTime) // exec first name, last name
	printLog("Server:", username, "has successfully registered!", " <", email, ">")
}

func sendPost(database *sql.DB, originalposter string, header string, content string) bool {
	if len(content) < 10 {
		Web.ErrorMsg = "Too few characters in content, spam detected!"
		return false
	}
	for i := 0; i < len(Web.Forum_data); i++ {
		if Web.Forum_data[i].Post_title == header {
			Web.ErrorMsg = "Post with the title '" + header + "' already exists!"
			return false
		}
	}

	currentTime := time.Now().Format("02.01.2006 15:04")
	var likes, dislikes int

	Web.Forum_data = append(Web.Forum_data, forumfamily{Originalposter: originalposter, Post_title: header, Post_content: content, Date_posted: currentTime, Post_likes: 0, Post_disLikes: 0})

	statement, _ := database.Prepare("INSERT INTO forum (originalposter, post_header, post_content, likes, dislikes, date) VALUES (?,?,?,?,?,?)")
	statement.Exec(originalposter, header, content, likes, dislikes, currentTime) // exec first name, last name

	return true
}

func sendComment(database *sql.DB, commenter string, forum_Commentbox string, forum_header string) {

	if len(forum_Commentbox) <= 3 {
		return
	}

	printLog(forum_header, ">", commenter, "lisas kommentaari:", forum_Commentbox)

	statement, _ := database.Prepare("INSERT INTO commentdb (commentor, forum_comments, post_header, likes, dislikes, date, likedbyusers, dislikedbyusers) VALUES (?,?,?,?,?,?,?,?)")
	var likes, dislikes int
	currentTime := time.Now().Format("02.01.2006 15:04")

	i := Web.tempint
	if Web.Forum_data[i].Post_title == forum_header {
		Web.Forum_data[i].Commentor_data = append(Web.Forum_data[i].Commentor_data, commentpandemic{Commentor: commenter, Forum_comment: forum_Commentbox, Post_header: forum_header, Comment_likes: 0, Comment_disLikes: 0, Date: currentTime, Likedby: "", Dislikedby: ""})
	}

	statement.Exec(commenter, forum_Commentbox, forum_header, likes, dislikes, currentTime, "", "") // exec first name, last name
}
