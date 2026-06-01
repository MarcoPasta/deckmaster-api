package controller

import (
	"log"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/storage"
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

func FindNarutoCardImages(ctx *gin.Context, presigner storage.Presigner) {
	keys, err := presigner.ListObjects(ctx, "cards/naruto/")
	if err != nil {
		log.Printf("Error listing objects: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list objects"})
		return
	}

	urls := make(map[string]string, len(keys))
	for _, key := range keys {
		signed, err := presigner.GetObject(ctx, key)
		if err != nil {
			continue
		}
		urls[key] = signed.URL
	}

	ctx.JSON(http.StatusOK, gin.H{"images": urls})
}
