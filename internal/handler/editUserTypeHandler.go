package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// ChangeUserTypeResultHandler handles changing of user type
func (m *Repository) ChangeUserTypeResultHandler(w http.ResponseWriter, r *http.Request) {
	sessionUserID := m.GetLoggedUser(w, r)
	if sessionUserID == 0 {
		setErrorAndRedirect(w, r, "unautorized", "/error-page")
		return
	}
	personalCabinetUserID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil {
		setErrorAndRedirect(w, r, "could not convert string to int: strconv.Atoi(r.URL.Query().Get(userID))", "/error-page")
		return
	}
	if sessionUserID != personalCabinetUserID {
		setErrorAndRedirect(w, r, "only owner of cabinet can submit secret code", "/error-page")
		return
	}
	inputPass := strings.TrimSpace(r.FormValue("changeUserType"))

	if r.Method == http.MethodPost && inputPass == m.App.ModeratorPass {
		user, _ := m.DB.GetUserByID(sessionUserID)
		user.Type = "moder"
		m.DB.EditUserType(user)

		data := make(map[string]interface{})
		data["userID"] = user.ID
		renderer.RendererTemplate(w, "edit_user_type_result.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "wrong secret code", http.StatusMethodNotAllowed)
	}

}
