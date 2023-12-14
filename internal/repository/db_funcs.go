package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// MakeDBTables creates DB tables and Guest User
func MakeDBTables(db *sql.DB) error {
	sqlQuerys := getQuerys()

	for _, sqlQuery := range sqlQuerys {
		statement, errPrepare := db.Prepare(sqlQuery)
		if errPrepare != nil {
			fmt.Println("errPrepare:", errPrepare)
			continue
		}

		_, err := statement.Exec()
		statement.Close()

		if err != nil {
			return fmt.Errorf("error executing query %s: %v", sqlQuery, err)
		}
	}

	userExist, _ := userExists(db, "guest")
	if !userExist {
		if _, err := db.Exec(guestUser); err != nil {
			return fmt.Errorf("error executing guestUser query: %v", err)
		}
	}

	return nil
}

// GetDB open or creates sqlite3 DataBase
func GetDB() (*DataBase, error) {
	database, err := sql.Open("sqlite3", "./mainDB.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	errDB := testDB(database)
	if errDB != nil {
		fmt.Printf("error Ping database: %v\n", errDB) // Change "err" to "errDB" here
		return nil, errDB
	}

	return &DataBase{SQL: database}, nil
}

// testDB tested DB connection
func testDB(db *sql.DB) error {
	// Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	err := db.Ping()
	return err
}

// userExists helper function to create Guest User
func userExists(db *sql.DB, username string) (bool, error) {
	// Check if the user with the given username exists.
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	var count int
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
