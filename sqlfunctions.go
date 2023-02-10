package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func DateCreated() {
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

// Getting posts by category aswell, leaving empty WILL GET THEM ALL
func AllPosts(category string) []Forumdata {
	var data []Forumdata
	var statement *sql.Stmt
	var rows *sql.Rows

	switch category {
	case "":
		statement, _ = DataBase.Prepare(`
		SELECT posts.id, users.username, posts.title, posts.content, posts.date, category.name as category, image
		FROM posts
		LEFT JOIN users on posts.userId = users.id
		LEFT JOIN category on posts.categoryId = category.id		
		ORDER BY date DESC
		`)
		rows, _ = statement.Query()

	default:
		statement, _ = DataBase.Prepare(`
		SELECT posts.id, users.username, posts.title, posts.content, posts.date, category.name as category, image
		FROM posts
		LEFT JOIN users on posts.userId = users.id
		LEFT JOIN category on posts.categoryId = category.id
		WHERE posts.categoryId = ? 
		ORDER BY date DESC
		`)
		rows, _ = statement.Query(category)

	}

	for rows.Next() {
		var post Forumdata
		rows.Scan(
			&post.Id,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.Date_posted,
			&post.Category,
			&post.Image,
		)
		data = append(data, post)
	}

	return data
}

func getCategories() []Category {
	var categories []Category
	rows, err := DataBase.Query("select * from category")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var category Category
		rows.Scan(
			&category.Id,
			&category.Name,
		)
		categories = append(categories, category)

	}

	return categories
}

func SavePost(title string, author int, content string, categoryId int, image string) bool {
	if image == "" {
		image = "false"
	}

	statement, _ := DataBase.Prepare("INSERT INTO posts (userId, title, content, date, categoryId, image) VALUES (?,?,?,?,?,?)")
	currentTime := time.Now().Format("02.01.2006 15:04")

	statement.Exec(author, title, content, currentTime, categoryId, image)

	return true
}

func GetPostById(postId int) (Forumdata, error) {
	var post Forumdata
	statement, _ := DataBase.Prepare(`SELECT 
	posts.id, posts.userId, posts.title, posts.content, posts.date,
	users.username, image,
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
		&post.Date_posted,
		&post.Author,
		&post.Image,
		&post.Likes,
		&post.Dislikes,
	)

	return post, err
}

func SaveComment(content string, userId, postId int) bool {
	statement, _ := DataBase.Prepare("INSERT INTO comments (userId, content, postId, datecommented) VALUES (?,?,?,?)")
	currentTime := time.Now().Format("02.01.2006 15:04")
	statement.Exec(userId, content, postId, currentTime)
	return true
}

func GetUsers() []Memberlist {
	var users []Memberlist
	rows, _ := DataBase.Query("SELECT id, username, email, datecreated from users")

	for rows.Next() {
		var user Memberlist
		rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.DateCreated,
		)
		users = append(users, user)
	}

	return users
}

func GetCommentsByPostId(id int) []Commentdata {
	var comments []Commentdata
	statement, _ := DataBase.Prepare(`
	SELECT 
  comments.id, comments.userId, comments.content, 
  users.username, comments.datecommented,
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
			&comment.Date_commented,
			&comment.Likes,
			&comment.Dislikes,
		)

		comments = append(comments, comment)
	}
	var rearrange []Commentdata

	for i := 0; i < len(comments); i++ {
		rearrange = append(rearrange, comments[len(comments)-i-1])
	}

	return rearrange
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
	statement, _ := DataBase.Prepare("INSERT OR REPLACE INTO session (key, userId) VALUES (?,?)")
	_, err := statement.Exec(key, userId)
	if err != nil {
		fmt.Println("one per user")
	}
}

func DeleteSession(key string, userId int) {
	statement, _ := DataBase.Prepare("DELETE FROM session WHERE key = ? AND userId = ?")
	_, err := statement.Exec(key, userId)
	if err != nil {
		fmt.Println("Error deleting record from session table:", err)
	}
}

func CanRegister(uid string, password string, email string, cpassword string, cemail string) bool {
	str := ""
	if cpassword != password {
		Web.ErrorMsg = "The passwords do not match!"
		return false
	} else if cemail != email {
		Web.ErrorMsg = "The emails do not match!"
		return false
	}

	for i := 0; i < len(Web.User_data); i++ {
		if Web.User_data[i].Username == uid {
			str += "u"
		} else if Web.User_data[i].Email == email {
			str += "e"
		}
	}
	if strings.Contains(str, "u") {
		if strings.Contains(str, "e") {
			Web.ErrorMsg = "This username and e-mail is already in use!"
			return false
		}
		Web.ErrorMsg = "This username is already taken!"
		return false
	} else if strings.Contains(str, "e") && !strings.Contains(str, "u") {
		Web.ErrorMsg = "This e-mail is already in use!"
		return false
	}

	return true
}

func InitDatabase() {
	DataBase.Exec(
		`
		BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "posts" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"userId"	INTEGER,
	"title"	TEXT,
	"content"	TEXT,
	"categoryId"	INTEGER,
	"date"	TEXT,
	"image"	TEXT
);
CREATE TABLE IF NOT EXISTS "category" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"name"	TEXT
);
CREATE TABLE IF NOT EXISTS "comments" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"userId"	INTEGER,
	"content"	TEXT,
	"postId"	INTEGER,
	"datecommented"	TEXT
);
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"username"	TEXT,
	"password"	TEXT,
	"email"	TEXT,
	"datecreated"	TEXT
);
CREATE TABLE IF NOT EXISTS "session" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"key"	TEXT UNIQUE,
	"userId"	INTEGER UNIQUE
);
CREATE TABLE IF NOT EXISTS "commentLikes" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"name"	TEXT,
	"userId"	INTEGER,
	"commentId"	INTEGER,
	UNIQUE("commentId","userId")
);
CREATE TABLE IF NOT EXISTS "postlikes" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"name"	TEXT,
	"userId"	INTEGER,
	"postId"	INTEGER,
	UNIQUE("postId","userId")
);
INSERT INTO "posts" VALUES (123,15,'jpgjpg','jpgjpgjpgjpgjpgjpg',1,'08.02.2023 19:59','upload-3437260623.jpg');
INSERT INTO "posts" VALUES (124,15,'pngpng','pngpngpng',1,'08.02.2023 19:59','upload-2606161593.png');
INSERT INTO "posts" VALUES (125,15,'gifgif','gifgifgif',1,'08.02.2023 19:59','upload-2191428446.gif');
INSERT INTO "category" VALUES (1,'Kosmos');
INSERT INTO "category" VALUES (2,'Märgatud Jõhvis');
INSERT INTO "comments" VALUES (174,15,'Niino',125,'09.02 2023 06:27');
INSERT INTO "users" VALUES (9,'sass','$2a$14$68nNeNBTdHQafzdQ0TXyKe4VSU7osrvRPlzF7RHGUz2nIrUX4mN8y','asd@gmail.com','03.02 2023');
INSERT INTO "users" VALUES (15,'joel','$2a$14$sqf5Stu0zBTfE9J4wBL47OeijFNu5rnfu/qcN3zOGEZGAwJ251udi','joelimeil@gmail.com','07.02 2023');
INSERT INTO "session" VALUES (221,'8f659668-9a51-45bd-b7f1-d15fdf240701',15);
INSERT INTO "commentLikes" VALUES (243,'like',15,174);
INSERT INTO "postlikes" VALUES (289,'like',15,128);
INSERT INTO "postlikes" VALUES (291,'like',15,125);
COMMIT;
`)
}
