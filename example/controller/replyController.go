package controller

import (
	"net/http"
	"strconv"

	"forum-authentication/types"
)

type ReplyController struct{}

func (_ *ReplyController) ReplyController(w http.ResponseWriter, r *http.Request) {
	var postReply types.PostReply

	switch r.Method {

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

		post_id_string := r.URL.Query().Get("post_id")

		post_id, err := strconv.Atoi(post_id_string)
		if err != nil {
			post_id = 0
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		content := r.FormValue("content")
		if user.Id == 0 || post_id == 0 || content == "" {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			return
		}

		postReply.CreatePostReply(post_id, user.Id, content)

		http.Redirect(w, r, "/?post="+post_id_string, http.StatusSeeOther)
	}
}
