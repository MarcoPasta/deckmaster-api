package controller

import (
	"log"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func FindAllPokemonCards(ctx *gin.Context) {

	cards, err := repository.FindAllPokemonCards(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

func FindPokemonCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	log.Printf("Fetching card with ID: %s", id)

	card, err := repository.FindPokemonCardByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": card})
}

func FindSamePokemonCardsByNames(ctx *gin.Context) {
	decklist := ctx.Query("decklist")

	parsedList := utils.ParseDecklist(ctx, decklist)
	if parsedList == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decklist format.",
		})
		return
	}

	cards, err := repository.FindSamePokemonCardsByNames(ctx, parsedList)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

func FindPokemonCardsByEnergyType(ctx *gin.Context) {
	var f Filter

	log.Printf("Filtering by energy types: %+v\n", f)
	if err := ctx.ShouldBind(&f); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters.",
		})
		return
	}

	log.Printf("Filtering by energy types: %+v\n", f.Colors)
	cards, err := repository.FindPokemonCardsByEnergyType(ctx, f.Colors)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error 404": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

func AddPokemonCardsBySets(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form.",
		})
		return
	}

	files := form.File["files"]
	err = repository.AddPokemonCardsBySets(ctx, files)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add cards.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cards added successfully.",
	})
}
