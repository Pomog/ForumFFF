package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Pomog/ForumFFF/pkg/config"
	"github.com/Pomog/ForumFFF/pkg/handler"
	"github.com/Pomog/ForumFFF/pkg/renderer"
)

const Port = ":8080"

func main() {
	// In the main function, an instance of the config.AppConfig struct is created and stored in the app variable.
	var app config.AppConfig

	// The script attempts to create a template cache using renderer.CreateTemplateCache() and stores it in the app.TemplateCache.
	tc, err := renderer.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}
	app.TemplateCache = tc
	//this var to set use cache true or false, when in Dev mode
	app.UseCache = false

	// The handler.NewRepo(&app) call creates a new repository (*Repository)
	// and passes the app configuration to it. This repository is responsible for holding application configuration and handling requests.
	// The NewRepo function in the handler package is used to create a new repository.
	repo := handler.NewRepo(&app)

	// The handler.NewHandlers(repo) function is called to set the global Repo variable to
	// the newly created repository. This allows other parts of the code to access the repository.
	handler.NewHandlers(repo)

	//The renderer.NewTemplate(&app) function is called to set up
	// templates with the app configuration. This prepares the application to render HTML templates.
	renderer.NewTemplate(&app)
	// http.HandleFunc("/", handler.Repo.MainHandler)
	// http.HandleFunc("/about", handler.Repo.AboutHandler)
	fmt.Println("Server started on Port:", Port)

	// err = http.ListenAndServe(Port, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	srv := &http.Server{
		Addr:    Port,
		Handler: rounts(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
