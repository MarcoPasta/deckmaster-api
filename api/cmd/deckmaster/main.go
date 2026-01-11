package main

import (
	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	_ "github.com/MarcoPasta/deckmaster/api/migrations"
	_ "github.com/lib/pq"
)

func main() {
	repository.ConnectToDatabase(internal.GetDBString())
	internal.SetupRouter()

	// Start server on port 8080 (default)
	internal.Router.Run()
}
