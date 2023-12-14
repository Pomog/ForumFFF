// Package helper provides common helper functions for handling HTTP requests and errors.
package helper

import (
	"net/http"
	"runtime/debug"

	"github.com/Pomog/ForumFFF/internal/config"
)

// app is a pointer to the AppConfig, which holds the application configuration.
var app *config.AppConfig

// NewHelper Set the config to the helper
func NewHelper(a *config.AppConfig) {
	app = a
}

// ClientError logs a client error and responds with the specified HTTP status code.
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client Error ", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError logs a server error, including the error message and stack trace,
// and responds with a 500 Internal Server Error status code
func ServerError(w http.ResponseWriter, err error) {
	app.ErrorLog.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
