package main

import (
	"database/sql"
	"net/http"
	"time"
)

// func checkIfPreviouslyLoggedin(username string) bool {
// 	// see on funktsioon mis võiks deleteda eelmise sessiooni, kui eelmine eksisteerib
// 	user := getUserFromSession(user.ID)
// 	if userid == -1 {
// 		return false
// 	}
// 	key := ""
// 	selector, _ := DataBase.Prepare("SELECT key FROM session WHERE userId=?")
// 	err := selector.QueryRow(Web.User_data[userid].ID).Scan(&key)
// 	if err != nil {
// 		//	fmt.Println("eelmist sessiooni pole!")
// 		return false
// 	}
// 	// delete eelmine sessioon

// 	statement, _ := DataBase.Prepare("DELETE FROM session WHERE key = ?")
// 	statement.Exec(key)
// 	return true
// }

func Login(username string, password string) (Memberlist, error) {
	loginStatement, _ := DataBase.Prepare("SELECT id, username, password, email FROM users WHERE username=?")
	sessionStatement, _ := DataBase.Prepare("SELECT key FROM session WHERE userId=?")

	var user Memberlist
	err := loginStatement.QueryRow(username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
	)

	if err == sql.ErrNoRows {
		return Memberlist{}, err
	}

	var key string
	err = sessionStatement.QueryRow(user.ID).Scan(&key)

	if err == nil {
		statement, _ := DataBase.Prepare("DELETE FROM session WHERE key = ?")
		statement.Exec(key)
	}

	return user, err
}

func Register(username, password, email string) {
	statement, _ := DataBase.Prepare("INSERT INTO users (username, password, email, datecreated) values (?,?,?,?)")
	currentTime := time.Now().Format("02.01.2006")
	statement.Exec(username, password, email, currentTime)
	Web.User_data = append(Web.User_data, Memberlist{ID: Web.User_data[len(Web.User_data)-1].ID + 1, Username: username, Email: email, DateCreated: currentTime})
}

func getUserFromSession(cookie string) Memberlist {
	var user Memberlist
	statement, _ := DataBase.Prepare(`SELECT userId, username, email, datecreated
	FROM session LEFT JOIN users
	ON session.userId = users.id 
	WHERE key = ?
	`)
	statement.QueryRow(cookie).Scan(&user)

	return user
}

func hasCookie(r *http.Request) bool {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		Web.LoggedUser = Memberlist{}
		return false
	}
	user := getUserFromSession(cookie.Value)

	Web.LoggedUser = user

	return true
}

func ClearCookies(w http.ResponseWriter, r *http.Request) {
	if !Web.Loggedin {
		http.SetCookie(w, &http.Cookie{
			Name:   "session-id",
			Value:  "",
			MaxAge: -1,
		})
	}
}
