package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//port := "8080" // webserver port
	database, _ := sql.Open("sqlite3", "./members.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?,?)")
	statement.Exec("Red", "Ferrari") //exec first name, last name
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")

	var id int
	var name string
	var password string
	for rows.Next() {
		rows.Scan(&id, &name, &password)
		fmt.Println(id, " ", name, " ", password)
	}
	// fmt.Println("statement: ", statement)
	// fmt.Println("database: ", database)

	// http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web")))) // handling web folder
	// http.HandleFunc("/", serverHandle)                                                // server handle
	// fmt.Printf("Starting server at port " + port + "\n")
	// err := http.ListenAndServe(":"+port, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func serverHandle(w http.ResponseWriter, r *http.Request) {
	homepage, err := template.ParseFiles("web/index.html") // home template
	errorCheck(err, true)
	errorpage, err2 := template.ParseFiles("web/error.html") //error template
	errorCheck(err2, true)

	if r.Method == "GET" {
		if r.URL.Path == "/" {
			homepage.Execute(w, "") // at homepage print homepage
		} else {
			errorpage.Execute(w, "")
		}

	}
}

func errorCheck(err error, exit bool) {
	if err != nil {
		fmt.Println(err)
		if exit {
			os.Exit(0)
		}
	}
}
