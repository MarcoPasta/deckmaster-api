package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateDigimonDecksRemoveAltArts, downUpdateDigimonDecksRemoveAltArts)
}

func upUpdateDigimonDecksRemoveAltArts(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	err := repository.DbMigrator.DropTable(&models.DigimonCardAltArt{})
	if err != nil {
		return err
	}

	return repository.DB.Migrator().AutoMigrate(&models.DigimonDecks{}, &models.DigimonDeckEntry{})
}

func downUpdateDigimonDecksRemoveAltArts(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return repository.DbMigrator.DropTable(&models.DigimonDecks{}, &models.DigimonDeckEntry{})
}
