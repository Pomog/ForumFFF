package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Pomog/ForumFFF/internal/helper"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// ThemeHandler handles the main functionality of the theme page.
func (m *Repository) ThemeHandler(w http.ResponseWriter, r *http.Request) {
	visitorID, err := getVisitorID(m, w, r)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get visitor\n"+err.Error(), "/error-page")
		return
	}
	visitor, err := m.DB.GetUserByID(visitorID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get visitor\n"+err.Error(), "/error-page")
		return
	}

	threadID := getThreadIDFromQuery(w, r)
	if threadID == 0 {
		setErrorAndRedirect(w, r, "Could not get thread or creator\n"+err.Error(), "/error-page")
		return
	}
	mainThread, creator, err := getThreadAndCreator(m, threadID)
	if err != nil || threadID == 0 {
		setErrorAndRedirect(w, r, "Could not get thread or creator\n"+err.Error(), "/error-page")
		return
	}

	handlePostActions(w, r, m, visitorID, visitor, mainThread)

	if r.Method == http.MethodPost && (r.FormValue("like") == "" && r.FormValue("dislike") == "") {
		handlePostCreation(w, r, m, visitorID, mainThread)
		return
	}

	postsInfo, err := getPostsInfo(m, w, r, threadID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get posts or user information \n"+err.Error(), "/error-page")
		return
	}

	data, err := prepareDataForThemePage(m, w, r, visitorID, postsInfo, mainThread, creator)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not data from prepareDataForThemePage \n"+err.Error(), "/error-page")
		return
	}

	data["loggedAsID"] = visitorID

	renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{
		Data: data,
	})
}

// getVisitorID retrieves the visitor ID from cookies or generates a guest ID.
func getVisitorID(m *Repository, w http.ResponseWriter, r *http.Request) (int, error) {
	visitorID, err := m.DB.GetGuestID()
	if err != nil {
		return 0, err
	}

	for _, cookie := range r.Cookies() {
		if cookie.Value == m.App.UserLogin.String() {
			userID, err := strconv.Atoi(cookie.Name)
			if err != nil {
				return 0, nil
			}
			if visitorID = userID; visitorID != 0 {
				break
			}
		}
	}
	return visitorID, nil
}

// getThreadAndCreator retrieves the main thread and its creator information.
func getThreadAndCreator(m *Repository, threadID int) (mainThread models.Thread, creator models.User, err error) {
	mainThread, err = m.DB.GetThreadByID(threadID)
	if err != nil {
		return mainThread, creator, err
	}

	creator, err = m.DB.GetUserByID(mainThread.UserID)
	return mainThread, creator, err
}

// handlePostActions handles like/dislike actions for posts.
func handlePostActions(w http.ResponseWriter, r *http.Request, m *Repository, visitorID int, visitor models.User, mainThread models.Thread) {
	like := r.FormValue("like")
	dislike := r.FormValue("dislike")

	// Handle like action
	if like != "" {
		handleLikeAction(w, r, m, visitorID, visitor, like)
		return
	}

	// Handle dislike action
	if dislike != "" {
		handleDislikeAction(w, r, m, visitorID, visitor, dislike)
		return
	}
}

// handleLikeAction handles the 'like' action for posts.
func handleLikeAction(w http.ResponseWriter, r *http.Request, m *Repository, visitorID int, visitor models.User, like string) {
	if visitor.UserName == "guest" {
		setErrorAndRedirect(w, r, guestRestiction, "/error-page")
		return
	}
	postID, _ := strconv.Atoi(like)
	err := m.DB.LikePostByUserIdAndPostId(visitorID, postID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not LikePostByUserIdAndPostId", "/error-page")
		return
	}
}

// handleDislikeAction handles the 'dislike' action for posts.
func handleDislikeAction(w http.ResponseWriter, r *http.Request, m *Repository, visitorID int, visitor models.User, dislike string) {
	if visitor.UserName == "guest" {
		setErrorAndRedirect(w, r, guestRestiction, "/error-page")
		return
	}
	postID, _ := strconv.Atoi(dislike)
	err := m.DB.DislikePostByUserIdAndPostId(visitorID, postID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not DislikePostByUserIdAndPostId", "/error-page")
		return
	}
}

// handlePostCreation handles the creation of a new post.
func handlePostCreation(w http.ResponseWriter, r *http.Request, m *Repository, visitorID int, mainThread models.Thread) {
	visitor, err := m.DB.GetUserByID(visitorID)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not get user", "/error-page")
		return
	}

	if visitor.UserName == "guest" || visitor.UserName == "" {
		setErrorAndRedirect(w, r, guestRestiction, "/error-page")
		return
	}

	err = r.ParseMultipartForm(m.App.FileSize << 20)
	if err != nil {
		setErrorAndRedirect(w, r, "Image is too large row 160 \n"+err.Error(), "/error-page")
		return
	}

	post, err := createPostFromRequest(m, w, r, visitorID, mainThread)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not create a post \n"+err.Error(), "/error-page")
		return
	}
	err = m.DB.CreatePost(post)
	if err != nil {
		setErrorAndRedirect(w, r, "Could not create a post \n"+err.Error(), "/error-page")
		return
	}

	path := fmt.Sprintf("/create_post_result?threadID=%v", post.ThreadId)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

// createPostFromRequest creates a post from the request data.
func createPostFromRequest(m *Repository, w http.ResponseWriter, r *http.Request, visitorID int, mainThread models.Thread) (models.Post, error) {
	post := models.Post{
		Subject:  ShortenerOfSubject(mainThread.Subject),
		Content:  r.FormValue("post-text"),
		UserID:   visitorID,
		ThreadId: mainThread.ID,
		Image:    r.FormValue("image"),
	}

	post.Content = strings.TrimSpace(post.Content)
	post.Content = helper.CorrectPunctuationsSpaces(post.Content)

	// Validation of the User info
	validationParameters := models.ValidationConfig{
		MinSubjectLen:    m.App.MinSubjectLen,
		MaxSubjectLen:    m.App.MaxSubjectLen,
		SingleWordMaxLen: len(m.App.LongestSingleWord),
	}

	validationsErrors := post.Validate(validationParameters)
	if len(validationsErrors) > 0 {
		// prepare error msg
		var errorMsg string
		for _, err := range validationsErrors {
			errorMsg += err.Error() + "\n"
		}
		return post, errors.New(errorMsg)
	}

	AttachFile(m, w, r, &post, nil)
	return post, nil
}

// getPostsInfo retrieves information for rendering posts.
func getPostsInfo(m *Repository, w http.ResponseWriter, r *http.Request, threadID int) ([]models.PostDataForThemePage, error) {
	posts, err := m.DB.GetAllPostsFromThread(threadID)
	var postsInfo []models.PostDataForThemePage

	if err != nil {
		return postsInfo, err
	}

	for _, post := range posts {
		// Populate postsInfo with post data
		var user models.User
		user, err = m.DB.GetUserByID(post.UserID)
		if err != nil {
			return postsInfo, err
		}

		userPostsAmount, err := m.DB.GetTotalPostsAmmountByUserID(post.UserID)
		if err != nil {
			return postsInfo, err
		}

		likes, dislikes, err := m.DB.CountLikesAndDislikesForPostByPostID(post.ID)
		if err != nil {
			return postsInfo, err
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
		info.Classification = post.Classification
		postsInfo = append(postsInfo, info)
	}

	return postsInfo, nil
}

// prepareDataForThemePage prepares data for rendering the theme page template.
func prepareDataForThemePage(m *Repository, w http.ResponseWriter, r *http.Request, visitorID int, postsInfo []models.PostDataForThemePage, mainThread models.Thread, creator models.User) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	// Populate data for rendering template
	data["posts"] = postsInfo

	//to get user in Nav bar ______________
	sessionUserID := m.GetLoggedUser(w, r)
	if sessionUserID == 0 {
		return data, errors.New("unautorized")
	}
	loggedUser, err := m.DB.GetUserByID(sessionUserID)
	if err != nil {
		return data, err
	}

	data["loggedAs"] = loggedUser.UserName
	data["loggedAsID"] = loggedUser.ID
	data["loggedUserType"] = loggedUser.Type
	//__________________________________
	creatorPostsAmount, err := m.DB.GetTotalPostsAmmountByUserID(mainThread.UserID)
	if err != nil {
		return data, err
	}

	data["creatorName"] = creator.UserName
	data["threadImg"] = mainThread.Image
	data["creatorID"] = creator.ID
	data["creatorRegistrationDate"] = creator.Created.Format("2006-01-02 15:04:05")
	data["creatorPostsAmount"] = creatorPostsAmount
	data["creatorImg"] = creator.Picture
	data["mainThreadName"] = mainThread.Subject
	data["mainThreadCategory"] = mainThread.Category
	data["mainThreadID"] = mainThread.ID
	data["mainThreadCreatedTime"] = mainThread.Created.Format("2006-01-02 15:04:05")
	data["games"] = m.App.GamesList

	return data, nil
}

// AttachFile attaches a file to a post or thread.
func AttachFile(m *Repository, w http.ResponseWriter, r *http.Request, post *models.Post, thread *models.Thread) {
	// ADD IMAGE TO STATIC_________________________
	// Get the file from the form data
	file, handler, errFileGet := r.FormFile("image")
	if errFileGet == nil {
		defer file.Close()

		// Validate file size (2 MB limit)
		if handler.Size > m.App.FileSize<<20 {
			setErrorAndRedirect(w, r, "File size should be below 2 MB", "/error-page")
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
		_, err := io.Copy(newFile, file)
		if err != nil {
			setErrorAndRedirect(w, r, fileSavingErrorMsg, "/error-page")
			return
		}
		if post != nil {
			post.Image = path.Join("/", newFilePath)
		} else if thread != nil {
			thread.Image = path.Join("/", newFilePath)
		}

	}
	//_______________________________________
}
