package utils

import (
	"fmt"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
)

// marshalling types
type CreateDigimonDeckRequest struct {
	UserId  string                    `json:"userId" binding:"required"`
	Name    string                    `json:"name" binding:"required"`
	Entries []models.DigimonDeckEntry `json:"entries" binding:"required"`
}

type CreatePokemonDeckRequest struct {
	UserId  string                    `json:"userId" binding:"required"`
	Name    string                    `json:"name" binding:"required"`
	Entries []models.PokemonDeckEntry `json:"entries" binding:"required"`
}

type CreateNarutoDeckRequest struct {
	UserId  string                   `json:"userId" binding:"required"`
	Name    string                   `json:"name" binding:"required"`
	Entries []models.NarutoDeckEntry `json:"entries" binding:"required"`
}

// marshalling functions
func MarshalDigimonDeckRequest(body CreateDigimonDeckRequest) models.DigimonDecks {
	return models.DigimonDecks{
		UserId:  body.UserId,
		Name:    body.Name,
		Entries: body.Entries,
	}
}

func MarshalPokemonDeckRequest(body CreatePokemonDeckRequest) models.PokemonDecks {
	fmt.Printf("Request body for marshalling: %+v\n", body)
	return models.PokemonDecks{
		UserId:  body.UserId,
		Name:    body.Name,
		Entries: body.Entries,
	}
}

func MarshalNarutoDeckRequest(body CreateNarutoDeckRequest) models.NarutoDecks {
	return models.NarutoDecks{
		UserId:  body.UserId,
		Name:    body.Name,
		Entries: body.Entries,
	}
}
