package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

func AllPostsRearrange(allposts []Forumdata) []Forumdata {
	var data []Forumdata
	for i := 0; i < len(allposts); i++ {
		data = append(data, allposts[len(allposts)-1-i])
	}

	return data
}

func AllPosts(category string) []Forumdata {
	var data []Forumdata
	// converting category name -> id
	// realCategoryID := 0
	realCategoryName := ""
	for i := 0; i < len(Web.Categories); i++ {
		if category == Web.Categories[i].Name {
			// realCategoryID = Web.Categories[i].Id
			realCategoryName = Web.Categories[i].Name
		}
	}
	// fmt.Println(realCategoryID, realCategoryName)

	rows, err := DataBase.Query(`
	SELECT posts.id, users.username, posts.title, posts.content, posts.date, category.name as category
	FROM posts
	LEFT JOIN users on posts.userId = users.id
	LEFT JOIN category on posts.categoryId = category.id
	`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var Id int
		var Category_name, Author, Title, Content, Date_posted string
		rows.Scan(
			&Id,
			&Author,
			&Title,
			&Content,
			&Date_posted,
			&Category_name,
		)
		if realCategoryName == Category_name && realCategoryName != "" {
			data = append(data, Forumdata{Id: Id, Author: Author, Title: Title, Content: Content, Date_posted: Date_posted, Category: category})
		} else if realCategoryName == "" {
			data = append(data, Forumdata{Id: Id, Author: Author, Title: Title, Content: Content, Date_posted: Date_posted, Category: Category_name})
		}
	}

	return data
}

func setupCategories() {
	rows, err := DataBase.Query("select * from category")
	if err != nil {
		log.Fatal(err)
	}

	if len(Web.Categories) != 0 {
		return
	}

	for rows.Next() {
		var id int
		var name string
		rows.Scan(
			&id,
			&name,
		)
		Web.Categories = append(Web.Categories, Category{Id: id, Name: name})
	}
}

func SavePost(title string, author int, content string, categoryId int) bool {
	if len(content) < 10 {
		return false
	}

	statement, _ := DataBase.Prepare("INSERT INTO posts (userId, title, content, date, categoryId) VALUES (?,?,?,?,?)")
	currentTime := time.Now().Format("02.01.2006 15:04")

	statement.Exec(author, title, content, currentTime, categoryId)

	return true
}

func GetPostById(postId int) (Forumdata, error) {
	var post Forumdata
	statement, _ := DataBase.Prepare(`SELECT 
	posts.id, posts.userId, posts.title, posts.content, posts.date,
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
		&post.Date_posted,
		&post.Author,
		&post.Likes,
		&post.Dislikes,
	)

	return post, err
}

func SaveComment(content string, userId, postId int) bool {
	statement, _ := DataBase.Prepare("INSERT INTO comments (userId, content, postId, datecommented) VALUES (?,?,?,?)")
	currentTime := time.Now().Format("02.01 2006 15:04")
	statement.Exec(userId, content, postId, currentTime)
	return true
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
