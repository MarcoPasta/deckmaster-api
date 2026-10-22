package migrations

import (
	"context"
	"database/sql"
	"log"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLanguageCol, downAddLanguageCol)
}

func upAddLanguageCol(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	// Check if column exists, only add if it doesn't (for DBs that ran before this migration)
	if !repository.DB.Migrator().HasColumn(&models.PokemonCard{}, "language") {
		if err := repository.DB.Exec("ALTER TABLE pokemon_cards ADD COLUMN language VARCHAR(2) NOT NULL DEFAULT 'de'").Error; err != nil {
			return err
		}
	}

	// Only add constraint if it doesn't already exist (idempotent for existing DBs)
	if !repository.DB.Migrator().HasConstraint(&models.PokemonCard{}, "pokemon_cards_code_number_language_unique") {
		if err := repository.DB.Exec("ALTER TABLE pokemon_cards ADD CONSTRAINT pokemon_cards_code_number_language_unique UNIQUE (code,card_number,language)").Error; err != nil {
			log.Printf("Warning: Could not add constraint pokemon_cards_code_number_language_unique: %v", err)
			// For fresh databases where AutoMigrate already created the table with language,
			// the constraint might already exist or fail for other reasons. Log but don't fail.
			// Try with IF NOT EXISTS syntax for PostgreSQL
			if err := repository.DB.Exec("ALTER TABLE pokemon_cards ADD CONSTRAINT IF NOT EXISTS pokemon_cards_code_number_language_unique UNIQUE (code,card_number,language)").Error; err != nil {
				log.Printf("Error adding constraint with IF NOT EXISTS: %v", err)
				// Continue anyway - constraint might already exist
			}
		}
	}

	return nil
}

func downAddLanguageCol(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	if err := repository.DB.Exec("ALTER TABLE pokemon_cards DROP CONSTRAINT IF EXISTS pokemon_cards_code_number_language_unique").Error; err != nil {
		log.Printf("Warning: Could not drop constraint: %v", err)
	}
	if err := repository.DB.Exec("ALTER TABLE pokemon_cards DROP COLUMN IF EXISTS language").Error; err != nil {
		log.Printf("Warning: Could not drop column: %v", err)
		return err
	}
	return nil
}
