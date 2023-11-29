package db_driver

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func MakeDBTables() error {
	database, err := GetDB()
	if err != nil {
		log.Println(err)
		return err
	}
	defer database.Close()

	sqlQuerys := getQuerys()

	for _, sqlQury := range sqlQuerys {
		statement, errPrepare := database.Prepare(sqlQury)
		if errPrepare != nil {
			fmt.Println("errPrepare:", errPrepare)
		}

		defer statement.Close()

		_, err = statement.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./mainDB.db")
	if err != nil {
		return nil, err
	}
	return database, nil
}
