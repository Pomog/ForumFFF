package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/Pomog/ForumFFF/internal/repository"
	"github.com/Pomog/ForumFFF/internal/repository/dbrepo"
	"github.com/google/uuid"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository handles the repository type, encapsulating the AppConfig and DatabaseInt dependencies.
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseInt
}

const (
	dbErrorUserPresent    = "DB Error func UserPresent"
	userAlreadyExistsMsg  = "User Already Exists"
	dbErrorCreateUser     = "DB Error func CreateUser"
	fileReceivingErrorMsg = "file receiving error"
	fileCreatingErrorMsg  = "Unable to create file"
	fileSavingErrorMsg    = "Unable to save file"
	guestRestiction       = "guests can not create Themes and Posts, Like or Dislike posts, please log in or register"
)

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *repository.DataBase) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewSQLiteRepo(a, db.SQL),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) ContactUsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "contactUs.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}

}

func (m *Repository) ForumRulesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "forumRules.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

func (m *Repository) HelpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "help.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

func (m *Repository) PrivatPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		renderer.RendererTemplate(w, "privatP.page.html", &models.TemplateData{})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

// PersonaCabinetHandler hanles the personal cabinet of selected user.
func (m *Repository) PersonaCabinetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID, _ := strconv.Atoi(r.URL.Query().Get("userID"))
		fmt.Println("userID ", userID)
		var personalInfo models.User
		user, errUser := m.DB.GetUserByID(userID)
		if errUser != nil {
			setErrorAndRedirect(w, r, "Could not get User from  GetUserByID(visitorID)", "/error-page")
			return
		}
		personalInfo.Email = user.Email
		personalInfo.ID = user.ID
		personalInfo.Created = user.Created
		personalInfo.FirstName = user.FirstName
		personalInfo.LastName = user.LastName
		personalInfo.Picture = user.Picture
		personalInfo.UserName = user.UserName
		totalPosts, _ := m.DB.GetTotalPostsAmmountByUserID(personalInfo.ID)
		data := make(map[string]interface{})
		data["personal"] = personalInfo
		data["totalPosts"] = totalPosts

		renderer.RendererTemplate(w, "personal.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

// GetAllThreadsForUserHandler gets all threads from user (user id)
func (m *Repository) GetAllThreadsForUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID, _ := strconv.Atoi(r.URL.Query().Get("userID"))
		user, errUser := m.DB.GetUserByID(userID)
		if errUser != nil {
			setErrorAndRedirect(w, r, "Could not get User from  GetUserByID(visitorID)", "/error-page")
			return
		}
		threads, err := m.DB.GetAllThreadsByUserID(user.ID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get threads from  GetAllThreadsByUserID(user.ID)", "/error-page")
			return
		}

		var threadsInfo []models.ThreadDataForMainPage
		for _, thread := range threads {
			var user models.User
			user, err = m.DB.GetUserByID(thread.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID", "/error-page")
				return
			}

			user, err = m.DB.GetUserByID(thread.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID", "/error-page")
				return
			}

			var info models.ThreadDataForMainPage
			info.ThreadID = thread.ID
			info.Subject = thread.Subject
			info.Created = thread.Created.Format("2006-01-02 15:04:05")
			info.Category = thread.Category

			info.PictureUserWhoCreatedThread = user.Picture
			info.UserNameWhoCreatedThread = user.UserName

			posts, err := m.DB.GetAllPostsFromThread(thread.ID)
			if err != nil {
				log.Fatal(err)
			}
			info.Posts = posts

			userIDwhoCreatedLastPost := getUserThatCreatedLastPost(posts)

			if userIDwhoCreatedLastPost != 0 || len(posts) != 0 {
				userWhoCreatedLastPost, err := m.DB.GetUserByID(userIDwhoCreatedLastPost)
				if err != nil {
					setErrorAndRedirect(w, r, "Could not get user as creator, m.DB.GetUserByID(getUserThatCreatedLastPost(posts)) 95", "/error-page")
					return
				}

				info.PictureUserWhoCreatedLastPost = userWhoCreatedLastPost.Picture
				info.UserNameWhoCreatedLastPost = userWhoCreatedLastPost.UserName
			} else if userIDwhoCreatedLastPost == 0 || len(posts) == 0 {
				info.Created = ""
			}

			threadsInfo = append(threadsInfo, info)
		}

		data := make(map[string]interface{})

		data["games"] = m.App.GamesList
		data["threads"] = threadsInfo

		renderer.RendererTemplate(w, "home.page.html", &models.TemplateData{
			Data: data,
		})

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
	}
}

// GetAllPostsForUserHandler gets all posts from user (user id)
func (m *Repository) GetAllPostsForUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID, _ := strconv.Atoi(r.URL.Query().Get("userID"))
		user, errUser := m.DB.GetUserByID(userID)
		if errUser != nil {
			setErrorAndRedirect(w, r, "Could not get User from  GetUserByID(visitorID)", "/error-page")
			return
		}

		posts, err := m.DB.GetAllPostsByUserID(user.ID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get all posts from user ID", "/error-page")
			return
		}

		var postsInfo []models.PostDataForThemePage

		for _, post := range posts {
			var user models.User
			user, err = m.DB.GetUserByID(post.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user by id", "/error-page")
				return
			}
			userPostsAmount, err := m.DB.GetTotalPostsAmmountByUserID(post.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get amount of Posts, GetTotalPostsAmountByUserID", "/error-page")
				return
			}

			likes, dislikes, err := m.DB.CountLikesAndDislikesForPostByPostID(post.ID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get Likes for Post, CountLikesAndDislikesForPostByPostID", "/error-page")
				return
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
		data["games"] = m.App.GamesList

		renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{
			Data: data,
		})
	}
}

func (m *Repository) GetAllLikedPostsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userID, _ := strconv.Atoi(r.URL.Query().Get("userID"))
		user, errUser := m.DB.GetUserByID(userID)
		if errUser != nil {
			setErrorAndRedirect(w, r, "Could not get User from  GetUserByID(visitorID)"+errUser.Error(), "/error-page")
			return
		}

		posts, err := m.DB.GetAllLikedPostsByUserID(user.ID)
		if err != nil {
			setErrorAndRedirect(w, r, "Could not get all posts from user ID GetAllLikedPostsByUserID(user.ID)"+err.Error(), "/error-page")
			return
		}

		var postsInfo []models.PostDataForThemePage

		for _, post := range posts {
			var user models.User
			user, err = m.DB.GetUserByID(post.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get user by id"+err.Error(), "/error-page")
				return
			}
			userPostsAmount, err := m.DB.GetTotalPostsAmmountByUserID(post.UserID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get amount of Posts, GetTotalPostsAmountByUserID"+err.Error(), "/error-page")
				return
			}

			likes, dislikes, err := m.DB.CountLikesAndDislikesForPostByPostID(post.ID)
			if err != nil {
				setErrorAndRedirect(w, r, "Could not get Likes for Post, CountLikesAndDislikesForPostByPostID"+err.Error(), "/error-page")
				return
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
		data["games"] = m.App.GamesList

		renderer.RendererTemplate(w, "theme.page.html", &models.TemplateData{
			Data: data,
		})
	}
}

// shortenerOfSubject helper function to squeeze theme name
func ShortenerOfSubject(input string) string {
	if len(input) <= 80 {
		return input
	}
	return "Topic:" + input[0:81] + "..."
}

func (m *Repository) GetLoggedUser(w http.ResponseWriter, r *http.Request) int {
	var UserID int
	loginUUID := m.App.UserLogin

	if loginUUID == uuid.Nil {
		m.App.InfoLog.Println("Could not get loginUUID in HomeHandler")
		return 0
	}

	for _, cookie := range r.Cookies() {
		if cookie.Value == loginUUID.String() {
			userID, _ := m.DB.GetUserIDForSessionID(cookie.Value)
			if UserID = userID; UserID != 0 {
				break
			}
		}
	}

	if UserID == 0 {
		return 0
	}
	return UserID
}
