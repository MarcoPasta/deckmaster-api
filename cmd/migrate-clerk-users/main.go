package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Preview migration without committing changes")
	flag.Parse()

	// Load environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "example"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "deckmaster"
	}
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is required")
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repository.DB = db

	// Set Clerk secret key
	clerk.SetKey(clerkSecretKey)

	fmt.Println("\n=== Clerk User Migration ===")
	if *dryRun {
		fmt.Println("DRY-RUN MODE: No changes will be committed")
	}
	fmt.Println("")

	ctx := context.Background()

	// List all users from Clerk
	fmt.Println("Fetching users from Clerk...")
	clerkUserList, err := clerkUser.List(ctx, &clerkUser.ListParams{})
	if err != nil {
		log.Fatalf("Failed to list users from Clerk: %v", err)
	}

	syncedCount := 0
	skippedCount := 0
	errorCount := 0

	now := time.Now()

	// Iterate through each user
	for _, clerkUserData := range clerkUserList.Users {
		userID := clerkUserData.ID

		// Check if user already exists
		var existingUser models.User
		result := db.Where("id = ?", userID).First(&existingUser)

		if result.Error == nil {
			email := ""
			if len(clerkUserData.EmailAddresses) > 0 {
				email = clerkUserData.EmailAddresses[0].EmailAddress
			}
			fmt.Printf("⊘ User %d: %s (%s) - Already exists, skipping\n", skippedCount+1, userID, email)
			skippedCount++
			continue
		}

		// Extract user data
		email := ""
		if len(clerkUserData.EmailAddresses) > 0 {
			email = clerkUserData.EmailAddresses[0].EmailAddress
		}

		phoneNumber := ""
		if len(clerkUserData.PhoneNumbers) > 0 {
			phoneNumber = clerkUserData.PhoneNumbers[0].PhoneNumber
		}

		imageURL := ""
		if clerkUserData.ImageURL != nil {
			imageURL = *clerkUserData.ImageURL
		}

		firstName := ""
		if clerkUserData.FirstName != nil {
			firstName = *clerkUserData.FirstName
		}

		lastName := ""
		if clerkUserData.LastName != nil {
			lastName = *clerkUserData.LastName
		}

		// Get role from public_metadata
		role := "user"
		if clerkUserData.PublicMetadata != nil {
			// Parse public_metadata as map to extract role
			var metadata map[string]interface{}
			if err := json.Unmarshal(clerkUserData.PublicMetadata, &metadata); err == nil {
				if roleValue, ok := metadata["role"]; ok {
					if roleStr, ok := roleValue.(string); ok {
						role = roleStr
					}
				}
			}
		}

		user := models.User{
			ID:              userID,
			Email:           email,
			FirstName:       firstName,
			LastName:        lastName,
			ProfileImageUrl: imageURL,
			PhoneNumber:     phoneNumber,
			Role:            role,
			ClerkSyncedAt:   &now,
		}

		if !*dryRun {
			// Insert into local database
			if err := db.Create(&user).Error; err != nil {
				fmt.Printf("✗ User %d: %s - Error: %v\n", syncedCount+errorCount+1, userID, err)
				errorCount++
				continue
			}
		}

		fmt.Printf("✓ User %d: %s (%s)\n", syncedCount+1, userID, email)
		syncedCount++
	}

	// Print summary
	fmt.Println("\n=== Migration Summary ===")
	fmt.Printf("✓ Successfully synced: %d users\n", syncedCount)
	fmt.Printf("⊘ Already existed: %d users\n", skippedCount)
	fmt.Printf("✗ Errors: %d users\n", errorCount)
	fmt.Printf("Total processed: %d users\n", syncedCount+skippedCount+errorCount)

	if *dryRun {
		fmt.Println("\nDRY-RUN: No changes were committed to the database")
	}
}
