package routes

import (
	"github.com/MarcoPasta/deckmaster-api/internal/controller"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func AddDigimonRoutes(api *gin.RouterGroup, admin *gin.RouterGroup, protected *gin.RouterGroup, presigner storage.Presigner) {

	AddDigimonAdminRoutes(admin)
	AddDigimonCardRoutes(api, presigner)
	AddDigimonDeckRoutes(protected, presigner)

}

func AddDigimonAdminRoutes(admin *gin.RouterGroup) {

	digimon := admin.Group("/digimon")
	{
		digimon.GET("/card/:id", func(ctx *gin.Context) {
			controller.GetDigimonCardByID(ctx)
		})
		digimon.PATCH("/card/:id", func(ctx *gin.Context) {
			controller.PatchDigimonCardByID(ctx)
		})
		digimon.DELETE("/card/:id", func(ctx *gin.Context) {
			controller.DeleteDigimonCardByID(ctx)
		})
		digimon.POST("/sets", func(ctx *gin.Context) {
			controller.AddDigimonCardsBySets(ctx)
		})
	}
}

func AddDigimonCardRoutes(api *gin.RouterGroup, presigner storage.Presigner) {
	cards := api.Group("/digimon/cards")
	{
		cards.GET("", func(ctx *gin.Context) {
			controller.FindAllDigimonCards(ctx)
		})
		// CardRoutes.GET("/id/:id", func(ctx *gin.Context) {
		// 	controller.FindDigimonCardByID(ctx)
		// })
		// CardRoutes.GET("/name/:name", func(ctx *gin.Context) {
		// 	controller.FindDigimonCardsByName(ctx)
		// })
		cards.GET("/images", func(ctx *gin.Context) {
			controller.FindDigimonCardImages(ctx, presigner)
		})
	}
}

func AddDigimonDeckRoutes(protected *gin.RouterGroup, presigner storage.Presigner) {
	decks := protected.Group("/digimon/decks")
	{
		decks.GET("/user", func(ctx *gin.Context) {
			controller.GetDigimonDeckByUserId(ctx)
		})
		decks.POST("/create", func(ctx *gin.Context) {
			controller.CreateDigimonDeck(ctx)
		})
		decks.GET("/:id/image", func(ctx *gin.Context) {
			controller.GetDigimonDeckImage(ctx, presigner)
		})
		decks.PATCH("/:id/image/update/:cardId", func(ctx *gin.Context) {
			controller.UpdateDigimonDeckImage(ctx)
		})
		decks.DELETE("/:id/delete", func(ctx *gin.Context) {
			controller.DeleteDigimonDeck(ctx)
		})
		decks.PUT("/:id/update", func(ctx *gin.Context) {
			controller.UpdateDigimonDeck(ctx)
		})
	}
}
