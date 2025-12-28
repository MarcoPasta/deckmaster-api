package main

import (
	"github.com/MarcoPasta/deckmaster/api/models"
	"github.com/MarcoPasta/deckmaster/api/routes"
)

func main() {
	models.ConnectToDatabase()
	routes.SetupRouter()

	// Start server on port 8080 (default)
	routes.Router.Run()
}
