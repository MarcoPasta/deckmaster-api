package repository

import (
	"fmt"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateDigimonDeck creates a new Digimon deck in the database
func CreateDigimonDeck(ctx *gin.Context, deck models.DigimonDecks) (models.DigimonDecks, error) {

	fmt.Println("Received Deck: ", deck)

	err := gorm.G[models.DigimonDecks](DB).Create(ctx, &deck)
	if err != nil {
		return models.DigimonDecks{}, err
	}
	return deck, nil
}

// GetDeckByUserId retrieves all Digimon decks for a given user ID
func GetDeckByUserId(ctx *gin.Context, userId string) ([]models.DigimonDecks, error) {
	return gorm.G[models.DigimonDecks](DB).Where("user_id = ?", userId).Find(ctx)
}
