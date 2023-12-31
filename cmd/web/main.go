package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Pomog/ForumFFF/internal/models"

	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/handler"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/Pomog/ForumFFF/internal/repository"
)

const Port = ":8080"

var app config.AppConfig
var infolog *log.Logger
var errorlog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

	app.InfoLog.Println("Trying to get DataBase connection")
	defer db.SQL.Close()

	app.InfoLog.Printf("Server starting on port %s\n", Port)

	app.PostLen = 2500    //post and topic size
	app.CategoryLen = 100 //category size

	app.FileSize = 2 //here we set 2mb of file size
	app.NameMinLen = 3
	app.NameMaxLen = 12
	app.PasswordMinLen = 6
	app.PasswordMaxLen = 20
	app.MinCategoryLen = 3
	app.MaxCategoryLen = 2500
	app.MinSubjectLen = 3
	app.MaxSubjectLen = 2500
	app.LongestSingleWord = "pneumonoultramicroscopicsilicovolcanoconiosis"
	app.GitHubClientID = "b383a969da8082007c4c"
	app.GitHubClientSecret = "2327f7f8677894e2726f70696d9321399efb6b61"
	app.GitHubRedirectURL = "http://localhost:8080/github-callback"

	app.GoogleClientID = "1045980296590-9c2dpdmdbp17kiih3oev5nsk55a33q29.apps.googleusercontent.com"
	app.GoogleClientSecret = "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"
	app.GoogleRedirectURL = "http://localhost:8080/google-callback"

	app.ModeratorPass = "123" //our secret pass for moders

	//the list of games that are represented and will be covered on site.
	app.GamesList = map[string]string{
		"Lineage 2":                "L2",
		"Lineage 3":                "L3",
		"World Of Warcraft":        "WOW",
		"New World":                "NW",
		"AION":                     "AION",
		"Conan exiles":             "Conan",
		"Guild Wars 2":             "GW2",
		"Archeage":                 "Archeage",
		"Black Desert Online":      "BDO",
		"World of Tanks":           "WOT",
		"EVE Online":               "EVE",
		"The Elder Scrolls Online": "ESO",
		"Throne and Liberty":       "TL",
	}

	srv := &http.Server{
		Addr:    Port,
		Handler: routes(&app),
	}

	// helper.SendEmail(app.ServerEmail, "The FFForum Server started at")
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*repository.DataBase, error) {
	fmt.Println("Starting application")

	// gob.Register() function is used to inform the encoding/gob package about custom types that may be encoded or decoded using the gob encoding format.
	gob.Register(models.User{})
	gob.Register(models.Thread{})
	gob.Register(models.Votes{})
	gob.Register(models.Post{})

	// change this to true when in production
	app.InProduction = false

	// set email adress for loggin
	app.ServerEmail = "ffforumadm@gmail.com"

	// info log
	infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog

	// error log
	errorlog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog

	db, err := repository.GetDB()
	if err != nil {
		log.Fatal("cannot get database connection :" + err.Error())
		return nil, err
	}

	err = repository.MakeDBTables(db.SQL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	repo := handler.NewRepo(&app, db)

	tc, err := renderer.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache : " + err.Error())
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	handler.NewHandlers(repo)

	renderer.NewTemplate(&app)

	return db, nil
}
