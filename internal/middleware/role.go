package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

func GetClerkUserRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user, ok := ctx.Get("usr")
		if !ok {
			log.Println("User not found in context")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		id := user.(*clerk.User).ID
		url := os.Getenv("CLERK_BACKEND_API") + "/users/" + id

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating request: ", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't create request"})
			return
		}

		req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error making request: ", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch user data for: " + id})
			return
		}
		defer res.Body.Close()

		var data map[string]any
		json.NewDecoder(res.Body).Decode(&data)

		ctx.Set("role", getMeta(data, "role"))
		ctx.Next()
	}
}

func NeedsAdminRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.MustGet("role")
		log.Println("User role: ", role)
		if role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User does not have the required role"})
			return
		}
		ctx.Next()
	}
}

func getMeta(data map[string]any, key string) any {
	meta, ok := data["public_metadata"].(map[string]any)
	if !ok {
		return nil
	}
	return meta[key]
}
