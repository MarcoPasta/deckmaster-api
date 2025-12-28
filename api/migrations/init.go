package migrations

import (
	"context"
	"log"

	"github.com/MarcoPasta/deckmaster/api/models"
	"gorm.io/gorm"
)

func init() {
	models.ConnectToDatabase()

	err := models.DB.AutoMigrate(&models.DigimonDecks{}, &models.DigimonCard{})
	if err != nil {
		log.Fatal(err)
	}

	card := models.DigimonCard{
		Name:  "Agumon",
		Color: "Red",
	}

	deck := models.DigimonDecks{
		Name:  "Mein Deck",
		Cards: []models.DigimonCard{card},
	}

	ctx := context.Background()
	err = gorm.G[models.DigimonDecks](models.DB).Create(ctx, &deck)
}
