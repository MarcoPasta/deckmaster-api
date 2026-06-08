package internal

import (
	"net/http"
	"time"

	"github.com/MarcoPasta/deckmaster/api/internal/controller"
	"github.com/MarcoPasta/deckmaster/api/internal/middleware"
	"github.com/MarcoPasta/deckmaster/api/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter(presigner storage.Presigner) {
	r := gin.Default()
	initRoutes(r, presigner)

	Router = r
}

func initRoutes(r *gin.Engine, presigner storage.Presigner) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://dev.deckmaster.mirage-platform.cards", "https://deckmaster.mirage-platform.cards", "http://localhost:5173", "http://localhost:4173", "http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
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

	clerkRoutes := api.Group("/clerk")

	clerkRoutes.GET("/session", func(ctx *gin.Context) {
		controller.GetSessionToken(ctx)
	})

	userRoutes := clerkRoutes.Group("/user")
	{
		userRoutes.GET("/info", func(ctx *gin.Context) {
			controller.GetUserInfo(ctx)
		})
		userRoutes.POST("/created", func(ctx *gin.Context) {
			controller.HandleClerkUserCreated(ctx)
		})
	}

	// Admin Routes
	adminRoutes := api.Group("/admin")
	{
		// Testing Endpoint
		adminRoutes.GET("/role", func(ctx *gin.Context) {
			role := ctx.MustGet("role")
			ctx.JSON(http.StatusOK, gin.H{"role": role})
		})
		// Testing Endpoint
		adminRoutes.GET("/admin", middleware.NeedsAdminRole(), func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Welcome, admin user!"})
		})
		// S3 Bucket
		bucketRoutes := adminRoutes.Group("/bucket")
		bucketRoutes.Use(middleware.VerifyJWT())
		bucketRoutes.Use(middleware.GetClerkUserRole())
		bucketRoutes.Use(middleware.NeedsAdminRole())
		{
			bucketRoutes.PUT("/presign-put", func(ctx *gin.Context) {
				objectKey := ctx.Query("key")
				if objectKey == "" {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing 'key' query parameter"})
					return
				}
				url, err := presigner.PutObject(ctx, objectKey)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"url": url.URL})
			})
		}
		//////////////
		digimonCard := adminRoutes.Group("/digimon")
		digimonCard.Use(middleware.VerifyJWT())
		digimonCard.Use(middleware.GetClerkUserRole())
		digimonCard.Use(middleware.NeedsAdminRole())
		{
			digimonCard.GET("/card/:id", func(ctx *gin.Context) {
				controller.GetDigimonCardByID(ctx)
			})
			digimonCard.PATCH("/card/:id", func(ctx *gin.Context) {
				controller.PatchDigimonCardByID(ctx)
			})
			digimonCard.DELETE("/card/:id", func(ctx *gin.Context) {
				controller.DeleteDigimonCardByID(ctx)
			})
			digimonCard.POST("/sets", func(ctx *gin.Context) {
				controller.AddDigimonCardsBySets(ctx)
			})
		}

		pokemonCard := adminRoutes.Group("/pokemon")
		pokemonCard.Use(middleware.VerifyJWT())
		pokemonCard.Use(middleware.GetClerkUserRole())
		pokemonCard.Use(middleware.NeedsAdminRole())
		{
			pokemonCard.GET("/card/:id", func(ctx *gin.Context) {
				controller.GetPokemonCardByID(ctx)
			})
			pokemonCard.PATCH("/card/:id", func(ctx *gin.Context) {
				controller.PatchPokemonCardByID(ctx)
			})
			pokemonCard.DELETE("/card/:id", func(ctx *gin.Context) {
				controller.DeletePokemonCardByID(ctx)
			})
			pokemonCard.POST("/sets", func(ctx *gin.Context) {
				controller.AddPokemonCardsBySets(ctx)
			})
		}

		narutoCard := adminRoutes.Group("/naruto")
		narutoCard.Use(middleware.VerifyJWT())
		narutoCard.Use(middleware.GetClerkUserRole())
		narutoCard.Use(middleware.NeedsAdminRole())
		{
			// narutoCard.GET("/card/:id", func(ctx *gin.Context) {
			// 	controller.GetNarutoCardByID(ctx)
			// })
			// narutoCard.PATCH("/card/:id", func(ctx *gin.Context) {
			// 	controller.PatchNarutoCardByID(ctx)
			// })
			// narutoCard.DELETE("/card/:id", func(ctx *gin.Context) {
			// 	controller.DeleteNarutoCardByID(ctx)
			// })
			narutoCard.POST("/sets", func(ctx *gin.Context) {
				controller.AddNarutoCardsBySets(ctx)
			})
		}

	}

	// ---------------------------------------------------------------------------------------------------------------------
	// ############### Cards ###############
	//
	// Digimon Cards
	CardRoutes := api.Group("/digimon/cards")
	{
		CardRoutes.GET("", func(ctx *gin.Context) {
			controller.FindAllDigimonCards(ctx)
		})
		// CardRoutes.GET("/id/:id", func(ctx *gin.Context) {
		// 	controller.FindDigimonCardByID(ctx)
		// })
		// CardRoutes.GET("/name/:name", func(ctx *gin.Context) {
		// 	controller.FindDigimonCardsByName(ctx)
		// })
		CardRoutes.GET("/images", func(ctx *gin.Context) {
			controller.FindDigimonCardImages(ctx, presigner)
		})
	}

	// Digimon Filtered
	// DigimonFilteredRoutes := CardRoutes.Group("/filter")
	// {
	// DigimonFilteredRoutes.GET("/color", func(ctx *gin.Context) {
	// 	controller.FindDigimonCardsByColor(ctx)
	// })
	// }

	// Pokemon Cards
	PokemonCardRoutes := api.Group("/pokemon/cards")
	{
		PokemonCardRoutes.GET("", func(ctx *gin.Context) {
			controller.FindAllPokemonCards(ctx)
		})
		PokemonCardRoutes.GET("/names", func(ctx *gin.Context) {
			controller.FindSamePokemonCardsByNames(ctx)
		})
		PokemonCardRoutes.GET("/images", func(ctx *gin.Context) {
			controller.FindPokemonCardImages(ctx, presigner)
		})
	}

	// Pokemon Filtered
	PokemonFilteredRoutes := PokemonCardRoutes.Group("/filter")
	{
		PokemonFilteredRoutes.GET("/energy-type", func(ctx *gin.Context) {
			controller.FindPokemonCardsByEnergyType(ctx)
		})
	}

	// Naruto Cards
	NarutoCardRoutes := api.Group("/naruto/cards")
	{
		NarutoCardRoutes.GET("", func(ctx *gin.Context) {
			controller.FindAllNarutoCards(ctx)
		})
		NarutoCardRoutes.GET("/names", func(ctx *gin.Context) {
			controller.FindNarutoCardsByNames(ctx)
		})
		NarutoCardRoutes.GET("/images", func(ctx *gin.Context) {
			controller.FindNarutoCardImages(ctx, presigner)
		})
	}

	// Authentication routes
	protectedRoutes := api.Group("/auth")
	protectedRoutes.Use(middleware.VerifyJWT())
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
			DigimondeckRoutes.GET("/:id/image", func(ctx *gin.Context) {
				controller.GetDigimonDeckImage(ctx, presigner)
			})
			DigimondeckRoutes.DELETE("/:id/delete", func(ctx *gin.Context) {
				controller.DeleteDigimonDeck(ctx)
			})
			DigimondeckRoutes.PUT("/:id/update", func(ctx *gin.Context) {
				controller.UpdateDigimonDeck(ctx)
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
			PokemondeckRoutes.GET("/:id/image", func(ctx *gin.Context) {
				controller.GetPokemonDeckImage(ctx, presigner)
			})
			PokemondeckRoutes.DELETE("/:id/delete", func(ctx *gin.Context) {
				controller.DeletePokemonDeck(ctx)
			})
			PokemondeckRoutes.PUT("/:id/update", func(ctx *gin.Context) {
				controller.UpdatePokemonDeck(ctx)
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
			NarutodeckRoutes.GET("/:id/image", func(ctx *gin.Context) {
				controller.GetNarutoDeckImage(ctx, presigner)
			})
			NarutodeckRoutes.DELETE("/:id/delete", func(ctx *gin.Context) {
				controller.DeleteNarutoDeck(ctx)
			})
			NarutodeckRoutes.PUT("/:id/update", func(ctx *gin.Context) {
				controller.UpdateNarutoDeck(ctx)
			})
		}
	}
}
