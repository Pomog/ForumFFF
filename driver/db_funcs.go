package driver

import (
	"database/sql"
	"net/http"

	"github.com/Pomog/ForumFFF/internal/helper"
	_ "github.com/mattn/go-sqlite3"
)

func makeUserTable() error {
	database, err := getDB()
	if err != nil {
		helper.ServerError(w, err)
		return err
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

	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func makeThreadTable() error {
	database, err := getDB()
	if err != nil {
		helper.ServerError(w, err)
		return err
	}
	defer database.Close()

	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS thread 
	(id INTEGER PRIMARY KEY, 
		subject TEXT,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		userID INTEGER)`)

	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func makePostTable() error {
	database, err := getDB()
	if err != nil {
		helper.ServerError(w, err)
		return err
	}
	defer database.Close()

	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS post 
	(id INTEGER PRIMARY KEY, 
		subject TEXT,
		content TEXT,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		threadID INTEGER,
		userID INTEGER)`)

	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func makeVotesTable() error {
	database, err := getDB()
	if err != nil {
		helper.ServerError(w, err)
		return err
	}
	defer database.Close()

	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS post 
	(id INTEGER PRIMARY KEY, 
		upCount INTEGER,
		downCount INTEGER,
		postID INTEGER)`)

	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func getDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./mainDB.db")
	if err != nil {
		return nil, err
	}
	return database, nil
}

type myWriter struct{}

var w *myWriter

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
