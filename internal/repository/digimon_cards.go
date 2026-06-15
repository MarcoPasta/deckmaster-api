package repository

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FindAllDigimonCards retrieves all Digimon cards from the database
func FindAllDigimonCards(ctx *gin.Context) ([]models.DigimonCard, error) {
	return gorm.G[models.DigimonCard](DB).Order("id asc").Find(ctx)
}

// FindDigimonCardByID retrieves a single Digimon card by its ID
func FindDigimonCardByID(ctx *gin.Context, id string) (*models.DigimonCard, error) {
	card, err := gorm.G[models.DigimonCard](DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func PatchDigimonCardByID(ctx *gin.Context, id int, body models.DigimonCard) (int, error) {

	log.Printf("Updating card with ID %d: %+v\n", id, body)
	rows, err := gorm.G[models.DigimonCard](DB).Where("id = ?", id).Updates(ctx, body)
	if err != nil {
		log.Printf("Error updating card with ID %d: %v\n", id, err)
		return 0, err
	}

	return rows, nil
}

func DeleteDigimonCardByID(ctx *gin.Context, id string) (int, error) {
	rows, err := gorm.G[models.DigimonCard](DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		log.Printf("Error deleting card with ID %s: %v\n", id, err)
		return 0, err
	}

	return rows, nil
}

///////////////////////////////////////////////////////////////////////////////

func AddDigimonCardsBySets(ctx *gin.Context, files []*multipart.FileHeader) error {
	for _, file := range files {
		cards, err := utils.DecodeCards(ctx, *file, models.DigimonCard{})
		if err != nil {
			log.Printf("Error decoding cards from file: %s, error: %v\n", file.Filename, err)
			return err
		}
		log.Printf("Decoded cards: %+v\n", cards)

		for _, card := range cards {
			log.Printf("Saving card to database: %s (ID: %s)\n", card.(models.DigimonCard).Name, card.(models.DigimonCard).ID)
			save := card.(models.DigimonCard)
			err := gorm.G[models.DigimonCard](DB).Create(ctx, &save)
			if err != nil {
				log.Printf("Error saving Digimon card to database: %v, error: %v\n", card.(models.DigimonCard).ID, err)
				return err
			}
		}
	}
	return nil
}

// FindDigimonCardsByName retrieves all Digimon cards that match the given name
// func FindDigimonCardsByName(ctx *gin.Context, name string) ([]models.DigimonCard, error) {
// 	return gorm.G[models.DigimonCard](DB).Where("name ILIKE ?", "%"+name+"%").Find(ctx)
// }

// FindDigimonCardsByColor retrieves all Digimon cards that have any of the specified colors
func FindDigimonCardsByColor(ctx *gin.Context, colors []string) ([]models.DigimonCard, error) {
	var cards []models.DigimonCard
	query := DB.WithContext(ctx)

	log.Printf("Querying by colors: %+v\n", colors)

	// Filter by colors - using containment operator for PostgreSQL array columns
	for _, color := range colors {
		query = query.Where("color LIKE ?", fmt.Sprintf("%%%v%%", color))
	}

	if err := query.Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}
