package main

import (
	"api/config"
	"api/server"
	"api/server/models"

	"log"
	"net/http"
)

// TODO: Restructure the whole project, it sucks right now

func main() {
	var err error
	log.Println("Initializing config structure")
	err = config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing database connections pool")
	err = models.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting local server")
	err = http.ListenAndServe(config.Port, server.New())
	if err != nil {
		log.Fatal(err)
	}
}
