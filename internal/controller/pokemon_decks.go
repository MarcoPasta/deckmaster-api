package controller

import (
	"fmt"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetPokemonDecksByUserId(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	decks, err := repository.GetPokemonDecksByUserId(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch decks"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": decks})
}

func CreatePokemonDeck(ctx *gin.Context) {
	body := utils.CreatePokemonDeckRequest{}
	err := ctx.ShouldBindBodyWithJSON(&body)

	if err != nil {
		fmt.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	saved, err := repository.CreatePokemonDeck(ctx, utils.MarshalPokemonDeckRequest(body))
	if err != nil {
		fmt.Println("Create deck error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, saved)
}
