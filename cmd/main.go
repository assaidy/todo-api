package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/assaidy/todo-api/repo"
	"github.com/assaidy/todo-api/router"
	"github.com/assaidy/todo-api/utils"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dbConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	pg, err := repo.NewPostgresDB(dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := repo.Migrate(pg.DB); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	defer pg.DB.Close()

	router := router.NewRouter(pg)

	log.Println("Running server on port %s", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))
}
