package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CheckOptions(c *gin.Context) {
	fmt.Println("### OPTIONS request received ###")
	c.JSON(http.StatusOK, gin.H{"message": "OPTIONS request successful"})
}

func GetUserInfo(c *gin.Context) {

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	url := os.Getenv("CLERK_BACKEND_API") + "/users/" + id
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": string(body)})

}
