package repository

import (
	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindAllDigimonCards(c *gin.Context) ([]models.DigimonCard, error) {
	// Implementation goes here
	cards, err := gorm.G[models.DigimonCard](DB).Find(c)
	return cards, err
}
