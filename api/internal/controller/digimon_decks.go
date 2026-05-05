package controller

import (
	"fmt"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"

	"github.com/gin-gonic/gin"
)

func CreateDigimonDeck(ctx *gin.Context) {

	body := utils.CreateDigimonDeckRequest{}
	err := ctx.ShouldBindBodyWithJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body blabla"})
		return
	}

	saved, err := repository.CreateDigimonDeck(ctx, utils.MarshalDigimonDeckRequest(body))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deck"})
		return
	}

	ctx.JSON(http.StatusOK, saved)
}

func GetDigimonDeckByUserId(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	fmt.Println("Received user_id: ", id)

	decks, err := repository.GetDigimonDeckByUserId(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch decks"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": decks})
}
