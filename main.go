package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var login bool = true

func main() {
	if login {
		Web.Loggedin = true
		Web.Currentuser = "Joel"
	}

	port := "8080" // webserver port
	database, err := sql.Open("sqlite3", "./database.db")
	errorCheck(err, true)
	Web.Sqlbase = database

	saveAllUsers(database)    // salvestame kõik kasutajad mällu, et hiljem oleks võimalik neid veebilehel lohistada
	saveAllPosts(database)    // salvestame kõik postitused mällu, et hiljem oleks võimalik neid veebilehel lohistada
	saveAllComments(database) // salvestame kõik kommentaarid mällu, et hiljem oleks võimalik neid veebilehel lohistada
	buildLikesStruct(Web.allcomments, len(Web.Userlist))
	buildDisLikesStruct(Web.allcomments, len(Web.Userlist))
	buildTopicLikesStruct(Web.allposts, len(Web.Userlist))
	buildTopicDisLikesStruct(Web.allposts, len(Web.Userlist))

	for i := 0; i < len(Web.TopicLikes); i++ {
		//	fmt.Println("links:", Web.TopicLikes[i].TopicID, Web.TopicLikes[i].UserID) //userid, roomid, status
	}
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web")))) // handling web folder
	http.HandleFunc("/", serverHandle)                                                // server handle
	fmt.Printf("Starting server at port " + port + "\n")

	if http.ListenAndServe(":"+port, nil) != nil {
		log.Fatal(err)
	}
}
