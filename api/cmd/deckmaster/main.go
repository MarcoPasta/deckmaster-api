package main

import (
	"log"
	"time"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/seeder"
	_ "github.com/MarcoPasta/deckmaster/api/migrations"
)

func main() {
	repository.ConnectToDatabase(internal.GetDBString())
	repository.RunMigrations()
	for i := range 3 {
		err := repository.CheckTableValue()
		if err == nil {
			break
		}
		log.Printf("run %d/3: seeding %v\n", i+1, err)
		if seedErr := seeder.Main("mode", err.Error()); seedErr != nil {
			log.Printf("seeder failed: %v, retrying in 5s...\n", seedErr)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("run %d/3: seeder finished successfully\n", i+1)
	}
	internal.SetupRouter()

	// Start server on port 8080 (default)
	internal.Router.Run()
}
