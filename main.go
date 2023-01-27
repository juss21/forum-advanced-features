package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	errorp := true

	if r.Method == "GET" {
		if r.URL.Path == "/" {
			errorp = true
			homepage.Execute(w, Web) // at homepage print homepage
		} else {
			//looping forum data
			for i := 0; i < len(Web.Forum_data); i++ {
				if Web.Forum_data[i].Post_title == r.URL.Path[1:] {
					errorp = false
					forumpage.Execute(w, Web.Forum_data[i])
				}
			}
			if r.URL.Path == "/login" {
				if Web.Loggedin && Web.Currentuser != "" {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				} else {
					loginpage.Execute(w, "")
				}
			} else if r.URL.Path == "/register" {
				if Web.Loggedin && Web.Currentuser != "" {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				} else {
					registerpage.Execute(w, "")
				}
			} else if r.URL.Path == "/members" {
				memberspage.Execute(w, Web)
			} else if r.URL.Path == "/logout" {
				Web.Currentuser = ""
				Web.Loggedin = false
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				for i := 0; i < len(Web.Forum_data); i++ {
					if strings.Contains(r.URL.Path, Web.Forum_data[i].Post_title) {
						Web.Currentpage = r.URL.Path[1:]
						like := r.FormValue("like")
						dislike := r.FormValue("dislike")
						if like != "" {
							sendLike(Web.Sqlbase, Web.Currentuser, Web.Currentpage, false)
							fmt.Println("postitus: ", r.URL.Path, " sai like!")
							Web.Sqlbase.Exec("UPDATE forum SET likes = Web.Forum_data[topicid].Post_likes WHERE post_header = title")
							Web.Sqlbase.Exec("UPDATE forum SET dislikes = Web.Forum_data[topicid].Post_dislikes WHERE post_header = title")
							Web.Sqlbase.Exec("UPDATE userdata SET likedcontent = Web.Userlist[userid].Likedcontent WHERE username = user")
							Web.Sqlbase.Exec("UPDATE userdata SET dislikedcontent = Web.Userlist[userid].Dislikedcontent WHERE username = user")
						}
						if dislike != "" {
							sendDisLike(Web.Sqlbase, Web.Currentuser, Web.Currentpage, false)
							fmt.Println("postitus: ", r.URL.Path, " sai dislike!")
							Web.Sqlbase.Exec("UPDATE forum SET likes = Web.Forum_data[topicid].Post_likes WHERE post_header = title")
							Web.Sqlbase.Exec("UPDATE forum SET dislikes = Web.Forum_data[topicid].Post_dislikes WHERE post_header = title")
							Web.Sqlbase.Exec("UPDATE userdata SET likedcontent = Web.Userlist[userid].Likedcontent WHERE username = user")
							Web.Sqlbase.Exec("UPDATE userdata SET dislikedcontent = Web.Userlist[userid].Dislikedcontent WHERE username = user")
						}
						http.Redirect(w, r, Web.Currentpage, http.StatusSeeOther)

						// //or j := 0; j < len(Web.Forum_data[i].Commentor_data); j++ {
						// comment_like := r.FormValue("comment_like")
						// comment_dislike := r.FormValue("comment_dislike")
						// hastag := 0
						// if comment_like != "" {
						// 	fmt.Println("kommentaar: ", r.URL.Path, " sai like!")

						// 	for url := 0; url < len(r.URL.Path); url++ {
						// 		if string(r.URL.Path[url]) == "#" {
						// 			fmt.Println("jah!")
						// 			hastag = url
						// 		}
						// 	}
						// 	fmt.Println(hastag, r.URL.Path[hastag:])
						// }
						// if comment_dislike != "" {
						// 	fmt.Println("kommentaar: ", r.URL.Path, " sai dislike!")
						// }
						//}
					}
				}

				if errorp {
					errorpage.Execute(w, "404 Page not found!")
				}
			}
		}
	} else if r.Method == "POST" {
		if r.URL.Path == "/" {
			post_header := r.FormValue("post_header")
			post_content := r.FormValue("post_content")
			if !Web.Loggedin {
				errorpage.Execute(w, "You must be logged in before you post!")
			} else {
				cansend, errormsg := sendPost(Web.Sqlbase, Web.Currentuser, post_header, post_content)
				if cansend {
					homepage.Execute(w, Web)
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
				Web.Currentuser = user_name
				Web.Loggedin = true
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
					sendRegister(Web.Sqlbase, user_name, user_password, user_email)
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				} else {
					registerpage.Execute(w, output)
				}
			}
		} else {
			if Web.Currentuser != "" {
				forum_commentbox := r.FormValue("forum_commentbox")
				currenturl := r.URL.Path[1:]

				sendComment(Web.Sqlbase, Web.Currentuser, forum_commentbox, currenturl)
				// fmt.Println(currenturl, commentor_data[0].Post_header)
				http.Redirect(w, r, currenturl, http.StatusSeeOther)

				forumpage.Execute(w, Web.Forum_data)
			} else {
				errorpage.Execute(w, "You must be logged in to do that!")
			}
		}
	}
}
