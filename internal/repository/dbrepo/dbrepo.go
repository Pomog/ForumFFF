package dbrepo

import (
	"database/sql"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/repository"
)

type SqliteBDRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseInt
}

func NewSQLiteRepo(a *config.AppConfig, conn *sql.DB) *SqliteBDRepo {
	return &SqliteBDRepo{
		App: a,
		DB:  conn,
	}
}
