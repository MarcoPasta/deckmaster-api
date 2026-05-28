package controller

import (
	"image/png"
	"net/http"
	"strconv"

	_ "embed"
	_ "image/jpeg"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/utils"
	"github.com/gin-gonic/gin"
	_ "golang.org/x/image/webp"
)

func GetPokemonDeckImage(ctx *gin.Context) {

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}

	deck, err := repository.GetPokemonDeckById(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}

	groups := utils.GroupPokemonCards(deck.Entries)
	img := utils.RenderDeckImage(groups, "pokemon")

	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", `attachment; filename="`+deck.Name+`.png"`)
	png.Encode(ctx.Writer, img)
}

func GetNarutoDeckImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}

	deck, err := repository.GetNarutoDeckById(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}

	groups := utils.GroupNarutoCards(deck.Entries)
	img := utils.RenderDeckImage(groups, "naruto")

	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", `attachment; filename="`+deck.Name+`.png"`)
	png.Encode(ctx.Writer, img)
}
