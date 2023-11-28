package dbrepo

import (
	"database/sql"

	"github.com/Pomog/ForumFFF/internal/config"
)

type sqliteBDRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewSQLiteRepo (conn *sql.DB, a *config.AppConfig) *sqliteBDRepo {
	return &sqliteBDRepo{
		App: a,
		DB: conn,
	}
}
