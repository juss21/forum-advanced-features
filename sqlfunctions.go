package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func test() {
	for i := 0; i < len(Web.User_data); i++ {
		if Web.User_data[i].ID == Web.LoggedUser.ID {
			Web.LoggedUser.DateCreated = Web.User_data[i].DateCreated
		}
	}
}

func LikesSent() {
	userid := Web.LoggedUser.ID

	rows, err := DataBase.Query(`SELECT posts.id, posts.title, users.username
	FROM postlikes
	LEFT join users on postlikes.userId = users.id
   LEFT JOIN posts on postlikes.postId = posts.id
	where postlikes.name = "like" and users.id = ?`, userid)
	if err != nil {
		fmt.Println("likessent()", err)
		os.Exit(0)
	}
	var id int
	var name, title string
	for rows.Next() {
		rows.Scan(
			&id,
			&title,
			&name,
		)
		Web.LikedComments = append(Web.LikedComments, Likedstuff{PostID: id, User: name, Title: title})
	}
}
func UserPosted() {
	userid := Web.LoggedUser.Username

	for i := 0; i < len(Web.Forum_data); i++ {
		if Web.Forum_data[i].Author == userid {
			Web.CreatedPosts = append(Web.CreatedPosts, Createdstuff{PostID: Web.Forum_data[i].Id, UserID: Web.LoggedUser.ID, PostTopic: Web.Forum_data[i].Title})
		}
	}
}

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

	statement, _ := DataBase.Prepare("INSERT INTO posts (userId, title, content) VALUES (?,?,?)")

	statement.Exec(author, title, content)

	return true
}

func GetPostById(postId int) Forumdata {
	var post Forumdata
	statement, _ := DataBase.Prepare(`SELECT 
	posts.id, posts.userId, posts.title, posts.content,
	users.username,
	COUNT(CASE WHEN postlikes.name = 'like' THEN 1 END) AS likes, 
	COUNT(CASE WHEN postlikes.name = 'dislike' THEN 1 END) AS dislikes
  FROM 
	posts 
	LEFT JOIN postlikes ON posts.id = postlikes.postId
	LEFT JOIN users ON posts.userId = users.id
  WHERE posts.id = ?
  GROUP by posts.id;
  `)
	err := statement.QueryRow(postId).Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.Author,
		&post.Likes,
		&post.Dislikes,
	)
	if err != nil {
		log.Fatal(err)
	}

	return post

}
func SaveComment(content string, userId, postId int) bool {

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
	statement, _ := DataBase.Prepare("INSERT INTO users (username, password, email, datecreated) values (?,?,?,?)")
	currentTime := time.Now().Format("02.01 2006")
	statement.Exec(username, password, email, currentTime)
}
func GetUsers() {
	rows, _ := DataBase.Query("SELECT id, username, email, datecreated from users")

	for rows.Next() {
		var ID int
		var Username, Email, DateC string
		rows.Scan(
			&ID,
			&Username,
			&Email,
			&DateC,
		)
		Web.User_data = append(Web.User_data, Memberlist{ID: ID, Username: Username, Email: Email, DateCreated: DateC})
	}
}

func GetCommentsByPostId(id int) []Commentdata {
	var comments []Commentdata
	statement, _ := DataBase.Prepare(`
	SELECT 
  comments.id, comments.userId, comments.content, 
  users.username,
  COUNT(CASE WHEN commentLikes.name = 'like' THEN 1 END) AS likes, 
  COUNT(CASE WHEN commentLikes.name = 'dislike' THEN 1 END) AS dislikes
FROM 
  comments 
  LEFT JOIN commentLikes ON comments.id = commentLikes.commentId
  LEFT JOIN users ON comments.userId = users.id
WHERE comments.postId= ?
GROUP by comments.id;
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
			&comment.Username,
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
