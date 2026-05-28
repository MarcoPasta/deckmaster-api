package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLanguageCol, downAddLanguageCol)
}

func upAddLanguageCol(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	if err := repository.DB.Exec("ALTER TABLE pokemon_cards ADD COLUMN language VARCHAR(2) NOT NULL DEFAULT 'de'").Error; err != nil {
		return err
	}

	if err := repository.DB.Exec("ALTER TABLE pokemon_cards ADD CONSTRAINT pokemon_cards_code_number_language_unique UNIQUE (code,card_number, language)").Error; err != nil {
		return err
	}

	return nil
}

func downAddLanguageCol(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	if err := repository.DB.Exec("ALTER TABLE pokemon_cards DROP CONSTRAINT pokemon_cards_code_number_language_unique").Error; err != nil {
		return err
	}
	if err := repository.DB.Exec("ALTER TABLE pokemon_cards DROP COLUMN language").Error; err != nil {
		return err
	}
	return nil
}
