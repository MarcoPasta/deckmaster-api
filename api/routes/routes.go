package routes

import (
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/controllers/cards"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter() {
	r := gin.Default()
	initRoutes(r)

	Router = r
}

func initRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	CardRoutes := r.Group("/cards")
	CardRoutes.GET("", func(ctx *gin.Context) {
		cards.Find(ctx)
	})
}
