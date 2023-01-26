package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type memberlist struct {
	ID       int
	Username string
	Password string
	Email    string
}
type forumfamily struct {
	Originalposter string
	Post_title     string
	Post_content   string
}

var loggedin bool = false
var forum_op string
var sqlbase *sql.DB
var userlist []memberlist
var forum_data []forumfamily

func getLogin(uid string, password string) bool {
	for i := 0; i < len(userlist); i++ {
		//if username or email correct
		if userlist[i].Username == uid || userlist[i].Email == uid {
			// if password correct
			if password == userlist[i].Password {
				return true
			}
		}
	}
	return false
}
func getRegister(uid string, password string, email string) bool {
	for i := 0; i < len(userlist); i++ {
		//TODO add email check aswell?
		//check if user already exists
		if userlist[i].Username == uid {
			return false
		}
	}
	return true
}

func sendPost(database *sql.DB, originalposter string, header string, content string) {
	statement, _ := database.Prepare("INSERT INTO forum (originalposter, post_header, post_content) VALUES (?,?,?)")
	statement.Exec(originalposter, header, content) //exec first name, last name
	forum_data = append(forum_data, forumfamily{Originalposter: originalposter, Post_title: header, Post_content: content})
	fmt.Println("Server:", header, "has been posted!", " by <", originalposter, ">")

	//foorumi sisu printimine konsooli
	for i := 0; i < len(forum_data); i++ {
		fmt.Println(forum_data[i])
	}
}

func sendRegister(database *sql.DB, username string, password string, email string) {
	statement, _ := database.Prepare("INSERT INTO userdata (username, password, email) VALUES (?,?,?)")
	statement.Exec(username, password, email) //exec first name, last name
	userlist = append(userlist, memberlist{ID: len(userlist) + 1, Username: username, Password: password, Email: email})
	fmt.Println("Server:", username, "has successfully registered!", " <", email, ">")

	//kasutajate printimine konsooli
	for i := 0; i < len(userlist); i++ {
		fmt.Println(userlist[i])
	}
}

func saveAllPosts(database *sql.DB) {
	rows, _ := database.Query("SELECT originalposter, post_header, post_content FROM forum")
	var originalposter string
	var post_header string
	var post_content string

	for rows.Next() {
		rows.Scan(&originalposter, &post_header, &post_content)
		forum_data = append(forum_data, forumfamily{Originalposter: originalposter, Post_title: post_header, Post_content: post_content})
	}

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
	}

}
func main() {
	/*IDEED ja TODO
	  ideed: login/register võivad pesitseda samal lehel

		todo:  teha forum topicu html page hiljem saadame sinna infot nagu groupie-tracker'is
		kommentaarideks ja like/dislike jaoks tuleb sqli täiendada aga selle jätaks hiljemaks | delete pole vist veel required?
		Hetkel: Loggedin bool muutub true'ks kui sisse logida see ja võiks ära kaotada html'is login/register nupu aga kuidas seda teha?

		kui panna uue postituse title'sse öäöü läheb lolliks vist? või oli lihtsalt mul liiga pikk TEXT

		Loggedin bool vist äkki ei tohiks olla koodis, kuidas seda paremini teha?
		või igale kasutajale lisada boolean loggedOn? ja siis lihtsalt uue posti postitamisel teha mitu checki?
		selleks tuleks currentuser kuidagi html'i salvestada vist

		Bugid:
		peale registreerimist uue postituse lisamine registreerib uue kasutaja, üldse see login süsteem suht broken atm
		r.Method == "POST"  mingi kamm ilmselt

	  TODO ja IDEED*/

	port := "8080" // webserver port
	database, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("ERROR: Faulty database! ", err)
	}
	sqlbase = database
	saveAllUsers(database) // salvestame kõik kasutajad mällu, et hiljem oleks võimalik neid veebilehel lohistada
	saveAllPosts(database) // salvestame kõik postitused mällu, et hiljem oleks võimalik neid veebilehel lohistada

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
	loginpage, err2 := template.ParseFiles("web/login.html") //login template
	errorCheck(err2, true)
	registerpage, err3 := template.ParseFiles("web/register.html") //register template
	errorCheck(err3, true)
	errorpage, err4 := template.ParseFiles("web/error.html") //error template
	errorCheck(err4, true)
	memberspage, err5 := template.ParseFiles("web/members.html") //memberlist template
	errorCheck(err5, true)
	if r.Method == "GET" {
		if r.URL.Path == "/" {
			homepage.Execute(w, forum_data) // at homepage print homepage
		} else if r.URL.Path == "/login" {
			loginpage.Execute(w, "")
		} else if r.URL.Path == "/register" {
			registerpage.Execute(w, "")
		} else if r.URL.Path == "/members" {
			memberspage.Execute(w, userlist)
		} else {
			errorpage.Execute(w, "")
		}

	} else if r.Method == "POST" {
		if r.URL.Path == "/login" {
			user_name := r.FormValue("user_name")         //text input
			user_password := r.FormValue("user_password") //font type
			if getLogin(user_name, user_password) {
				homepage.Execute(w, forum_data)
				fmt.Println("Server:", user_name, "has logged in!")
				forum_op = user_name
				loggedin = true
			} else {
				loginpage.Execute(w, "Please check your password and account name and try again.")
			}
		}
		if r.URL.Path == "/register" {
			user_name := r.FormValue("user_name")         //text input
			user_password := r.FormValue("user_password") //font type
			password_confirmation := r.FormValue("user_password_confirmation")
			user_email := r.FormValue("user_email")
			email_confirmation := r.FormValue("user_email_confirmation")
			fmt.Println(user_email)
			if user_password == password_confirmation && user_email == email_confirmation {
				if getRegister(user_name, user_password, user_email) {
					sendRegister(sqlbase, user_name, user_password, user_email)
					homepage.Execute(w, forum_data)
				}
			} else {
				registerpage.Execute(w, "Password or E-mail does not match!")
			}
		}
		if r.URL.Path == "/" {
			post_header := r.FormValue("post_header")
			post_content := r.FormValue("post_content")
			fmt.Println(loggedin, forum_op, post_header, post_content)
			if loggedin {
				sendPost(sqlbase, forum_op, post_header, post_content)
				homepage.Execute(w, forum_data)
			} else {
				errorpage.Execute(w, "You must be logged in before you post!")
			}
		}
	}
}
