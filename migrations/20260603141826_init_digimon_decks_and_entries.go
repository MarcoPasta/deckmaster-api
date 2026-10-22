package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitDigimonDecksAndEntries, downInitDigimonDecksAndEntries)
}

func upInitDigimonDecksAndEntries(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return repository.DB.Migrator().AutoMigrate(&models.DigimonDecks{}, &models.DigimonDeckEntry{})
}

func downInitDigimonDecksAndEntries(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return repository.DB.Migrator().DropTable(&models.DigimonDecks{}, &models.DigimonDeckEntry{})
}
