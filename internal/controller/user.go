package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MarcoPasta/deckmaster-api/internal/models"

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

// TODO: implement future proof error code handling, this is only for now
var errorCodes = map[string]int{
	"success":                    200,
	"request_not_succesful":      400,
	"authentication_invalid":     401,
	"resource_not_found":         404,
	"invalid_request_parameters": 412,
}

func HandleClerkUserCreated(ctx *gin.Context) {

	var hook models.ClerkWebhook
	err := ctx.ShouldBindJSON(&hook)
	if err != nil {
		fmt.Println(err)
	}

	payload := strings.NewReader(`{
		"public_metadata": {
			"role": "user"
		},
		"private_metadata": {},
		"unsafe_metadata": {}
	}`)

	req, err := http.NewRequest("PATCH", os.Getenv("CLERK_BACKEND_API")+fmt.Sprintf("/users/%v/metadata", hook.Data.Id), payload)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	var codes models.ClerkErrorResponse
	err = json.Unmarshal(body, &codes)
	if err != nil {
		log.Println(err)
	}

	if len(codes.Errors) > 0 {
		ctx.JSON(errorCodes[codes.Errors[0].Code], codes.Errors[0].Message)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "user role created successfully"})
	}

}

// GetUserProfile returns the authenticated user's profile with social links
func GetUserProfile(ctx *gin.Context) {}

// GetSocialLinks returns all social links for the authenticated user
func GetSocialLinks(ctx *gin.Context) {}

// CreateSocialLink creates a new social link for the authenticated user
func CreateSocialLink(ctx *gin.Context) {}

// UpdateSocialLink updates an existing social link
func UpdateSocialLink(ctx *gin.Context) {}

// DeleteSocialLink deletes a social link
func DeleteSocialLink(ctx *gin.Context) {}
