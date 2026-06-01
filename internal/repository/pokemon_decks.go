package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"gorm.io/gorm"
)

func CreatePokemonDeck(ctx context.Context, deck models.PokemonDecks) (models.PokemonDecks, error) {

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
		DB.Model(&models.PokemonCard{}).Where("id IN ?", cardIds).Count(&count)
		if count != int64(len(cardIds)) {
			return models.PokemonDecks{}, fmt.Errorf("one or more card IDs do not exist in the database")
		}
	}

	err := DB.Omit("Entries.Card").Create(&deck).Error
	if err != nil {
		return models.PokemonDecks{}, err
	}
	return deck, nil
}

func GetPokemonDecksByUserId(ctx context.Context, userId string) ([]models.PokemonDecks, error) {
	var decks []models.PokemonDecks
	err := DB.Preload("Entries.Card").Where("user_id = ?", userId).Find(&decks).Error
	if err != nil {
		return nil, err
	}
	return decks, nil
}

func UpdatePokemonDeck(ctx context.Context, id string, updatedDeck models.PokemonDecks) (int, error) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// Delete existing entries for this deck
		if err := tx.Where("deck_id = ?", id).Delete(&models.PokemonDeckEntry{}).Error; err != nil {
			return err
		}

		// Update the deck's own fields
		if err := tx.Model(&models.PokemonDecks{}).Where("id = ?", id).Updates(map[string]interface{}{
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

func DeletePokemonDeck(ctx context.Context, id int) (int, error) {
	// err := DB.Transaction(func(tx *gorm.DB) error {
	// 	// Delete existing entries for this deck
	// 	if err := tx.Where("deck_id = ?", id).Delete(&models.PokemonDeckEntry{}).Error; err != nil {
	// 		return err
	// 	}

	// 	// Delete the deck itself
	// 	if err := tx.Where("id = ?", id).Delete(&models.PokemonDecks{}).Error; err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	rows, err := gorm.G[models.PokemonDecks](DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete deck: %w", err)
	}
	fmt.Printf("Deleted %d rows\n", rows)
	return rows, nil
}
