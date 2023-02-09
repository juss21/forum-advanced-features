package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ParseFiles(filename string) *template.Template {
	temp, err := template.ParseFiles(filename) // home template
	errorCheck(err, true)

	return temp
}

func errorCheck(err error, fatal bool) {
	if err != nil {
		fmt.Println("Error():", err)
		if fatal {
			os.Exit(0)
		}
	}
}

func homePageHandle(w http.ResponseWriter, r *http.Request) {
	errorpage := ParseFiles("web/templates/error.html")
	header := ParseFiles("web/templates/header.html")
	homepage := ParseFiles("web/templates/index.html")

	Web.Forum_data = AllPosts(Web.SelectedFilter)
	Web.Categories = getCategories()
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result

	ClearCookies(w, r)
	switch r.Method {
	case "GET":	
		data := Web
		header.Execute(w, data)
		homepage.Execute(w, data)
	case "POST":
		title := r.FormValue("post_header")
		content := r.FormValue("post_content")
		category, _ := strconv.Atoi(r.FormValue("category"))
		filterstatus := r.FormValue("categoryfilter")

		Web.SelectedFilter = filterstatus
		if title == "" || content == "" {
			header.Execute(w, Web)
			errorpage.Execute(w, "Error! Post title/content cannot be empty!")
			return
		}

		if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
			header.Execute(w, Web)
			errorpage.Execute(w, "You must be logged in before you post!")
			return
		}

		imageName, err := uploadFile(w, r)
		if err != nil {
			header.Execute(w, Web)
			errorpage.Execute(w, "File size too big")
			return
		}

		if !SavePost(title, Web.LoggedUser.ID, content, category, imageName) {
			header.Execute(w, Web)
			errorpage.Execute(w, Web.ErrorMsg)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func forumPageHandler(w http.ResponseWriter, r *http.Request) {
	header := ParseFiles("web/templates/header.html")
	forumpage := ParseFiles("web/templates/forumpage.html")
	errorpage := ParseFiles("web/templates/error.html")

	postId, _ := strconv.Atoi(path.Base(r.URL.Path))
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result

	post, err := GetPostById(postId) // TODO implementeerida error kui pole ühtegi posti
	if err != nil {
		header.Execute(w, Web)
		errorpage.Execute(w, "Post not Found")
		return
	}

	post.Comments = GetCommentsByPostId(postId)
	Web.CurrentPost = post
	ClearCookies(w, r)

	switch r.Method {
	case "GET":
		post.Loggedin = Web.Loggedin
		header.Execute(w, Web)
		forumpage.Execute(w, post)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	loginpage := ParseFiles("web/templates/login.html")
	header := ParseFiles("web/templates/header.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	switch r.Method {
	case "GET":
		if Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		header.Execute(w, Web)
		loginpage.Execute(w, "")
	case "POST":
		user_name := r.FormValue("user_name")
		user_password := r.FormValue("user_password")

		user, err := Login(user_name, user_password)
		match := CheckPasswordHash(user_password, user.Password)

		if err != nil || !match {
			header.Execute(w, Web)
			loginpage.Execute(w, "Please check your password and account name and try again.")
			return
		}

		id, _ := uuid.NewV4()
		cookie := &http.Cookie{
			Name:    "session-id",
			Value:   id.String(),
			Expires: time.Now().Add(30 * time.Minute),
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		SaveSession(cookie.Value, user.ID)

		Web.Loggedin = hasCookie(r)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func logOutHandler(w http.ResponseWriter, r *http.Request) {
	ClearCookies(w, r)
	if !Web.Loggedin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	userId := Web.LoggedUser.ID

	cookie, _ := r.Cookie("session-id")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result

	http.SetCookie(w, &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	})

	Web.Loggedin = false
	DeleteSession(cookie.Value, userId)
	Web.LoggedUser = Memberlist{}
	Web.CreatedPosts = []Createdstuff{}
	Web.LikedComments = []Likedstuff{}
	http.Redirect(w, r, "/", http.StatusSeeOther) // TODO lisada sõnum, et on edukal välja logitud
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	registerpage := ParseFiles("web/templates/register.html")
	header := ParseFiles("web/templates/header.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	switch r.Method {
	case "GET":
		if Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		header.Execute(w, Web)
		registerpage.Execute(w, Web.ErrorMsg)
	case "POST":
		user_name := r.FormValue("user_name")         // text input
		user_password := r.FormValue("user_password") // font type
		user_email := r.FormValue("user_email")

		hash, _ := HashPassword(user_password)
		if CanRegister(user_name, hash, user_email, hash, user_email) {
			Register(user_name, hash, user_email) // TODO lisada sõnum, et edukalt registreeritud
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			header.Execute(w, Web)
			registerpage.Execute(w, Web.ErrorMsg)
		}
	}
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
	memberspage := ParseFiles("web/templates/members.html")
	header := ParseFiles("web/templates/header.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	switch r.Method {
	case "GET":
		header.Execute(w, Web)
		memberspage.Execute(w, Web.User_data)
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO mingi imelik Faviconi bug, template korda tegemisega saaks valmis vist
	errorpage := ParseFiles("web/templates/error.html")
	header := ParseFiles("web/templates/header.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	switch r.Method {
	case "POST":
		comment := r.FormValue("forum_commentbox") // TODO Kuvada kommentaari, teha like/dislike süsteem

		if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
			header.Execute(w, Web)
			errorpage.Execute(w, "You must be logged in before you comment!")
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
	errorpage := ParseFiles("web/templates/error.html")
	header := ParseFiles("web/templates/header.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
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
	errorpage := ParseFiles("web/templates/error.html")
	header := ParseFiles("web/templates/header.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	commentId, _ := strconv.Atoi(path.Base(r.URL.Path))

	if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
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
	accountpage := ParseFiles("web/templates/account.html")
	header := ParseFiles("web/templates/header.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	Web.CreatedPosts = []Createdstuff{}
	Web.LikedComments = []Likedstuff{}
	ClearCookies(w, r)
	switch r.Method {
	case "GET":
		if !Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		UserPosted()
		LikesSent()
		DateCreated()
		header.Execute(w, Web)
		accountpage.Execute(w, Web)
	}
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	filterstatus := r.FormValue("categoryfilter")
	ClearCookies(w, r)
	switch r.Method {
	case "GET":
		Web.SelectedFilter = filterstatus
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func uploadFile(w http.ResponseWriter, r *http.Request) (string, error) {
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		return "", nil
	}

	fileSize := 20971520
	if handler.Size > int64(fileSize) {
		return "", errors.New("cant be over 20mb")
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	tempFile, err := ioutil.TempFile("web/temp-images", "upload-*"+ext)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	return strings.Split(tempFile.Name(), "/")[2], nil
}
