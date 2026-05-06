package controller

import (
	"fmt"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetNarutoDecksByUserId(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	fmt.Println("Received user_id: ", id)

	decks, err := repository.GetNarutoDecksByUserId(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch decks"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": decks})
}

func CreateNarutoDeck(ctx *gin.Context) {
	body := utils.CreateNarutoDeckRequest{}
	err := ctx.ShouldBindBodyWithJSON(&body)

	if err != nil {
		fmt.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	saved, err := repository.CreateNarutoDeck(ctx, utils.MarshalNarutoDeckRequest(body))
	if err != nil {
		fmt.Println("Create deck error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, saved)
}
