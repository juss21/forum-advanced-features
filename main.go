package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func main() {
	/*IDEED ja TODO
	  ideed: login/register võivad pesitseda samal lehel
	todo:
		Loggedin bool muutub hetkel true'ks kui sisse logida see ja võiks ära kaotada html'is login/register nupu aga kuidas seda teha?
		Loggedin bool vist äkki ei tohiks olla koodis, kuidas seda paremini teha?
		selleks tuleks currentuser kuidagi html'i salvestada vist
	TODO ja IDEED*/
	port := "8080" // webserver port
	database, err := sql.Open("sqlite3", "./database.db")
	errorCheck(err, true)
	Web.Sqlbase = database

	saveAllUsers(database)    // salvestame kõik kasutajad mällu, et hiljem oleks võimalik neid veebilehel lohistada
	saveAllPosts(database)    // salvestame kõik postitused mällu, et hiljem oleks võimalik neid veebilehel lohistada
	saveAllComments(database) // salvestame kõik kommentaarid mällu, et hiljem oleks võimalik neid veebilehel lohistada

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web")))) // handling web folder
	http.HandleFunc("/", serverHandle)                                                // server handle
	fmt.Printf("Starting server at port " + port + "\n")

	if http.ListenAndServe(":"+port, nil) != nil {
		log.Fatal(err)
	}
}
