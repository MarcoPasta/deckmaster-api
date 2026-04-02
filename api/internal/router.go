package internal

import (
	"net/http"

	"github.com/MarcoPasta/deckmaster/api/internal/controller"
	"github.com/MarcoPasta/deckmaster/api/internal/middleware"
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

	// Cards
	CardRoutes := r.Group("/cards")
	CardRoutes.GET("", func(ctx *gin.Context) {
		controller.FindAllDigimonCards(ctx)
	})
	CardRoutes.GET("/id/:id", func(ctx *gin.Context) {
		controller.FindDigimonCardByID(ctx)
	})
	CardRoutes.GET("/name/:name", func(ctx *gin.Context) {
		controller.FindDigimonCardsByName(ctx)
	})

	// Filtered
	FilteredRoutes := r.Group("/filter")
	FilteredRoutes.GET("/color", func(ctx *gin.Context) {
		controller.FindDigimonCardsByColor(ctx)
	})

	authRoutes := r.Group("/auth")
	authRoutes.Use(middleware.ClerkMiddleware())
	authRoutes.GET("/clerk", func(ctx *gin.Context) {
		middleware.ClerkListUsers(ctx)
	})
}
