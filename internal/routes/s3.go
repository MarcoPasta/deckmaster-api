package routes

import (
	"net/http"

	"github.com/MarcoPasta/deckmaster-api/internal/middleware"
	"github.com/MarcoPasta/deckmaster-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func AddS3Routes(api *gin.RouterGroup, presigner storage.Presigner) {
	bucketRoutes := api.Group("/bucket")
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
}
