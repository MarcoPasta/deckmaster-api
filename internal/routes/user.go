package routes

import (
	"github.com/MarcoPasta/deckmaster-api/internal/controller"
	"github.com/gin-gonic/gin"
)

func AddUserRoutes(protected *gin.RouterGroup) {
	user := protected.Group("/user")
	{
		user.GET("/info", func(ctx *gin.Context) {
			controller.GetUserInfo(ctx)
		})
		user.POST("/created", func(ctx *gin.Context) {
			controller.HandleClerkUserCreated(ctx)
		})
		user.GET("/profile", func(ctx *gin.Context) {
			controller.GetUserProfile(ctx)
		})
	}

}
