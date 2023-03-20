package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"forum/app"
)

func main() {
	file, err := os.Stat("database.db")
	port := "8080" // webserver port
	// if .db file deleted, it will create new one and populate with data
	if errors.Is(err, os.ErrNotExist) {
		app.DataBase, _ = sql.Open("sqlite3", "database.db")
		app.InitDatabase()
		fmt.Println("New database created ", file)
	} else {
		fmt.Println("masiin")
		app.DataBase, _ = sql.Open("sqlite3", "./database.db")
	}

	app.Server(port)
}
