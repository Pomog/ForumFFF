package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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

	sendEmail()

	db, err := run()
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

	app.InfoLog.Println("Trying to get DataBase connection")
	defer db.SQL.Close()

	app.InfoLog.Printf("Server starting on port %s\n", Port)

	srv := &http.Server{
		Addr:    Port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func sendEmail() {
	// test mail

	from := "ffforumadm@gmail.com"
	password := "wwop dffy xhnz rdbv"
	to := "denver1033@gmail.com"
	subject := "Test Email"
	body := "Hello, this is a test email from Golang."

	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "Test" + "\r\n" +
		"\r\n" + body

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Email sent successfully.")

	}
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

	// info log
	infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog

	// error log
	errorlog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog

	db, err := repository.GetDB()
	if err != nil {
		log.Fatal("cannot get database connection")
		return nil, err
	}

	err = repository.MakeDBTables(db.SQL)
	if err != nil {
		log.Fatal("cannot create database tables")
		return nil, err
	}

	repo := handler.NewRepo(&app, db)

	tc, err := renderer.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	handler.NewHandlers(repo)

	renderer.NewTemplate(&app)

	return db, nil
}
