package handler

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

	visitor, err := m.DB.GetUserByID(visitorID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get visitor ID, m.DB.GetUserByID(visitorID)", "/error-page")
	}

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
		if visitor.UserName == "guest" || visitor.UserName == "" {
			setErrorAndRedirect(w, r, guestRestiction, "/error-page")
			return
		}

		// Parse the form data, including files Need to Set Upper limit for DATA
		err := r.ParseMultipartForm(1 << 20)
		if err != nil {
			setErrorAndRedirect(w, r, "Image is too large", "/error-page")
			return
		}

		post := models.Post{
			Subject:  ShortenerOfSubject(mainThread.Subject),
			Content:  r.FormValue("post-text"),
			UserID:   visitorID,
			ThreadId: mainThread.ID,
			Image:    r.FormValue("image"),
		}

		// checking if there is a text before thread creation
		if post.Content == "" {
			setErrorAndRedirect(w, r, "Empty post can not be created", "/error-page")
			return
		}

		// checking text length
		if len(post.Content) > 500 {
			setErrorAndRedirect(w, r, "Only 500 symbols allowed", "/error-page")
			return
		}
		// ADD IMAGE TO STATIC
		// Get the file from the form data
		file, handler, errFileGet := r.FormFile("image")
		if errFileGet != nil {
			setErrorAndRedirect(w, r, fileReceivingErrorMsg, "/error-page")
			return
		}
		defer file.Close()

		// Validate file size (1 MB limit)
		if handler.Size > 1<<20 {
			setErrorAndRedirect(w, r, "File size should be below 1 MB", "/error-page")
			return
		}

		// Validate file type (must be an image)
		contentType := handler.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
			setErrorAndRedirect(w, r, "Wrong File Formate, allowed jpeg, png, gif ", "/error-page")
			return
		}

		// Create a new file in the "static/post_images" directory
		newFilePath := filepath.Join("static/post_images", handler.Filename)
		newFile, errFileCreate := os.Create(newFilePath)
		if errFileCreate != nil {
			setErrorAndRedirect(w, r, fileCreatingErrorMsg, "/error-page")
			return
		}
		defer newFile.Close()

		// Copy the uploaded file to the new file
		_, err = io.Copy(newFile, file)
		if err != nil {
			setErrorAndRedirect(w, r, fileSavingErrorMsg, "/error-page")
			return
		}

		post.Image = path.Join("/", newFilePath)

		err = m.DB.CreatePost(post)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not create a post"+err.Error(), "/error-page")
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
		info.Image = post.Image
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
