package handler

import (
	"fmt"
	"log"
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
			log.Fatal(err1)
		}
		post, err2 := m.DB.GetPostByID(postID)
		if err2 != nil {
			log.Fatal(err2)
		}
		initialFormData.Content = post.Content
		initialFormData.Subject = post.Subject
		initialFormData.UserID = post.UserID
		initialFormData.ThreadId = post.ThreadId
		data := make(map[string]interface{})
		data["content"] = initialFormData
		renderer.RendererTemplate(w, "edit_topic.page.html", &models.TemplateData{
			Data: data,
		})

		fmt.Println(r.FormValue("post-text"))
		fmt.Println("printing")

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}
