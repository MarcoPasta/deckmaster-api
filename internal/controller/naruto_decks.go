package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/MarcoPasta/deckmaster-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetNarutoDecksByUserId(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

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

func UpdateNarutoDeck(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deck id is required"})
		return
	}
	body := utils.UpdateNarutoDeckRequest{}
	err := ctx.ShouldBindBodyWithJSON(&body)
	if err != nil {
		fmt.Println("Binding error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := repository.UpdateNarutoDeck(ctx, id, utils.MarshalUpdateNarutoDeckRequest(body))
	if err != nil {
		fmt.Println("Update deck error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"updated": updated})

}

func UpdateNarutoDeckImage(ctx *gin.Context) {
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

	err := repository.UpdateNarutoDeckImage(ctx, id, cardId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cover image updated successfully"})
}

func DeleteNarutoDeck(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "deck id is required"})
		return
	}

	deckId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}

	rows, err := repository.DeleteNarutoDeck(ctx, deckId)
	if err != nil {
		fmt.Println("Delete deck error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deck deleted successfully, %d rows affected", rows)})
}
