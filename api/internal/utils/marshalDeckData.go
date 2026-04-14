package utils

import "github.com/MarcoPasta/deckmaster/api/internal/models"

// mashalling types
type CreateDeckRequest struct {
	UserId string               `json:"userId" binding:"required"`
	Name   string               `json:"name" binding:"required"`
	Cards  []models.DigimonCard `json:"cards" binding:"required"`
}

// marshalling functions
func MarshalDigimonDeckRequest(body CreateDeckRequest) models.DigimonDecks {
	return models.DigimonDecks{
		UserId: body.UserId,
		Name:   body.Name,
		Cards:  body.Cards,
	}
}
