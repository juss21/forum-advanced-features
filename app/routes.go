package app

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func createAndExecute(w http.ResponseWriter, fileName string) {
	page, err := template.New("index.html").ParseFiles("web/templates/index.html", fmt.Sprint("web/templates/", fileName))
	if err != nil {
		fmt.Println(err.Error())

		createAndExecuteError(w, "500 Internal Server Error")
		return
	}
	err = page.Execute(w, Web)
	if err != nil {
		createAndExecuteError(w, "500 Internal Server Error")
		fmt.Println(err.Error())
		return
	}
}

func createAndExecuteError(w http.ResponseWriter, msg string) {
	page, _ := template.New("index.html").ParseFiles("web/templates/index.html", fmt.Sprint("web/templates/", "error.html"))
	Web.ErrorMsg = msg
	page.Execute(w, Web)
	Web.ErrorMsg = ""
}

func homePageHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		createAndExecuteError(w, "404! Page not found")
		return
	}

	Web.Forum_data = AllPosts(Web.SelectedFilter)
	Web.Categories = getCategories()
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)
	Web.Notifications = GetNotifications()
	Web.LoggedUser.Notifications = len(Web.Notifications)
	switch r.Method {
	case "GET":
		createAndExecute(w, "homepage.html")
	case "POST":
		title := r.FormValue("post_header")
		content := r.FormValue("post_content")
		category := (r.Form["category"])
		filterstatus := r.FormValue("categoryfilter")

		Web.SelectedFilter = filterstatus
		if title == "" || content == "" {
			w.WriteHeader(400)
			createAndExecuteError(w, "Error! Post title/content cannot be empty!")
			return
		}

		if !Web.Loggedin {
			w.WriteHeader(400)
			createAndExecuteError(w, "You must be logged in before you post!")
			return
		}

		imageName, err := uploadFile(w, r)
		if err != nil {
			w.WriteHeader(400)
			createAndExecuteError(w, "File size too big")
			return
		}

		for _, v := range category {
			convert, _ := strconv.Atoi(v)
			if !SavePost(title, Web.LoggedUser.ID, content, convert, imageName) {
				Web.ErrorMsg = "You must be logged in before you post!"
				createAndExecute(w, "homepage.html")
				Web.ErrorMsg = ""
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func forumPageHandler(w http.ResponseWriter, r *http.Request) {
	postId, _ := strconv.Atoi(path.Base(r.URL.Path))
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)

	DeleteNoticfication(Web.LoggedUser.ID, postId)
	Web.Notifications = GetNotifications()
	Web.LoggedUser.Notifications = len(Web.Notifications)
	post, err := GetPostById(postId)
	if err != nil {
		w.WriteHeader(400)
		createAndExecuteError(w, "Post not Found")
		return
	}

	post.Comments = GetCommentsByPostId(postId)
	Web.CurrentPost = post

	switch r.Method {
	case "GET":
		createAndExecute(w, "forumpage.html")
	case "POST":

		if !Web.Loggedin {
			createAndExecuteError(w, "You must be logged in, you 🦀")
			return

		}
		if r.FormValue("deletePost") != "" {

			DeletePostById(strconv.Itoa(postId))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		if r.FormValue("editPost") != "" {
			Web.CurrentPost.Edit = true

			createAndExecute(w, "forumpage.html")
		}

		if r.FormValue("cancelPost") != "" {
			Web.CurrentPost.Edit = false
			createAndExecute(w, "forumpage.html")
		}

		if r.FormValue("savePost") != "" {
			title, content := r.FormValue("post_header"), r.FormValue("post_content")

			Web.CurrentPost.Edit = false
			imageName, err := uploadFile(w, r)
			if err != nil {
				fmt.Println(err)
				createAndExecuteError(w, "505 INTERNAL ERROR")
			}
			EditPostById(postId, title, imageName, content)
			http.Redirect(w, r, "/post/"+strconv.Itoa(postId), http.StatusSeeOther)

		}
		if r.FormValue("deleteComment") != "" {
			a, _ := strconv.ParseInt(r.FormValue("deleteComment"), 10, 64)

			DeleteCommentById(strconv.Itoa(int(a)), postId)
			http.Redirect(w, r, "/post/"+strconv.Itoa(postId), http.StatusSeeOther)
		}

		if r.FormValue("editComment") != "" {
			Web.CurrentComment.Edit = true
			Web.CurrentComment.Id, _ = strconv.Atoi(r.FormValue("editComment"))
			createAndExecute(w, "forumpage.html")
		}

		if r.FormValue("cancel") != "" {
			Web.CurrentComment.Edit = false
			createAndExecute(w, "forumpage.html")
			http.Redirect(w, r, "/post/"+strconv.Itoa(postId), http.StatusSeeOther)
		}

		if r.FormValue("Save") != "" {

			content := r.FormValue("comment_content")
			Web.CurrentComment.Edit = false
			EditCommentById(Web.CurrentComment.Id, content)

			http.Redirect(w, r, "/post/"+strconv.Itoa(postId), http.StatusSeeOther)

		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)

	switch r.Method {
	case "GET":
		if Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		createAndExecute(w, "login.html")
		if Web.ErrorMsg == "You have successfully registered! Please log in." {
			Web.ErrorMsg = ""
		}
	case "POST":
		user_name := r.FormValue("user_name")
		user_password := r.FormValue("user_password")

		user, err := Login(user_name, user_password)
		match := CheckPasswordHash(user_password, user.Password)

		if err != nil || !match {
			w.WriteHeader(400)
			Web.ErrorMsg = "Please check you password and username, might be incorrect"
			createAndExecute(w, "login.html")
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
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)
	http.SetCookie(w, &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	})

	Web.Loggedin = false

	Web.LoggedUser = Memberlist{}
	Web.CreatedPosts = []CreatedPosts{}
	Web.LikedPosts = []LikedPosts{}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)

	switch r.Method {
	case "GET":
		if Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		createAndExecute(w, "register.html")
	case "POST":
		user_name := r.FormValue("user_name")
		user_password := r.FormValue("user_password")
		user_password_confirmation := r.FormValue("user_password_confirmation")

		user_email := r.FormValue("user_email")
		user_email_confirmation := r.FormValue("user_email_confirmation")

		if user_password != user_password_confirmation {
			Web.ErrorMsg = "Passwords must be same"
		} else if user_email != user_email_confirmation {
			Web.ErrorMsg = "Emails must be same"
		}

		if Web.ErrorMsg != "" {
			createAndExecute(w, "register.html")
			Web.ErrorMsg = ""
			return
		}

		hash, _ := HashPassword(user_password)

		err := Register(user_name, hash, user_email)
		if err != nil {
			Web.ErrorMsg = err.Error()
			createAndExecute(w, "register.html")
			Web.ErrorMsg = ""
			return
		}
		Web.ErrorMsg = "You have successfully registered! Please log in."
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func membersHandler(w http.ResponseWriter, r *http.Request) {
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)

	switch r.Method {
	case "GET":
		Web.Notifications = GetNotifications()
		Web.LoggedUser.Notifications = len(Web.Notifications)
		Web.User_data = GetUsers()
		createAndExecute(w, "members.html")
	}
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		createAndExecuteError(w, "We know where you live")
	case "POST":
		Web.LoggedUser, Web.Loggedin = getUserFromSession(r)

		comment := r.FormValue("forum_commentbox")

		if !Web.Loggedin {
			w.WriteHeader(400)
			createAndExecuteError(w, "You must be logged in before you comment!")
			return
		}
		if SaveComment(comment, Web.LoggedUser.ID, Web.CurrentPost.Id, Web.CurrentPost.Title, Web.LoggedUser.Username, Web.CurrentPost.UserId) {
			postId := strconv.Itoa(Web.CurrentPost.Id)
			http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
		} else {
			createAndExecuteError(w, "You must be logged in before you comment!")
			return
		}
	}
}

func postLikeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		createAndExecuteError(w, "We know where you live")
	case "POST":
		Web.LoggedUser, Web.Loggedin = getUserFromSession(r)

		if !Web.Loggedin {
			w.WriteHeader(400)
			createAndExecuteError(w, "You must be logged in before you Like!")
			return
		}
		postLike := r.FormValue("button")
		postId := strconv.Itoa(Web.CurrentPost.Id)
		SavePostLike(postLike, Web.LoggedUser.ID, Web.CurrentPost.Id, Web.CurrentPost.Title)
		http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
	}
}

func commentLikeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		createAndExecuteError(w, "We know where you live")
	case "POST":
		Web.LoggedUser, Web.Loggedin = getUserFromSession(r)
		postLike := r.FormValue("button")

		commentId, _ := strconv.Atoi(path.Base(r.URL.Path))

		if !Web.Loggedin {
			w.WriteHeader(400)
			createAndExecuteError(w, "You must be logged in before you Like!")
			return
		}

		postId := strconv.Itoa(Web.CurrentPost.Id)
		SaveCommentLike(postLike, Web.LoggedUser.ID, commentId)
		http.Redirect(w, r, "/post/"+postId, http.StatusSeeOther)
	}
}

func accountDetails(w http.ResponseWriter, r *http.Request) {
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)
	Web.CreatedComments = []CreatedComments{}
	Web.CreatedPosts = []CreatedPosts{}
	Web.LikedPosts = []LikedPosts{}
	Web.DisLikedPosts = []LikedPosts{}
	Web.LikedComments = []LikedComments{}
	Web.DisLikedComments = []LikedComments{}

	switch r.Method {
	case "GET":
		if !Web.Loggedin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		Web.Notifications = GetNotifications()
		Web.LoggedUser.Notifications = len(Web.Notifications)
		UserPosted()
		GetLikedComments("like")
		GetLikedComments("dislike")
		GetCreatedComments()
		GetLikedPosts("like")
		GetLikedPosts("dislike")
		DateCreated()
		createAndExecute(w, "account.html")
	}
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)
	filterstatus := r.FormValue("categoryfilter")

	switch r.Method {
	case "GET":
		Web.SelectedFilter = filterstatus
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func postEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || !Web.Loggedin {
		createAndExecuteError(w, "We know where you live")
		return
	}
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

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	Web.LoggedUser, Web.Loggedin = getUserFromSession(r)
	Web.Notifications = GetNotifications()
	Web.LoggedUser.Notifications = len(Web.Notifications)

	switch r.Method {
	case "GET":
		createAndExecute(w, "activity.html")
	}
}
