package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Pomog/ForumFFF/internal/helper"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// HomeHandler handles both GET and POST requests for the home page.
func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	sessionUserID := m.GetLoggedUser(w, r)

	if r.Method == http.MethodGet {
		handleGetRequest(w, r, m, sessionUserID)

	} else if r.Method == http.MethodPost {
		handlePostRequest(w, r, m, sessionUserID)
	}
}

// handleGetRequest handles GET requests for the home page.
func handleGetRequest(w http.ResponseWriter, r *http.Request, m *Repository, sessionUserID int) {
	search := r.FormValue("search")
	category := r.URL.Query().Get("searchCategory")

	threads, err := getThreadsBySearchOrCategory(m, search, category)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get threads"+err.Error(), "/error-page")
		return
	}

	threadsInfo, err := processThreads(m, threads)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get threadsInfo from: processThreads(m, threads)"+err.Error(), "/error-page")
		return
	}

	data, err := prepareDataForTemplate(w, r, m, sessionUserID, threadsInfo)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get data from: prepareDataForTemplate(w, r, m, sessionUserID, threadsInfo)"+err.Error(), "/error-page")
		return
	}
	renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{
		Data: data,
	})
}

// getThreadsBySearchOrCategory retrieves threads based on search or category after trimming spaces.
func getThreadsBySearchOrCategory(m *Repository, search, category string) ([]models.Thread, error) {
	var threads []models.Thread
	var err error

	// Trim leading and trailing spaces from search and category
	search = strings.TrimSpace(search)
	category = strings.TrimSpace(category)

	if search != "" {
		threads, err = m.DB.GetSearchedThreads(search)
	} else if category != "" {
		threads, err = m.DB.GetSearchedThreadsByCategory(category)
	} else {
		threads, err = m.DB.GetAllThreads()
	}
	return threads, err
}

// processThreads processes the retrieved threads to generate necessary info.
func processThreads(m *Repository, threads []models.Thread) ([]models.ThreadDataForMainPage, error) {
	var threadsInfo []models.ThreadDataForMainPage

	for _, thread := range threads {
		info, err := processThreadInfo(m, thread)
		if err != nil {
			return threadsInfo, err
		}
		threadsInfo = append(threadsInfo, info)
	}

	return threadsInfo, nil
}

// processThreadInfo processes individual thread info for display.
func processThreadInfo(m *Repository, thread models.Thread) (models.ThreadDataForMainPage, error) {
	info := models.ThreadDataForMainPage{}

	user, err := m.DB.GetUserByID(thread.UserID)
	if err != nil {
		return info, err
	}

	// Populate info with thread and user data
	info.ThreadID = thread.ID
	info.Subject = thread.Subject
	info.Created = thread.Created.Format("2006-01-02 15:04:05")
	info.Category = thread.Category

	info.PictureUserWhoCreatedThread = user.Picture
	info.UserNameWhoCreatedThread = user.UserName
	posts, err := m.DB.GetAllPostsFromThread(thread.ID)
	if err != nil {
		log.Fatal(err)
	}
	info.Posts = posts

	userIDwhoCreatedLastPost := getUserThatCreatedLastPost(posts)
	if userIDwhoCreatedLastPost != 0 || len(posts) != 0 {
		userWhoCreatedLastPost, err := m.DB.GetUserByID(userIDwhoCreatedLastPost)
		if err != nil {
			return info, err
		}

		info.PictureUserWhoCreatedLastPost = userWhoCreatedLastPost.Picture
		info.UserNameWhoCreatedLastPost = userWhoCreatedLastPost.UserName
	} else if userIDwhoCreatedLastPost == 0 || len(posts) == 0 {
		info.Created = ""
	}

	return info, nil
}

// prepareDataForTemplate prepares data for rendering the template.
func prepareDataForTemplate(w http.ResponseWriter, r *http.Request, m *Repository, sessionUserID int, threadsInfo []models.ThreadDataForMainPage) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	loggedUser, err := m.DB.GetUserByID(sessionUserID)
	if err != nil {
		return data, err
	}

	data["games"] = m.App.GamesList
	data["threads"] = threadsInfo
	data["loggedAs"] = loggedUser.UserName
	data["loggedAsID"] = loggedUser.ID

	return data, nil
}

// handlePostRequest handles POST requests for creating new threads.
func handlePostRequest(w http.ResponseWriter, r *http.Request, m *Repository, sessionUserID int) {
	loggedUser, err := m.DB.GetUserByID(sessionUserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get user by ID: m.DB.GetUserByID(sessionUserID)", "/error-page")
		return
	}

	userName := loggedUser.UserName
	if userName == "guest" || strings.TrimSpace(userName) == "" {
		setErrorAndRedirect(w, r, guestRestiction, "/error-page")
		return
	}

	thread := createThreadFromRequest(m, w, r, sessionUserID)
	thread.Category = strings.TrimSpace(helper.CorrectPunctuationsSpaces(thread.Category))
	thread.Subject = strings.TrimSpace(helper.CorrectPunctuationsSpaces(thread.Subject))

	// Validation of the thread info
	validationParameters := models.ValidationConfig{
		MinCategoryLen:   m.App.MinCategoryLen,
		MaxCategoryLen:   m.App.MaxCategoryLen,
		MinSubjectLen:    m.App.MinSubjectLen,
		MaxSubjectLen:    m.App.MaxSubjectLen,
		SingleWordMaxLen: len(m.App.LongestSingleWord),
	}

	validationsErrors := thread.Validate(validationParameters)
	if len(validationsErrors) > 0 {
		// prepare error msg
		var errorMsg string
		for _, err := range validationsErrors {
			errorMsg += err.Error() + "\n"
		}
		setErrorAndRedirect(w, r, errorMsg, "/error-page")
		return
	}

	id, err := m.DB.CreateThread(thread)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not create thread: m.DB.CreateThread(thread)"+err.Error(), "/error-page")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/theme?threadID=%d", id), http.StatusPermanentRedirect)
}

// createThreadFromRequest creates a thread from the HTTP request.
func createThreadFromRequest(m *Repository, w http.ResponseWriter, r *http.Request, sessionUserID int) models.Thread {
	thread := models.Thread{
		Subject:  r.FormValue("message-text"),
		Category: r.FormValue("category-text"),
		UserID:   sessionUserID,
	}
	AttachFile(m, w, r, nil, &thread)
	return thread
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

func getThreadIDFromQuery(w http.ResponseWriter, r *http.Request) int {
	threadID, err := strconv.Atoi(r.URL.Query().Get("threadID"))
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get all posts from thread", "/error-page")
		return 0
	}
	return threadID
}
