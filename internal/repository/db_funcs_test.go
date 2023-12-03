package repository

import (
	"os"
	"testing"
)

var ObligatoryTables = []string{
	"users", "thread", "votes", "post",
}

func Test_GetDB(t *testing.T) {

	db, err := MakeDBTables()

	if err != nil {
		t.Errorf("Could not %s", "sql.Open(sqlite3, ./mainDB.db)")
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()

	var tableName string
	var allTables []string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			t.Error(err)
		}
		allTables = append(allTables, tableName)
	}

	if err := rows.Err(); err != nil {
		t.Error(err)
	}

	if !compareTableNames(ObligatoryTables, allTables) {
		t.Errorf("\nExpected %v\nreceived %v", (ObligatoryTables), (allTables))
	}
	filepath := "./mainDB.db"
	if _, err := os.Stat(filepath); err == nil {
		// Delete the file
		err := os.Remove(filepath)
		if err != nil {
			t.Error("Error:", err)
		}
	}

}

func compareTableNames(want, get []string) bool {
	if len(want) != len(get) {
		return false
	}
	for _, elem := range want {
		if !containsString(get, elem) {
			return false
		}
	}
	return true
}

// containsString checks if a string is present in a slice of strings
func containsString(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}