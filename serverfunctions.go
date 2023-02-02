package main

import (
	"net/http"
	"strings"
)

func getCommentFormValue(r *http.Request, forum_header string) (datafound bool, like string, dislike string) {
	i := Web.tempint
	r.FormValue("like")
	if Web.Forum_data[i].Post_title == forum_header {
		//Web.Forum_data[i].Commentor_data -> peaks selle forumi commentor_data olema!

		return true, "", ""
	}

	return false, "", ""
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

func getRegister(uid string, password string, email string, cpassword string, cemail string) bool {
	str := ""
	if cpassword != password {
		Web.ErrorMsg = "The passwords do not match!"
		return false
	} else if cemail != email {
		Web.ErrorMsg = "The emails do not match!"
		return false
	}

	for i := 0; i < len(Web.Userlist); i++ {
		if Web.Userlist[i].Username == uid {
			str += "u"
		} else if Web.Userlist[i].Email == email {
			str += "e"
		}
	}
	if strings.Contains(str, "u") {
		if strings.Contains(str, "e") {
			Web.ErrorMsg = "This username and e-mail is already in use!"
			return false
		}
		Web.ErrorMsg = "This username is already taken!"
		return false
	} else if strings.Contains(str, "e") && !strings.Contains(str, "u") {
		Web.ErrorMsg = "This e-mail is already in use!"
		return false
	}

	return true
}