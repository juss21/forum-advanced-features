package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type memberlist struct {
	ID       int
	Username string
	Password string
	Email    string
}

var userlist []memberlist

func sendRegister(database *sql.DB, username string, password string, email string) {
	statement, _ := database.Prepare("INSERT INTO userdata (username, password, email) VALUES (?,?,?)")
	statement.Exec(username, password, email) //exec first name, last name
}
func saveAllUsers(database *sql.DB) {
	rows, _ := database.Query("SELECT id, username, password, email FROM userdata")
	var id int
	var username string
	var password string
	var email string

	for rows.Next() {
		rows.Scan(&id, &username, &password, &email)
		userlist = append(userlist, memberlist{ID: id, Username: username, Password: password, Email: email})
		//fmt.Println(id, " ", username, " ", password, " ", email)
	}

}

func main() {
	port := "8080" // webserver port
	database, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("ERROR: Faulty database! ", err)
	}
	saveAllUsers(database) // salvestame kõik kasutajad mällu, et hiljem oleks võimalik neid veebilehel lohistada
	//sendRegister(database, "first", "last", "first.last@mail.ee") // function that will be used later in life

	//kasutajate printimine
	for i := 0; i < len(userlist); i++ {
		fmt.Println(userlist[i])
	}

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web")))) // handling web folder
	http.HandleFunc("/", serverHandle)                                                // server handle
	fmt.Printf("Starting server at port " + port + "\n")

	if http.ListenAndServe(":"+port, nil) != nil {
		log.Fatal(err)
	}
}

func serverHandle(w http.ResponseWriter, r *http.Request) {
	homepage, err := template.ParseFiles("web/index.html") // home template
	errorCheck(err, true)
	loginpage, err2 := template.ParseFiles("web/login.html") //error template
	errorCheck(err2, true)
	registerpage, err3 := template.ParseFiles("web/register.html") //error template
	errorCheck(err3, true)
	errorpage, err4 := template.ParseFiles("web/error.html") //error template
	errorCheck(err4, true)
	memberspage, err5 := template.ParseFiles("web/members.html") //error template
	errorCheck(err5, true)
	if r.Method == "GET" {
		if r.URL.Path == "/" {
			homepage.Execute(w, "") // at homepage print homepage
		} else if r.URL.Path == "/login" {
			loginpage.Execute(w, "")
		} else if r.URL.Path == "/register" {
			registerpage.Execute(w, "")
		} else if r.URL.Path == "/members" {
			memberspage.Execute(w, userlist)
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
