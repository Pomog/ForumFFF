package main

import (
	"fmt"
	"net/http"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/handler"
)

func routes(a *config.AppConfig) http.Handler {

	mux := http.NewServeMux()

	statics := []string{"css", "logo", "ava"}

	for _, static := range statics {
		mux.Handle(fmt.Sprintf("/static/%s/", static), http.StripPrefix(fmt.Sprintf("/static/%s/", static), http.FileServer(http.Dir(fmt.Sprintf("static/%s", static)))))
	}

	// fileServer := http.FileServer(http.Dir("./static/"))
	// mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", handler.Repo.HomeHandler)
	mux.HandleFunc("/login", handler.Repo.LoginHandler)
	mux.HandleFunc("/register", handler.Repo.RegisterHandler)
	mux.HandleFunc("/theme", handler.Repo.ThemeHandler)
	return mux
}
