package main

import (
	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/storage"
	_ "github.com/MarcoPasta/deckmaster/api/migrations"
)

func main() {
	repository.ConnectToDatabase(internal.GetDBString())
	repository.RunMigrations()
	internal.SetupRouter(storage.CreateS3Client())

	// Start server on port 8080 (default)
	internal.Router.Run()
}
