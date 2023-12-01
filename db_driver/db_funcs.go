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
		return database, err
	}

	sqlQuerys := getQuerys()

	for _, sqlQury := range sqlQuerys {
		statement, errPrepare := database.Prepare(sqlQury)
		if errPrepare != nil {
			fmt.Println("errPrepare:", errPrepare)
		}

		defer statement.Close()

		_, err = statement.Exec()
		if err != nil {
			return database, err
		}
	}

	return database, nil
}

func GetDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./mainDB.db")
	if err != nil {
		return nil, err
	}
	errDB := testDB(database)
	if errDB != nil {
		return nil, errDB
	}
	return database, nil
}

func testDB(db *sql.DB) error {
	//Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	err := db.Ping()
	return err
}
