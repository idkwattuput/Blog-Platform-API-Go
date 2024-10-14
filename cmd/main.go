package main

import (
	"database/sql"
	"log"

	"github.com/idkwattuput/blogging-platform-api-go/cmd/api"
	"github.com/idkwattuput/blogging-platform-api-go/config"
	"github.com/idkwattuput/blogging-platform-api-go/db"
)

func main() {
	cnnString := config.Envs.DBUrl

	db, err := db.NewPostgreSQLStorage(cnnString)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
