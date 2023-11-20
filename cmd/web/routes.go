package main

import (
	"net/http"

	"github.com/Pomog/ForumFFF/pkg/config"
	"github.com/Pomog/ForumFFF/pkg/handler"
)

func routes(a *config.AppConfig) http.Handler {

	mux := http.NewServeMux()

	//connect CSS and JS files
	mux.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	mux.Handle("/static/logo/", http.StripPrefix("/static/logo/", http.FileServer(http.Dir("static/logo"))))
	mux.Handle("/static/ava/", http.StripPrefix("/static/ava/", http.FileServer(http.Dir("static/ava"))))

	mux.HandleFunc("/", handler.Repo.MainHandler)
	mux.HandleFunc("/about", handler.Repo.AboutHandler)
	mux.HandleFunc("/theme", handler.Repo.ThemeHandler)
	return mux
}
