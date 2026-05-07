package main

import (
	"log"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/seeder"
	_ "github.com/MarcoPasta/deckmaster/api/migrations"
)

func main() {
	repository.ConnectToDatabase(internal.GetDBString())
	repository.RunMigrations()
	if err := repository.CheckTableValue(); err != nil {
		log.Println("Error checking table values:", err)
		seeder.Main("mode", err.Error())
	}
	internal.SetupRouter()

	// Start server on port 8080 (default)
	internal.Router.Run()
}
