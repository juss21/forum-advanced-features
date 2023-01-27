package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

type memberlist struct {
	ID              int
	Username        string
	Password        string
	Email           string
	Likedcontent    string
	Dislikedcontent string
}
type forumfamily struct {
	Originalposter string
	Post_title     string
	Post_content   string
	Commentor_data []commentpandemic
	Date_posted    string
	Post_likes     int
	Post_disLikes  int
}
type commentpandemic struct {
	Date             string
	Commentor        string
	Forum_comment    string
	Post_header      string
	Comment_likes    int
	Comment_disLikes int
}

type webstuff struct {
	Loggedin    bool
	Currentuser string
	Currentpage string
	Sqlbase     *sql.DB
	Userlist    []memberlist
	Forum_data  []forumfamily
}

var Web webstuff

func errorCheck(err error, exit bool) {
	if err != nil {
		fmt.Println(err)
		if exit {
			os.Exit(0)
		}
	}
}

func getLogin(uid string, password string) bool {
	for i := 0; i < len(Web.Userlist); i++ {
		// if username or email correct
		if Web.Userlist[i].Username == uid || Web.Userlist[i].Email == uid {
			// if password correct
			if password == Web.Userlist[i].Password {
				return true
			}
		}
	}
	return false
}

func getRegister(uid string, password string, email string) (bool, string) {
	str := ""
	for i := 0; i < len(Web.Userlist); i++ {
		// TODO add email check aswell?
		// check if user already exists
		if Web.Userlist[i].Username == uid {
			str += "u"
		} else if Web.Userlist[i].Email == email {
			str += "e"
		}
	}
	if strings.Contains(str, "u") {
		if strings.Contains(str, "e") {
			return false, "This username and e-mail is already in use!"
		}
		return false, "This username is already taken!"
	} else if strings.Contains(str, "e") && !strings.Contains(str, "u") {
		return false, "This e-mail is already in use!"
	}

	return true, "Account has been created!"
}
