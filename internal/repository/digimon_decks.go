package repository

import (
	"fmt"
	"strconv"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateDigimonDeck creates a new Digimon deck in the database
func CreateDigimonDeck(ctx *gin.Context, deck models.DigimonDecks) (models.DigimonDecks, error) {

	// Validate that all card IDs exist
	if len(deck.Entries) > 0 {
		uniqueIds := make(map[uint]struct{})
		for _, e := range deck.Entries {
			uniqueIds[e.CardId] = struct{}{}
		}
		cardIds := make([]uint, 0, len(uniqueIds))
		for id := range uniqueIds {
			cardIds = append(cardIds, id)
		}
		var count int64
		DB.Model(&models.DigimonCard{}).Where("id IN ?", cardIds).Count(&count)
		if count != int64(len(cardIds)) {
			return models.DigimonDecks{}, fmt.Errorf("one or more card IDs do not exist in the database")
		}
	}

	err := DB.Omit("Entries.Card").Create(&deck).Error
	if err != nil {
		return models.DigimonDecks{}, err
	}
	return deck, nil
}

// GetDigimonDeckByUserId retrieves all Digimon decks for a given user ID
func GetDigimonDeckByUserId(ctx *gin.Context, userId string) ([]models.DigimonDecks, error) {
	var decks []models.DigimonDecks
	err := DB.Preload("Entries.Card").Where("user_id = ?", userId).Find(&decks).Error
	if err != nil {
		return nil, err
	}
	return decks, nil
}

func UpdateDigimonDeck(ctx *gin.Context, id string, updatedDeck models.DigimonDecks) (int, error) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// Delete existing entries for this deck
		if err := tx.Where("deck_id = ?", id).Delete(&models.DigimonDeckEntry{}).Error; err != nil {
			return err
		}

		// Update the deck's own fields
		if err := tx.Model(&models.DigimonDecks{}).Where("id = ?", id).Updates(map[string]interface{}{
			"name": updatedDeck.Name,
		}).Error; err != nil {
			return err
		}

		// Create the new entries with the correct deck ID
		if len(updatedDeck.Entries) > 0 {
			deckId, _ := strconv.ParseUint(id, 10, 64)
			for i := range updatedDeck.Entries {
				updatedDeck.Entries[i].DeckId = uint(deckId)
			}
			if err := tx.Omit("Card").Create(&updatedDeck.Entries).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to update deck: %w", err)
	}
	return len(updatedDeck.Entries), nil
}

func UpdateDigimonDeckImage(ctx *gin.Context, id string, cardId string) error {
	deckId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid deck id: %w", err)
	}

	err = DB.Model(&models.DigimonDecks{}).Where("id = ?", deckId).Update("cover_card", cardId).Error
	if err != nil {
		return fmt.Errorf("failed to update deck image: %w", err)
	}

	return nil
}

func DeleteDigimonDeck(ctx *gin.Context, id string) (int, error) {
	rows, err := gorm.G[models.DigimonDecks](DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete deck: %w", err)
	}
	fmt.Printf("Deleted %d rows\n", rows)
	return rows, nil
}
