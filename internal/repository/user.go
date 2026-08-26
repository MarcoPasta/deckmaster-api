package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/google/uuid"
)

// SyncUserFromClerkOnFirstLogin fetches user data from Clerk and syncs to local DB
func SyncUserFromClerkOnFirstLogin(clerkUserID string) (*models.User, error) {
	ctx := context.Background()

	// Fetch user from Clerk API
	clerkUser, err := user.Get(ctx, clerkUserID)
	if err != nil {
		log.Printf("Failed to fetch user from Clerk: %v", err)
		return nil, err
	}

	// Extract email
	email := ""
	if len(clerkUser.EmailAddresses) > 0 {
		email = clerkUser.EmailAddresses[0].EmailAddress
	}

	// Extract phone number
	phoneNumber := ""
	if len(clerkUser.PhoneNumbers) > 0 {
		phoneNumber = clerkUser.PhoneNumbers[0].PhoneNumber
	}

	// Extract image URL
	imageURL := ""
	if clerkUser.ImageURL != nil {
		imageURL = *clerkUser.ImageURL
	}

	// Get first name
	firstName := ""
	if clerkUser.FirstName != nil {
		firstName = *clerkUser.FirstName
	}

	// Get last name
	lastName := ""
	if clerkUser.LastName != nil {
		lastName = *clerkUser.LastName
	}

	// Get role from public_metadata (default to "user")
	role := "user"
	if clerkUser.PublicMetadata != nil {
		// Parse public_metadata as map to extract role
		var metadata map[string]interface{}
		if err := json.Unmarshal(clerkUser.PublicMetadata, &metadata); err == nil {
			if roleValue, ok := metadata["role"]; ok {
				if roleStr, ok := roleValue.(string); ok {
					role = roleStr
				}
			}
		}
	}

	now := time.Now()

	// Create local user record
	localUser := &models.User{
		ID:              clerkUser.ID,
		Email:           email,
		FirstName:       firstName,
		LastName:        lastName,
		ProfileImageUrl: imageURL,
		PhoneNumber:     phoneNumber,
		Role:            role,
		ClerkSyncedAt:   &now,
	}

	// Insert into local DB
	if err := DB.Create(localUser).Error; err != nil {
		log.Printf("Failed to sync user to local DB: %v", err)
		return nil, err
	}

	log.Printf("Successfully synced user %s to local DB", clerkUserID)
	return localUser, nil
}

// GetOrSyncUser retrieves user from local DB or syncs from Clerk if not exists
func GetOrSyncUser(clerkUserID string) (*models.User, error) {
	var localUser models.User
	result := DB.Where("id = ?", clerkUserID).First(&localUser)

	// User exists in local DB
	if result.Error == nil {
		return &localUser, nil
	}

	// User not found, sync from Clerk
	if result.Error.Error() == "record not found" {
		return SyncUserFromClerkOnFirstLogin(clerkUserID)
	}

	// Other error occurred
	return nil, result.Error
}

// UpdateUserFromClerk updates user data from Clerk
func UpdateUserFromClerk(clerkUserID string) (*models.User, error) {
	ctx := context.Background()

	// Fetch user from Clerk API
	clerkUser, err := user.Get(ctx, clerkUserID)
	if err != nil {
		log.Printf("Failed to fetch user from Clerk: %v", err)
		return nil, err
	}

	// Extract email
	email := ""
	if len(clerkUser.EmailAddresses) > 0 {
		email = clerkUser.EmailAddresses[0].EmailAddress
	}

	// Extract phone number
	phoneNumber := ""
	if len(clerkUser.PhoneNumbers) > 0 {
		phoneNumber = clerkUser.PhoneNumbers[0].PhoneNumber
	}

	// Extract image URL
	imageURL := ""
	if clerkUser.ImageURL != nil {
		imageURL = *clerkUser.ImageURL
	}

	// Get first name
	firstName := ""
	if clerkUser.FirstName != nil {
		firstName = *clerkUser.FirstName
	}

	// Get last name
	lastName := ""
	if clerkUser.LastName != nil {
		lastName = *clerkUser.LastName
	}

	now := time.Now()

	// Update local user record
	updateData := map[string]interface{}{
		"email":               email,
		"first_name":         firstName,
		"last_name":          lastName,
		"profile_image_url":  imageURL,
		"phone_number":       phoneNumber,
		"clerk_synced_at":    &now,
	}

	var localUser models.User
	if err := DB.Model(&models.User{}).Where("id = ?", clerkUserID).Updates(updateData).First(&localUser, "id = ?", clerkUserID).Error; err != nil {
		log.Printf("Failed to update user in local DB: %v", err)
		return nil, err
	}

	log.Printf("Successfully updated user %s in local DB", clerkUserID)
	return &localUser, nil
}

// LogPasswordChange logs a password change event for audit trail
func LogPasswordChange(userID string, ipAddress string, userAgent string) error {
	logEntry := &models.PasswordChangeLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		ChangedAt: time.Now(),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	if err := DB.Create(logEntry).Error; err != nil {
		log.Printf("Failed to log password change for user %s: %v", userID, err)
		return err
	}

	return nil
}
