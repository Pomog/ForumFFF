package server

import (
	"fmt"
	"net/http"

	"forum-authentication/controller"
)

var (
	userController         controller.UserController
	postController         controller.PostController
	categoryController     controller.CategoryController
	homePageController     controller.HomePageController
	handleRatingController controller.RatingController
	handleReplyController  controller.ReplyController
)

type Server struct {
	ListenAddress string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddress: listenAddr,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/signup", userController.CreateUser)

	http.HandleFunc("/login", userController.Login)

	http.HandleFunc("/logout", userController.Logout)

	http.HandleFunc("/me", userController.ProfilePage)

	http.HandleFunc("/create", postController.CreatePost)

	http.HandleFunc("/category", categoryController.CategoryController)

	http.HandleFunc("/handle-rating", handleRatingController.RatingController)

	http.HandleFunc("/handle-reply", handleReplyController.ReplyController)

	http.HandleFunc("/auth/google", userController.GoogleAuth)

	http.HandleFunc("/auth/google/callback", userController.AuthCallback)

	http.HandleFunc("/auth/github", userController.GithubAuth)

	http.HandleFunc("/auth/github/callback", userController.AuthCallback)

	http.HandleFunc("/", homePageController.HomePage)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets"))))

	fmt.Println("Server running at port", s.ListenAddress)
	return http.ListenAndServe(s.ListenAddress, nil)
}
