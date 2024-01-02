package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum-authentication/types"
)

type CategoryController struct{}

func (_ *CategoryController) CategoryController(w http.ResponseWriter, r *http.Request) {

	category_slug := r.URL.Query().Get("slug")

	if r.Method == "GET" {
		categories := []types.Categories{}
		var err error

		if category_slug == "" {
			categories, err = category.GetCategories()
		} else {
			categories, err = category.GetCategoryBySlug(category_slug)
		}

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error getting categories", http.StatusInternalServerError)
			return
		}

		categoriesJson, err := json.Marshal(categories)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(categoriesJson)

	}
}
