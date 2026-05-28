package repository

import (
	"context"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
)

func GetPokemonDeckById(ctx context.Context, deckId uint) (models.PokemonDecks, error) {
	var deck models.PokemonDecks
	err := DB.Preload("Entries.Card").First(&deck, deckId).Error
	return deck, err
}

func GetNarutoDeckById(ctx context.Context, deckId uint) (models.NarutoDecks, error) {
	var deck models.NarutoDecks
	err := DB.Preload("Entries.Card").First(&deck, deckId).Error
	return deck, err
}
