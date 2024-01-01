package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"forum-authentication/types"
)

type PostController struct{}

var post types.Post

func (_ *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	_, err := ValidateSession(w, r)

	categories, err := category.GetCategories()
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		SessionValid    bool
		Categories      []types.Categories
		CurrentCategory types.Categories
	}{
		SessionValid:    err == nil,
		Categories:      categories,
		CurrentCategory: category,
	}

	switch r.Method {
	case "GET":

		RenderPage(w, "ui/templates/createPost.html", data)

	case "POST":

		user, err := ValidateSession(w, r)

		referer := r.Header.Get("referer")

		if err != nil {
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}

		if (user == types.User{}) {
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryIDs := r.PostForm["categories"]

		for _, categoryIDStr := range categoryIDs {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}

			// Create a separate Post for each selected category
			post := &types.Post{
				Title:   title,
				Content: content,
				UserId:  user.Id,
			}

			postID, err := post.CreatePost(*post)

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Error creating post", http.StatusInternalServerError)
				return
			}

			postsCategory := &types.PostCategories{
				CategoryId: categoryID,
				PostId:     int(postID),
			}
			_, err = category.CreatePostCategory(*&postsCategory)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Error creating posts category", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/?post="+strconv.Itoa(int(postID)), http.StatusSeeOther)
		}

	}
}
