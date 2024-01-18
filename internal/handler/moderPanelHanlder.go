package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
	} else if r.Method == http.MethodPost {
		handlePostRequestModerPage(w, r, m, sessionUserID)
	}
}

func handlePostRequestModerPage(w http.ResponseWriter, r *http.Request, m *Repository, sessionUserID int) {
	err := r.ParseForm()
	if err != nil {
		setErrorAndRedirect(w, r, "Could not parse form "+err.Error(), "/error-page")
		return
	}

	topicID := r.FormValue("topicID")
	selectedCategory := r.FormValue("btnradio" + topicID)

	if r.FormValue("postID") != "" {
		postIDF := r.FormValue("postID")
		selectedCategory := r.FormValue("btnradio" + postIDF)

		postID, err := strconv.Atoi(r.FormValue("postID"))
		if err != nil {
			setErrorAndRedirect(w, r, "Could not convert string into int "+err.Error(), "/error-page")
			return
		}
		post, err := m.DB.GetPostByID(postID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get post by id "+err.Error(), "/error-page")
			return
		}
		cat := models.TextClassification(selectedCategory)

		err = m.DB.EditPostClassification(post, cat)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not edit post classification "+err.Error(), "/error-page")
			return
		}
	}

	if r.FormValue("topicID") != "" {
		topicID, err := strconv.Atoi(r.FormValue("topicID"))
		if err != nil {
			setErrorAndRedirect(w, r, "Could not convert string into int "+err.Error(), "/error-page")
			return
		}
		topic, err := m.DB.GetThreadByID(topicID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get topic by id "+err.Error(), "/error-page")
			return
		}
		cat := models.TextClassification(selectedCategory)
		fmt.Println("selectedCategory", cat)

		err = m.DB.EditTopicClassification(topic, cat)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not edit topic classification "+err.Error(), "/error-page")
			return
		}
	}

	// Redirect back to the previous page (referer)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

// handleGetRequest handles GET requests for the home page.
func handleGetRequestModerPage(w http.ResponseWriter, r *http.Request, m *Repository, sessionUserID int) {
	topicCat := r.URL.Query().Get("topic")
	postCat := r.URL.Query().Get("post")
	var topics []models.Thread
	var posts []models.Post
	var err error
	if topicCat != "" && postCat == "" {
		topics, err = m.DB.GetAllThreadsByClassification(models.TextClassification(topicCat))
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get topics by category "+err.Error(), "/error-page")
			return
		}
	} else if topicCat == "" && postCat != "" {
		posts, err = m.DB.GetAllPostsByClassification(models.TextClassification(postCat))
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get posts by category "+err.Error(), "/error-page")
			return
		}
	}

	data := make(map[string]interface{})
	loggedUser, err := m.DB.GetUserByID(sessionUserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get logged user from: GetUserByID(sessionUserID) -"+err.Error(), "/error-page")
		return
	}

	data["loggedAs"] = loggedUser.UserName
	data["loggedAsID"] = loggedUser.ID
	data["loggedUserType"] = loggedUser.Type
	data["categories"] = models.Classifications
	data["posts"] = posts
	data["topics"] = topics

	renderer.RendererTemplate(w, "moderMain.page.html", &models.TemplateData{
		Data: data,
	})
}
