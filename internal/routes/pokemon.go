package routes

import (
	"log"

	"github.com/MarcoPasta/deckmaster-api/internal/controller"
	"github.com/MarcoPasta/deckmaster-api/internal/middleware"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func AddPokemonRoutes(api *gin.RouterGroup, admin *gin.RouterGroup, protected *gin.RouterGroup, presigner storage.Presigner) {

	AddPokemonAdminRoutes(admin)
	AddPokemonCardRoutes(api, presigner)
	AddPokemonDeckRoutes(protected, presigner)

}

func AddPokemonAdminRoutes(admin *gin.RouterGroup) {
	pokemon := admin.Group("/pokemon")
	pokemon.Use(middleware.VerifyJWT())
	pokemon.Use(middleware.GetClerkUserRole())
	pokemon.Use(middleware.NeedsAdminRole())
	{
		pokemon.GET("/card/:id", func(ctx *gin.Context) {
			controller.GetPokemonCardByID(ctx)
		})
		pokemon.PATCH("/card/:id", func(ctx *gin.Context) {
			controller.PatchPokemonCardByID(ctx)
		})
		pokemon.DELETE("/card/:id", func(ctx *gin.Context) {
			controller.DeletePokemonCardByID(ctx)
		})
		pokemon.POST("/sets", func(ctx *gin.Context) {
			controller.AddPokemonCardsBySets(ctx)
		})
	}
}

func AddPokemonCardRoutes(api *gin.RouterGroup, presigner storage.Presigner) {
	cards := api.Group("/pokemon/cards")
	{
		cards.GET("", func(ctx *gin.Context) {
			controller.FindAllPokemonCards(ctx)
		})
		cards.GET("/names", func(ctx *gin.Context) {
			controller.FindSamePokemonCardsByNames(ctx)
		})
		cards.GET("/images", func(ctx *gin.Context) {
			controller.FindPokemonCardImages(ctx, presigner)
		})
	}

	// Pokemon Filtered
	PokemonFilteredRoutes := cards.Group("/filter")
	{
		PokemonFilteredRoutes.GET("/energy-type", func(ctx *gin.Context) {
			controller.FindPokemonCardsByEnergyType(ctx)
		})
	}
}

func AddPokemonDeckRoutes(protected *gin.RouterGroup, presigner storage.Presigner) {
	decks := protected.Group("/pokemon/decks")
	{
		decks.GET("/user", func(ctx *gin.Context) {
			controller.GetPokemonDecksByUserId(ctx)
		})
		decks.POST("/create", func(ctx *gin.Context) {
			controller.CreatePokemonDeck(ctx)
		})
		decks.GET("/:id/image", func(ctx *gin.Context) {
			controller.GetPokemonDeckImage(ctx, presigner)
		})
		decks.PATCH("/:id/image/update/:cardId", func(ctx *gin.Context) {
			controller.UpdatePokemonDeckImage(ctx)
		})
		decks.DELETE("/:id/delete", func(ctx *gin.Context) {
			controller.DeletePokemonDeck(ctx)
		})
		decks.PUT("/:id/update", func(ctx *gin.Context) {
			controller.UpdatePokemonDeck(ctx)
		})
		decks.POST("/cardmarket", func(ctx *gin.Context) {
			log.Println("Received request for Pokemon Deck Card Market List")
			controller.GetPokemonDeckCardMarketList(ctx)
		})
	}
}
