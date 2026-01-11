package controller

import (
	"log"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/gin-gonic/gin"
)

type Filter struct {
	Colors []string `form:"colors[]" collection_format:"csv"`
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
	log.Println("ID from context: ", id)

}

func FindDigimonCardsByName(c *gin.Context) {
	name := c.Param("name")
	log.Println("Name from context: ", name)
}

func FindDigimonCardsByColor(c *gin.Context) {
	var f Filter
	if err := c.ShouldBind(&f); err != nil {
		log.Println("Error binding query parameters")
	} else {
		log.Println("Colors from query: ", f.Colors)
	}
}
