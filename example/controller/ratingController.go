package controller

import (
	"net/http"
	"strconv"

	"forum-authentication/types"
)

type RatingController struct{}

func (_ *RatingController) RatingController(w http.ResponseWriter, r *http.Request) {

	var replyRating types.ReplyRating

	if r.Method == "POST" {
		user, err := ValidateSession(w, r)
		referer := r.Header.Get("Referer")

		if err != nil {
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}

		if (user == types.User{}) {
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}

		post_id_string := r.URL.Query().Get("post_id")
		rating_id_string := r.URL.Query().Get("rating_id")

		// Define a common function to process ratings.
		processRating := func(id int, rating string, handleFunc func(int, int, string)) {
			if id == 0 {
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}

			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			ratingValue := r.FormValue("rating")
			if user.Id == 0 || ratingValue == "" {
				http.Error(w, "Missing parameters", http.StatusBadRequest)
				return
			}

			handleFunc(id, user.Id, ratingValue)
			http.Redirect(w, r, referer, http.StatusSeeOther)
		}

		if post_id_string != "" {
			postID, err := strconv.Atoi(post_id_string)

			if err != nil {
				http.Error(w, "Invalid post_id", http.StatusBadRequest)
				return
			}
			processRating(postID, "rating", postRating.HandlePostRating)

		} else if rating_id_string != "" {
			ratingID, err := strconv.Atoi(rating_id_string)
			if err != nil {
				http.Error(w, "Invalid rating_id", http.StatusBadRequest)
				return
			}
			processRating(ratingID, "rating", replyRating.HandleReplyRating)
		}
	}
}
