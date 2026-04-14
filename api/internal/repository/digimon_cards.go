package repository

import (
	"fmt"
	"log"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FindAllDigimonCards retrieves all Digimon cards from the database
func FindAllDigimonCards(ctx *gin.Context) ([]models.DigimonCard, error) {
	return gorm.G[models.DigimonCard](DB).Find(ctx)
}

// FindDigimonCardByID retrieves a single Digimon card by its ID
func FindDigimonCardByID(ctx *gin.Context, id string) (*models.DigimonCard, error) {
	card, err := gorm.G[models.DigimonCard](DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// FindDigimonCardsByName retrieves all Digimon cards that match the given name
func FindDigimonCardsByName(ctx *gin.Context, name string) ([]models.DigimonCard, error) {
	return gorm.G[models.DigimonCard](DB).Where("name ILIKE ?", "%"+name+"%").Find(ctx)
}

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
