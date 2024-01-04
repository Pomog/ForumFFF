package handler

import (
	"net/http"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

func (m *Repository) ModerPanelHandler(w http.ResponseWriter, r *http.Request) {
	sessionUserID := m.GetLoggedUser(w, r)
	if sessionUserID == 0 {
		setErrorAndRedirect(w, r, "unautorized", "/error-page")
		return
	}
	user, err := m.DB.GetUserByID(sessionUserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get user from: GetUserByID(sessionUserID) "+err.Error(), "/error-page")
		return
	}
	if user.Type != "moder" {
		setErrorAndRedirect(w, r, "Unauthorized access, only Moderator can visit this page -"+err.Error(), "/error-page")
		return
	}

	if r.Method == http.MethodGet {
		handleGetRequestModerPage(w, r, m, sessionUserID)
	}
}

// handleGetRequest handles GET requests for the home page.
func handleGetRequestModerPage(w http.ResponseWriter, r *http.Request, m *Repository, sessionUserID int) {
	// topics := r.URL.Query().Get("topics")
	// posts := r.URL.Query().Get("posts")
	data := make(map[string]interface{})
	loggedUser, err := m.DB.GetUserByID(sessionUserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get logged user from: GetUserByID(sessionUserID) -"+err.Error(), "/error-page")
		return
	}

	data["loggedAs"] = loggedUser.UserName
	data["loggedAsID"] = loggedUser.ID
	data["loggedUserType"] = loggedUser.Type

	renderer.RendererTemplate(w, "moderMain.page.html", &models.TemplateData{
		Data: data,
	})
}
