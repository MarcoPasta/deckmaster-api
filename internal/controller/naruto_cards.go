package controller

import (
	"log"
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/storage"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"
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

func FindNarutoCardsByNames(ctx *gin.Context) {
	log.Println("Received request to find Naruto cards by names")
	decklist := ctx.Query("decklist")

	parsedList := utils.ParseNarutoDecklist(decklist)
	log.Printf("Parsed decklist: %+v\n", parsedList)
	if parsedList == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decklist format.",
		})
		return
	}

	cards, err := repository.FindNarutoCardsByNames(ctx, parsedList)
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

func AddNarutoCardsBySets(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form.",
		})
		return
	}

	files := form.File["files"]
	err = repository.AddNarutoCardsBySets(ctx, files)
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
