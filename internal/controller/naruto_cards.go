package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
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

// ########## NARUTO CARDS BY ID ##########

func GetNarutoCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	card, err := repository.GetNarutoCardByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": card})
}

func PatchNarutoCardByID(ctx *gin.Context) {

	var body models.NarutoCard

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid card ID.",
		})
		return
	}

	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Error binding JSON:", body, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body.",
		})
		return
	}

	log.Printf("Body: %+v\n", body)

	rows, err := repository.PatchNarutoCardByID(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update card.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Card updated successfully: " + strconv.Itoa(rows) + " rows affected.",
	})
}

func DeleteNarutoCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	rows, err := repository.DeleteNarutoCardByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete card.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Card deleted successfully: " + strconv.Itoa(rows) + " rows affected.",
	})
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
