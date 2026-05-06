package main

import (
	"log"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	_ "github.com/MarcoPasta/deckmaster/api/internal/repository/migrations"
	_ "github.com/lib/pq"
)

func main() {
	repository.ConnectToDatabase(internal.GetDBString())
	repository.RunMigrations()
	if err := repository.CheckTableValue(); err != nil {
		log.Fatal(err)
	}
	internal.SetupRouter()

	// Start server on port 8080 (default)
	internal.Router.Run()
}
