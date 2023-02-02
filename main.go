package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error
	port := "8080" // webserver port
	DataBase, _ = sql.Open("sqlite3", "./database.db")

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	http.HandleFunc("/", homePageHandle)
	http.HandleFunc("/post/", forumPageHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logOutHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/members", membersHandler)
	http.HandleFunc("/comment", commentHandler)
	http.HandleFunc("/likePost", postLikeHandler)
	http.HandleFunc("/likeComment/", commentLikeHandler)
	http.HandleFunc("/account", accountDetails)

	fmt.Printf("Starting server at port " + port + "\n")
	if http.ListenAndServe(":"+port, nil) != nil {
		log.Fatal(err)
	}
}
