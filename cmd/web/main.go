package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/Pomog/ForumFFF/db_driver"
	"github.com/Pomog/ForumFFF/internal/config"
	"github.com/Pomog/ForumFFF/internal/handler"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
	"github.com/Pomog/ForumFFF/internal/repository/dbrepo"
	"github.com/google/uuid"
	"github.com/Pomog/ForumFFF/internal/repository/dbrepo"
)

const Port = ":8080"

var app config.AppConfig
var sessionID uuid.UUID

func main() {

	// gob.Register() function is used to inform the encoding/gob package about custom types that may be encoded or decoded using the gob encoding format.
	gob.Register(models.User{})
	gob.Register(models.Thread{})
	gob.Register(models.Votes{})
	gob.Register(models.Post{})

	sessionID = uuid.New()

	app.UserLogin = &sessionID

	newUser := models.User{
		UserName:  "test1",
		Password:  "123",
		FirstName: "testFirstName1",
		LastName:  "testLastName1",
		Email:     "test1@mail.com",
	}

	db, _ := db_driver.GetDB()

	repo := dbrepo.NewSQLiteRepo(db, &app)

	repo.CreatetUser(newUser)

	userNameToFind := "test1"
	userNameToFind2 := "noSuchUser"
	userEmailToFind := "test1@mail.com"

	isPresent, _ := repo.UserPresent(userNameToFind, userEmailToFind)
	log.Printf("User %s present - %v", userNameToFind, isPresent)
	fmt.Println("--------------------------------------------------")
	isPresent2, _ := repo.UserPresent(userNameToFind2, userEmailToFind)
	log.Printf("User %s present - %v", userNameToFind2, isPresent2)

	// -------------------------------------------------------------

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server starting on port %s\n", Port)

	srv := &http.Server{
		Addr:    Port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func createDB() {
	db_driver.MakeDBTables()

	newUser := models.User{
		UserName:  "test1",
		Password:  "123",
		FirstName: "testFirstName1",
		LastName:  "testLastName1",
		Email:     "test1@mail.com",
	}

	db, _ := db_driver.GetDB()

	repo := dbrepo.NewSQLiteRepo(db, &app)

	repo.CreatetUser(newUser)

	userNameToFind := "test1"
	userNameToFind2 := "noSuchUser"
	userEmailToFind := "test1@mail.com"

	isPresent, _ := repo.UserPresent(userNameToFind, userEmailToFind)
	log.Printf("User %s present - %v", userNameToFind, isPresent)
	fmt.Println("--------------------------------------------------")
	isPresent2, _ := repo.UserPresent(userNameToFind2, userEmailToFind)
	log.Printf("User %s present - %v", userNameToFind2, isPresent2)

}

func run() error {
	fmt.Println("Starting application")

	// change this to true when in production
	app.InProduction = false

	//cookies should be set

	tc, err := renderer.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handler.NewRepo(&app)
	handler.NewHandlers(repo)

	renderer.NewTemplate(&app)

	return nil
}
