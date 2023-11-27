package config

import (
	"html/template"
	"log"
)

// AppConfig holds the app config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool
	ErrorLog      *log.Logger
}
