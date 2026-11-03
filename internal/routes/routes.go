package routes

import (
	"net/http"
	"time"

	"github.com/MarcoPasta/deckmaster-api/internal/middleware"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, presigner storage.Presigner) {
	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://prod.deckmaster.mirage-platform.cards", "https://deckmaster.mirage-platform.cards", "http://localhost:5173", "http://localhost:4173", "http://localhost:8081", "https://aloe-afoot-shingle.ngrok-free.dev"},
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

	// routing groups
	api := r.Group("/api")
	admin := api.Group("/admin")
	protected := api.Group("/auth")

	// assign middleware
	admin.
		Use(middleware.VerifyJWT()).
		Use(middleware.GetClerkUserRole()).
		Use(middleware.NeedsAdminRole())
	protected.
		Use(middleware.VerifyJWT())

	// init routes
	AddClerkRoutes(api)
	AddUserRoutes(protected)
	AddAdminRoutes(admin)
	AddS3Routes(api, presigner)
	AddDigimonRoutes(api, admin, protected, presigner)
	AddPokemonRoutes(api, admin, protected, presigner)
	AddNarutoRoutes(api, admin, protected, presigner)

}
