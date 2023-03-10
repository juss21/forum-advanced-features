package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getPostID(commentID int) (int, string) {
	rows, err := DataBase.Query(`SELECT id, postId, content
	FROM comments WHERE id = ?`, commentID)
	if err != nil {
		fmt.Println("getPostID()")
		os.Exit(0)
	}

	var cid, pid int
	var content string
	for rows.Next() {
		rows.Scan(
			&cid,
			&pid,
			&content,
		)
		if cid == commentID {
			return pid, content
		}
	}
	return 0, "..."
}
func GetCreatedComments() {
	userid := Web.LoggedUser.ID

	rows, err := DataBase.Query(`SELECT id, content, postId FROM comments WHERE userId = ?`, userid)
	if err != nil {
		fmt.Println("GetCreatedComments()", err)
		os.Exit(0)
	}

	var commentID, postID int
	var content string
	for rows.Next() {
		rows.Scan(
			&commentID,
			&content,
			&postID,
		)

		Web.CreatedComments = append(Web.CreatedComments, CreatedComments{CommentID: commentID, Content: content, PostID: postID})
	}
}

func GetLikedComments(likeswitch string) {

	rows, err := DataBase.Query(`SELECT comments.id, content, postId
	FROM comments
	LEFT join commentLikes on comments.id = commentLikes.commentId
	where commentLikes.name = ?`, likeswitch)
	if err != nil {
		fmt.Println("GetLikedComments()", err)
		os.Exit(0)
	}

	var commentID, postID int
	var content string
	for rows.Next() {
		rows.Scan(
			&commentID,
			&content,
			&postID,
		)
		if likeswitch == "like" {
			Web.LikedComments = append(Web.LikedComments, LikedComments{CommentID: commentID, PostId: postID, Content: content})
		} else if likeswitch == "dislike" {
			Web.DisLikedComments = append(Web.DisLikedComments, LikedComments{CommentID: commentID, PostId: postID, Content: content})
		}
	}
}

func GetLikedPosts(likeswitch string) {
	userid := Web.LoggedUser.ID

	rows, err := DataBase.Query(`SELECT posts.id, posts.title, users.username
	FROM postlikes
	LEFT join users on postlikes.userId = users.id
    LEFT JOIN posts on postlikes.postId = posts.id
	where postlikes.name = ? and users.id = ?`, likeswitch, userid)
	if err != nil {
		fmt.Println("GetPostLikes()", err)
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
		if likeswitch == "like" {
			Web.LikedPosts = append(Web.LikedPosts, LikedPosts{PostID: id, User: name, Title: title})
		} else if likeswitch == "dislike" {
			Web.DisLikedPosts = append(Web.DisLikedPosts, LikedPosts{PostID: id, User: name, Title: title})
		}
	}
}

func DateCreated() {
	for i := 0; i < len(Web.User_data); i++ {
		if Web.User_data[i].ID == Web.LoggedUser.ID {
			Web.LoggedUser.DateCreated = Web.User_data[i].DateCreated
		}
	}
}

func UserPosted() { // TODO get user data
	userid := Web.LoggedUser.Username

	for i := 0; i < len(Web.Forum_data); i++ {
		if Web.Forum_data[i].Author == userid {
			Web.CreatedPosts = append(Web.CreatedPosts, CreatedPosts{PostID: Web.Forum_data[i].Id, UserID: Web.LoggedUser.ID, PostTopic: Web.Forum_data[i].Title})
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
	users.username, image, category.name, 
	COUNT(CASE WHEN postlikes.name = 'like' THEN 1 END) AS likes, 
	COUNT(CASE WHEN postlikes.name = 'dislike' THEN 1 END) AS dislikes
  FROM 
	posts 
	LEFT JOIN postlikes ON posts.id = postlikes.postId
	LEFT JOIN users ON posts.userId = users.id
	LEFT JOIN category on posts.categoryId = category.id
  WHERE posts.id = ?
  GROUP by posts.id
  `)
	err := statement.QueryRow(postId).Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.Date_posted,
		&post.Author,
		&post.Image,
		&post.Category,
		&post.Likes,
		&post.Dislikes,
	)

	return post, err
}
func DeleteNoticfication(UserID int, PostID int, User string, TargetID int, Content string) {
	_, err := DataBase.Exec(`DELETE FROM notifications WHERE UserID= ? AND PostID= ? AND User= ?
	 AND TargetID= ? AND Content=?`, UserID, PostID, User, TargetID, Content)
	if err != nil {
		fmt.Println(err)
	}
}

func DeletePostById(id string) {
	DataBase.Exec("DELETE FROM posts WHERE id= ?", id)
	_, err := DataBase.Exec("DELETE FROM comments WHERE postid= ?", id)
	if err != nil {
		fmt.Println(err)
	}
	_, err2 := DataBase.Exec("DELETE FROM notifications WHERE PostID= ?", id)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func DeleteCommentById(id string, postid int) {
	_, err := DataBase.Exec("DELETE FROM comments WHERE id= ?", id)
	if err != nil {
		fmt.Println(err)
	}
	commentID, erro := strconv.Atoi(id)
	if erro != nil {
		fmt.Println(erro)
	}

	_, err2 := DataBase.Exec("DELETE FROM notifications WHERE targetID = ?", commentID)
	if err2 != nil {
		fmt.Println(err)
	}
}

func SaveComment(content string, userId int, postId int, title string, user string, authorUserId int) bool {
	statement, _ := DataBase.Prepare("INSERT INTO comments (userId, content, postId, datecommented) VALUES (?,?,?,?)")
	currentTime := time.Now().Format("02.01.2006 15:04")
	statement.Exec(userId, content, postId, currentTime)

	rows, err1 := DataBase.Query("SELECT id from comments WHERE userId = ? AND content = ? AND postId = ?", userId, content, postId)
	if err1 != nil {
		fmt.Println(err1)
	}
	//authorUDIz
	var commentID int

	for rows.Next() {
		rows.Scan(&commentID)
	}
	//fmt.Println(commentID)
	activity := "commented"
	notification := "on your post: " + title
	//statement.Exec("INSERT INTO notifications VALUES(?,?,?,?)", authorUserId, postId, user, title)
	_, err := DataBase.Exec("INSERT INTO notifications VALUES(?,?,?,?,?,?)", authorUserId, postId, user, commentID, activity, notification)
	if err != nil {
		fmt.Println(err)
	}
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

func GetNotifications() []Notifications {
	var notifications []Notifications

	rows, _ := DataBase.Query("SELECT * FROM notifications WHERE UserID= ? and User!= ?",
		Web.LoggedUser.ID, Web.LoggedUser.Username)

	for rows.Next() {
		var notification Notifications
		rows.Scan(
			&notification.UserID,
			&notification.PostID,
			&notification.User,
			&notification.TargetID,
			&notification.Activity,
			&notification.Content,
		)
		notifications = append(notifications, notification)
	}

	return notifications
}

func SavePostLike(like string, userId, postId int, title string) {
	var dbLike string
	statement, _ := DataBase.Prepare("SELECT name FROM postlikes WHERE  userId = ? and postId = ?")
	statement.QueryRow(userId, postId).Scan(&dbLike)
	if like == dbLike {
		toggleLike, _ := DataBase.Prepare("DELETE FROM postlikes WHERE  userId = ? and postId = ?")
		toggleLike.Exec(userId, postId)

		notifyLike, _ := DataBase.Prepare(`DELETE FROM notifications WHERE activity != 'commented' and user = (SELECT username FROM users WHERE id = ?) and postId = ?`)
		notifyLike.Exec(userId, postId)
	} else {
		toggleLike, _ := DataBase.Prepare("DELETE FROM postlikes WHERE  userId = ? and postId = ?")
		toggleLike.Exec(userId, postId)
		saving, _ := DataBase.Prepare("INSERT INTO postlikes (name, userId, postId) VALUES (?,?,?)")
		_, err := saving.Exec(like, userId, postId)

		notifyLike, _ := DataBase.Prepare(`DELETE FROM notifications WHERE activity != 'commented' and user = (SELECT username FROM users WHERE id = ?) and postId = ?`)
		notifyLike.Exec(userId, postId)
		notifyUpdate, _ := DataBase.Prepare(`INSERT INTO notifications VALUES((SELECT userid from posts where id = ? ),
		(SELECT id FROM posts WHERE id = ? ), (SELECT username  FROM users WHERE id = ?),?,?,?)`)

		activity := like + "s"
		content := "your post: " + title
		_, err2 := notifyUpdate.Exec(postId, postId, userId, 0, activity, content)

		if err == nil || err2 != nil {
			return
		}
	}

}

func opposite(str string) string {
	if str == "like" {
		return "dislike"
	} else if str == "dislike" {
		return "like"
	}
	return ""
}

func SaveCommentLike(like string, userId, commentId int) {
	var dbLike string
	statement, _ := DataBase.Prepare("SELECT name FROM commentLikes WHERE  userId = ? and commentId = ?")
	statement.QueryRow(userId, commentId).Scan(&dbLike)
	if like == dbLike {
		toggleLike, _ := DataBase.Prepare("DELETE FROM commentLikes WHERE  userId = ? and commentId = ?")
		toggleLike.Exec(userId, commentId)

		notifyLike, _ := DataBase.Prepare(`DELETE FROM notifications WHERE activity != 'commented' and user = (SELECT username FROM users WHERE id = ?) and targetID = ?`)
		notifyLike.Exec(userId, commentId)
	} else {
		toggleLike, _ := DataBase.Prepare("DELETE FROM commentLikes WHERE  userId = ? and commentId = ?")
		toggleLike.Exec(userId, commentId)
		saving, _ := DataBase.Prepare("INSERT INTO commentLikes (name, userId, commentId) VALUES (?,?,?)")
		_, err := saving.Exec(like, userId, commentId)

		notifyLike, _ := DataBase.Prepare(`DELETE FROM notifications WHERE activity != 'commented' and user = (SELECT username FROM users WHERE id = ?) and targetID = ?`)
		notifyLike.Exec(userId, commentId)
		notifyUpdate, _ := DataBase.Prepare(`INSERT INTO notifications VALUES((SELECT userid from comments where id = ? ),
		(SELECT postId FROM comments WHERE id = ? ), (SELECT username  FROM users WHERE id = ?),?,?,?)`)
		rows, _ := DataBase.Query("SELECT content from comments where id = ?", commentId)
		var comment string
		for rows.Next() {
			rows.Scan(
				&comment,
			)
		}
		activity := like + "s"
		content := "your comment: " + comment

		_, err2 := notifyUpdate.Exec(commentId, commentId, userId, commentId, activity, content)

		if err == nil || err2 != nil {
			return
		}
	}
}

func SaveSession(key string, userId int) {
	DataBase.Exec("DELETE FROM session WHERE userId = ?", userId)

	statement, _ := DataBase.Prepare("INSERT INTO session (key, userId) VALUES (?,?)")
	_, err := statement.Exec(key, userId)
	if err != nil {
		fmt.Println("one per user")
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

func EditPostById(postId int, title, image, content string) error {
	if image == "" {
		image = Web.CurrentPost.Image
	}
	currentTime := time.Now().Format("02.01.2006 15:04")

	_, err := DataBase.Exec("UPDATE posts SET title=?, content = ?,  date = ?, image=?   where id = ? ", title, content, currentTime, image, postId)

	return err
}
func EditCommentById(commentId int, content string) error {
	// if image == "" {
	// 	image = "false"
	// }
	currentTime := time.Now().Format("02.01.2006 15:04")

	_, err := DataBase.Exec("UPDATE comments SET content = ?,  datecommented = ?   where id = ? ", content, currentTime, commentId)

	return err
}

func InitDatabase() {
	DataBase.Exec(
		`
		BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"username"	TEXT UNIQUE,
	"password"	TEXT,
	"email"	TEXT UNIQUE,
	"datecreated"	TEXT
);
CREATE TABLE IF NOT EXISTS "postlikes" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"name"	TEXT,
	"userId"	INTEGER,
	"postId"	INTEGER,
	UNIQUE("postId","userId")
);
CREATE TABLE IF NOT EXISTS "commentLikes" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"name"	TEXT,
	"userId"	INTEGER,
	"commentId"	INTEGER,
	UNIQUE("commentId","userId")
);
CREATE TABLE IF NOT EXISTS "session" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"key"	TEXT UNIQUE,
	"userId"	INTEGER UNIQUE
);
CREATE TABLE IF NOT EXISTS "comments" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"userId"	INTEGER,
	"content"	TEXT,
	"postId"	INTEGER,
	"datecommented"	TEXT
);
CREATE TABLE IF NOT EXISTS "category" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"name"	TEXT
);
CREATE TABLE IF NOT EXISTS "posts" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"userId"	INTEGER,
	"title"	TEXT,
	"content"	TEXT,
	"categoryId"	INTEGER,
	"date"	TEXT,
	"image"	TEXT
);
INSERT INTO "users" VALUES (2,'isabella','$2a$14$5VY414NmXYll0cNJVk71l.vEpj2/DF/JZ/vCfr8PRuQZkeU9N5BBO','isabella@gmail.com','12.02.2023');
INSERT INTO "users" VALUES (4,'sinisterObtuse','$2a$14$3ceCJAGSpb813jupNxTZSOghkTVDZ7/j32zFT9WmrnePGDKWxGvEC','sinister@gmail.com','12.02.2023');
INSERT INTO "users" VALUES (5,'andrei','$2a$14$xSYbFdGIgX5Pe5svwjA7KOFKllaFuVgrEvsAbNV3hJnHaYDMDuq3u','andrei@koodJõhvi.com','12.02.2023');
INSERT INTO "postlikes" VALUES (2,'like',4,3);
INSERT INTO "postlikes" VALUES (3,'dislike',4,4);
INSERT INTO "postlikes" VALUES (4,'like',2,5);
INSERT INTO "postlikes" VALUES (5,'like',5,4);
INSERT INTO "commentLikes" VALUES (1,'like',2,1);
INSERT INTO "commentLikes" VALUES (3,'like',5,3);
INSERT INTO "comments" VALUES (1,4,'Woah so interesting',3,'12.02.2023 13:32');
INSERT INTO "comments" VALUES (2,4,'Well it''s pure spamming. Reported!',4,'12.02.2023 13:32');
INSERT INTO "comments" VALUES (3,2,'I know her, she lives on Ädala street. I can send her number',5,'12.02.2023 13:34');
INSERT INTO "category" VALUES (1,'Kosmos');
INSERT INTO "category" VALUES (2,'Märgatud Jõhvis');
INSERT INTO "posts" VALUES (1,2,'Jesse Marcel and Roswell conspiracy theories','In February 1978, UFO researcher Stanton Friedman interviewed Jesse Marcel, the only person known to have accompanied the Roswell debris from where it was recovered to Fort Worth where reporters saw material that was claimed to be part of the recovered object. Marcel''s statements contradicted those he made to the press in 1947.[79]

In November 1979, Marcel''s first filmed interview was featured in a documentary titled "UFO''s Are Real", co-written by Friedman.[80] The film had a limited release but was later syndicated for broadcasting. On February 28, 1980, sensationalist tabloid the National Enquirer brought large-scale attention to the Marcel story.[81] On September 20, 1980, the TV series In Search of... aired an interview where Marcel described his participation in the 1947 press conference',1,'12.02.2023 13:22','false');
INSERT INTO "posts" VALUES (2,2,'Linda Moulton Howe and cattle mutilations','Linda Moulton Howe is an advocate of conspiracy theories that cattle mutilations are of extraterrestrial origin and speculations that the U.S. government is involved with aliens.',1,'12.02.2023 13:23','false');
INSERT INTO "posts" VALUES (3,2,'In popular fiction','Works of popular fiction have included premises and scenes in which a government intentionally prevents disclosure to its populace of the discovery of non-human, extraterrestrial intelligence. Motion picture examples include 2001: A Space Odyssey (as well as the earlier novel by Arthur C. Clarke),[136][137] Easy Rider,[138] the Steven Spielberg films Close Encounters of the Third Kind and E.T. the Extra-Terrestrial, Hangar 18, Total Recall, Men in Black, and Independence Day. Television series and films including The X-Files, Dark Skies, and Stargate have also featured efforts by governments to conceal information about extraterrestrial beings. The plot of the Sidney Sheldon novel The Doomsday Conspiracy involves a UFO conspiracy.[139]

In March 2001, former astronaut and United States Senator John Glenn appeared on an episode of the TV series Frasier playing a fictional version of himself who confesses to a UFO coverup.[140] ',1,'12.02.2023 13:24','false');
INSERT INTO "posts" VALUES (4,2,'👽👽👽👽👽👽👽👽👽','Aliens are frikkin 😎 ',1,'12.02.2023 13:25','false');
INSERT INTO "posts" VALUES (5,4,'Found ID-card in Jõhvi Konserdimaja','Found ID-card. Contact me by email, if it might be  yours. ',2,'12.02.2023 13:32','upload-2033429674.jpg');
COMMIT;
`)
}
