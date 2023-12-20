package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/google/uuid"
)

// HomeHandler handles both GET and POST requests for the registration page.
func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	var UserID int

	loginUUID := m.App.UserLogin

	if loginUUID == uuid.Nil {
		m.App.InfoLog.Println("Could not get loginUUID in HomeHandler")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	for _, cookie := range r.Cookies() {
		if cookie.Value == loginUUID.String() {
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
		search := r.FormValue("search")
		var threads []models.Thread
		var err error
		if search != "" {
			threads, err = m.DB.GetSearchedThreads(search)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get Threads m.DB.GetSearchedThreads", "/error-page")
			}
		} else {
			threads, err = m.DB.GetAllThreads()
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get Threads m.DB.GetAllThreads", "/error-page")
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
			userWhoCreatedLastPost, err := m.DB.GetUserByID(getUserThatCreatedLastPost(posts))
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID(getUserThatCreatedLastPost(posts))", "/error-page")
				return
			}

			info.PictureUserWhoCreatedLastPost = userWhoCreatedLastPost.Picture
			info.UserNameWhoCreatedLastPost = userWhoCreatedLastPost.UserName
			threadsInfo = append(threadsInfo, info)
		}

		data := make(map[string]interface{})
		loggedUser, err := m.DB.GetUserByID(UserID)
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

		loggedUser, err := m.DB.GetUserByID(UserID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID(UserID), HomeHandler", "/error-page")
			return
		}
		userName := loggedUser.UserName
		if userName == "guest" {
			setErrorAndRedirect(w, r, guestRestiction, "/error-page")
		}
		thread := models.Thread{
			Subject: r.FormValue("message-text"),
			UserID:  UserID,
		}

		// checking if there is a text before thread creation
		if thread.Subject == "" {
			setErrorAndRedirect(w, r, "Empty thread can not be created", "/error-page")
			return
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

func getThreadIDFromQuery(w http.ResponseWriter, r *http.Request) int {
	threadID, err := strconv.Atoi(r.URL.Query().Get("threadID"))
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get all posts from thread", "/error-page")
	}
	return threadID
}
