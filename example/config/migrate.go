package config

import (
	"database/sql"
	"log"
)

func Run() {
	migrate(DB, UserTable)
	migrate(DB, PostTable)
	migrate(DB, CategoryTable)
	migrate(DB, PostCategoryTable)
	migrate(DB, PostRatingTable)
	migrate(DB, PostRepliesTable)
	migrate(DB, PostRepliesRatingTable)
	migrate(DB, SessionTable)
}

func migrate(db *sql.DB, query string) {
	statement, err := db.Prepare(query)
	if err == nil {
		_, creationError := statement.Exec()
		if creationError != nil {
			log.Println(creationError.Error())
		}
	} else {
		log.Println(err.Error())
	}
}
