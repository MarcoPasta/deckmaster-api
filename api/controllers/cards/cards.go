package cards

import (
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Find(c *gin.Context) {
	cards, err := gorm.G[models.DigimonCard](models.DB).Find(c)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
	}

	c.JSON(http.StatusOK, gin.H{"content": cards})
}
