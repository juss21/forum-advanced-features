package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func sendTopicLike(database *sql.DB, topic string, like bool) {
	userid := getUserID(Web.Currentuser)
	topicID := getTopicID(topic)
	statement_tl, rror := database.Prepare("UPDATE forumlikes SET status = ?")
	statement_tdl, rror := database.Prepare("UPDATE forumdislikes SET status = ?")
	if rror != nil {
		fmt.Println(rror)
		os.Exit(-1)
	}
	for i := 0; i < len(Web.TopicLikes); i++ {
		if Web.TopicLikes[i].TopicID == topicID && Web.TopicLikes[i].UserID == userid {
			if like {
				if Web.TopicLikes[i].Status == 0 {
					if Web.TopicDisLikes[i].Status == 1 {
						Web.TopicDisLikes[i].Status = 0
						Web.Forum_data[topicID].Post_disLikes -= 1
					}
					Web.TopicLikes[i].Status = 1
					Web.Forum_data[topicID].Post_likes += 1

					statement_tl.Exec(Web.TopicLikes[i].Status)     // exec
					statement_tdl.Exec(Web.TopicDisLikes[i].Status) // exec
					continue
				}
				if Web.TopicLikes[i].Status == 1 {

					Web.TopicLikes[i].Status = 0
					Web.Forum_data[topicID].Post_likes -= 1

					statement_tl.Exec(Web.TopicLikes[i].Status)     // exec
					statement_tdl.Exec(Web.TopicDisLikes[i].Status) // exec
					continue
				}
			} else {
				if Web.TopicDisLikes[i].Status == 0 {
					if Web.TopicLikes[i].Status == 1 {
						Web.TopicLikes[i].Status = 0
						Web.Forum_data[topicID].Post_likes -= 1
					}

					Web.TopicDisLikes[i].Status = 1
					Web.Forum_data[topicID].Post_disLikes += 1

					statement_tl.Exec(Web.CommentLikes[i].Status)     // exec
					statement_tdl.Exec(Web.CommentDisLikes[i].Status) // exec
					continue
				}
				if Web.TopicDisLikes[i].Status == 1 {

					Web.TopicDisLikes[i].Status = 0
					Web.Forum_data[topicID].Post_disLikes -= 1

					statement_tl.Exec(Web.CommentLikes[i].Status)     // exec
					statement_tdl.Exec(Web.CommentDisLikes[i].Status) // exec
					continue
				}
			}
		}
	}
}

func sendCommentLike(database *sql.DB, commentid int, like bool) {
	userid := getUserID(Web.Currentuser)
	statement_cl, rror := database.Prepare("UPDATE commentlikes SET status = ?")
	statement_cdl, rror := database.Prepare("UPDATE commentdislikes SET status = ?")

	if rror != nil {
		fmt.Println(rror)
		os.Exit(-1)
	}
	//var comment_id, status int

	topicID := getTopicID(Web.Currentpage)
	temp := 0
	for test := 0; test < len(Web.Forum_data[topicID].Commentor_data); test++ {
		if Web.Forum_data[topicID].Commentor_data[test].ID == commentid {
			temp = test
		}
	}

	for i := 0; i < len(Web.CommentLikes); i++ {
		if Web.CommentLikes[i].CommentID == commentid && Web.CommentLikes[i].UserID == userid {

			if like {
				// if commentid and userid match
				if Web.CommentLikes[i].Status == 0 {
					if Web.CommentDisLikes[i].Status == 1 {
						Web.CommentDisLikes[i].Status = 0
						Web.Forum_data[topicID].Commentor_data[temp].Comment_disLikes -= 1
					}
					Web.CommentLikes[i].Status = 1
					Web.Forum_data[topicID].Commentor_data[temp].Comment_likes += 1
					statement_cl.Exec(Web.CommentLikes[i].Status)     // exec
					statement_cdl.Exec(Web.CommentDisLikes[i].Status) // exec
					continue
				}
				if Web.CommentLikes[i].Status == 1 {
					Web.CommentLikes[i].Status = 0
					Web.Forum_data[topicID].Commentor_data[temp].Comment_likes -= 1
					statement_cl.Exec(Web.CommentLikes[i].Status)     // exec
					statement_cdl.Exec(Web.CommentDisLikes[i].Status) // exec
					continue
				}
			} else {
				// if commentid and userid match
				if Web.CommentDisLikes[i].Status == 0 {
					if Web.CommentLikes[i].Status == 1 {
						Web.CommentLikes[i].Status = 0
						Web.Forum_data[topicID].Commentor_data[temp].Comment_likes -= 1
					}
					Web.CommentDisLikes[i].Status = 1
					Web.Forum_data[topicID].Commentor_data[temp].Comment_disLikes += 1
					statement_cdl.Exec(Web.CommentDisLikes[i].Status) // exec
					statement_cl.Exec(Web.CommentLikes[i].Status)     // exec
					continue
				}
				if Web.CommentDisLikes[i].Status == 1 {
					Web.CommentDisLikes[i].Status = 0
					Web.Forum_data[topicID].Commentor_data[temp].Comment_disLikes -= 1
					statement_cdl.Exec(Web.CommentDisLikes[i].Status) // exec
					statement_cl.Exec(Web.CommentLikes[i].Status)     // exec
					continue
				}
			}
		}
	}

}

func sendRegister(database *sql.DB, username string, password string, email string) {
	statement, _ := database.Prepare("INSERT INTO userdata (username, password, email, likedcontent, dislikedcontent, datecreated) VALUES (?,?,?,?,?,?)")
	currentTime := time.Now().Format("02.01 2006")
	Web.Userlist = append(Web.Userlist, memberlist{ID: len(Web.Userlist) + 1, Username: username, Password: password, DateCreated: currentTime, Email: email})

	statement.Exec(username, password, email, "", "", currentTime) // exec first name, last name
	printLog("Server:", username, "has successfully registered!", " <", email, ">")
}

func sendPost(database *sql.DB, originalposter string, header string, content string) bool {
	if len(content) < 10 {
		Web.ErrorMsg = "Too few characters in content, spam detected!"
		return false
	}
	for c := 0; c < len(header); c++ {
		if (rune(header[c]) == 35 || rune(header[c]) == 37 || rune(header[c]) == 64) || rune(header[c]) < 32 && rune(header[c]) > 122 {
			Web.ErrorMsg = "Unfriendly characters found in header!"
			return false
		}
	}

	for i := 0; i < len(Web.Forum_data); i++ {
		if Web.Forum_data[i].Post_title == header {
			Web.ErrorMsg = "Post with the title '" + header + "' already exists!"
			return false
		}
	}

	currentTime := time.Now().Format("02.01.2006 15:04")
	var likes, dislikes int
	Web.allposts += 1
	nextid := Web.allposts

	Web.Forum_data = append(Web.Forum_data, forumfamily{ID: nextid, Originalposter: originalposter, Post_title: header, Post_content: content, Date_posted: currentTime, Post_likes: 0, Post_disLikes: 0})

	statement, _ := database.Prepare("INSERT INTO forum (originalposter, post_header, post_content, likes, dislikes, date) VALUES (?,?,?,?,?,?)")
	statement.Exec(originalposter, header, content, likes, dislikes, currentTime) // exec first name, last name

	return true
}

func sendComment(database *sql.DB, commenter string, forum_Commentbox string, forum_header string) {

	printLog(forum_header, ">", commenter, "lisas kommentaari:", forum_Commentbox)

	statement, _ := database.Prepare("INSERT INTO commentdb (commentor, comments, post_header, likes, dislikes, date) VALUES (?,?,?,?,?,?)")
	var likes, dislikes int
	currentTime := time.Now().Format("02.01.2006 15:04")

	i := Web.tempint

	Web.allcomments += 1
	nextid := Web.allcomments

	if Web.Forum_data[i].Post_title == forum_header {
		Web.Forum_data[i].Commentor_data = append(Web.Forum_data[i].Commentor_data, commentpandemic{ID: nextid, Commentor: commenter, Forum_comment: forum_Commentbox, Post_header: forum_header, Comment_likes: 0, Comment_disLikes: 0, Date: currentTime})
	}

	statement.Exec(commenter, forum_Commentbox, forum_header, likes, dislikes, currentTime) // exec first name, last name
	printLog("appended data: ", Web.Forum_data[i].Commentor_data)

}
