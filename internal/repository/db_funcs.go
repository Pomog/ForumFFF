package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var dbConnSQL = &DataBase{}

func MakeDBTables() (*sql.DB, error) {
	dbConn, err := GetDB()
	if err != nil {
		log.Println(err)
		return dbConn.SQL, err
	}

	sqlQuerys := getQuerys()

	for _, sqlQury := range sqlQuerys {
		statement, errPrepare := dbConn.SQL.Prepare(sqlQury)
		if errPrepare != nil {
			fmt.Println("errPrepare:", errPrepare)
		}

		defer statement.Close()

		_, err = statement.Exec()
		if err != nil {
			return dbConn.SQL, err
		}
	}

	return dbConn.SQL, nil
}

func GetDB() (*DataBase, error) {
	database, err := sql.Open("sqlite3", "./mainDB.db")
	if err != nil {
		return nil, err
	}
	errDB := testDB(database)
	if errDB != nil {
		return nil, errDB
	}

	dbConnSQL.SQL = database

	return dbConnSQL, nil
}

func testDB(db *sql.DB) error {
	//Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	err := db.Ping()
	return err
}
