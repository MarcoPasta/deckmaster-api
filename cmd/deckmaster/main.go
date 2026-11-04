package main

import (
	"os"

	"github.com/MarcoPasta/deckmaster-api/internal"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	_ "github.com/MarcoPasta/deckmaster-api/migrations"
	"github.com/clerk/clerk-sdk-go/v2"
)

func main() {
	repository.ConnectToDatabase(internal.GetDBString())
	repository.RunMigrations()
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
	internal.SetupRouter(storage.CreateS3Client())

	// Start server on port 8080 (default)
	internal.Router.Run()
}
