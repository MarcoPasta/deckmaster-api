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
	"gorm.io/gorm/clause"
)

// FindAllPokemonCards retrieves all Pokemon cards from the database
func FindAllPokemonCards(ctx *gin.Context) ([]models.PokemonCard, error) {
	return gorm.G[models.PokemonCard](DB).Where("language = ?", "de").Order(utils.PokemonSetOrderCase()).Order("CAST(NULLIF(SPLIT_PART(card_number, '/', 1), '') AS INTEGER) asc").Find(ctx)
}

// PokemonCardByID
func GetPokemonCardByID(ctx *gin.Context, id string) (*models.PokemonCard, error) {
	card, err := gorm.G[models.PokemonCard](DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &card, nil
}

func PatchPokemonCardByID(ctx *gin.Context, id int, body models.PokemonCard) (int, error) {

	log.Printf("Updating card with ID %d: %+v\n", id, body)
	rows, err := gorm.G[models.PokemonCard](DB).Where("id = ?", id).Updates(ctx, body)
	if err != nil {
		log.Printf("Error updating card with ID %d: %v\n", id, err)
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
func FindSamePokemonCardsByNames(ctx *gin.Context, parsedList []utils.ImportPokemonDeck) ([]models.PokemonCard, error) {

	var cards []models.PokemonCard
	var queue []utils.ImportPokemonDeck

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
			save := card.(models.PokemonCard)
			err := DB.WithContext(ctx).
				Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "code"},
						{Name: "card_number"},
						{Name: "language"},
					},
					DoUpdates: clause.AssignmentColumns([]string{
						"card_type", "trainer_type", "phase", "name", "kp",
						"energy_type", "ability", "attacks", "weakness",
						"resistance", "retreat_cost", "set", "tags", "updated_at",
					}),
				}).
				Create(&save).Error
			if err != nil {
				log.Printf("Error upserting Pokemon card to database: %s, error: %v\n", save.Name, err)
				return err
			}
		}
	}
	return nil
}
