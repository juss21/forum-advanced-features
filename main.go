package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/web/", http.StripPrefix("/web/", fs))

	mux.HandleFunc("/", rateLimiter(homePageHandle))
	mux.HandleFunc("/post/", rateLimiter(forumPageHandler))
	mux.HandleFunc("/login", rateLimiter(loginHandler))
	mux.HandleFunc("/logout", rateLimiter(logOutHandler))
	mux.HandleFunc("/register", rateLimiter(registerHandler))
	mux.HandleFunc("/members", rateLimiter(membersHandler))
	mux.HandleFunc("/comment", rateLimiter(commentHandler))
	mux.HandleFunc("/likePost", rateLimiter(postLikeHandler))
	mux.HandleFunc("/likeComment/", rateLimiter(commentLikeHandler))
	mux.HandleFunc("/account", rateLimiter(accountDetails))
	mux.HandleFunc("/changefilter", rateLimiter(filterHandler))

	fmt.Printf("Starting server at port " + port + "\n")
	if http.ListenAndServe(":"+port, mux) != nil {
		log.Fatal(err)
	}
}

/* func rateLimiterMiddleware(next func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	errorpage := ParseFiles("web/templates/error.html")
	header := ParseFiles("web/templates/header.html")
	user := request.RemoteAddr
	count := userRequestAmounts[user]
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !limiter.Allow() {
			writer.Write([]byte("rate limit exceeded "))
			return
		} else {
			endpointExample(writer, request)
		}
	})
} */

func rateLimiter(page http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.RemoteAddr
		go func() {
			time.Sleep(time.Minute)
			delete(userRequestAmounts, user)
		}()

		count := userRequestAmounts[user]
		if count > 50 {
			w.WriteHeader(429)
			createAndExecuteError(w, "Too many requests! Please wait a minute.")
			return
		}
		userRequestAmounts[user] += 1
		page.ServeHTTP(w, r)
	})
}
