package handler

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/forms"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/Pomog/ForumFFF/internal/repository"
	"github.com/Pomog/ForumFFF/internal/repository/dbrepo"
	"github.com/google/uuid"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repositroy is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseInt
}

const (
	dbErrorUserPresent   = "DB Error func UserPresent"
	userAlreadyExistsMsg = "User Already Exists"
	dbErrorCreateUser    = "DB Error func CreateUser"
	fileRecivingErrorMsg = "file receiving error"
	fileCreatingErrorMsg = "Unable to create file"
	fileSaveingErrorMsg  = "Unable to save file"
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
		userID, _ := m.DB.UserPresentLogin(loginData.Email, loginData.Password)
		if userID != 0 {
			m.App.UserLogin = uuid.New()
			m.DB.InsertSessionintoDB(m.App.UserLogin.String(), userID)

			cookie := &http.Cookie{
				Name:  strconv.Itoa(userID),
				Value: m.App.UserLogin.String(),
			}
			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			setErrorAndRedirect(w, r, "Wrong email or password", "/error-page")
		}
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
		// Parse the form data, including files Need to Set Upper limit for DATA
		err := r.ParseMultipartForm((1 << 20))
		if err != nil {
			setErrorAndRedirect(w, r, dbErrorUserPresent, "/error-page")
			return
		}

		// Create a User struct with data from the HTTP request form
		registrationData := models.User{
			FirstName: r.FormValue("firstName"),
			LastName:  r.FormValue("lastName"),
			UserName:  r.FormValue("nickName"),
			Email:     r.FormValue("emailRegistr"),
			Password:  r.FormValue("passwordReg"),
			Picture:   r.FormValue("avatar"),
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
		userAlreadyExist, err := m.DB.UserPresent(registrationData.UserName, registrationData.Email)
		if err != nil {
			setErrorAndRedirect(w, r, "DB Error func UserPresent", "/error-page")
			return
		}

		if userAlreadyExist {
			setErrorAndRedirect(w, r, "User AlreadyExists", "/error-page")
		} else {
			// Get the file from the form data
			file, handler, errfileGet := r.FormFile("avatar")
			if errfileGet != nil {
				setErrorAndRedirect(w, r, fileRecivingErrorMsg, "/error-page")
				return
			}
			defer file.Close()

			// Create a new file in the "static/ava" directory
			newFilePath := filepath.Join("static/ava", handler.Filename)
			newFile, errfileCreate := os.Create(newFilePath)
			if errfileCreate != nil {
				setErrorAndRedirect(w, r, fileCreatingErrorMsg, "/error-page")
				return
			}
			defer newFile.Close()

			// Copy the uploaded file to the new file
			_, err = io.Copy(newFile, file)
			if err != nil {
				setErrorAndRedirect(w, r, fileSaveingErrorMsg, "/error-page")
				return
			}

			registrationData.Picture = path.Join("/", newFilePath)

			err := m.DB.CreateUser(registrationData)
			if err != nil {
				setErrorAndRedirect(w, r, "DB Error func CreateUser", "/error-page")
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}

// MainHandler is a method of the Repository struct that handles requests to the main page.
// It renders the "home.page.html" template to the provided HTTP response writer.
func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	var UserID int
	for _, cookie := range r.Cookies() {
		if cookie.Value != "" {
			userID, err := m.DB.GetUserIDForSessionID(cookie.Value)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get UserID from Cookies m.DB.GetUserIDForSessionID", "/error-page")
			}
			UserID = userID
		}
	}

	if UserID == 0 {
		setErrorAndRedirect(w, r, "Could not verify User, Please LogIN", "/error-page")
	}

	if r.Method == http.MethodGet {
		threads, err := m.DB.GetAllThreads()
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get Threads m.DB.GetAllThreads", "/error-page")
		}

		var threadsInfo []models.ThreadDataForMainPage
		for _, thread := range threads {
			var user models.User
			user, _ = m.DB.GetUserByID(thread.UserID)
			var info models.ThreadDataForMainPage
			info.Subject = thread.Subject
			info.Created = thread.Created.Format("2006-01-02 15:04:05")

			info.PictureUserWhoCreatedThread = user.Picture
			info.UserNameWhoCreatedThread = user.UserName

			posts, err := m.DB.GetAllPostsFromThread(thread.ID)
			if err != nil {
				log.Fatal(err)
			}
			info.Posts = posts
			userWhoCreatedLastPost, _ := m.DB.GetUserByID(getUserThatCreatedLastPost(posts))
			info.PictureUserWhoCreatedLastPost = userWhoCreatedLastPost.Picture
			info.UserNameWhoCreatedLastPost = userWhoCreatedLastPost.UserName
			threadsInfo = append(threadsInfo, info)
		}

		data := make(map[string]interface{})

		data["threads"] = threadsInfo

		renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{
			Data: data,
		})
	} else if r.Method == http.MethodPost {
		thread := models.Thread{
			Subject: r.FormValue("message-text"),
			UserID:  UserID,
		}

		err := m.DB.CreateThread(thread)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not create a thread", "/error-page")
		}
		http.Redirect(w, r, "/theme", http.StatusPermanentRedirect)

	}

}

func getUserThatCreatedLastPost(posts []models.Post) int {
	var id int
	timeOfLastPost, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	for _, post := range posts {
		if post.Created.After(timeOfLastPost) {
			timeOfLastPost = post.Created
			id = post.UserID
		}

	}
	return id
}

func getThreadIDFromCookies(r *http.Request) int {
	return 2
}

// MainHandler is a method of the Repository struct that handles requests to the main page.
// It renders the "home.page.html" template to the provided HTTP response writer.
func (m *Repository) ThemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		threadID := getThreadIDFromCookies(r)
		posts, err := m.DB.GetAllPostsFromThread(threadID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get all posts from thread", "/error-page")
		}

		var postsInfo []models.PostDataForThemePage

		for _, post := range posts {
			var user models.User
			user, err = m.DB.GetUserByID(post.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user by id", "/error-page")
			}
			var info models.PostDataForThemePage
			info.Subject = post.Subject
			info.Created = post.Created.Format("2006-01-02 15:04:05")
			info.Content = post.Content
			info.PictureUserWhoCreatedPost = user.Picture
			info.UserNameWhoCreatedPost = user.UserName
			postsInfo = append(postsInfo, info)
		}

		data := make(map[string]interface{})

		data["posts"] = postsInfo

		mainThread, err := m.DB.GetThreadByID(getThreadIDFromCookies(r))
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get thread by id", "/error-page")
		}

		creator, err := m.DB.GetUserByID(mainThread.UserID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get user as creator", "/error-page")
		}

		data["creatorName"] = creator.UserName
		data["creatorImg"] = creator.Picture
		data["mainThreadName"] = mainThread.Subject
		data["mainThreadCreatedTime"] = mainThread.Created.Format("2006-01-02 15:04:05")
		renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{
			Data: data,
		})
	}
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
			<p>An error occurred: ` + errorMessage + `</p>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlContent))
}

// setErrorContext sets the error message in the context and adds it to the redirect URL
func setErrorAndRedirect(w http.ResponseWriter, r *http.Request, errorMessage string, redirectURL string) {
	log.Printf("Error: %s", errorMessage)
	// Append the error message as a query parameter in the redirect URL
	redirectURL += "?error=" + url.QueryEscape(errorMessage)

	// Perform the redirect
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
