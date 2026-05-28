package repository

import (
	"context"
	"fmt"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
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
