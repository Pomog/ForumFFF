package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		setErrorAndRedirect(w, r, "Could not get threads", "/error-page")
		return
	}

	threadsInfo := processThreads(m, threads)

	data := prepareDataForTemplate(w, r, m, sessionUserID, threadsInfo)
	renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{
		Data: data,
	})
}

// getThreadsBySearchOrCategory retrieves threads based on search or category.
func getThreadsBySearchOrCategory(m *Repository, search, category string) ([]models.Thread, error) {
	var threads []models.Thread
	var err error

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
func processThreads(m *Repository, threads []models.Thread) []models.ThreadDataForMainPage {
	var threadsInfo []models.ThreadDataForMainPage

	for _, thread := range threads {
		info := processThreadInfo(m, thread)
		threadsInfo = append(threadsInfo, info)
	}

	return threadsInfo
}

// processThreadInfo processes individual thread info for display.
func processThreadInfo(m *Repository, thread models.Thread) models.ThreadDataForMainPage {
	info := models.ThreadDataForMainPage{}

	user, err := m.DB.GetUserByID(thread.UserID)
	if err != nil {
		// setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID", "/error-page")
		// return info
		log.Fatal(err)
	}

	// Populate info with thread and user data
	info.ThreadID = thread.ID
	info.Subject = thread.Subject
	info.Created = thread.Created.Format("2006-01-02 15:04:05")
	info.Category = thread.Category

	info.PictureUserWhoCreatedThread = user.Picture
	info.UserNameWhoCreatedThread = user.UserName

	return info
}

// prepareDataForTemplate prepares data for rendering the template.
func prepareDataForTemplate(w http.ResponseWriter, r *http.Request, m *Repository, sessionUserID int, threadsInfo []models.ThreadDataForMainPage) map[string]interface{} {
	data := make(map[string]interface{})
	loggedUser, err := m.DB.GetUserByID(sessionUserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get user by ID: m.DB.GetUserByID(sessionUserID)", "/error-page")
		return data // returning empty data
	}

	data["games"] = m.App.GamesList
	data["threads"] = threadsInfo
	data["loggedAs"] = loggedUser.UserName
	data["loggedAsID"] = loggedUser.ID

	return data
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
	// checking if there is a text before thread creation
	if strings.TrimSpace(thread.Subject) == "" {
		setErrorAndRedirect(w, r, "Empty thread can not be created", "/error-page")
		return
	}

	// checking text length
	if len(thread.Subject) > m.App.PostLen {
		setErrorAndRedirect(w, r, fmt.Sprintf("the post is to long, %d syblos allowed", m.App.PostLen), "/error-page")
		return
	}

	// checking if there is a category before thread creation
	if strings.TrimSpace(thread.Category) == "" {
		setErrorAndRedirect(w, r, "Empty category can not be created", "/error-page")
		return
	}

	// checking category length
	if len(thread.Category) > m.App.CategoryLen {
		setErrorAndRedirect(w, r, fmt.Sprintf("the category is to long, %d syblos allowed", m.App.CategoryLen), "/error-page")
		return
	}

	id, err := m.DB.CreateThread(thread)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not create thread: m.DB.CreateThread(thread)", "/error-page")
	}

	// r.Form.Del("message-text")
	// r.Form.Del("category-text")

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
	}
	return threadID
}
