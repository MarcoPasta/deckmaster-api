package controller

import (
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/gin-gonic/gin"
)

func FindAllNarutoCards(ctx *gin.Context) {

	cards, err := repository.FindAllNarutoCards(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}
