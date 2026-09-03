package utils

import (
	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"gorm.io/gorm"
)

// marshalling types
type CreateDigimonDeckRequest struct {
	UserId    string                    `json:"userId" binding:"required"`
	Name      string                    `json:"name" binding:"required"`
	CoverCard string                    `json:"coverCard" binding:"required"`
	Entries   []models.DigimonDeckEntry `json:"entries" binding:"required"`
}

type CreatePokemonDeckRequest struct {
	UserId    string                    `json:"userId" binding:"required"`
	Name      string                    `json:"name" binding:"required"`
	CoverCard string                    `json:"coverCard" binding:"required"`
	Entries   []models.PokemonDeckEntry `json:"entries" binding:"required"`
}

type CreateNarutoDeckRequest struct {
	UserId    string                   `json:"userId" binding:"required"`
	Name      string                   `json:"name" binding:"required"`
	CoverCard string                   `json:"coverCard" binding:"required"`
	Entries   []models.NarutoDeckEntry `json:"entries" binding:"required"`
}

type UpdateDigimonDeckRequest struct {
	DeckId    uint                      `json:"deckId" binding:"required"`
	Name      string                    `json:"name" binding:"required"`
	CoverCard string                    `json:"coverCard" binding:"required"`
	Entries   []models.DigimonDeckEntry `json:"entries" binding:"required"`
}

type UpdatePokemonDeckRequest struct {
	DeckId    uint                      `json:"deckId" binding:"required"`
	Name      string                    `json:"name" binding:"required"`
	CoverCard string                    `json:"coverCard" binding:"required"`
	Entries   []models.PokemonDeckEntry `json:"entries" binding:"required"`
}

type UpdateNarutoDeckRequest struct {
	DeckId    uint                     `json:"deckId" binding:"required"`
	Name      string                   `json:"name" binding:"required"`
	CoverCard string                   `json:"coverCard" binding:"required"`
	Entries   []models.NarutoDeckEntry `json:"entries" binding:"required"`
}

type GetPokemonDeckCardMarketListResponse struct {
	Id         uint   `json:"id"`
	Code       string `json:"code"`
	CardNumber string `json:"card_number"`
}

type CardMarketRequest struct {
	Cards []GetPokemonDeckCardMarketListResponse `json:"cards"`
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

func MarshalUpdateDigimonDeckRequest(body UpdateDigimonDeckRequest) models.DigimonDecks {
	return models.DigimonDecks{
		Model:     gorm.Model{ID: uint(body.DeckId)},
		Name:      body.Name,
		CoverCard: body.CoverCard,
		Entries:   body.Entries,
	}
}

func MarshalUpdatePokemonDeckRequest(body UpdatePokemonDeckRequest) models.PokemonDecks {
	return models.PokemonDecks{
		Model:     gorm.Model{ID: uint(body.DeckId)},
		Name:      body.Name,
		CoverCard: body.CoverCard,
		Entries:   body.Entries,
	}
}

func MarshalUpdateNarutoDeckRequest(body UpdateNarutoDeckRequest) models.NarutoDecks {
	return models.NarutoDecks{
		Model:     gorm.Model{ID: uint(body.DeckId)},
		Name:      body.Name,
		CoverCard: body.CoverCard,
		Entries:   body.Entries,
	}
}

func MarshalGetPokemonDeckCardMarketListResponse(cards CardMarketRequest) []GetPokemonDeckCardMarketListResponse {
	var response []GetPokemonDeckCardMarketListResponse
	for _, card := range cards.Cards {
		response = append(response, GetPokemonDeckCardMarketListResponse{
			Id:         card.Id,
			Code:       card.Code,
			CardNumber: card.CardNumber,
		})
	}
	return response
}
