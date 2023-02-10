package main

import (
	"database/sql"
	"errors"
	"fmt"
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

	return user, err
}

func Register(username, password, email string) error {
	insertStatement, _ := DataBase.Prepare("INSERT INTO users (username, password, email, datecreated) values (?,?,?,?)")

	currentTime := time.Now().Format("02.01.2006")
	_, err := insertStatement.Exec(username, password, email, currentTime)
	if err != nil {
		switch err.Error() {
		case "UNIQUE constraint failed: users.username":
			return errors.New("Username " + username + " is taken")
		case "UNIQUE constraint failed: users.email":
			return errors.New("Email " + email + " is taken")
		default:
			return err
		}
	}

	return nil
}

func getUserFromSession(r *http.Request) (Memberlist, bool) {
	cookie, _ := r.Cookie("session-id")
	if cookie == nil {
		return Memberlist{}, false
	}
	var user Memberlist
	statement, _ := DataBase.Prepare(`SELECT userId, username, email, datecreated
	FROM session LEFT JOIN users
	ON session.userId = users.id 
	WHERE key = ?
	`)
	err := statement.QueryRow(cookie.Value).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DateCreated,
	)
	if err != nil {
		return Memberlist{}, false
	}

	return user, true
}

func hasCookie(r *http.Request) bool {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		return false
	}
	var count int
	err = DataBase.QueryRow("SELECT COUNT(*) FROM session WHERE key = ?", cookie.Value).Scan(&count)

	if err != nil {
		fmt.Println(err)
	}

	return count == 1
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
