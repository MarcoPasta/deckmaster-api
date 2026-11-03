package internal

import (
	"github.com/MarcoPasta/deckmaster-api/internal/routes"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter(presigner storage.Presigner) {
	r := gin.Default()
	routes.InitRoutes(r, presigner)

	Router = r
}
