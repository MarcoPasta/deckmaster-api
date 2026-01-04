package main

import (
	_ "github.com/MarcoPasta/deckmaster/api/migrations"
	"github.com/MarcoPasta/deckmaster/api/models"
	"github.com/MarcoPasta/deckmaster/api/routes"
	_ "github.com/lib/pq"
)

func main() {
	models.ConnectToDatabase()
	routes.SetupRouter()

	// Start server on port 8080 (default)
	routes.Router.Run()
}
