package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// RegisterHandler handles both GET and POST requests for the registration page.
func (m *Repository) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	sessionUserID := m.GetLoggedUser(w, r)
	user, _ := m.DB.GetUserByID(sessionUserID)

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

		if user.UserName == "guest" || user.UserName == "" {
			setErrorAndRedirect(w, r, "Guests can not edit/delete posts", "/error-page")
		} else if user.ID != post.UserID {
			setErrorAndRedirect(w, r, "Only Admin or Creator of the Post can Edit / Delete it", "/error-page")
		}
		initialFormData.Content = post.Content
		initialFormData.Subject = post.Subject
		initialFormData.UserID = post.UserID
		initialFormData.ThreadId = post.ThreadId
		initialFormData.ID = post.ID
		data := make(map[string]interface{})
		data["content"] = initialFormData
		renderer.RendererTemplate(w, "edit_post.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}

func (m *Repository) EditPostResultHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if strings.TrimSpace(r.FormValue("post-text")) == "" || len(r.FormValue("post-text")) > m.App.PostLen {
			setErrorAndRedirect(w, r, "The post is empty or too long", "/error-page")
			return
		}

		postID, err1 := strconv.Atoi(r.URL.Query().Get("postID"))
		if err1 != nil {
			setErrorAndRedirect(w, r, "Could not convert postID into integer: "+err1.Error(), "/error-page")
		}
		post, err2 := m.DB.GetPostByID(postID)
		if err2 != nil {
			setErrorAndRedirect(w, r, "Could not get post from GetPostByID: "+err2.Error(), "/error-page")
		}
		post.Content = r.FormValue("post-text")
		err3 := m.DB.EditPost(post)

		if err3 != nil {
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

func (m *Repository) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	sessionUserID := m.GetLoggedUser(w, r)
	user, _ := m.DB.GetUserByID(sessionUserID)

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
		err3 := m.DB.DeletePost(post)

		if err3 != nil {
			setErrorAndRedirect(w, r, "Could not m.DB.DeletePost(post): "+err3.Error(), "/error-page")
		}

		message := fmt.Sprintf("Post ID - %v deleted by User %s with email %s", post.ID, user.UserName, user.Email)
		fmt.Println(message)
		// helper.SendEmail(m.App.ServerEmail, message)

		data := make(map[string]interface{})
		data["post"] = post.Content
		data["threadID"] = post.ThreadId

		renderer.RendererTemplate(w, "edit_post_result.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}
