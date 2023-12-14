package handler

import (
	"fmt"
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

	uuid := m.App.UserLogin
	if uuid.String() == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		m.App.InfoLog.Println("Could not get uuid in HomeHandler")
	}

	for _, cookie := range r.Cookies() {
		if cookie.Value == uuid.String() {
			userID, _ := m.DB.GetUserIDForSessionID(cookie.Value)
			if UserID = userID; UserID != 0 {
				break
			}

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
			info.ThreadID = thread.ID
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
		loggedUser, _ := m.DB.GetUserByID(UserID)
		data["threads"] = threadsInfo
		data["loggedAs"] = loggedUser.UserName

		renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{
			Data: data,
		})
	} else if r.Method == http.MethodPost {
		loggedUser, _ := m.DB.GetUserByID(UserID)
		userName := loggedUser.UserName
		if userName == "guest" {
			setErrorAndRedirect(w, r, "Guests can not create Themes and Posts, please log in or register!", "/error-page")
		}
		thread := models.Thread{
			Subject: r.FormValue("message-text"),
			UserID:  UserID,
		}

		id, err := m.DB.CreateThread(thread)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not create a thread", "/error-page")
		}

		http.Redirect(w, r, fmt.Sprintf("/theme?threadID=%d", id), http.StatusPermanentRedirect)
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

func getThreadIDFromURLquery(w http.ResponseWriter, r *http.Request) int {
	threadID, err := strconv.Atoi(r.URL.Query().Get("threadID"))
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get all posts from thread", "/error-page")
	}
	return threadID
}

// MainHandler is a method of the Repository struct that handles requests to the main page.
// It renders the "home.page.html" template to the provided HTTP response writer.
func (m *Repository) ThemeHandler(w http.ResponseWriter, r *http.Request) {

	visitorID, _ := m.DB.GetGuestID()

	for _, cookie := range r.Cookies() {
		if cookie.Value == m.App.UserLogin.String() {
			userID, err := strconv.Atoi(cookie.Name)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get visitor ID", "/error-page")
			}
			if visitorID = userID; visitorID != 0 {
				break
			}

		}
	}

	threadID := getThreadIDFromURLquery(w, r)

	mainThread, err := m.DB.GetThreadByID(threadID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get thread by id", "/error-page")
	}

	creator, err := m.DB.GetUserByID(mainThread.UserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get user as creator", "/error-page")
	}

	like := r.FormValue("like")
	dislike := r.FormValue("dislike")
	if like != "" {
		postID, _ := strconv.Atoi(like)
		m.DB.LikePostByUserIdAndPostId(visitorID, postID)
	}
	if dislike != "" {
		postID, _ := strconv.Atoi(dislike)
		m.DB.DislikePostByUserIdAndPostId(visitorID, postID)
	}
	//new post
	if r.Method == http.MethodPost && len(r.FormValue("post-text")) != 0 {
		post := models.Post{
			Subject:  shortenerOfSubject(mainThread.Subject),
			Content:  r.FormValue("post-text"),
			UserID:   visitorID,
			ThreadId: mainThread.ID,
		}
		err = m.DB.CreatePost(post)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not create a post", "/error-page")
		}
	}
	//-------
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
		userPostsAmmount, err := m.DB.GetTotalPostsAmmountByUserID(post.UserID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get ammount of Posts, GetTotalPostsAmmountByUserID", "/error-page")
		}

		likes, dislikes, err := m.DB.CountLikesAndDislikesForPostByPostID(post.ID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get Likes for Post, CountLikesAndDislikesForPostByPostID", "/error-page")
		}

		var info models.PostDataForThemePage
		info.ID = post.ID
		info.Subject = post.Subject
		info.Created = post.Created.Format("2006-01-02 15:04:05")
		info.Content = post.Content
		info.PictureUserWhoCreatedPost = user.Picture
		info.UserNameWhoCreatedPost = user.UserName
		info.UserRegistrationDate = user.Created.Format("2006-01-02 15:04:05")
		info.UserPostsAmmount = userPostsAmmount
		info.Likes = likes
		info.Dislikes = dislikes
		postsInfo = append(postsInfo, info)
	}

	data := make(map[string]interface{})

	data["posts"] = postsInfo

	creatorPostsAmmount, err := m.DB.GetTotalPostsAmmountByUserID(mainThread.UserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get ammount of Posts, GetTotalPostsAmmountByUserID", "/error-page")
	}

	data["creatorName"] = creator.UserName
	data["creatorRegistrationDate"] = creator.Created.Format("2006-01-02 15:04:05")
	data["creatorPostsAmmount"] = creatorPostsAmmount
	data["creatorImg"] = creator.Picture
	data["mainThreadName"] = mainThread.Subject
	data["mainThreadCreatedTime"] = mainThread.Created.Format("2006-01-02 15:04:05")

	renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{
		Data: data,
	})

}

func shortenerOfSubject(input string) string {
	if len(input) <= 20 {
		return input
	}
	return ("Topic:" + input[0:21] + "...")
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

func (m *Repository) PrivatPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "privatP.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

func (m *Repository) PersonaCabinetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		visitorID, _ := m.DB.GetGuestID()

		for _, cookie := range r.Cookies() {
			if cookie.Value == m.App.UserLogin.String() {
				userID, err := strconv.Atoi(cookie.Name)
				if err != nil {
					setErrorAndRedirect(w, r, "Could not get visitor ID", "/error-page")
				}
				if visitorID = userID; visitorID != 0 {
					break
				}

			}
		}
		user, errUser := m.DB.GetUserByID(visitorID)
		if errUser != nil {
			setErrorAndRedirect(w, r, "Could not get User from  GetUserByID(visitorID)", "/error-page")
		}
		var personalInfo models.User
		personalInfo.Email = user.Email
		personalInfo.Created = user.Created
		personalInfo.FirstName = user.FirstName
		personalInfo.LastName = user.LastName
		personalInfo.Picture = user.Picture
		personalInfo.UserName = user.UserName
		data := make(map[string]interface{})
		data["personal"] = personalInfo

		renderer.RendererTemplate(w, "personal.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}
