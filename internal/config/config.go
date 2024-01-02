package config

import (
	"html/template"
	"log"

	"github.com/google/uuid"
)

// AppConfig holds the app config
type AppConfig struct {
	UseCache           bool
	TemplateCache      map[string]*template.Template
	InfoLog            *log.Logger
	InProduction       bool
	ErrorLog           *log.Logger
	UserLogin          uuid.UUID // This field can hold a UUID value
	ServerEmail        string
	PostLen            int //this parameter limits post and topic size
	CategoryLen        int //this parameter limits category size
	FileSize           int64
	GamesList          (map[string]string)
	LongestSingleWord  string
	NameMinLen         int
	NameMaxLen         int
	PasswordMinLen     int
	PasswordMaxLen     int
	MinSubjectLen      int
	MaxSubjectLen      int
	MinCategoryLen     int
	MaxCategoryLen     int
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}
