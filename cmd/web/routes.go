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

	mux.HandleFunc("/", handler.Repo.MainHandler)
	mux.HandleFunc("/theme", handler.Repo.ThemeHandler)
	return mux
}
