package routes

import (
	"github.com/MarcoPasta/deckmaster-api/internal/controller"
	"github.com/gin-gonic/gin"
)

func AddUserRoutes(api *gin.RouterGroup, protected *gin.RouterGroup) {

	api.POST("/user/created", func(ctx *gin.Context) {
		controller.HandleClerkUserCreated(ctx)
	})

	user := protected.Group("/user")
	{
		// test route?
		user.GET("/info", func(ctx *gin.Context) {
			controller.GetUserInfo(ctx)
		})
		user.GET("/profile/:id", func(ctx *gin.Context) {
			controller.GetUserProfile(ctx)
		})
	}

}
