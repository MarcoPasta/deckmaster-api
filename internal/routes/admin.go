package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAdminRoutes(admin *gin.RouterGroup) {
	{
		// Testing Endpoint
		admin.GET("/role", func(ctx *gin.Context) {
			role := ctx.MustGet("role")
			ctx.JSON(http.StatusOK, gin.H{"role": role})
		})
		// Testing Endpoint
		admin.GET("/admin", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Welcome, admin user!"})
		})
	}
}
