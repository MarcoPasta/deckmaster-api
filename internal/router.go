package internal

import (
	"log"
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
	// ---------------------------------------------------------------------------------------------------------------------
	// ############################## CORS config ##############################
	//
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://prod.deckmaster.mirage-platform.cards", "https://deckmaster.mirage-platform.cards", "http://localhost:5173", "http://localhost:4173", "http://localhost:8081"},
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

	// ---------------------------------------------------------------------------------------------------------------------
	// ############################## Routes Grouping ##############################
	//
	api := r.Group("/api")

	// ---------------------------------------------------------------------------------------------------------------------
	// ############################## Clerk ##############################
	//
	clerkRoutes := api.Group("/clerk")
	clerkRoutes.GET("/session", func(ctx *gin.Context) {
		controller.GetSessionToken(ctx)
	})

	// ---------------------------------------------------------------------------------------------------------------------
	// ############################## User ##############################
	//
	userRoutes := clerkRoutes.Group("/user")
	{
		userRoutes.GET("/info", func(ctx *gin.Context) {
			controller.GetUserInfo(ctx)
		})
		userRoutes.POST("/created", func(ctx *gin.Context) {
			controller.HandleClerkUserCreated(ctx)
		})
	}

	// ---------------------------------------------------------------------------------------------------------------------
	// ############################## Admin ##############################
	//
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

		//
		// ############################## Admin / S3 Bucket ##############################
		//
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
		//
		// ############################## Admin / Digimon ##############################
		//
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

		// ############################## Admin / Pokemon ##############################
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

		//
		// ############################## Admin / Naruto ##############################
		//
		narutoCard := adminRoutes.Group("/naruto")
		narutoCard.Use(middleware.VerifyJWT())
		narutoCard.Use(middleware.GetClerkUserRole())
		narutoCard.Use(middleware.NeedsAdminRole())
		{
			narutoCard.GET("/card/:id", func(ctx *gin.Context) {
				controller.GetNarutoCardByID(ctx)
			})
			narutoCard.PATCH("/card/:id", func(ctx *gin.Context) {
				controller.PatchNarutoCardByID(ctx)
			})
			narutoCard.DELETE("/card/:id", func(ctx *gin.Context) {
				controller.DeleteNarutoCardByID(ctx)
			})
			narutoCard.POST("/sets", func(ctx *gin.Context) {
				controller.AddNarutoCardsBySets(ctx)
			})
		}

<<<<<<< HEAD
=======
<<<<<<< HEAD
	}

	// ---------------------------------------------------------------------------------------------------------------------
	// ############################## Cards ##############################
	//
=======
>>>>>>> main
		// User resync routes (admin only)
		adminUsersRoutes := adminRoutes.Group("/users")
		adminUsersRoutes.Use(middleware.VerifyJWT())
		adminUsersRoutes.Use(middleware.GetClerkUserRole())
		adminUsersRoutes.Use(middleware.NeedsAdminRole())
		{
			adminUsersRoutes.POST("/resync/:user_id", func(ctx *gin.Context) {
				controller.AdminResyncUser(ctx)
			})
			adminUsersRoutes.POST("/resync-all", func(ctx *gin.Context) {
				controller.AdminResyncAllUsers(ctx)
			})
		}
	}
<<<<<<< HEAD
=======
>>>>>>> b3ffc9b (add user sync features)
>>>>>>> main
	// ############################## Digimon ##############################
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

	// ############################## Pokemon ##############################
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

	// ############################## Naruto ##############################
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

	// ---------------------------------------------------------------------------------------------------------------------
	// ############################## Protected Routes ##############################
	protectedRoutes := api.Group("/auth")
	protectedRoutes.Use(middleware.VerifyJWT())
	{
<<<<<<< HEAD
=======
<<<<<<< HEAD
		protectedRoutes.GET("/user", func(ctx *gin.Context) {
			controller.GetUserInfo(ctx)
		})
=======
>>>>>>> main
		// User profile routes
		protectedRoutes.GET("/profile", func(ctx *gin.Context) {
			controller.GetUserProfile(ctx)
		})
		protectedRoutes.PUT("/profile", func(ctx *gin.Context) {
			controller.UpdateUserProfile(ctx)
		})

		// Password change logging endpoint
		protectedRoutes.POST("/profile/log-password-change", func(ctx *gin.Context) {
			controller.LogPasswordChange(ctx)
		})

		// User social links routes
		socialLinksRoutes := protectedRoutes.Group("/profile/social-links")
		{
			socialLinksRoutes.GET("", func(ctx *gin.Context) {
				controller.GetSocialLinks(ctx)
			})
			socialLinksRoutes.POST("", func(ctx *gin.Context) {
				controller.CreateSocialLink(ctx)
			})
			socialLinksRoutes.PUT("/:id", func(ctx *gin.Context) {
				controller.UpdateSocialLink(ctx)
			})
			socialLinksRoutes.DELETE("/:id", func(ctx *gin.Context) {
				controller.DeleteSocialLink(ctx)
			})
		}
<<<<<<< HEAD
=======
>>>>>>> b3ffc9b (add user sync features)
>>>>>>> main

		// ############################## Digimon Decks ##############################
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
			DigimondeckRoutes.PATCH("/:id/image/update/:cardId", func(ctx *gin.Context) {
				controller.UpdateDigimonDeckImage(ctx)
			})
			DigimondeckRoutes.DELETE("/:id/delete", func(ctx *gin.Context) {
				controller.DeleteDigimonDeck(ctx)
			})
			DigimondeckRoutes.PUT("/:id/update", func(ctx *gin.Context) {
				controller.UpdateDigimonDeck(ctx)
			})
		}

		// ############################## Pokemon Decks ##############################
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
			PokemondeckRoutes.PATCH("/:id/image/update/:cardId", func(ctx *gin.Context) {
				controller.UpdatePokemonDeckImage(ctx)
			})
			PokemondeckRoutes.DELETE("/:id/delete", func(ctx *gin.Context) {
				controller.DeletePokemonDeck(ctx)
			})
			PokemondeckRoutes.PUT("/:id/update", func(ctx *gin.Context) {
				controller.UpdatePokemonDeck(ctx)
			})
			PokemondeckRoutes.POST("/cardmarket", func(ctx *gin.Context) {
				log.Println("Received request for Pokemon Deck Card Market List")
				controller.GetPokemonDeckCardMarketList(ctx)
			})
		}
		// ############################## Naruto Decks ##############################
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
			NarutodeckRoutes.PATCH("/:id/image/update/:cardId", func(ctx *gin.Context) {
				controller.UpdateNarutoDeckImage(ctx)
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
