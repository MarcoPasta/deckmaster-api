package routes

import (
	"fmt"

	"github.com/MarcoPasta/deckmaster-api/internal/controller"
	"github.com/gin-gonic/gin"
)

func AddClerkRoutes(api *gin.RouterGroup) {

	clerkRoutes := api.Group("/clerk")
	clerkRoutes.POST("/", func(ctx *gin.Context) {
		fmt.Println("Webhook arrived")
	})
	clerkRoutes.GET("/session", func(ctx *gin.Context) {
		controller.GetSessionToken(ctx)
	})

}
