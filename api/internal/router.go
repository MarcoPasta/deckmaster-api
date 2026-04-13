package internal

import (
	"net/http"
	"time"

	"github.com/MarcoPasta/deckmaster/api/internal/controller"
	"github.com/MarcoPasta/deckmaster/api/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter() {
	r := gin.Default()
	initRoutes(r)

	Router = r
}

func initRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")

	// Cards
	CardRoutes := api.Group("/cards")
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
	FilteredRoutes := CardRoutes.Group("/filter")
	FilteredRoutes.GET("/color", func(ctx *gin.Context) {
		controller.FindDigimonCardsByColor(ctx)
	})

	// Authentication routes
	authRoutes := api.Group("/auth")
	authRoutes.Use(middleware.ClerkAuthorizeHeader())
	authRoutes.GET("/user", func(ctx *gin.Context) {
		controller.GetUserInfo(ctx)
	})

}
