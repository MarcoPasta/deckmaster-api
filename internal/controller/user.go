package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckOptions(ctx *gin.Context) {
	fmt.Println("### OPTIONS request received ###")
	ctx.JSON(http.StatusOK, gin.H{"message": "OPTIONS request successful"})
}

// Retrieve session token for bruno testing
func GetSessionToken(ctx *gin.Context) {
	url := "https://api.clerk.dev/v1/sessions/"
	user_id := "user_3BlhFfnu05v9ZdYqtLPmEw6IsGL"
	payload := strings.NewReader("{\"user_id\": \"" + user_id + "\"}")
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body"})
		return
	}

	url = "https://api.clerk.dev/v1/sessions/" + data["id"].(string) + "/tokens/Bruno_Testing_Template"
	req, _ = http.NewRequest("POST", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session token"})
		return
	}

	defer res.Body.Close()
	var tokenData map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&tokenData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode token response body"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"session_token": tokenData["jwt"]})
}

func GetUserInfo(ctx *gin.Context) {

	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	url := os.Getenv("CLERK_BACKEND_API") + "/users/" + id
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"content": string(body)})

}

func HandleClerkUserCreated(ctx *gin.Context) {
	// Receive request body with new user
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("### HandleClerkUserCreated error ### --> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	// Decode JSON data, assign new map with role data
	var data map[string]map[string]string
	json.NewDecoder(bytes.NewReader(body)).Decode(&data)
	log.Println("### Received Data For ### --> ", data["data"]["id"])

	url := os.Getenv("CLERK_BACKEND_API") + "/users/" + data["data"]["id"]

	req, _ := http.NewRequest("PATCH", url, bytes.NewBufferString(`{"public_metadata": {"role": "user"}}`))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("### HandleClerkUserCreated error ### --> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change user role"})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("### HandleClerkUserCreated error ### --> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	fmt.Println(res)
	fmt.Println(string(resBody))

	ctx.JSON(http.StatusOK, gin.H{"message": "User role changed successfully" + string(resBody)})

}
