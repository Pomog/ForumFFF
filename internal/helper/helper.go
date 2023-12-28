// Package helper provides common helper functions for handling HTTP requests and errors.
package helper

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

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

var punctuationsSlice = []string{".", ",", "!", "?", ":", ";"}       // TODO: add configurable punctuation, by config file or env var

/*
CorrectPunctuationsSpaces takes an input string and adds spaces around punctuation marks
based on the specified rules duscribed by inner function applyPunctuationRules.
Punctuation marks are placed close to the previous word
and separated by a space from the following content. Exceptions include groups of
punctuation like '...' or '!?'.
*/
func CorrectPunctuationsSpaces(input string) string {
	punctuationPattern := strings.Join(punctuationsSlice, "|\\")
	strPattern := fmt.Sprintf(`(\b[\w]+)\s*([%s]+)\s*`, punctuationPattern) // shoud be two groups, where the first element is the word, and the second is the punctuation mark

	re := regexp.MustCompile(strPattern)

	return correctString(re, input)
}

/*
applies the punctuation spacing rules to the input string.
It uses a regular expression to match and replace substrings based on the rules defined in applyPunctuationRules.
*/
func correctString(re *regexp.Regexp, input string) string {
	return re.ReplaceAllStringFunc(input, func(match string) string {
		return applyPunctuationRules(re, match, input)
	})
}

/*
applies the punctuation spacing rules to the matched substring.
It takes the matched substring, extracts the word and punctuation mark preseeded by it,
and places them together with or without following space based on whether it's the end of the string.
*/
func applyPunctuationRules(re *regexp.Regexp, match string, input string) string {
	submatches := re.FindStringSubmatch(match)
	if strings.HasSuffix(input, match) {
		return fmt.Sprintf("%s%s", submatches[1], submatches[2]) // end of the string
	}
	return fmt.Sprintf("%s%s ", submatches[1], submatches[2]) // not the end of the string
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
