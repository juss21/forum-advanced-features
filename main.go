package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func sendRegister(database *sql.DB, username string, password string, email string) {
	statement, _ := database.Prepare("INSERT INTO userdata (username, password, email) VALUES (?,?,?)")
	statement.Exec(username, password, email) //exec first name, last name
}
func printAllUsers(database *sql.DB) {
	rows, _ := database.Query("SELECT id, username, password, email FROM userdata")
	var id int
	var name string
	var password string
	var email string

	for rows.Next() {
		rows.Scan(&id, &name, &password, &email)
		fmt.Println(id, " ", name, " ", password, " ", email)
	}

}

func main() {
	//port := "8080" // webserver port
	database, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("ERROR: Faulty database! ", err)
	}
	//sendRegister(database, "first", "last", "first.last@mail.ee") // function that will be used later in life
	printAllUsers(database)

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
