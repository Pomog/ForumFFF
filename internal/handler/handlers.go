package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/forms"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/Pomog/ForumFFF/internal/repository"
	"github.com/Pomog/ForumFFF/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repositroy is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseInt
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *repository.DataBase) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewSQLiteRepo(a, db.SQL),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var emtpyLogin models.User
		data := make(map[string]interface{})
		data["loginData"] = emtpyLogin
		renderer.RendererTemplate(w, "login.page.html", &models.TemplateData{
			Form: forms.NewForm(nil),
			Data: data,
		})
	} else if r.Method == http.MethodPost {
		// Parse the raw request body into r.Form
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		// Create a User struct with data from the HTTP request form
		loginData := models.User{
			Email:    r.FormValue("emailLogIn"),
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

		// Check if User is Presaent in the DB, ERR should be handled
		result, _ := m.DB.UserPresent(loginData.UserName, loginData.Email)
		fmt.Println("UserPresent: ", result)
		if result {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}

		// if there is no error, we upload Form data into our Session
		//WHAT to use here?
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (m *Repository) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var emptyRegistration models.User
		data := make(map[string]interface{})
		data["registrationData"] = emptyRegistration
		renderer.RendererTemplate(w, "register.page.html", &models.TemplateData{
			Form: forms.NewForm(nil),
			Data: data,
		})

	} else if r.Method == http.MethodPost {
		// Parse the raw request body into r.Form
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		// Create a User struct with data from the HTTP request form
		registrationData := models.User{
			FirstName: r.FormValue("firstName"),
			LastName:  r.FormValue("lastName"),
			UserName:  r.FormValue("nickName"),
			Email:     r.FormValue("emailRegistr"),
			Password:  r.FormValue("passwordReg"),
		}

		// Create a new form instance based on the HTTP request's PostForm
		form := forms.NewForm(r.PostForm)

		// Validation checks for required fields and their specific formats and lengths
		form.Required("firstName", "lastName", "nickName", "emailRegistr", "passwordReg")
		form.First_LastName_Min_Max_Len("firstName", 3, 12, r)
		form.First_LastName_Min_Max_Len("lastName", 3, 12, r)
		form.First_LastName_Min_Max_Len("nickName", 3, 12, r)
		form.EmailFormat("emailRegistr", r)
		form.PassFormat("passwordReg", 6, 15, r)

		// Check if the form data is valid; if not, render the home page with error messages
		if !form.Valid() {
			data := make(map[string]interface{})
			data["registrationData"] = registrationData
			renderer.RendererTemplate(w, "register.page.html", &models.TemplateData{
				Form: form,
				Data: data,
			})
			return
		}

		// Check if User is Presaent in the DB, ERR should be handled
		result, _ := m.DB.UserPresent(registrationData.UserName, registrationData.Email)
		fmt.Println("UserPresent: ", result)
		if result {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}

		// if there is no error, we upload Form data into our Session
		//WHAT to use here?
	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}

// MainHandler is a method of the Repository struct that handles requests to the main page.
// It renders the "home.page.html" template to the provided HTTP response writer.
func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {

	renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{})

}

// MainHandler is a method of the Repository struct that handles requests to the main page.
// It renders the "home.page.html" template to the provided HTTP response writer.
func (m *Repository) ThemeHandler(w http.ResponseWriter, r *http.Request) {
	renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{})
}
