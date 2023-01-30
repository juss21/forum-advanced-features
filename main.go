package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func main() {
	Web.Loggedin = true
	Web.Currentuser = "Joel"
	port := "8081" // webserver port
	database, err := sql.Open("sqlite3", "./database.db")
	errorCheck(err, true)
	Web.Sqlbase = database

	saveAllUsers(database)    // salvestame kõik kasutajad mällu, et hiljem oleks võimalik neid veebilehel lohistada
	saveAllPosts(database)    // salvestame kõik postitused mällu, et hiljem oleks võimalik neid veebilehel lohistada
	saveAllComments(database) // salvestame kõik kommentaarid mällu, et hiljem oleks võimalik neid veebilehel lohistada

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web")))) // handling web folder
	http.HandleFunc("/", serverHandle)                                                // server handle
	fmt.Printf("Starting server at port " + port + "\n")

	fmt.Println(Web.allcomments)
	if http.ListenAndServe(":"+port, nil) != nil {
		log.Fatal(err)
	}
}
