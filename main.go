package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
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

	sqlbase = database
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

func ParseFiles(filename string) *template.Template {
	temp, err := template.ParseFiles(filename) // home template
	errorCheck(err, true)

	return temp
}

func serverHandle(w http.ResponseWriter, r *http.Request) {
	homepage := ParseFiles("web/index.html")
	loginpage := ParseFiles("web/login.html")
	registerpage := ParseFiles("web/register.html")
	errorpage := ParseFiles("web/error.html")
	memberspage := ParseFiles("web/members.html")
	forumpage := ParseFiles("web/forumpage.html")

	forum_op = "Joel"
	loggedin = true

	errorp := true

	if r.Method == "GET" {
		if r.URL.Path == "/" {
			errorp = true
			homepage.Execute(w, forum_data) // at homepage print homepage
		} else {
			//looping forum data
			for i := 0; i < len(forum_data); i++ {
				if forum_data[i].Post_title == r.URL.Path[1:] {
					errorp = false
					forumpage.Execute(w, forum_data[i])
				}
			}
			if r.URL.Path == "/login" {
				loginpage.Execute(w, "")
			} else if r.URL.Path == "/register" {
				registerpage.Execute(w, "")
			} else if r.URL.Path == "/members" {
				memberspage.Execute(w, userlist)
			} else {
				if errorp {
					errorpage.Execute(w, "404 Page not found!")
				}
			}
		}
	} else if r.Method == "POST" {
		if r.URL.Path == "/" {
			post_header := r.FormValue("post_header")
			post_content := r.FormValue("post_content")
			if !loggedin {
				errorpage.Execute(w, "You must be logged in before you post!")
			} else {
				cansend, errormsg := sendPost(sqlbase, forum_op, post_header, post_content)
				if cansend {
					homepage.Execute(w, forum_data)
				} else {
					errorpage.Execute(w, errormsg)
				}
			}
		} else if r.URL.Path == "/login" {
			user_name := r.FormValue("user_name")         // text input
			user_password := r.FormValue("user_password") // font type
			str := "Please check your password and account name and try again."
			if getLogin(user_name, user_password) {
				str = "Sign in was successful"
				fmt.Println("Server:", user_name, "has logged in!")
				forum_op = user_name
				loggedin = true
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				loginpage.Execute(w, str)
			}
		} else if r.URL.Path == "/register" {
			user_name := r.FormValue("user_name")         // text input
			user_password := r.FormValue("user_password") // font type
			password_confirmation := r.FormValue("user_password_confirmation")
			user_email := r.FormValue("user_email")
			email_confirmation := r.FormValue("user_email_confirmation")
			//str := "Password or E-mail does not match!"
			isValid := false
			output := ""
			if user_password == password_confirmation && user_email == email_confirmation {
				isValid, output = getRegister(user_name, user_password, user_email)
				if isValid {
					sendRegister(sqlbase, user_name, user_password, user_email)
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				} else {
					registerpage.Execute(w, output)
				}
			}
		} else {
			commentor := forum_op
			if commentor != "" {
				forum_commentbox := r.FormValue("forum_commentbox")
				currenturl := r.URL.Path[1:]

				sendComment(sqlbase, forum_op, forum_commentbox, currenturl)
				// fmt.Println(currenturl, commentor_data[0].Post_header)
				http.Redirect(w, r, currenturl, http.StatusSeeOther)

				forumpage.Execute(w, forum_data)

			} else {
				errorpage.Execute(w, "You must be logged in to do that!")
			}
		}
	}
}
