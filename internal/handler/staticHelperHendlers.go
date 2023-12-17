package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/Pomog/ForumFFF/internal/repository"
	"github.com/Pomog/ForumFFF/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository handles the repository type, encapsulating the AppConfig and DatabaseInt dependencies.
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseInt
}

const (
	dbErrorUserPresent    = "DB Error func UserPresent"
	userAlreadyExistsMsg  = "User Already Exists"
	dbErrorCreateUser     = "DB Error func CreateUser"
	fileReceivingErrorMsg = "file receiving error"
	fileCreatingErrorMsg  = "Unable to create file"
	fileSavingErrorMsg    = "Unable to save file"
	guestRestiction       = "Guests can not create Themes and Posts, please log in or register!"

	emptyUUID = "00000000-0000-0000-0000-000000000000"
)

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
			<p>An error occurred: <strong>` + errorMessage + `</strong></p>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(htmlContent))
	if err != nil {
		setErrorAndRedirect(w, r, err.Error(), "/error-page")
	}
}

// setErrorContext sets the error message in the context and adds it to the redirect URL
func setErrorAndRedirect(w http.ResponseWriter, r *http.Request, errorMessage string, redirectURL string) {
	// Append the error message as a query parameter in the redirect URL
	redirectURL += "?error=" + url.QueryEscape(errorMessage)

	// Perform the redirect
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (m *Repository) ContactUsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "contactUs.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}

func (m *Repository) ForumRulesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "forumRules.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

func (m *Repository) HelpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "help.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

func (m *Repository) PrivatPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "privatP.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

func (m *Repository) PersonaCabinetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID, _ := strconv.Atoi(r.URL.Query().Get("userID"))
		fmt.Println("UserID", userID)
		var personalInfo models.User
		user, errUser := m.DB.GetUserByID(userID)
		if errUser != nil {
			setErrorAndRedirect(w, r, "Could not get User from  GetUserByID(visitorID)", "/error-page")
		}

		// visitorID, _ := m.DB.GetGuestID()

		// for _, cookie := range r.Cookies() {
		// 	if cookie.Value == m.App.UserLogin.String() {
		// 		userID, err := strconv.Atoi(cookie.Name)
		// 		if err != nil {
		// 			setErrorAndRedirect(w, r, "Could not get visitor ID", "/error-page")
		// 		}
		// 		if visitorID = userID; visitorID != 0 {
		// 			break
		// 		}

		// 	}
		// }
		// user, errUser := m.DB.GetUserByID(visitorID)
		// if errUser != nil {
		// 	setErrorAndRedirect(w, r, "Could not get User from  GetUserByID(visitorID)", "/error-page")
		// }
		// var personalInfo models.User
		personalInfo.Email = user.Email
		personalInfo.ID = user.ID
		personalInfo.Created = user.Created
		personalInfo.FirstName = user.FirstName
		personalInfo.LastName = user.LastName
		personalInfo.Picture = user.Picture
		personalInfo.UserName = user.UserName
		totalPosts, _ := m.DB.GetTotalPostsAmmountByUserID(personalInfo.ID)
		data := make(map[string]interface{})
		data["personal"] = personalInfo
		data["totalPosts"] = totalPosts

		renderer.RendererTemplate(w, "personal.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}
