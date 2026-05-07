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
	gin.SetMode(gin.ReleaseMode)
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

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")

	// Digimon Cards
	CardRoutes := api.Group("/digimon/cards")
	CardRoutes.GET("", func(ctx *gin.Context) {
		controller.FindAllDigimonCards(ctx)
	})
	CardRoutes.GET("/id/:id", func(ctx *gin.Context) {
		controller.FindDigimonCardByID(ctx)
	})
	CardRoutes.GET("/name/:name", func(ctx *gin.Context) {
		controller.FindDigimonCardsByName(ctx)
	})

	// Digimon Filtered
	DigimonFilteredRoutes := CardRoutes.Group("/filter")
	DigimonFilteredRoutes.GET("/color", func(ctx *gin.Context) {
		controller.FindDigimonCardsByColor(ctx)
	})

	// Pokemon Cards
	PokemonCardRoutes := api.Group("/pokemon/cards")
	PokemonCardRoutes.GET("", func(ctx *gin.Context) {
		controller.FindAllPokemonCards(ctx)
	})
	PokemonCardRoutes.GET("/id/:id", func(ctx *gin.Context) {
		controller.FindPokemonCardByID(ctx)
	})
	PokemonCardRoutes.GET("/name/:name", func(ctx *gin.Context) {
		controller.FindPokemonCardsByName(ctx)
	})

	// Pokemon Filtered
	PokemonFilteredRoutes := PokemonCardRoutes.Group("/filter")
	PokemonFilteredRoutes.GET("/energy-type", func(ctx *gin.Context) {
		controller.FindPokemonCardsByEnergyType(ctx)
	})

	// Naruto Cards
	NarutoCardRoutes := api.Group("/naruto/cards")
	NarutoCardRoutes.GET("", func(ctx *gin.Context) {
		controller.FindAllNarutoCards(ctx)
	})

	// Authentication routes
	protectedRoutes := api.Group("/auth")
	protectedRoutes.Use(middleware.ClerkAuthorizeHeader())
	{
		protectedRoutes.GET("/user", func(ctx *gin.Context) {
			controller.GetUserInfo(ctx)
		})
		// Digimon
		DigimondeckRoutes := protectedRoutes.Group("/digimon/decks")
		{
			DigimondeckRoutes.GET("/user", func(ctx *gin.Context) {
				controller.GetDigimonDeckByUserId(ctx)
			})
			DigimondeckRoutes.POST("/create", func(ctx *gin.Context) {
				controller.CreateDigimonDeck(ctx)
			})
		}
		PokemondeckRoutes := protectedRoutes.Group("/pokemon/decks")
		{
			PokemondeckRoutes.GET("/user", func(ctx *gin.Context) {
				controller.GetPokemonDecksByUserId(ctx)
			})
			PokemondeckRoutes.POST("/create", func(ctx *gin.Context) {
				controller.CreatePokemonDeck(ctx)
			})
		}
		NarutodeckRoutes := protectedRoutes.Group("/naruto/decks")
		{
			NarutodeckRoutes.GET("/user", func(ctx *gin.Context) {
				controller.GetNarutoDecksByUserId(ctx)
			})
			NarutodeckRoutes.POST("/create", func(ctx *gin.Context) {
				controller.CreateNarutoDeck(ctx)
			})
		}

	}

}
