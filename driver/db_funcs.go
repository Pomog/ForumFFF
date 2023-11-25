package driver

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func MakeDB() {
	database, err := sql.Open("sqlite3", "./mainDB.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer database.Close()

	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS users 
	(id INTEGER PRIMARY KEY, 
		username varchar(100) DEFAULT "", 
		password varchar(100) DEFAULT "",
		first_name varchar(100) DEFAULT "",
		last_name varchar(100) DEFAULT "",
		email varchar(254) DEFAULT "",
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		picture  TEXT DEFAULT "static/ava/pomog_ava.png", 
		last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)

	defer statement.Close()

	_, _ = statement.Exec()

}
