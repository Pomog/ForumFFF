package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum-authentication/types"
)

type HomePageController struct{}

var (
	category   types.Categories
	postRating types.PostRating
	postReply  types.PostReply
)

func (_ *HomePageController) HomePage(w http.ResponseWriter, r *http.Request) {
	user, err := ValidateSession(w, r)
	fmt.Println(user)

	data := struct {
		SessionValid        bool
		Categories          []types.Categories
		CurrentCategory     types.Categories
		CurrentPost         types.Post
		CurrentPostReplies  []types.PostReply
		CurrentPostDislikes int
		CurrentPostLikes    int

		Posts []types.Post
	}{
		SessionValid:        err == nil,
		Categories:          []types.Categories{},
		CurrentCategory:     types.Categories{},
		CurrentPost:         types.Post{},
		CurrentPostReplies:  []types.PostReply{},
		CurrentPostDislikes: 0,
		CurrentPostLikes:    0,
		Posts:               []types.Post{},
	}

	// Check if the URL path is not root, return not found template
	if r.URL.Path != "/" {
		tmpl, err := template.ParseGlob("ui/templates/notFound.html")
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(w, r)
		return
	}

	// Get all categories for the topics sidebar
	categories, err := category.GetCategories()
	if err != nil {
		log.Println(err)
	}
	data.Categories = categories

	postID := r.URL.Query().Get("post")
	filter := r.URL.Query().Get("filter")
	categorySlug := r.URL.Query().Get("category")

	if postID != "" {
		postID = r.URL.Query().Get("post")
		category, err := category.GetCurrentCategory(categorySlug)
		if err != nil {
			renderNotFoundTemplate(w, r)
			return
		}

		data.CurrentCategory = category
		currentPost, err := post.GetPostById(postID)
		if err != nil {
			renderNotFoundTemplate(w, r)
			return
		}

		dislikes, likes, err := postRating.GetPostRatings(postID)
		if err != nil {
			renderNotFoundTemplate(w, r)
			return
		}
		content, err := postReply.GetPostReplies(postID)

		if err != nil {
			renderNotFoundTemplate(w, r)
			return
		}

		data.CurrentPostReplies = content
		data.CurrentPostDislikes = dislikes
		data.CurrentPostLikes = likes
		data.CurrentPost = currentPost
		renderTemplate("ui/templates/post.html", w, data)

	} else if categorySlug != "" {

		category, err := category.GetCurrentCategory(categorySlug)

		if err != nil || category.Id == 0 {
			renderNotFoundTemplate(w, r)
			return
		}

		data.CurrentCategory = category

		var posts []types.Post

		switch filter {
		case "liked-posts":
			user, err := ValidateSession(w, r)

			referer := r.Header.Get("referer")

			if err != nil {
				http.Redirect(w, r, referer, http.StatusSeeOther)
				return
			}
			posts, err = post.GetCategoryLikedPosts(category, user.Id)
			if err != nil {
				log.Println(err)
			}
			break

		case "created-posts":
			user, err := ValidateSession(w, r)

			referer := r.Header.Get("referer")

			if err != nil {
				http.Redirect(w, r, referer, http.StatusSeeOther)
				return
			}
			posts, err = post.GetCategoryCreatedPosts(category, user.Id)
			if err != nil {
				log.Println(err)
			}
			break

		default:
			posts, err = post.GetCategoryPosts(category)
		}

		if err != nil || len(posts) == 0 {
			log.Println(err)
		}

		data.Posts = posts
		renderTemplate("ui/templates/home.html", w, data)

	} else {
		// when category or post id is not provided, return first category from the database

		data.CurrentCategory = categories[0]
		data.Posts, err = post.GetCategoryPosts(categories[0])

		if err != nil {
			log.Println(err)
		}
		var posts []types.Post
		category, err := category.GetCurrentCategory("python")
		if err != nil {
			log.Println(err)
		}

		switch filter {
		case "liked-posts":
			user, err := ValidateSession(w, r)

			referer := r.Header.Get("referer")

			if err != nil {
				http.Redirect(w, r, referer, http.StatusSeeOther)
				return
			}
			posts, err = post.GetCategoryLikedPosts(category, user.Id)
			if err != nil {
				log.Println(err)
			}
			break

		case "created-posts":
			user, err := ValidateSession(w, r)

			referer := r.Header.Get("referer")

			if err != nil {
				http.Redirect(w, r, referer, http.StatusSeeOther)
				return
			}
			posts, err = post.GetCategoryCreatedPosts(category, user.Id)
			if err != nil {
				log.Println(err)
			}
			break

		default:
			posts, err = post.GetCategoryPosts(category)
		}

		if err != nil || len(posts) == 0 {
			log.Println(err)
		}

		data.Posts = posts
		renderTemplate("ui/templates/home.html", w, data)
	}
}

func renderNotFoundTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("ui/templates/notFound.html")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, r)
}

func renderTemplate(templatePath string, w http.ResponseWriter, data interface{}) {
	tmpl, err := template.ParseGlob(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, data)
}
