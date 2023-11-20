package main

import (
	"net/http"

	"github.com/Pomog/ForumFFF/pkg/config"
	"github.com/Pomog/ForumFFF/pkg/handler"
)

func routes(a *config.AppConfig) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Repo.MainHandler)
	mux.HandleFunc("/about", handler.Repo.AboutHandler)
	return mux
}
