package handler

import (
	"fmt"
	"net/http"

	"github.com/Pomog/ForumFFF/pkg/config"
	"github.com/Pomog/ForumFFF/pkg/models"
	"github.com/Pomog/ForumFFF/pkg/renderer"
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
	emailLog := r.FormValue("emailLogIn")
	passLog := r.FormValue("passwordLogIn")

	nickname := r.FormValue("nickName")
	emailReg := r.FormValue("emailRegistr")
	passwordReg := r.FormValue("passwordReg")
	passwordRep := r.FormValue("passwordRep")

	if emailLog != "" {
		fmt.Println("log:", emailLog)
		fmt.Println("pass:", passLog)
	} else if nickname != "" {
		fmt.Println("nickname:", nickname)
		fmt.Println("emailReg:", emailReg)
		fmt.Println("passwordReg:", passwordReg)
		fmt.Println("passwordRep:", passwordRep)
	}

	renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{})
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
