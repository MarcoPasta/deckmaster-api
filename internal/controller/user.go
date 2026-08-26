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

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userID := data["data"]["id"]
	log.Println("### Received Data For ### --> ", userID)

	// Sync user from Clerk to local DB
	localUser, err := repository.SyncUserFromClerkOnFirstLogin(userID)
	if err != nil {
		log.Println("### HandleClerkUserCreated error syncing user ### --> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync user to local database"})
		return
	}

	// Also set role in Clerk for external compatibility
	url := os.Getenv("CLERK_BACKEND_API") + "/users/" + userID

	req, _ := http.NewRequest("PATCH", url, bytes.NewBufferString(`{"public_metadata": {"role": "user"}}`))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("CLERK_SECRET_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("### HandleClerkUserCreated error setting Clerk role ### --> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set user role in Clerk"})
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("### HandleClerkUserCreated error reading response ### --> ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	fmt.Println(res)
	fmt.Println(string(resBody))

	ctx.JSON(http.StatusOK, gin.H{"message": "User created and synced successfully", "user": localUser})

}

// GetUserProfile returns the authenticated user's profile with social links
func GetUserProfile(ctx *gin.Context) {
	localUser, exists := ctx.Get("localUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := localUser.(*models.User)

	// Load social links
	var socialLinks []models.UserSocialLink
	if err := repository.DB.Where("user_id = ?", user.ID).Find(&socialLinks).Error; err != nil {
		log.Println("Failed to fetch social links:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social links"})
		return
	}

	user.SocialLinks = socialLinks

	ctx.JSON(http.StatusOK, user)
}

// UpdateUserProfile updates the authenticated user's profile
func UpdateUserProfile(ctx *gin.Context) {
	localUser, exists := ctx.Get("localUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := localUser.(*models.User)

	var updateData struct {
		FirstName   string                `json:"first_name"`
		LastName    string                `json:"last_name"`
		Phone       string                `json:"phone_number"`
		SocialLinks []models.UserSocialLink `json:"social_links"`
	}

	if err := ctx.BindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate fields
	validationErrors := validateProfileUpdate(updateData.FirstName, updateData.LastName, updateData.Phone, updateData.SocialLinks)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Update user fields
	if err := repository.DB.Model(user).Updates(map[string]interface{}{
		"first_name":   updateData.FirstName,
		"last_name":    updateData.LastName,
		"phone_number": updateData.Phone,
	}).Error; err != nil {
		log.Println("Failed to update user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// If social links provided, update them
	if len(updateData.SocialLinks) > 0 {
		// Delete existing social links
		if err := repository.DB.Where("user_id = ?", user.ID).Delete(&models.UserSocialLink{}).Error; err != nil {
			log.Println("Failed to delete existing social links:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social links"})
			return
		}

		// Create new social links
		for _, link := range updateData.SocialLinks {
			link.UserID = user.ID
			link.ID = uuid.New().String()
			if err := repository.DB.Create(&link).Error; err != nil {
				log.Println("Failed to create social link:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social link"})
				return
			}
		}
	}

	// Load updated user with social links
	var updatedUser models.User
	if err := repository.DB.Preload("SocialLinks").First(&updatedUser, "id = ?", user.ID).Error; err != nil {
		log.Println("Failed to fetch updated user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

// validateProfileUpdate validates profile update data
func validateProfileUpdate(firstName, lastName, phone string, socialLinks []models.UserSocialLink) []string {
	var errors []string

	// Validate first name
	if len(firstName) > 0 && len(firstName) > 100 {
		errors = append(errors, "First name must be 100 characters or less")
	}

	// Validate last name
	if len(lastName) > 0 && len(lastName) > 100 {
		errors = append(errors, "Last name must be 100 characters or less")
	}

	// Validate phone number
	if len(phone) > 0 && len(phone) > 20 {
		errors = append(errors, "Phone number must be 20 characters or less")
	}

	// Validate social links
	for i, link := range socialLinks {
		if link.URL == "" {
			errors = append(errors, fmt.Sprintf("Social link %d: URL is required", i+1))
		} else if len(link.URL) > 500 {
			errors = append(errors, fmt.Sprintf("Social link %d: URL must be 500 characters or less", i+1))
		}

		if link.Platform == "" {
			errors = append(errors, fmt.Sprintf("Social link %d: Platform is required", i+1))
		} else if link.Platform != "twitter" && link.Platform != "discord" && link.Platform != "youtube" && link.Platform != "custom" {
			errors = append(errors, fmt.Sprintf("Social link %d: Invalid platform. Must be one of: twitter, discord, youtube, custom", i+1))
		}

		if link.Platform == "custom" && link.Label == "" {
			errors = append(errors, fmt.Sprintf("Social link %d: Label is required for custom platforms", i+1))
		}
	}

	return errors
}

// GetSocialLinks returns all social links for the authenticated user
func GetSocialLinks(ctx *gin.Context) {
	localUser, exists := ctx.Get("localUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := localUser.(*models.User)

	var socialLinks []models.UserSocialLink
	if err := repository.DB.Where("user_id = ?", user.ID).Find(&socialLinks).Error; err != nil {
		log.Println("Failed to fetch social links:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social links"})
		return
	}

	ctx.JSON(http.StatusOK, socialLinks)
}

// CreateSocialLink creates a new social link for the authenticated user
func CreateSocialLink(ctx *gin.Context) {
	localUser, exists := ctx.Get("localUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := localUser.(*models.User)

	var input struct {
		Platform string `json:"platform" binding:"required"`
		URL      string `json:"url" binding:"required"`
		Label    string `json:"label"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	link := models.UserSocialLink{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Platform:  input.Platform,
		URL:       input.URL,
		Label:     input.Label,
	}

	if err := repository.DB.Create(&link).Error; err != nil {
		log.Println("Failed to create social link:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social link"})
		return
	}

	ctx.JSON(http.StatusCreated, link)
}

// UpdateSocialLink updates an existing social link
func UpdateSocialLink(ctx *gin.Context) {
	localUser, exists := ctx.Get("localUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := localUser.(*models.User)
	linkID := ctx.Param("id")

	var link models.UserSocialLink
	if err := repository.DB.Where("id = ? AND user_id = ?", linkID, user.ID).First(&link).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Social link not found"})
		return
	}

	var input struct {
		Platform string `json:"platform"`
		URL      string `json:"url"`
		Label    string `json:"label"`
	}

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updateData := map[string]interface{}{}
	if input.Platform != "" {
		updateData["platform"] = input.Platform
	}
	if input.URL != "" {
		updateData["url"] = input.URL
	}
	updateData["label"] = input.Label

	if err := repository.DB.Model(&link).Updates(updateData).Error; err != nil {
		log.Println("Failed to update social link:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social link"})
		return
	}

	ctx.JSON(http.StatusOK, link)
}

// DeleteSocialLink deletes a social link
func DeleteSocialLink(ctx *gin.Context) {
	localUser, exists := ctx.Get("localUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := localUser.(*models.User)
	linkID := ctx.Param("id")

	var link models.UserSocialLink
	if err := repository.DB.Where("id = ? AND user_id = ?", linkID, user.ID).First(&link).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Social link not found"})
		return
	}

	if err := repository.DB.Delete(&link).Error; err != nil {
		log.Println("Failed to delete social link:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social link"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Social link deleted successfully"})
}

// AdminResyncUser resyncs a user from Clerk (admin only)
func AdminResyncUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	user, err := repository.UpdateUserFromClerk(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resync user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"synced": true, "user": user})
}

// AdminResyncAllUsers resyncs all users from Clerk (admin only)
func AdminResyncAllUsers(ctx *gin.Context) {
	var users []models.User
	if err := repository.DB.Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	syncedCount := 0
	errorCount := 0

	for _, user := range users {
		_, err := repository.UpdateUserFromClerk(user.ID)
		if err != nil {
			errorCount++
		} else {
			syncedCount++
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"synced": syncedCount,
		"errors": errorCount,
		"total":  len(users),
	})
}

// LogPasswordChange logs a password change event (called from frontend after Clerk handles the change)
func LogPasswordChange(ctx *gin.Context) {
	localUser, exists := ctx.Get("localUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user := localUser.(*models.User)

	// Get client IP and user agent
	clientIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	// Log the password change
	if err := repository.LogPasswordChange(user.ID, clientIP, userAgent); err != nil {
		log.Println("Failed to log password change:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log password change"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password change logged successfully"})
}
