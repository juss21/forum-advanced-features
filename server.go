package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createTemplate(fileName string) *template.Template {
	page, err := template.New("index.html").ParseFiles("web/templates/index.html", fmt.Sprint("web/templates/", fileName))
	if err != nil {
		log.Fatal(err)
	}
	return page
}

func homePageHandle(w http.ResponseWriter, r *http.Request) {
	errorpage := createTemplate("error.html")
	homepage := createTemplate("homepage.html")

	if r.URL.Path != "/" {
		w.WriteHeader(404)
		Web.ErrorMsg = "404! Page not found"
		errorpage.Execute(w, Web)
		Web.ErrorMsg = ""
		return
	}

	Web.Forum_data = AllPosts(Web.SelectedFilter)
	Web.Categories = getCategories()
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)

	// ClearCookies(w, r)
	switch r.Method {
	case "GET":
		data := Web
		err := homepage.Execute(w, data)
		if err != nil {
			fmt.Println(err)
		}
	case "POST":
		title := r.FormValue("post_header")
		content := r.FormValue("post_content")
		category, _ := strconv.Atoi(r.FormValue("category"))
		filterstatus := r.FormValue("categoryfilter")

		Web.SelectedFilter = filterstatus
		if title == "" || content == "" {
			w.WriteHeader(400)
			errorpage.Execute(w, "Error! Post title/content cannot be empty!")
			return
		}

		if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
			w.WriteHeader(400)
			errorpage.Execute(w, "You must be logged in before you post!")
			return
		}

		imageName, err := uploadFile(w, r)
		if err != nil {
			w.WriteHeader(400)
			errorpage.Execute(w, "File size too big")
			return
		}

		if !SavePost(title, Web.LoggedUser.ID, content, category, imageName) {
			errorpage.Execute(w, Web.ErrorMsg)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func forumPageHandler(w http.ResponseWriter, r *http.Request) {
	forumpage := createTemplate("forumpage.html")
	errorpage := createTemplate("error.html")

	postId, _ := strconv.Atoi(path.Base(r.URL.Path))
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r) // setting loggedin bool status depending on hasCookie result

	post, err := GetPostById(postId)
	if err != nil {
		w.WriteHeader(400)
		errorpage.Execute(w, "Post not Found")
		return
	}

	post.Comments = GetCommentsByPostId(postId)
	Web.CurrentPost = post

	switch r.Method {
	case "GET":
		forumpage.Execute(w, Web)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	loginpage := createTemplate("login.html")

	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result

	switch r.Method {
	case "GET":
		if Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		loginpage.Execute(w, Web)
	case "POST":
		user_name := r.FormValue("user_name")
		user_password := r.FormValue("user_password")

		user, err := Login(user_name, user_password)
		match := CheckPasswordHash(user_password, user.Password)

		if err != nil || !match {
			w.WriteHeader(400)
			Web.ErrorMsg = "Please check you password and username, might be incorrect"
			loginpage.Execute(w, Web)
			Web.ErrorMsg = ""
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

		Web.Loggedin = true
		Web.LoggedUser = user

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
	Web.CreatedPosts = []Createdstuff{} // TODO vaadata mida see värk siin teeb
	Web.LikedComments = []Likedstuff{}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	registerpage := createTemplate("register.html")
	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	switch r.Method {
	case "GET":
		if Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		registerpage.Execute(w, Web)
	case "POST":
		user_name := r.FormValue("user_name")
		user_password := r.FormValue("user_password")
		user_password_confirmation := r.FormValue("user_password_confirmation")

		user_email := r.FormValue("user_email")
		user_email_confirmation := r.FormValue("user_email_confirmation") // TODO Teha miskit kui molemad oleksid katki

		if user_password != user_password_confirmation {
			Web.ErrorMsg = "Passwords must be same"
			registerpage.Execute(w, Web)
			Web.ErrorMsg = ""
			return
		}

		if user_email != user_email_confirmation {
			Web.ErrorMsg = "Emails must be same"
			registerpage.Execute(w, Web)
			Web.ErrorMsg = ""
			return
		}

		hash, _ := HashPassword(user_password)

		err := Register(user_name, hash, user_email)
		if err != nil {
			Web.ErrorMsg = err.Error()
			registerpage.Execute(w, Web)
			Web.ErrorMsg = ""
			return
		}

		Web.ErrorMsg = "You have successfully registered! Please log in."
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
	memberspage := createTemplate("members.html")

	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	switch r.Method {
	case "GET":
		Web.User_data = GetUsers()
		memberspage.Execute(w, Web)
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	errorpage := createTemplate("error.html")

	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	switch r.Method {
	case "POST":
		comment := r.FormValue("forum_commentbox")

		if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
			w.WriteHeader(400)
			errorpage.Execute(w, "You must be logged in before you comment!")
			return
		}
		if SaveComment(comment, Web.LoggedUser.ID, Web.CurrentPost.Id) {
			postId := strconv.Itoa(Web.CurrentPost.Id)
			http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
		} else {
			errorpage.Execute(w, Web.ErrorMsg)
			return
		}
	}
}

func postLikeHandler(w http.ResponseWriter, r *http.Request) {
	errorpage := createTemplate("error.html")

	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
		w.WriteHeader(400)
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
	errorpage := createTemplate("error.html")

	Web.Loggedin = hasCookie(r) // setting loggedin bool status depending on hasCookie result
	ClearCookies(w, r)
	commentId, _ := strconv.Atoi(path.Base(r.URL.Path))

	if !Web.Loggedin { // kui objekt on tühi, siis pole keegi sisse loginud
		w.WriteHeader(400)
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
	accountpage := createTemplate("account.html")

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
