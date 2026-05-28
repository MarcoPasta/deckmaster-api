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

func FindAllDigimonCards(ctx *gin.Context) {

	cards, err := repository.FindAllDigimonCards(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

func FindDigimonCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	log.Printf("Fetching card with ID: %s", id)

	card, err := repository.FindDigimonCardByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": card})
}

func FindDigimonCardsByName(ctx *gin.Context) {
	name := ctx.Param("name")

	cards, err := repository.FindDigimonCardsByName(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

func FindDigimonCardsByColor(ctx *gin.Context) {
	var f Filter

	log.Printf("Filtering by colors: %+v\n", f)
	if err := ctx.ShouldBind(&f); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters.",
		})
		return
	}

	log.Printf("Filtering by colors: %+v\n", f.Colors)
	cards, err := repository.FindDigimonCardsByColor(ctx, f.Colors)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error 404": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}
