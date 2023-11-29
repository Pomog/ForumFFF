package handler

import (
	"log"
	"net/http"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/forms"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// TemplateData holds data sent from handlers to templates

// Repo the repository used by the handlers
var Repo *Repository

// Repositroy is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// MainHandler is a method of the Repository struct that handles requests to the main page.
// It renders the "home.page.html" template to the provided HTTP response writer.
func (m *Repository) MainHandler(w http.ResponseWriter, r *http.Request) {

	// Parse the raw request body into r.Form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	// Create a User struct with data from the HTTP request form
	loginData := models.User{
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
		data["loginData"] = loginData
		renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

}

// AboutHandler is a method of the Repository struct that handles requests to the about page.
// It renders the "about.page.html" template to the provided HTTP response writer.
func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	//perform some logic

	// stringData := models.TemplateData{
	// 	StringMap: map[string]string{"test": "this is test data!"},
	// }

	//send data to the template
	renderer.RendererTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: map[string]string{"test": "this is test data!"},
	})
}

// MainHandler is a method of the Repository struct that handles requests to the main page.
// It renders the "home.page.html" template to the provided HTTP response writer.
func (m *Repository) ThemeHandler(w http.ResponseWriter, r *http.Request) {
	renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{})
}
