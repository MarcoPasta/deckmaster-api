package middleware

import (
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"

	"github.com/gin-gonic/gin"
)

func VerifyJWT() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
		// using middleware to retrieve claims from token
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed, no valid session claims"})
				return
			}

			// another clerk middleware to fetch current user informations
			usr, err := user.Get(r.Context(), claims.Subject)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Authorization failed, no valid user found in claims"})
				return
			}

			// write claims and user in gins context
			ctx.Set("claims", claims)
			ctx.Set("usr", usr)
		})

		// clerk middleware that parses the token
		clerkhttp.WithHeaderAuthorization()(handler).ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Next()
	}
}
