package repository

import (
	"fmt"
	"log"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FindAllPokemonCards retrieves all Pokemon cards from the database
func FindAllPokemonCards(ctx *gin.Context) ([]models.PokemonCard, error) {
	return gorm.G[models.PokemonCard](DB).Find(ctx)
}

// FindPokemonCardByID retrieves a single Pokemon card by its ID
func FindPokemonCardByID(ctx *gin.Context, id string) (*models.PokemonCard, error) {
	card, err := gorm.G[models.PokemonCard](DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// FindPokemonCardsByName retrieves all Pokemon cards that match the given name
func FindPokemonCardsByName(ctx *gin.Context, name string) ([]models.PokemonCard, error) {
	return gorm.G[models.PokemonCard](DB).Where("name ILIKE ?", "%"+name+"%").Find(ctx)
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
