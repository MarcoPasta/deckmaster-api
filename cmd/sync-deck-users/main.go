package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fix := flag.Bool("fix", false, "Fix orphaned deck records")
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

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repository.DB = db

	fmt.Println("\n=== Deck User Consistency Check ===")
	if *fix {
		fmt.Println("FIX MODE: Orphaned records will be deleted")
	}
	fmt.Println("")

	// Check Pokemon decks
	fmt.Println("Checking Pokemon decks...")
	pokemonValidCount, pokemonOrphanedCount := checkAndFixPokeDecks(db, *fix)

	// Check Digimon decks
	fmt.Println("Checking Digimon decks...")
	digimonValidCount, digimonOrphanedCount := checkAndFixDigimonDecks(db, *fix)

	// Check Naruto decks
	fmt.Println("Checking Naruto decks...")
	narutoValidCount, narutoOrphanedCount := checkAndFixNarutoDecks(db, *fix)

	// Print summary
	fmt.Println("\n=== Consistency Check Summary ===")
	fmt.Printf("Pokemon: %d valid users, %d orphaned\n", pokemonValidCount, pokemonOrphanedCount)
	fmt.Printf("Digimon: %d valid users, %d orphaned\n", digimonValidCount, digimonOrphanedCount)
	fmt.Printf("Naruto: %d valid users, %d orphaned\n", narutoValidCount, narutoOrphanedCount)

	totalOrphaned := pokemonOrphanedCount + digimonOrphanedCount + narutoOrphanedCount

	if totalOrphaned == 0 {
		fmt.Println("\n✓ All deck users are valid!")
	} else {
		if *fix {
			fmt.Printf("\n✓ Cleaned up %d orphaned deck records\n", totalOrphaned)
		} else {
			fmt.Printf("\n⚠ Found %d orphaned deck records\nRun with --fix flag to clean them up\n", totalOrphaned)
		}
	}
}

func checkAndFixPokeDecks(db *gorm.DB, fix bool) (int, int) {
	type PokemonDecks struct {
		ID        uint
		UserId    string
		Name      string
		CoverCard string
	}

	var decks []PokemonDecks
	db.Table("pokemon_decks").Find(&decks)

	validCount := 0
	orphanedCount := 0

	for _, deck := range decks {
		var user models.User
		result := db.Where("id = ?", deck.UserId).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			orphanedCount++
			if fix {
				db.Table("pokemon_decks").Delete(&deck)
				fmt.Printf("  ✗ Deleted orphaned Pokemon deck: %d (user_id: %s)\n", deck.ID, deck.UserId)
			} else {
				fmt.Printf("  ⚠ Orphaned Pokemon deck: %d (user_id: %s)\n", deck.ID, deck.UserId)
			}
		} else {
			validCount++
		}
	}

	fmt.Printf("  ✓ Pokemon: %d valid, %d orphaned\n", validCount, orphanedCount)
	return validCount, orphanedCount
}

func checkAndFixDigimonDecks(db *gorm.DB, fix bool) (int, int) {
	type DigimonDecks struct {
		ID        uint
		UserId    string
		Name      string
		CoverCard string
	}

	var decks []DigimonDecks
	db.Table("digimon_decks").Find(&decks)

	validCount := 0
	orphanedCount := 0

	for _, deck := range decks {
		var user models.User
		result := db.Where("id = ?", deck.UserId).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			orphanedCount++
			if fix {
				db.Table("digimon_decks").Delete(&deck)
				fmt.Printf("  ✗ Deleted orphaned Digimon deck: %d (user_id: %s)\n", deck.ID, deck.UserId)
			} else {
				fmt.Printf("  ⚠ Orphaned Digimon deck: %d (user_id: %s)\n", deck.ID, deck.UserId)
			}
		} else {
			validCount++
		}
	}

	fmt.Printf("  ✓ Digimon: %d valid, %d orphaned\n", validCount, orphanedCount)
	return validCount, orphanedCount
}

func checkAndFixNarutoDecks(db *gorm.DB, fix bool) (int, int) {
	type NarutoDecks struct {
		ID        uint
		UserId    string
		Name      string
		CoverCard string
	}

	var decks []NarutoDecks
	db.Table("naruto_decks").Find(&decks)

	validCount := 0
	orphanedCount := 0

	for _, deck := range decks {
		var user models.User
		result := db.Where("id = ?", deck.UserId).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			orphanedCount++
			if fix {
				db.Table("naruto_decks").Delete(&deck)
				fmt.Printf("  ✗ Deleted orphaned Naruto deck: %d (user_id: %s)\n", deck.ID, deck.UserId)
			} else {
				fmt.Printf("  ⚠ Orphaned Naruto deck: %d (user_id: %s)\n", deck.ID, deck.UserId)
			}
		} else {
			validCount++
		}
	}

	fmt.Printf("  ✓ Naruto: %d valid, %d orphaned\n", validCount, orphanedCount)
	return validCount, orphanedCount
}
