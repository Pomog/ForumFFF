package db_driver

import (
	"database/sql"
	"net/http"

	"github.com/Pomog/ForumFFF/internal/helper"
	_ "github.com/mattn/go-sqlite3"
)

func MakeDBTables() error {
	database, err := getDB()
	if err != nil {
		helper.ServerError(w, err)
		return err
	}
	defer database.Close()

	for _, sqlQury := range sqlQurys {

		statement, _ := database.Prepare(sqlQury)

		defer statement.Close()

		_, err = statement.Exec()
		if err != nil {
			return err
		}
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
