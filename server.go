package main

import (
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/gofrs/uuid"
)

func ParseFiles(filename string) *template.Template {
	temp, _ := template.ParseFiles(filename) // home template
	// errorCheck(err, true)

	return temp
}
func homePageHandle(w http.ResponseWriter, r *http.Request) {
	errorpage := ParseFiles("web/error.html")
	header := ParseFiles("web/header.html")
	homepage := ParseFiles("web/index.html")

	Web.Forum_data = AllPosts()

	switch r.Method {
	case "GET":
		header.Execute(w, Web)
		homepage.Execute(w, Web)
	case "POST":
		title := r.FormValue("post_header")
		content := r.FormValue("post_content")

		if Web.LoggedUser == (Memberlist{}) { // kui objekt on tühi, siis pole keegi sisse loginud
			Web.ErrorMsg = "You must be logged in before you post!"
			header.Execute(w, Web)
			errorpage.Execute(w, Web.ErrorMsg)
			return
		}
		if !SavePost(title, Web.LoggedUser.ID, content) {
			header.Execute(w, Web)
			errorpage.Execute(w, Web.ErrorMsg)
			return
		}

		Web.Forum_data = AllPosts()
		header.Execute(w, Web)
		homepage.Execute(w, Web)
	}
}

func forumPageHandler(w http.ResponseWriter, r *http.Request) {
	header := ParseFiles("web/header.html")
	forumpage := ParseFiles("web/forumpage.html")

	// errorpage := ParseFiles("web/error.html")
	postId, _ := strconv.Atoi(path.Base(r.URL.Path))

	post := GetPostById(postId) // TODO implementeerida error kui pole ühtegi posti
	post.Comments = GetCommentsByPostId(postId)
	Web.CurrentPost = post

	switch r.Method {
	case "GET":
		post.Loggedin = Web.Loggedin

		header.Execute(w, Web)
		forumpage.Execute(w, post)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO lisada tagasi osa,kus peidab kui on logitud ja näitab kui vaja välja logida
	loginpage := ParseFiles("web/login.html")
	//	header := ParseFiles("web/header.html")

	switch r.Method {
	case "GET":
		loginpage.Execute(w, "")
	case "POST":
		user_name := r.FormValue("user_name")
		user_password := r.FormValue("user_password")
		user, err := Login(user_name, user_password)
		if err != nil {
			feedback := "Please check your password and account name and try again."
			loginpage.Execute(w, feedback)
			return
		}

		Web.LoggedUser = Memberlist{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		Web.Loggedin = true
		cookie, err := r.Cookie("session-id")

		if err != nil {
			id, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:  "session-id",
				Value: id.String(),
			}

			http.SetCookie(w, cookie)
			SaveSession(cookie.Value, user.ID)

		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func logOutHandler(w http.ResponseWriter, r *http.Request) {
	Web.LoggedUser = Memberlist{}
	Web.Loggedin = false
	Web.CreatedPosts = []Createdstuff{}
	http.Redirect(w, r, "/", http.StatusSeeOther) // TODO lisada sõnum, et on edukal välja logitud
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	registerpage := ParseFiles("web/register.html")
	header := ParseFiles("web/header.html")

	switch r.Method {
	case "GET":
		header.Execute(w, Web)
		registerpage.Execute(w, "")
	case "POST":
		user_name := r.FormValue("user_name")         // text input
		user_password := r.FormValue("user_password") // font type
		user_email := r.FormValue("user_email")

		Register(user_name, user_password, user_email) // TODO lisada sõnum, et edukalt registreeritud
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
	memberspage := ParseFiles("web/members.html")
	header := ParseFiles("web/header.html")

	switch r.Method {
	case "GET":
		members := GetUsers()
		header.Execute(w, Web)
		memberspage.Execute(w, members)
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO mingi imelik Faviconi bug, template korda tegemisega saaks valmis vist
	// forumpage := ParseFiles("web/forumpage.html")
	errorpage := ParseFiles("web/error.html")
	header := ParseFiles("web/header.html")

	// postId := path.Base(r.URL.Path)

	switch r.Method {
	case "POST":
		comment := r.FormValue("forum_commentbox") // TODO Kuvada kommentaari, teha like/dislike süsteem

		if Web.LoggedUser == (Memberlist{}) { // kui objekt on tühi, siis pole keegi sisse loginud
			Web.ErrorMsg = "You must be logged in before you comment!"
			header.Execute(w, Web)
			errorpage.Execute(w, Web.ErrorMsg)
			return
		}
		if SaveComment(comment, Web.LoggedUser.ID, Web.CurrentPost.Id) {
			postId := strconv.Itoa(Web.CurrentPost.Id)
			http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
		} else {
			header.Execute(w, Web)
			errorpage.Execute(w, Web.ErrorMsg)
			return
		}
	}
}

func postLikeHandler(w http.ResponseWriter, r *http.Request) {
	errorpage := ParseFiles("web/error.html")
	header := ParseFiles("web/header.html")

	if Web.LoggedUser == (Memberlist{}) { // kui objekt on tühi, siis pole keegi sisse loginud
		header.Execute(w, Web)
		errorpage.Execute(w, "You must be logged in before you Like!")
		return
	}
	postLike := r.FormValue("button")

	switch r.Method {
	case "POST":

		postId := strconv.Itoa(Web.CurrentPost.Id)
		SavePostLike(postLike, Web.LoggedUser.ID, Web.CurrentPost.Id)
		http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
	}
}

func commentLikeHandler(w http.ResponseWriter, r *http.Request) {
	errorpage := ParseFiles("web/error.html")
	header := ParseFiles("web/header.html")

	commentId, _ := strconv.Atoi(path.Base(r.URL.Path))

	if Web.LoggedUser == (Memberlist{}) { // kui objekt on tühi, siis pole keegi sisse loginud
		header.Execute(w, Web)
		errorpage.Execute(w, "You must be logged in before you Like!")
		return
	}
	postLike := r.FormValue("button")

	switch r.Method {
	case "POST":
		postId := strconv.Itoa(Web.CurrentPost.Id)
		SaveCommentLike(postLike, Web.LoggedUser.ID, commentId)
		http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
	}
}

func accountDetails(w http.ResponseWriter, r *http.Request) {
	accountpage := ParseFiles("web/account.html")
	header := ParseFiles("web/header.html")
	Web.CreatedPosts = []Createdstuff{}
	Web.LikedComments = []Likedstuff{}
	switch r.Method {
	case "GET":
		UserPosted()
		LikesSent()

		header.Execute(w, Web)
		accountpage.Execute(w, Web)
	}
}
