package main

import (
	"fmt"
	"net/http"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/handler"
)

// routes returns an HTTP handler that routes requests to the appropriate handlers.
func routes(a *config.AppConfig) http.Handler {

	mux := http.NewServeMux()

	registerStaticHandlers(mux)

	// Register handlers for application-specific routes.
	mux.HandleFunc("/", handler.Repo.LoginHandler)
	mux.HandleFunc("/home", handler.Repo.HomeHandler)
	mux.HandleFunc("/theme", handler.Repo.ThemeHandler)
	mux.HandleFunc("/registration", handler.Repo.RegisterHandler)
	mux.HandleFunc("/error-page", handler.Repo.ErrorPage)
	mux.HandleFunc("/contactUs", handler.Repo.ContactUsHandler)
	mux.HandleFunc("/forum_rules", handler.Repo.ForumRulesHandler)
	mux.HandleFunc("/help", handler.Repo.HelpHandler)
	mux.HandleFunc("/privat_policy", handler.Repo.PrivatPolicyHandler)
	mux.HandleFunc("/personal_cabinet", handler.Repo.PersonaCabinetHandler)
	mux.HandleFunc("/edit_topic", handler.Repo.EditTopicHandler)
	mux.HandleFunc("/edit_topic_result", handler.Repo.EditTopicResultHandler)

	return mux
}

// registerStaticHandlers registers handlers for static content.
func registerStaticHandlers(mux *http.ServeMux) {
	//Define a list of static directories (e.g., "css", "logo", "ava").
	statics := []string{"css", "logo", "ava"}

	// Register handlers for static content.
	for _, static := range statics {
		mux.Handle(fmt.Sprintf("/static/%s/", static), http.StripPrefix(fmt.Sprintf("/static/%s/", static), http.FileServer(http.Dir(fmt.Sprintf("static/%s", static)))))
	}
}
