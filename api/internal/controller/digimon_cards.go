package controller

import (
	"log"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/gin-gonic/gin"
)

type Filter struct {
	Colors []string `form:"color[]" collection_format:"csv"`
}

func FindAllDigimonCards(c *gin.Context) {

	cards, err := repository.FindAllDigimonCards(c)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
	}

	c.JSON(http.StatusOK, gin.H{"content": cards})
}

func FindDigimonCardByID(c *gin.Context) {
	id := c.Param("id")

	log.Printf("Fetching card with ID: %s", id)

	card, err := repository.FindDigimonCardByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": card})
}

func FindDigimonCardsByName(c *gin.Context) {
	name := c.Param("name")

	cards, err := repository.FindDigimonCardsByName(c, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": cards})
}

func FindDigimonCardsByColor(c *gin.Context) {
	var f Filter

	log.Printf("Filtering by colors: %+v\n", f)
	if err := c.ShouldBind(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters.",
		})
		return
	}

	log.Printf("Filtering by colors: %+v\n", f.Colors)
	cards, err := repository.FindDigimonCardsByColor(c, f.Colors)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error 404": "Cards not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": cards})
}
