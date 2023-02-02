package main

import (
	"database/sql"
	"fmt"
	"log"
)

func AllPosts() []Forumdata {
	var data []Forumdata
	rows, err := DataBase.Query("SELECT posts.id, users.username, posts.title, posts.content, date FROM posts LEFT JOIN users on posts.userId = users.id")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var post Forumdata
		rows.Scan(
			&post.Id,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.Date_posted,
		)
		data = append(data, post)
	}

	return data
}

func SavePost(title string, author int, content string) bool {
	if len(title) < 5 {
		Web.ErrorMsg = "Too few arguments in title, spam?"
		return false
	}
	if len(content) < 10 {
		Web.ErrorMsg = "Too few arguments in content! Spam!"
		return false
	}
	statement, _ := DataBase.Prepare("INSERT INTO posts (userId, title, content) VALUES (?,?,?)")

	statement.Exec(author, title, content)

	return true
}

func GetPostById(postId int) Forumdata {
	var post Forumdata
	statement, _ := DataBase.Prepare(`SELECT 
	posts.id, posts.userId, posts.title, posts.content,  
	COUNT(CASE WHEN postlikes.name = 'like' THEN 1 END) AS likes, 
	COUNT(CASE WHEN postlikes.name = 'dislike' THEN 1 END) AS dislikes  
  FROM 
	posts 
	LEFT JOIN postlikes ON posts.id = postlikes.postId	
	WHERE posts.id = ?
	GROUP by posts.id
  `)
	err := statement.QueryRow(postId).Scan(
		&post.Id,
		&post.Author,
		&post.Title,
		&post.Content,
		&post.Likes,
		&post.Dislikes,
	)
	if err != nil {
		log.Fatal(err)
	}

	return post

}
func SaveComment(content string, userId, postId int) bool {
	if len(content) < 10 {
		Web.ErrorMsg = "Too few arguments in commentbox! SPAM"
		return false
	}
	statement, _ := DataBase.Prepare("INSERT INTO comments (userId, content, postId) VALUES (?,?,?)")
	statement.Exec(userId, content, postId)
	return true
}

func Login(username, password string) (Memberlist, error) {
	var user Memberlist
	statement, _ := DataBase.Prepare("SELECT id, username, email FROM users WHERE username=? and password=?")
	err := statement.QueryRow(username, password).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)

	if err == sql.ErrNoRows {
		return Memberlist{}, err
	}

	return user, err
}

func Register(username, password, email string) {
	statement, _ := DataBase.Prepare("INSERT INTO users (username, password, email) values (?,?,?)")
	statement.Exec(username, password, email)
}
func GetUsers() []Memberlist {
	var data []Memberlist
	rows, _ := DataBase.Query("SELECT id, username, email from users")

	for rows.Next() {
		var member Memberlist
		rows.Scan(
			&member.ID,
			&member.Username,
			&member.Email,
		)
		data = append(data, member)
	}

	return data
}

func GetCommentsByPostId(id int) []Commentdata {
	var comments []Commentdata
	statement, _ := DataBase.Prepare(`
	SELECT 
	comments.id, comments.userId, comments.content, 
  COUNT(CASE WHEN commentLikes.name = 'like' THEN 1 END) AS likes, 
  COUNT(CASE WHEN commentLikes.name = 'dislike' THEN 1 END) AS dislikes  
FROM 
  comments 
  LEFT JOIN commentLikes ON comments.id = commentLikes.commentId  
  WHERE comments.postId= ?
  GROUP by comments.id
	`)
	rows, err := statement.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var comment Commentdata
		rows.Scan(
			&comment.Id,
			&comment.UserId,
			&comment.Content,
			&comment.Likes,
			&comment.Dislikes,
		)
		comments = append(comments, comment)
	}

	return comments
}

func SavePostLike(like string, userId, postId int) {
	var dbLike string
	statement, _ := DataBase.Prepare("SELECT name FROM postlikes WHERE  userId = ? and postId = ?")
	statement.QueryRow(userId, postId).Scan(&dbLike)
	if like == dbLike {
		toggleLike, _ := DataBase.Prepare("DELETE FROM postlikes WHERE  userId = ? and postId = ?")
		toggleLike.Exec(userId, postId)
	} else {
		toggleLike, _ := DataBase.Prepare("DELETE FROM postlikes WHERE  userId = ? and postId = ?")
		toggleLike.Exec(userId, postId)
		saving, _ := DataBase.Prepare("INSERT INTO postlikes (name, userId, postId) VALUES (?,?,?)")
		_, err := saving.Exec(like, userId, postId)
		if err == nil {
			return
		}
	}
}

func SaveCommentLike(like string, userId, commentId int) {
	var dbLike string
	statement, _ := DataBase.Prepare("SELECT name FROM commentLikes WHERE  userId = ? and commentId = ?")
	statement.QueryRow(userId, commentId).Scan(&dbLike)
	if like == dbLike {
		toggleLike, _ := DataBase.Prepare("DELETE FROM commentLikes WHERE  userId = ? and commentId = ?")
		toggleLike.Exec(userId, commentId)
	} else {
		toggleLike, _ := DataBase.Prepare("DELETE FROM commentLikes WHERE  userId = ? and commentId = ?")
		toggleLike.Exec(userId, commentId)
		saving, _ := DataBase.Prepare("INSERT INTO commentLikes (name, userId, commentId) VALUES (?,?,?)")
		_, err := saving.Exec(like, userId, commentId)
		if err == nil {
			return
		}
	}
}

func SaveSession(key string, userId int) {
	statement, _ := DataBase.Prepare("INSERT INTO session (key, userId) VALUES (?,?)")
	_, err := statement.Exec(key, userId)

	if err != nil {
		fmt.Println("one per user")
	}
}
