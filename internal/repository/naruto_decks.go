package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"gorm.io/gorm"
)

func CreateNarutoDeck(ctx context.Context, deck models.NarutoDecks) (models.NarutoDecks, error) {

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
		DB.Model(&models.NarutoCard{}).Where("id IN ?", cardIds).Count(&count)
		if count != int64(len(cardIds)) {
			return models.NarutoDecks{}, fmt.Errorf("one or more card IDs do not exist in the database")
		}
	}

	err := DB.Omit("Entries.Card").Create(&deck).Error
	if err != nil {
		return models.NarutoDecks{}, err
	}
	return deck, nil
}

func GetNarutoDecksByUserId(ctx context.Context, userId string) ([]models.NarutoDecks, error) {
	var decks []models.NarutoDecks
	err := DB.Preload("Entries.Card").Where("user_id = ?", userId).Find(&decks).Error
	if err != nil {
		return nil, err
	}
	return decks, nil
}

func UpdateNarutoDeck(ctx context.Context, id string, updatedDeck models.NarutoDecks) (int, error) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// Delete existing entries for this deck
		if err := tx.Where("deck_id = ?", id).Delete(&models.NarutoDeckEntry{}).Error; err != nil {
			return err
		}

		// Update the deck's own fields
		if err := tx.Model(&models.NarutoDecks{}).Where("id = ?", id).Updates(map[string]interface{}{
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

func UpdateNarutoDeckImage(ctx context.Context, id string, cardId string) error {
	deckId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid deck id: %w", err)
	}

	err = DB.Model(&models.NarutoDecks{}).Where("id = ?", deckId).Update("cover_card", cardId).Error
	if err != nil {
		return fmt.Errorf("failed to update deck image: %w", err)
	}

	return nil
}

func DeleteNarutoDeck(ctx context.Context, id int) (int, error) {
	rows, err := gorm.G[models.NarutoDecks](DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete deck: %w", err)
	}
	fmt.Printf("Deleted %d rows\n", rows)
	return rows, nil
}
