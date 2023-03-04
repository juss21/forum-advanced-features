package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Server() {
	var err error
	port := "8080" // webserver port

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

	mux.HandleFunc("/editPost", rateLimiter(postEditHandler))
	mux.HandleFunc("/editComment", rateLimiter(postEditHandler))
	mux.HandleFunc("/activity", rateLimiter(notificationHandler))
	fmt.Printf("Starting server at port " + port + "\n")
	if http.ListenAndServe(":"+port, mux) != nil {
		log.Fatal(err)
	}
}

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
