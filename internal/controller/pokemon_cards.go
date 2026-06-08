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

func FindPokemonCardImages(ctx *gin.Context, presigner storage.Presigner) {
	keys, err := presigner.ListObjects(ctx, "cards/pokemon/")
	if err != nil {
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

// ####### ALL POKEMON CARDS ####### //
func FindAllPokemonCards(ctx *gin.Context) {

	cards, err := repository.FindAllPokemonCards(ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

// ####### POKEMON CARD BY ID ####### //

func GetPokemonCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	card, err := repository.GetPokemonCardByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": card})
}

func PatchPokemonCardByID(ctx *gin.Context) {

	var body models.PokemonCard

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

	rows, err := repository.PatchPokemonCardByID(ctx, id, body)
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

func DeletePokemonCardByID(ctx *gin.Context) {
	id := ctx.Param("id")

	rows, err := repository.DeletePokemonCardByID(ctx, id)
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

// ####### POKEMON CARDS BY NAME ####### //
func FindSamePokemonCardsByNames(ctx *gin.Context) {
	decklist := ctx.Query("decklist")

	parsedList := utils.ParsePokemonDecklist(ctx, decklist)
	if parsedList == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid decklist format.",
		})
		return
	}

	cards, err := repository.FindSamePokemonCardsByNames(ctx, parsedList)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

// ####### POKEMON CARDS BY ENERGY TYPE ####### //
func FindPokemonCardsByEnergyType(ctx *gin.Context) {
	var f Filter

	log.Printf("Filtering by energy types: %+v\n", f)
	if err := ctx.ShouldBind(&f); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters.",
		})
		return
	}

	log.Printf("Filtering by energy types: %+v\n", f.Colors)
	cards, err := repository.FindPokemonCardsByEnergyType(ctx, f.Colors)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error 404": "Cards not found.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": cards})
}

// ####### ADD POKEMON CARDS BY SETS ####### //
func AddPokemonCardsBySets(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse multipart form.",
		})
		return
	}

	files := form.File["files"]
	err = repository.AddPokemonCardsBySets(ctx, files)
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
