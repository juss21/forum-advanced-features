package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error
	port := "8080" // webserver port

	file, err := os.Stat("database.db")

	// if .db file deleted, it will create new one and populate with data
	if errors.Is(err, os.ErrNotExist) {
		DataBase, _ = sql.Open("sqlite3", "database.db")
		InitDatabase()
		fmt.Println("New database created ", file)
	} else {
		DataBase, _ = sql.Open("sqlite3", "./database.db")
	}

	GetUsers()
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/web/", http.StripPrefix("/web/", fs))

	mux.HandleFunc("/", homePageHandle)
	mux.HandleFunc("/post/", forumPageHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logOutHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/members", membersHandler)
	mux.HandleFunc("/comment", commentHandler)
	mux.HandleFunc("/likePost", postLikeHandler)
	mux.HandleFunc("/likeComment/", commentLikeHandler)
	mux.HandleFunc("/account", accountDetails)
	mux.HandleFunc("/changefilter", filterHandler)

	fmt.Printf("Starting server at port " + port + "\n")
	if http.ListenAndServe(":"+port, mux) != nil {
		log.Fatal(err)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}
}
