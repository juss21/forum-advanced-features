package main

import (
	"database/sql"
	"net/http"
	"time"
)

func checkIfPreviouslyLoggedin(username string) bool {
	// see on funktsioon mis võiks deleteda eelmise sessiooni, kui eelmine eksisteerib
	userid := getUserLoopValueSTR(username)
	if userid == -1 {
		return false
	}
	key := ""
	selector, _ := DataBase.Prepare("SELECT key FROM session WHERE userId=?")
	err := selector.QueryRow(Web.User_data[userid].ID).Scan(&key)
	if err != nil {
		//	fmt.Println("eelmist sessiooni pole!")
		return false
	}
	// delete eelmine sessioon
	statement, _ := DataBase.Prepare("DELETE FROM session WHERE key = ?")
	statement.Exec(key)
	return true
}

func Login(username, password string) (Memberlist, error) {
	checkIfPreviouslyLoggedin(username)
	var user Memberlist
	statement, _ := DataBase.Prepare("SELECT id, username, password, email FROM users WHERE username=?")
	err := statement.QueryRow(username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
	)

	if err == sql.ErrNoRows {
		return Memberlist{}, err
	}

	return user, err
}

func Register(username, password, email string) {
	statement, _ := DataBase.Prepare("INSERT INTO users (username, password, email, datecreated) values (?,?,?,?)")
	currentTime := time.Now().Format("02.01 2006")
	statement.Exec(username, password, email, currentTime)
	Web.User_data = append(Web.User_data, Memberlist{ID: Web.User_data[len(Web.User_data)-1].ID + 1, Username: username, Email: email, DateCreated: currentTime})
}

func getUserIDFromSession(cookie string) (index int) {
	index = -1
	statement, _ := DataBase.Prepare("SELECT userId FROM session WHERE key=?")
	statement.QueryRow(cookie).Scan(&index)

	return index
}
func getUserLoopValue(userId int) int {
	for i := 0; i < len(Web.User_data); i++ {
		if Web.User_data[i].ID == userId {
			return i
		}
	}
	return -1
}
func getUserLoopValueSTR(username string) int {
	for i := 0; i < len(Web.User_data); i++ {
		if Web.User_data[i].Username == username {
			return i
		}
	}
	return -1
}
func hasCookie(r *http.Request) bool {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		Web.LoggedUser = Memberlist{}
		return false
	}

	userid := getUserIDFromSession(cookie.Value)
	if userid == -1 {
		Web.LoggedUser = Memberlist{}
		return false
	}
	userLid := getUserLoopValue(userid)

	Web.LoggedUser = Memberlist{
		ID:       userid,
		Username: Web.User_data[userLid].Username,
		Email:    Web.User_data[userLid].Email,
	}
	//fmt.Println("session:", cookie.Value, "userid:", userid, userLid, Web.User_data[userLid].Username)

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