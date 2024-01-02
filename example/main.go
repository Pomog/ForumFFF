package main

import (
	"log"

	"forum-authentication/config"
	"forum-authentication/server"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	_, err := config.InitializeDB()
	if err != nil {
		log.Println("Driver creation failed", err.Error())
	}

	config.Run()
}

func main() {

	server := server.NewServer(":8080")
	log.Fatal(server.Start())

}
