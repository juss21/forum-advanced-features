package main

import (
	"fmt"
	"net/http"
)

func getUserIDFromSession(cookie string) (index int) {
	fmt.Println(cookie)
	statement, _ := DataBase.Prepare("SELECT userId FROM session WHERE key=?")
	err := statement.QueryRow(cookie).Scan(&index)
	if err != nil {
		fmt.Println("viga siin!")
		errorCheck(err, true)
	}

	return index
}
func getUserLoopValue(userId int) int {
	for i := 0; i < len(Web.User_data); i++ {
		if Web.User_data[i].ID == userId {
			return i
		}
	}
	return 0
}

func hasCookie(r *http.Request) bool {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		Web.LoggedUser = Memberlist{}
		return false
	}

	userid := getUserIDFromSession(cookie.Value)
	userLid := getUserLoopValue(userid)

	Web.LoggedUser = Memberlist{
		ID:       userid,
		Username: Web.User_data[userLid].Username,
		Email:    Web.User_data[userLid].Email,
	}
	fmt.Println("session:", cookie.Value, "userid:", userid, userLid, Web.User_data[userLid].Username)

	return true
}

// func GetSessionId(username string) (session string, userid int) {
// 	if username == "" {
// 		fmt.Println("GetSessionId:", username, "cannot be empty!")
// 		os.Exit(0)
// 	}

// 	for i := 0; i < len(Web.User_data); i++ {
// 		if Web.User_data[i].Username == username {
// 			userid = i
// 		}
// 	}

// 	statement, _ := DataBase.Prepare("SELECT key FROM session WHERE userId=?")
// 	err := statement.QueryRow(Web.User_data[userid].ID).Scan(&session)

// 	errorCheck(err, true)

// 	return session, userid
// }

// func GetSessionKeyMatch(r *http.Request) bool {
// 	cookie, _ := r.Cookie("session-id")

// 	isTrue, _ := DataBase.Prepare("SELECT key FROM session WHERE key=?")
// 	err := isTrue.QueryRow(cookie.Value).Scan(&Web.LoggedUser.Session)
// 	errorCheck(err, true)
// 	return true
// }
