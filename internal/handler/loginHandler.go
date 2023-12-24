package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Pomog/ForumFFF/internal/forms"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/google/uuid"
)

// LoginHandler handles both GET and POST requests for the login page.
func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var emptyLogin models.User
		data := make(map[string]interface{})
		data["loginData"] = emptyLogin
		renderer.RendererTemplate(w, "login.page.html", &models.TemplateData{
			Form: forms.NewForm(nil),
			Data: data,
		})
	} else if r.Method == http.MethodPost {
		// Parse the raw request body into r.Form
		err := r.ParseForm()
		if err != nil {
			setErrorAndRedirect(w, r, err.Error(), "/error-page")
			return
		}
		// Create a User struct with data from the HTTP request form
		loginData := models.User{
			Email:    strings.ToLower(r.FormValue("emailLogIn")),
			Password: r.FormValue("passwordLogIn"),
		}

		// Create a new form instance based on the HTTP request's PostForm
		form := forms.NewForm(r.PostForm)

		// Validation checks for required fields and their specific formats and lengths
		form.Required("emailLogIn", "passwordLogIn")

		// Check if the form data is valid; if not, render the home page with error messages
		if !form.Valid() {
			data := make(map[string]interface{})
			data["loginData"] = loginData
			renderer.RendererTemplate(w, "login.page.html", &models.TemplateData{
				Form: form,
				Data: data,
			})
			return
		}

		// Check if User is Present in the DB, ERR should be handled
		userID, _ := m.DB.UserPresentLogin(loginData.Email, loginData.Password)
		if userID != 0 {
			m.App.UserLogin = uuid.New()
			err := m.DB.InsertSessionintoDB(m.App.UserLogin.String(), userID)
			if err != nil {
				setErrorAndRedirect(w, r, err.Error(), "/error-page")
				return
			}

			cookie := &http.Cookie{
				Name:  strconv.Itoa(userID),
				Value: m.App.UserLogin.String(),
			}
			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			setErrorAndRedirect(w, r, "Wrong email or password", "/error-page")
			return
		}
	}

}
