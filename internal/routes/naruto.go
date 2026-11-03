package routes

import (
	"github.com/MarcoPasta/deckmaster-api/internal/controller"
	"github.com/MarcoPasta/deckmaster-api/internal/middleware"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func AddNarutoRoutes(api *gin.RouterGroup, admin *gin.RouterGroup, protected *gin.RouterGroup, presigner storage.Presigner) {

	AddNarutoAdminRoutes(admin)
	AddNarutoCardRoutes(api, presigner)
	AddNarutoDeckRoutes(protected, presigner)

}

func AddNarutoAdminRoutes(admin *gin.RouterGroup) {
	naruto := admin.Group("/naruto")
	naruto.Use(middleware.VerifyJWT())
	naruto.Use(middleware.GetClerkUserRole())
	naruto.Use(middleware.NeedsAdminRole())
	{
		naruto.GET("/card/:id", func(ctx *gin.Context) {
			controller.GetNarutoCardByID(ctx)
		})
		naruto.PATCH("/card/:id", func(ctx *gin.Context) {
			controller.PatchNarutoCardByID(ctx)
		})
		naruto.DELETE("/card/:id", func(ctx *gin.Context) {
			controller.DeleteNarutoCardByID(ctx)
		})
		naruto.POST("/sets", func(ctx *gin.Context) {
			controller.AddNarutoCardsBySets(ctx)
		})
	}
}

func AddNarutoCardRoutes(api *gin.RouterGroup, presigner storage.Presigner) {
	cards := api.Group("/naruto/cards")
	{
		cards.GET("", func(ctx *gin.Context) {
			controller.FindAllNarutoCards(ctx)
		})
		cards.GET("/names", func(ctx *gin.Context) {
			controller.FindNarutoCardsByNames(ctx)
		})
		cards.GET("/images", func(ctx *gin.Context) {
			controller.FindNarutoCardImages(ctx, presigner)
		})
	}
}

func AddNarutoDeckRoutes(protected *gin.RouterGroup, presigner storage.Presigner) {
	decks := protected.Group("/naruto/decks")
	{
		decks.GET("/user", func(ctx *gin.Context) {
			controller.GetNarutoDecksByUserId(ctx)
		})
		decks.POST("/create", func(ctx *gin.Context) {
			controller.CreateNarutoDeck(ctx)
		})
		decks.GET("/:id/image", func(ctx *gin.Context) {
			controller.GetNarutoDeckImage(ctx, presigner)
		})
		decks.PATCH("/:id/image/update/:cardId", func(ctx *gin.Context) {
			controller.UpdateNarutoDeckImage(ctx)
		})
		decks.DELETE("/:id/delete", func(ctx *gin.Context) {
			controller.DeleteNarutoDeck(ctx)
		})
		decks.PUT("/:id/update", func(ctx *gin.Context) {
			controller.UpdateNarutoDeck(ctx)
		})
	}
}
