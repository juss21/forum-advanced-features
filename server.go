package main

import (
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

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

	errorp := true //error found boolean

	if r.Method == "GET" {
		if r.URL.Path == "/" {
			errorp = true
			homepage.Execute(w, Web) // at homepage launch homepage
		} else if r.URL.Path == "/login" {
			if Web.Loggedin && Web.Currentuser != "" {
				http.Redirect(w, r, "/", http.StatusSeeOther) //if already logged in redirect back to the homepage
			} else {
				loginpage.Execute(w, "") //at loginpage launch loginpage
			}
		} else if r.URL.Path == "/register" {
			if Web.Loggedin && Web.Currentuser != "" {
				http.Redirect(w, r, "/", http.StatusSeeOther) //if already logged in redirect back to the homepage
			} else {
				registerpage.Execute(w, "") //at registerpage launch registerpage
			}
		} else if r.URL.Path == "/members" {
			memberspage.Execute(w, Web) //at registerpage launch memberpage
		} else if r.URL.Path == "/logout" {
			if Web.Currentuser != "" {
				printLog("Server:", Web.Currentuser, "has logged out!")
				Web.Currentuser = ""
				Web.Loggedin = false
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther) //if logging out return the user back to the login page
		} else {
			// this one is for all the different forum pages:
			for i := 0; i < len(Web.Forum_data); i++ {
				if r.URL.Path[1:] == Web.Forum_data[i].Post_title {
					errorp = false //since a page was found set the boolean to false
					Web.Forum_data[i].Currentuser = Web.Currentuser
					Web.Forum_data[i].Loggedin = Web.Loggedin
					Web.tempint = i
					forumpage.Execute(w, Web.Forum_data[i]) //at "/Web.Forum_data[i].Post_title" launch forumpage with certain information
				}
			}
			if errorp {
				Web.ErrorMsg = "404 Page not Found!"
				errorpage.Execute(w, Web)
			}
		}
	} else if r.Method == "POST" {
		if r.URL.Path == "/" {
			post_header := r.FormValue("post_header")
			post_content := r.FormValue("post_content")
			if !Web.Loggedin {
				Web.ErrorMsg = "You must be logged on before you can post!"
				errorpage.Execute(w, Web) // if user is not logged in and is trying to post
			} else {
				//checks whether new post can be sent or not
				if sendPost(Web.Sqlbase, Web.Currentuser, post_header, post_content) {
					homepage.Execute(w, Web) //opening homepage with fresh data
				} else {
					errorpage.Execute(w, Web) //opening errorpage with errormessage
				}
			}
		} else if r.URL.Path == "/login" {
			user_name := r.FormValue("user_name")         // get username input
			user_password := r.FormValue("user_password") // get password input

			//login attempt #1
			if getLogin(user_name, user_password) {
				printLog("Server:", user_name, "has logged in!")
				Web.Currentuser = user_name
				Web.Loggedin = true
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				loginpage.Execute(w, "Please check your password and account name and try again.") //error output
			}
		} else if r.URL.Path == "/register" {
			user_name := r.FormValue("user_name")         // get username input
			user_password := r.FormValue("user_password") // get password input
			password_confirmation := r.FormValue("user_password_confirmation")
			user_email := r.FormValue("user_email") // get email input
			email_confirmation := r.FormValue("user_email_confirmation")

			isValid := getRegister(user_name, user_password, user_email, password_confirmation, email_confirmation)
			//if everything is correct
			if isValid {
				sendRegister(Web.Sqlbase, user_name, user_password, user_email) //registering account
				http.Redirect(w, r, "/login", http.StatusSeeOther)              //redirecting to login page
			} else {
				registerpage.Execute(w, Web.ErrorMsg) //send user back to register page with error output
			}
		} else {
			//this is for adding comments to the forumpage
			if Web.Currentuser != "" && Web.Loggedin {
				forum_commentbox := r.FormValue("forum_commentbox") // commentbox data
				Web.Currentpage = r.URL.Path[1:]                    //current page

				//sending a comment
				sendComment(Web.Sqlbase, Web.Currentuser, forum_commentbox, Web.Currentpage)

				//getting whether buttonclick was "like" or "dislike"
				like := r.FormValue("like")
				dislike := r.FormValue("dislike")

				if like != "" {
					//in case it was "like" sendLike
					sendLike(Web.Sqlbase, Web.Currentuser, Web.Currentpage, false)
					printLog("postitus: ", r.URL.Path, " sai like!")
				}
				if dislike != "" {
					//in case it was "dislike" sendDisLike
					sendDisLike(Web.Sqlbase, Web.Currentuser, Web.Currentpage, false)
					printLog("postitus: ", r.URL.Path, " sai dislike!")
				}
				//refreshing forumpage
				http.Redirect(w, r, Web.Currentpage, http.StatusSeeOther)
				forumpage.Execute(w, Web.Forum_data[Web.tempint])
			} else {
				Web.ErrorMsg = "You must be logged on to do that!"
				errorpage.Execute(w, Web) //opening errorpage with errormessage
			}
		}
	}
}
