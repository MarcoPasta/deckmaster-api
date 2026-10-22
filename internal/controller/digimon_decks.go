package controller

import (
	"log"
	"net/http"

	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/MarcoPasta/deckmaster-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetDigimonDeckByUserId(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	decks, err := repository.GetDigimonDeckByUserId(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch decks"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": decks})
}

func CreateDigimonDeck(ctx *gin.Context) {
	body := utils.CreateDigimonDeckRequest{}
	err := ctx.ShouldBindBodyWithJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Received CreateDigimonDeck request: %+v\n", body)

	saved, err := repository.CreateDigimonDeck(ctx, utils.MarshalDigimonDeckRequest(body))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deck"})
		return
	}

	ctx.JSON(http.StatusOK, saved)
}

func UpdateDigimonDeck(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deck id is required"})
		return
	}
	body := utils.UpdateDigimonDeckRequest{}
	err := ctx.ShouldBindBodyWithJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Received UpdateDigimonDeck request for deck ID %s: %+v\n", id, body)

	updated, err := repository.UpdateDigimonDeck(ctx, id, utils.MarshalUpdateDigimonDeckRequest(body))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deck"})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func UpdateDigimonDeckImage(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deck id is required"})
		return
	}

	cardId := ctx.Param("cardId")
	if cardId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cardId is required"})
		return
	}

	err := repository.UpdateDigimonDeckImage(ctx, id, cardId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cover image updated successfully"})
}

func DeleteDigimonDeck(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deck id is required"})
		return
	}

	log.Printf("Received DeleteDigimonDeck request for deck ID %s\n", id)

	rows, err := repository.DeleteDigimonDeck(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deck"})
		return
	}

	log.Printf("Deleted %d rows for deck ID %s\n", rows, id)

	ctx.JSON(http.StatusOK, gin.H{"message": "Deck deleted successfully"})
}
