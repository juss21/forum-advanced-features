package main

import (
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

func AllPostsRearrange(allposts []Forumdata) []Forumdata {
	var data []Forumdata
	for i := 0; i < len(allposts); i++ {
		data = append(data, allposts[len(allposts)-1-i])
	}

	return data
}

func AllPosts(category string) []Forumdata {
	var data []Forumdata

	realCategoryName := ""
	for i := 0; i < len(Web.Categories); i++ {
		if category == Web.Categories[i].Name {

			realCategoryName = Web.Categories[i].Name
		}
	}

	rows, err := DataBase.Query(`
	SELECT posts.id, users.username, posts.title, posts.content, posts.date, category.name as category, image
	FROM posts
	LEFT JOIN users on posts.userId = users.id
	LEFT JOIN category on posts.categoryId = category.id
	`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var Id int
		var Category_name, Author, Title, Content, Date_posted, Image string
		rows.Scan(
			&Id,
			&Author,
			&Title,
			&Content,
			&Date_posted,
			&Category_name,
			&Image,
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
	"date"	TEXT
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
INSERT INTO "posts" VALUES (87,15,'SUUR Postitus!!!','Tere see on väga suur postitus!',2,'03.02.2023 04:52');
INSERT INTO "posts" VALUES (88,13,'Osdan deleskoob','Tahan naha guud',1,'05.02.2023 19:28');
INSERT INTO "posts" VALUES (89,14,'Pikk postitus','Lühikokkuvõte
Käesoleva bakalaureuse töö eesmärgiks oli uurida prügila nõrgvee moodustamist, omadusi ja
selle erinevaid töötlusmeetodeid. Antud bakalaureusetöös antakse ülevaade bioloogilistele,
füüsikalistele kui ka keemilistele töötlusprotsessidele.
Ohutu ja usaldusväärne jäätmete ladestamine on kõige tähtsam ülesanne jäätmekäitluses.
Prügilate rajamisel keskendutakse kõige enam faktile, et prügila ei tekitaks negatiivset mõju
keskkonnale. Enamik tänapäevased prügilad on varustatud nõrgvee ja prügilagaasi
kogumissüsteemidega, mille kohta tehakse regulaarset seiret. Prügila valed hooldusvõtted
võivad tekitada nii õhu-, pinnase- kui ka veereostuse.
Nõrgvesi on prügist läbi imbunud vedelik, mis sisaldab saasteaineid. Nõrgvesi tekib prügilas
kõige enam sademete tõttu, samuti tekib jäätmete niiskusest. Nõrgvee omadused sõltuvad
jäätmete koostisest ja prügila vanusest. Mida noorem on prügila, seda saastatum on selle
prügila nõrgvesi, sest jäätmed ei ole jõudnud täielikult laguneda. Nõrgvee biolagundatavus on
ajas muutuv. Kui on tegemist noore prügilaga, siis on nõrgvee biolagundatavus väike, prügila
vananedes biolagundatavus suureneb.
Prügilas tekkiv nõrgvesi on toksiline ning kui seda valesti töödelda võib see reostada
ümbritsevat pinnast ja põhjavett. Prügila nõrgvee töötlemisel kasutatakse erinevaid
bioloogilisi, keemilisi ja füüsikalisi reovee puhastusmeetodeid. Bioloogilised töötlusmeetodid
nagu aktiivmuda protsess, nitrifikatsioon-denitrifikatsioon protsess ning anaeroobsed ja
aeroobsed laguunid põhinevad mikroorganismide kasutamisel, kes kasutavad orgaanilist ja
anorgaanilist materjali sünteesimaks endale eluks vajaliku energiat. Bioloogilised
töötlusmeetodid on efektiivsed vaid noorema prügila nõrgvee töötlemiseks. Põhiliselt
kasutatakse bioloogilisi töötlusmeetodeid enne füüsikalisi ja keemilisi protsesse. Füüsikalised
töötlusmeetodid nagu flotatsioon, filtratsioon ja settimine põhinevad füüsika seadustel.
Keemilistes töötlusprotsessides kasutatakse kemikaalide lisamist või muud keemilist
reaktsiooni saasteainete vähendamiseks. Keemilisteks töötlusmeetoditeks on näiteks
sadestamine, neutraliseerimine ja oksüdatsioon. Parima tulemuse annab aga töötlusmeetodite
erinev kombineerimine.
Selleks, et vältida prügila nõrgvee ohtliku mõju keskkonnale, tuleb jäätmeid ohutult ja
vastavalt nõuetele ladestada ning käidelda. Samuti on väga oluline nõrgvee kokku kogumine
ning nõuetele vastav puhastamine.',2,'07.02.2023 01:57');
INSERT INTO "posts" VALUES (90,15,'uus postitus','terekst teuele kogiile',1,'08.02.2023 10:04');
INSERT INTO "posts" VALUES (92,15,'uus postitud 2','paelgaplglapgplapgelapg',2,'08.02.2023 10:09');
INSERT INTO "posts" VALUES (93,16,'uue kasutaja uus postist','uue kasutaja uus postist',2,'08.02.2023 10:11');
INSERT INTO "category" VALUES (1,'Kosmos');
INSERT INTO "category" VALUES (2,'Märgatud Jõhvis');
INSERT INTO "comments" VALUES (155,13,'Minu kommentaar!',87,'05.02 2023 19:27');
INSERT INTO "comments" VALUES (156,13,'tere',80,'06.02 2023 13:40');
INSERT INTO "comments" VALUES (157,13,'Hello gud frend',79,'06.02 2023 13:41');
INSERT INTO "comments" VALUES (158,14,'eeaeasd',88,'07.02 2023 01:51');
INSERT INTO "comments" VALUES (159,15,'pühas römps',89,'07.02 2023 03:35');
INSERT INTO "comments" VALUES (160,15,'sitt kokkuvõte!',89,'07.02 2023 05:10');
INSERT INTO "comments" VALUES (161,14,':/',89,'07.02 2023 05:47');
INSERT INTO "comments" VALUES (162,15,'Lühikokkuvõte Käesoleva bakalaureuse töö eesmärgiks oli uurida prügila nõrgvee moodustamist, omadusi ja selle erinevaid töötlusmeetodeid. Antud bakalaureusetöös antakse ülevaade bioloogilistele, füüsikalistele kui ka keemilistele töötlusprotsessidele. Ohutu ja usaldusväärne jäätmete ladestamine on kõige tähtsam ülesanne jäätmekäitluses. Prügilate rajamisel keskendutakse kõige enam faktile, et prügila ei tekitaks negatiivset mõju keskkonnale. Enamik tänapäevased prügilad on varustatud nõrgvee ja prügilagaasi kogumissüsteemidega, mille kohta tehakse regulaarset seiret. Prügila valed hooldusvõtted võivad tekitada nii õhu-, pinnase- kui ka veereostuse. Nõrgvesi on prügist läbi imbunud vedelik, mis sisaldab saasteaineid. Nõrgvesi tekib prügilas kõige enam sademete tõttu, samuti tekib jäätmete niiskusest. Nõrgvee omadused sõltuvad jäätmete koostisest ja prügila vanusest. Mida noorem on prügila, seda saastatum on selle prügila nõrgvesi, sest jäätmed ei ole jõudnud täielikult laguneda. Nõrgvee biolagundatavus on ajas muutuv. Kui on tegemist noore prügilaga, siis on nõrgvee biolagundatavus väike, prügila vananedes biolagundatavus suureneb. Prügilas tekkiv nõrgvesi on toksiline ning kui seda valesti töödelda võib see reostada ümbritsevat pinnast ja põhjavett. Prügila nõrgvee töötlemisel kasutatakse erinevaid bioloogilisi, keemilisi ja füüsikalisi reovee puhastusmeetodeid. Bioloogilised töötlusmeetodid nagu aktiivmuda protsess, nitrifikatsioon-denitrifikatsioon protsess ning anaeroobsed ja aeroobsed laguunid põhinevad mikroorganismide kasutamisel, kes kasutavad orgaanilist ja anorgaanilist materjali sünteesimaks endale eluks vajaliku energiat. Bioloogilised töötlusmeetodid on efektiivsed vaid noorema prügila nõrgvee töötlemiseks. Põhiliselt kasutatakse bioloogilisi töötlusmeetodeid enne füüsikalisi ja keemilisi protsesse. Füüsikalised töötlusmeetodid nagu flotatsioon, filtratsioon ja settimine põhinevad füüsika seadustel. Keemilistes töötlusprotsessides kasutatakse kemikaalide lisamist või muud keemilist reaktsiooni saasteainete vähendamiseks. Keemilisteks töötlusmeetoditeks on näiteks sadestamine, neutraliseerimine ja oksüdatsioon. Parima tulemuse annab aga töötlusmeetodite erinev kombineerimine. Selleks, et vältida prügila nõrgvee ohtliku mõju keskkonnale, tuleb jäätmeid ohutult ja vastavalt nõuetele ladestada ning käidelda. Samuti on väga oluline nõrgvee kokku kogumine ning nõuetele vastav puhastamine.',89,'07.02 2023 06:17');
INSERT INTO "comments" VALUES (163,15,'asd',89,'07.02 2023 14:34');
INSERT INTO "comments" VALUES (164,15,'aaaaa',89,'07.02 2023 14:36');
INSERT INTO "comments" VALUES (165,15,'lahe',89,'07.02 2023 14:52');
INSERT INTO "comments" VALUES (166,0,'norm',89,'08.02 2023 06:12');
INSERT INTO "comments" VALUES (167,15,'egagagaegag',92,'08.02 2023 10:09');
INSERT INTO "comments" VALUES (168,15,'uus komment',93,'08.02 2023 10:11');
INSERT INTO "users" VALUES (9,'sass','$2a$14$68nNeNBTdHQafzdQ0TXyKe4VSU7osrvRPlzF7RHGUz2nIrUX4mN8y','asd@gmail.com','03.02 2023');
INSERT INTO "users" VALUES (10,'asd2','$2a$14$T9JaA1fXvHty0vvPPMZJ/ehRoDaBkXPKgvouL2uvqTl2IsuTvynJW','asd2@gmail.com','03.02 2023');
INSERT INTO "users" VALUES (12,'sass1','$2a$14$BaMvZfDTXrJlEehpkjFKkeikm4Xi4nLmf8wtmCh7OVbJl/AfaTUbu','asd5@gmail.com','03.02 2023');
INSERT INTO "users" VALUES (13,'viies','$2a$14$P6wDEhhzn3u17HH1BRbORuV.LVbPRsKLkgQQGTIjL1l7ab3g39OYO','asd@gmail.com','05.02 2023');
INSERT INTO "users" VALUES (14,'Huh','$2a$14$/t8js11JZMVr7s2a.052LezUSKUiheC5fDWM8gFDCKCgDprk28DQ2','huh@huh.huh','07.02 2023');
INSERT INTO "users" VALUES (15,'joel','$2a$14$sqf5Stu0zBTfE9J4wBL47OeijFNu5rnfu/qcN3zOGEZGAwJ251udi','joelimeil@gmail.com','07.02 2023');
INSERT INTO "users" VALUES (16,'uuskasutaja','$2a$14$b0pddKkpWytbpb4EtbHnueh4GueVZ.ZrYu8yP5orBGrFFH6hLzt6e','uuskasutaja7@mail.ee','08.02 2023');
INSERT INTO "session" VALUES (163,'55f07ebd-20ca-405b-8cc7-7b064dbb5389',15);
INSERT INTO "commentLikes" VALUES (184,'like',13,152);
INSERT INTO "commentLikes" VALUES (187,'like',13,153);
INSERT INTO "commentLikes" VALUES (188,'dislike',13,151);
INSERT INTO "commentLikes" VALUES (191,'like',13,154);
INSERT INTO "commentLikes" VALUES (192,'like',13,155);
INSERT INTO "commentLikes" VALUES (193,'like',13,149);
INSERT INTO "commentLikes" VALUES (194,'dislike',13,150);
INSERT INTO "commentLikes" VALUES (202,'dislike',13,145);
INSERT INTO "commentLikes" VALUES (206,'dislike',13,146);
INSERT INTO "commentLikes" VALUES (211,'like',13,156);
INSERT INTO "commentLikes" VALUES (212,'like',13,157);
INSERT INTO "commentLikes" VALUES (226,'like',14,158);
INSERT INTO "commentLikes" VALUES (228,'like',15,159);
INSERT INTO "commentLikes" VALUES (230,'like',15,160);
INSERT INTO "commentLikes" VALUES (231,'dislike',14,161);
INSERT INTO "commentLikes" VALUES (234,'like',15,162);
INSERT INTO "commentLikes" VALUES (236,'like',16,168);
INSERT INTO "commentLikes" VALUES (238,'like',15,168);
INSERT INTO "postlikes" VALUES (149,'like',13,87);
INSERT INTO "postlikes" VALUES (150,'like',13,78);
INSERT INTO "postlikes" VALUES (152,'like',13,77);
INSERT INTO "postlikes" VALUES (155,'like',13,81);
INSERT INTO "postlikes" VALUES (184,'like',14,79);
INSERT INTO "postlikes" VALUES (236,'like',14,88);
INSERT INTO "postlikes" VALUES (264,'like',14,89);
INSERT INTO "postlikes" VALUES (267,'dislike',15,88);
INSERT INTO "postlikes" VALUES (280,'like',15,89);
INSERT INTO "postlikes" VALUES (283,'dislike',0,89);
INSERT INTO "postlikes" VALUES (284,'like',0,88);
COMMIT;`)
}
