package helper

import (
	"net/http"
	"runtime/debug"

	"github.com/Pomog/ForumFFF/internal/config"
)

var app *config.AppConfig

// Set the config to the helper
func NewHelper(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client Error ", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	app.ErrorLog.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
