package middleware

import (
	"fmt"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
)

// clerk handler

func CheckClerkAuth(c *gin.Context) (error, *clerk.User) {

	fmt.Println("Calling Clerk Handler")

	details, err := user.Get(c, "user_3BlhFfnu05v9ZdYqtLPmEw6IsGL")
	if err != nil {
		return err, nil
	}

	return nil, details
}

func ClerkListUsers(c *gin.Context) {

	fmt.Println("Listing Clerk Users")
	userList, err := user.List(c, &user.ListParams{})
	if err != nil {
		fmt.Println("Error listing users:", err)
		return
	}

	fmt.Println("User list:", userList)

}

func ClerkMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Middleware logic here
		clerk.SetKey("sk_test_TEYkt1b0ykCc55Y03YFkFiqKOZx6MgKEIzvzmqV4ZA")
		fmt.Println("Setting up Clerk Middleware")
		c.Next()
	}
}
