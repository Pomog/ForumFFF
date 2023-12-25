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

// HomeHandler handles both GET and POST requests for the registration page.
func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	sessionUserID := m.GetLoggedUser(w, r)

	if r.Method == http.MethodGet {
		search := r.FormValue("search")
		var threads []models.Thread
		var err error
		if search != "" {
			threads, err = m.DB.GetSearchedThreads(search)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get Threads m.DB.GetSearchedThreads", "/error-page")
				return
			}
		} else {
			threads, err = m.DB.GetAllThreads()
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get Threads m.DB.GetAllThreads", "/error-page")
				return
			}
		}

		var threadsInfo []models.ThreadDataForMainPage
		for _, thread := range threads {
			var user models.User
			user, err = m.DB.GetUserByID(thread.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID", "/error-page")
				return
			}

			user, err = m.DB.GetUserByID(thread.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID", "/error-page")
				return
			}

			var info models.ThreadDataForMainPage
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
					setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID(getUserThatCreatedLastPost(posts)) 95", "/error-page")
					return
				}

				info.PictureUserWhoCreatedLastPost = userWhoCreatedLastPost.Picture
				info.UserNameWhoCreatedLastPost = userWhoCreatedLastPost.UserName
			} else if userIDwhoCreatedLastPost == 0 || len(posts) == 0 {
				info.Created = ""
			}

			threadsInfo = append(threadsInfo, info)
		}

		data := make(map[string]interface{})
		loggedUser, err := m.DB.GetUserByID(sessionUserID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID(UserID)", "/error-page")
			return
		}
		data["threads"] = threadsInfo
		data["loggedAs"] = loggedUser.UserName
		data["loggedAsID"] = loggedUser.ID

		renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{
			Data: data,
		})
	} else if r.Method == http.MethodPost {

		loggedUser, err := m.DB.GetUserByID(sessionUserID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID(UserID), HomeHandler", "/error-page")
			return
		}

		userName := loggedUser.UserName
		if userName == "guest" || strings.TrimSpace(userName) == "" {
			setErrorAndRedirect(w, r, guestRestiction, "/error-page")
			return
		}
		thread := models.Thread{
			Subject:  r.FormValue("message-text"),
			Category: r.FormValue("category-text"),
			UserID:   sessionUserID,
		}
		AttachFile(m, w, r, nil, &thread)

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
			setErrorAndRedirect(w, r, "Could not create a thread: "+err.Error(), "/error-page")
			return
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

func getThreadIDFromQuery(w http.ResponseWriter, r *http.Request) int {
	threadID, err := strconv.Atoi(r.URL.Query().Get("threadID"))
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get all posts from thread", "/error-page")
	}
	return threadID
}
