package main

import (
	"fmt"
	"net/http"

	"github.com/Pomog/ForumFFF/pkg/config"
	"github.com/Pomog/ForumFFF/pkg/handler"
)

func routes(a *config.AppConfig) http.Handler {

	mux := http.NewServeMux()

	statics := []string{"css", "logo", "ava"}

	for _, static := range statics {
		mux.Handle(fmt.Sprintf("/static/%s/", static), http.StripPrefix(fmt.Sprintf("/static/%s/", static), http.FileServer(http.Dir(fmt.Sprintf("static/%s", static)))))
	}

	//connect CSS and JS files
	// mux.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	// mux.Handle("/static/logo/", http.StripPrefix("/static/logo/", http.FileServer(http.Dir("static/logo"))))
	// mux.Handle("/static/ava/", http.StripPrefix("/static/ava/", http.FileServer(http.Dir("static/ava"))))

	mux.HandleFunc("/", handler.Repo.MainHandler)
	mux.HandleFunc("/about", handler.Repo.AboutHandler)
	mux.HandleFunc("/theme", handler.Repo.ThemeHandler)
	return mux
}
