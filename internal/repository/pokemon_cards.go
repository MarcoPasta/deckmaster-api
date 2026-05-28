package repository

import (
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FindAllPokemonCards retrieves all Pokemon cards from the database
func FindAllPokemonCards(ctx *gin.Context) ([]models.PokemonCard, error) {
	return gorm.G[models.PokemonCard](DB).Where("language = ?", "de").Find(ctx)
}

// PokemonCardByID
func GetPokemonCardByID(ctx *gin.Context, id string) (*models.PokemonCard, error) {
	card, err := gorm.G[models.PokemonCard](DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func PatchPokemonCardByID(ctx *gin.Context, body models.PokemonCard) (int, error) {
	rows, err := gorm.G[models.PokemonCard](DB).Where("id = ?", body.ID).Updates(ctx, body)
	if err != nil {
		log.Printf("Error updating card with ID %d: %v\n", body.ID, err)
		return 0, err
	}

	return rows, nil
}

func DeletePokemonCardByID(ctx *gin.Context, id string) (int, error) {
	rows, err := gorm.G[models.PokemonCard](DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		log.Printf("Error deleting card with ID %s: %v\n", id, err)
		return 0, err
	}

	return rows, nil
}

// PokemonCardsByName
func FindSamePokemonCardsByNames(ctx *gin.Context, parsedList []utils.ImportDeck) ([]models.PokemonCard, error) {

	var cards []models.PokemonCard
	var queue []utils.ImportDeck

	for _, item := range parsedList {

		log.Println("Searching Card: ", item)

		card, err := gorm.G[models.PokemonCard](DB).Where("card_number like ? AND code = ?", item.Number+"%", strings.ToLower(item.Code)).First(ctx)
		if err != nil {
			log.Printf("Error fetching card: %s, appending to queue", item.Name)
			queue = append(queue, item)
			continue
		}
		for i := 0; i < item.Count; i++ {
			cards = append(cards, card)
		}
	}
	log.Println("Cards in queue: ", queue)
	for _, item := range queue {
		log.Printf("Attempting to find card by name: %s\n", item.Name)
		card, err := gorm.G[models.PokemonCard](DB).Where("name = ?", item.Name).First(ctx)
		if err != nil {
			log.Printf("Error fetching card by name: %s, error: %v\n", item.Name, err)
			continue
		}
		for i := 0; i < item.Count; i++ {
			cards = append(cards, card)
		}
	}

	return cards, nil
}

// FindPokemonCardsByEnergyType retrieves all Pokemon cards that have any of the specified energy types
func FindPokemonCardsByEnergyType(ctx *gin.Context, energyTypes []string) ([]models.PokemonCard, error) {
	var cards []models.PokemonCard
	query := DB.WithContext(ctx)

	log.Printf("Querying by energy types: %+v\n", energyTypes)

	for _, energyType := range energyTypes {
		query = query.Where("energy_type LIKE ?", fmt.Sprintf("%%%v%%", energyType))
	}

	if err := query.Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func AddPokemonCardsBySets(ctx *gin.Context, files []*multipart.FileHeader) error {
	for _, file := range files {
		cards, err := utils.DecodeCards(ctx, *file, models.PokemonCard{})
		if err != nil {
			log.Printf("Error decoding cards from file: %s, error: %v\n", file.Filename, err)
			return err
		}
		log.Printf("Decoded cards: %+v\n", cards)

		for _, card := range cards {
			log.Printf("Saving card to database: %s (ID: %s)\n", card.(models.PokemonCard).Name, card.(models.PokemonCard).ID)
			save := card.(models.PokemonCard)
			err := gorm.G[models.PokemonCard](DB).Create(ctx, &save)
			if err != nil {
				log.Printf("Error saving Pokemon card to database: %v, error: %v\n", card.(models.PokemonCard).ID, err)
				return err
			}
		}
	}
	return nil
}
