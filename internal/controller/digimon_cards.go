package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	"github.com/gin-gonic/gin"
)

type Filter struct {
	Colors []string `form:"color[]" collection_format:"csv"`
}

/*
 * GET ALL DIGIMON CARDS AND IMAGES
 */
func FindAllDigimonCards(ctx *gin.Context) {

	cards, err := repository.FindAllDigimonCards(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

func FindDigimonCardImages(ctx *gin.Context, presigner storage.Presigner) {
	keys, err := presigner.ListObjects(ctx, "cards/digimon/")
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

/////////////////////////////////////////////////////////////////////////////

/*
 * CRUD DIGIMON CARD BY ID
 */
func GetDigimonCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	log.Printf("Fetching card with ID: %s", id)

	card, err := repository.FindDigimonCardByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": card})
}

func PatchDigimonCardByID(ctx *gin.Context) {

	var body models.DigimonCard

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

	rows, err := repository.PatchDigimonCardByID(ctx, id, body)
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

func DeleteDigimonCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	rows, err := repository.DeleteDigimonCardByID(ctx, id)
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

/*
 * UPLOAD DIGIMON CARDS BY SETS
 */

func AddDigimonCardsBySets(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form.",
		})
		return
	}

	files := form.File["files"]
	err = repository.AddDigimonCardsBySets(ctx, files)
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

// func FindDigimonCardsByColor(ctx *gin.Context) {
// 	var f Filter

// 	log.Printf("Filtering by colors: %+v\n", f)
// 	if err := ctx.ShouldBind(&f); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid query parameters.",
// 		})
// 		return
// 	}

// 	log.Printf("Filtering by colors: %+v\n", f.Colors)
// 	cards, err := repository.FindDigimonCardsByColor(ctx, f.Colors)
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"error 404": "Cards not found.",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"content": cards})
// }

// func FindDigimonCardsByName(ctx *gin.Context) {
// 	name := ctx.Param("name")

// 	cards, err := repository.FindDigimonCardsByName(ctx, name)
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"error": "Cards not found.",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"content": cards})
// }
