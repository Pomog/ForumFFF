package handler

import (
	"net/http"
	"strconv"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// RegisterHandler handles both GET and POST requests for the registration page.
func (m *Repository) EditTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var initialFormData models.Post

		postID, err1 := strconv.Atoi(r.URL.Query().Get("postID"))
		if err1 != nil {
			setErrorAndRedirect(w, r, "Could not convert postID into integer: "+err1.Error(), "/error-page")
		}
		post, err2 := m.DB.GetPostByID(postID)
		if err2 != nil {
			setErrorAndRedirect(w, r, "Could not get post from GetPostByID: "+err2.Error(), "/error-page")
		}
		initialFormData.Content = post.Content
		initialFormData.Subject = post.Subject
		initialFormData.UserID = post.UserID
		initialFormData.ThreadId = post.ThreadId
		initialFormData.ID = post.ID
		data := make(map[string]interface{})
		data["content"] = initialFormData
		renderer.RendererTemplate(w, "edit_topic.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}

func (m *Repository) EditTopicResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postID, err1 := strconv.Atoi(r.URL.Query().Get("postID"))
		if err1 != nil {
			setErrorAndRedirect(w, r, "Could not convert postID into integer: "+err1.Error(), "/error-page")
		}
		post, err2 := m.DB.GetPostByID(postID)
		if err2 != nil {
			setErrorAndRedirect(w, r, "Could not get post from GetPostByID: "+err2.Error(), "/error-page")
		}
		post.Content = r.FormValue("post-text")
		err3:=m.DB.EditPost(post)
		if err3!=nil{
			setErrorAndRedirect(w, r, "Could not edit post using EditPost(post): "+err3.Error(), "/error-page")
		}
		data := make(map[string]interface{})
		data["post"] = post.Content
		data["threadID"] = post.ThreadId
	
		renderer.RendererTemplate(w, "edit_topic_result.page.html", &models.TemplateData{
			Data: data,
		})


	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}
