package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// RegisterHandler handles both GET and POST requests for the registration page.
func (m *Repository) EditTopicHandler(w http.ResponseWriter, r *http.Request) {
	sessionUserID := m.GetLoggedUser(w, r)
	user, _ := m.DB.GetUserByID(sessionUserID)

	if r.Method == http.MethodPost {

		var initialFormData models.Thread

		topicID, err1 := strconv.Atoi(r.URL.Query().Get("topicID"))
		if err1 != nil {
			setErrorAndRedirect(w, r, "Could not convert topicID into integer: "+err1.Error(), "/error-page")
			return
		}
		topic, err2 := m.DB.GetThreadByID(topicID)
		if err2 != nil {
			setErrorAndRedirect(w, r, "Could not get topic from m.DB.GetThreadByID(topicID): "+err2.Error(), "/error-page")
			return
		}

		if user.UserName == "guest" || user.UserName == "" {
			setErrorAndRedirect(w, r, "Guests can not edit/delete topics", "/error-page")
			return

		} else if user.ID != topic.UserID {
			setErrorAndRedirect(w, r, "Only Admin or Creator of the Topic can Edit / Delete it", "/error-page")
			return
		}
		initialFormData.Subject = topic.Subject
		initialFormData.UserID = topic.UserID
		initialFormData.ID = topic.ID
		data := make(map[string]interface{})
		data["content"] = initialFormData
		data["creatorName"] = user.UserName
		renderer.RendererTemplate(w, "edit_topic.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}

func (m *Repository) EditTopicResultHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if strings.TrimSpace(r.FormValue("post-text")) == "" || len(r.FormValue("post-text")) > m.App.PostLen {
			setErrorAndRedirect(w, r, "The post is empty or too long", "/error-page")
			return
		}

		topicID, err1 := strconv.Atoi(r.URL.Query().Get("topicID"))
		if err1 != nil {
			setErrorAndRedirect(w, r, "Could not convert topicID into integer: "+err1.Error(), "/error-page")
			return
		}
		topic, err2 := m.DB.GetThreadByID(topicID)
		if err2 != nil {
			setErrorAndRedirect(w, r, "Could not get post from GetPostByID: "+err2.Error(), "/error-page")
			return
		}
		topic.Subject = r.FormValue("post-text")
		err3 := m.DB.EditTopic(topic)

		if err3 != nil {
			setErrorAndRedirect(w, r, "Could not edit post using EditPost(post): "+err3.Error(), "/error-page")
			return
		}

		data := make(map[string]interface{})
		data["topic"] = topic.Subject
		data["threadID"] = topic.ID

		renderer.RendererTemplate(w, "edit_topic_result.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}
