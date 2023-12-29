package handler

import (
	"net/http"
	"net/url"
	"strings"
)

// ErrorPage handles the "/error-page" route
func (m *Repository) ErrorPage(w http.ResponseWriter, r *http.Request) {
	// Retrieve the error value from the query parameter
	errorMessage := r.URL.Query().Get("error")

	if errorMessage == "" {
		// If the error value is not present, handle it accordingly
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Error Page</title>
		</head>
		<body>
			<h1>Error</h1>
			<p>An error occurred:</p>
			<p><strong>` + strings.ReplaceAll(errorMessage, "\n", "<br>") + `</strong></p>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(htmlContent))
	if err != nil {
		setErrorAndRedirect(w, r, err.Error(), "/error-page")
		return
	}
}

// setErrorAndRedirect sets the error message in the context and adds it to the redirect URL
func setErrorAndRedirect(w http.ResponseWriter, r *http.Request, errorMessage string, redirectURL string) {
	// Append the error message as a query parameter in the redirect URL
	redirectURL += "?error=" + url.QueryEscape(errorMessage)

	// Perform the redirect and return immediately
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
