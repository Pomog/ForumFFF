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

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	loginData := models.User{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		UserName:  r.FormValue("nickName"),
		Email:     r.FormValue("emailRegistr"),
		Password:  r.FormValue("passwordReg"),
	}

	form := forms.NewForm(r.PostForm)

	form.Required("firstName", "lastName", "nickName", "emailRegistr", "passwordReg")

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
