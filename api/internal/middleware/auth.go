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
	return func(c *gin.Context) {
		clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

		handler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			fmt.Println("### Clerk Authorization Middleware ###")
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No Session Claims found"})
				c.Abort()
				return
			}
			usr, err := clerkUser.Get(r.Context(), claims.Subject)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not Authorized"})
				c.Abort()
				return
			}
			fmt.Printf("Authenticated user: %s (ID: %s)\n", *usr.Username, usr.ID)
		}))

		handler.ServeHTTP(c.Writer, c.Request)
		c.Next()
	}
}
