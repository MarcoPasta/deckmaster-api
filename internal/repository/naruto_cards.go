package repository

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"github.com/MarcoPasta/deckmaster-api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindAllNarutoCards(ctx context.Context) ([]models.NarutoCard, error) {
	return gorm.G[models.NarutoCard](DB).Order("id asc").Find(ctx)
}

func FindNarutoCardsByNames(ctx context.Context, decklist []utils.ImportNarutoDeck) ([]models.NarutoCard, error) {
	var cards []models.NarutoCard
	for _, item := range decklist {
		var card models.NarutoCard
		if item.Set == "" || item.Number == "" || item.Rarity == "" {
			continue
		}

		card, err := gorm.G[models.NarutoCard](DB).Where("card_number = ?", item.Number).Where("rarity = ?", item.Rarity).First(ctx)
		if err != nil {
			log.Printf("Error finding card with number %s and rarity %s: %v\n", item.Number, item.Rarity, err)
			return nil, err
		}
		switch item.Count {
		case 1:
			cards = append(cards, card)
		case 2:
			cards = append(cards, card, card)
		default:
			log.Printf("Unexpected count %d for card %s\n", card)
			cards = append(cards, card) // Default to adding one copy if count is unexpected
		}
	}
	return cards, nil
}

func GetNarutoCardByID(ctx *gin.Context, id string) (*models.NarutoCard, error) {
	card, err := gorm.G[models.NarutoCard](DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func PatchNarutoCardByID(ctx *gin.Context, id int, body models.NarutoCard) (int, error) {

	log.Printf("Updating card with ID %d: %+v\n", id, body)
	rows, err := gorm.G[models.NarutoCard](DB).Where("id = ?", id).Updates(ctx, body)
	if err != nil {
		log.Printf("Error updating card with ID %d: %v\n", id, err)
		return 0, err
	}

	return rows, nil
}

func DeleteNarutoCardByID(ctx *gin.Context, id string) (int, error) {
	rows, err := gorm.G[models.NarutoCard](DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		log.Printf("Error deleting card with ID %s: %v\n", id, err)
		return 0, err
	}

	return rows, nil
}

func AddNarutoCardsBySets(ctx *gin.Context, files []*multipart.FileHeader) error {
	for _, file := range files {
		cards, err := utils.DecodeCards(ctx, *file, models.NarutoCard{})
		if err != nil {
			log.Printf("Error decoding cards from file: %s, error: %v\n", file.Filename, err)
			return err
		}
		log.Printf("Decoded cards: %+v\n", cards)

		for _, card := range cards {
			log.Printf("Saving card to database: %s (ID: %s)\n", card.(models.NarutoCard).Name, card.(models.NarutoCard).ID)
			save := card.(models.NarutoCard)
			err := gorm.G[models.NarutoCard](DB).Create(ctx, &save)
			if err != nil {
				log.Printf("Error saving Naruto card to database: %v, error: %v\n", card.(models.NarutoCard).ID, err)
				return err
			}
		}
	}
	return nil
}
