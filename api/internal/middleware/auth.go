package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"

	"github.com/gin-gonic/gin"
)

// clerk handler
func ClerkAuthorizeHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

		fmt.Println("### Clerk Authorization Middleware ###")
		handler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No Session Claims found"})
				ctx.Abort()
				return
			}

			usr, err := clerkUser.Get(r.Context(), claims.Subject)
			if err != nil {
				fmt.Println("Claims: ", claims.Subject)
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User is not Authorized"})
				ctx.Abort()
				return
			}
			fmt.Printf("Authenticated user: %s (ID: %s)\n", *usr.Username, usr.ID)
		}))

		handler.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Next()
	}
}
