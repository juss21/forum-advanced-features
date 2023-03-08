package app

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

func ClearCookies(w http.ResponseWriter, r *http.Request) {
	if !Web.Loggedin {
		http.SetCookie(w, &http.Cookie{
			Name:   "session-id",
			Value:  "",
			MaxAge: -1,
		})
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
