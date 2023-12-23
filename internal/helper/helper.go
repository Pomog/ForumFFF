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

// // sendEmail to emailAdress with starting server time
// func SendEmail(emailAdress, message string) {
// 	// test mail
// 	from := "ffforumadm@gmail.com"
// 	password := ""
// 	to := emailAdress
// 	subject := "Test Email"
// 	time := time.Now().Format("2006-01-02 15:04:05")
// 	body := fmt.Sprintf("%s at : %s", message, time)

// 	msg := "To: " + to + "\r\n" +
// 		"Subject: " + subject + "Test" + "\r\n" +
// 		"\r\n" + body

// 	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

// 	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
// 	if err != nil {
// 		app.ErrorLog.Println(err)
// 	} else {
// 		fmt.Println("Email sent successfully.")
// 	}
// }
