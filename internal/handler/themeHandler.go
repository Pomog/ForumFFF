package handler

import (
	"net/http"
	"strconv"

	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// ThemeHandler handles both GET and POST requests for the theme page
func (m *Repository) ThemeHandler(w http.ResponseWriter, r *http.Request) {

	visitorID, _ := m.DB.GetGuestID()

	for _, cookie := range r.Cookies() {
		if cookie.Value == m.App.UserLogin.String() {
			userID, err := strconv.Atoi(cookie.Name)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get visitor ID", "/error-page")
			}
			if visitorID = userID; visitorID != 0 {
				break
			}
		}
	}

	visitor, _ := m.DB.GetUserByID(visitorID)

	threadID := getThreadIDFromQuery(w, r)

	mainThread, err := m.DB.GetThreadByID(threadID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get thread by id", "/error-page")
	}

	creator, err := m.DB.GetUserByID(mainThread.UserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get user as creator", "/error-page")
		return
	}

	like := r.FormValue("like")
	dislike := r.FormValue("dislike")
	if like != "" {
		if visitor.UserName == "guest" {
			setErrorAndRedirect(w, r, guestRestiction, "/error-page")
		}
		postID, _ := strconv.Atoi(like)
		err := m.DB.LikePostByUserIdAndPostId(visitorID, postID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not LikePostByUserIdAndPostId", "/error-page")
		}
	}
	if dislike != "" {
		if visitor.UserName == "guest" {
			setErrorAndRedirect(w, r, guestRestiction, "/error-page")
			return
		}
		postID, _ := strconv.Atoi(dislike)
		err := m.DB.DislikePostByUserIdAndPostId(visitorID, postID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not DislikePostByUserIdAndPostId", "/error-page")
		}
	}
	//new post
	if r.Method == http.MethodPost && len(r.FormValue("post-text")) != 0 {
		if visitor.UserName == "guest" {
			setErrorAndRedirect(w, r, guestRestiction, "/error-page")
			return
		}
		post := models.Post{
			Subject:  ShortenerOfSubject(mainThread.Subject),
			Content:  r.FormValue("post-text"),
			UserID:   visitorID,
			ThreadId: mainThread.ID,
		}
		err = m.DB.CreatePost(post)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not create a post", "/error-page")
		}
	}
	//-------
	posts, err := m.DB.GetAllPostsFromThread(threadID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get all posts from thread", "/error-page")
	}

	var postsInfo []models.PostDataForThemePage

	for _, post := range posts {
		var user models.User
		user, err = m.DB.GetUserByID(post.UserID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get user by id", "/error-page")
		}
		userPostsAmount, err := m.DB.GetTotalPostsAmmountByUserID(post.UserID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get amount of Posts, GetTotalPostsAmountByUserID", "/error-page")
		}

		likes, dislikes, err := m.DB.CountLikesAndDislikesForPostByPostID(post.ID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get Likes for Post, CountLikesAndDislikesForPostByPostID", "/error-page")
		}

		var info models.PostDataForThemePage
		info.ID = post.ID
		info.Subject = post.Subject
		info.Created = post.Created.Format("2006-01-02 15:04:05")
		info.Content = post.Content
		info.PictureUserWhoCreatedPost = user.Picture
		info.UserNameWhoCreatedPost = user.UserName
		info.UserIDWhoCreatedPost = user.ID
		info.UserRegistrationDate = user.Created.Format("2006-01-02 15:04:05")
		info.UserPostsAmmount = userPostsAmount
		info.Likes = likes
		info.Dislikes = dislikes
		postsInfo = append(postsInfo, info)
	}

	data := make(map[string]interface{})

	data["posts"] = postsInfo

	creatorPostsAmount, err := m.DB.GetTotalPostsAmmountByUserID(mainThread.UserID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get amount of Posts, GetTotalPostsAmountByUserID", "/error-page")
	}

	data["creatorName"] = creator.UserName
	data["creatorID"] = creator.ID
	data["creatorRegistrationDate"] = creator.Created.Format("2006-01-02 15:04:05")
	data["creatorPostsAmount"] = creatorPostsAmount
	data["creatorImg"] = creator.Picture
	data["mainThreadName"] = mainThread.Subject
	data["mainThreadCreatedTime"] = mainThread.Created.Format("2006-01-02 15:04:05")

	renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{
		Data: data,
	})
}


