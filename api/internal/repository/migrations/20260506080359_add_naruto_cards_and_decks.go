package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddNarutoCardsAndDecks, downAddNarutoCardsAndDecks)
}

func upAddNarutoCardsAndDecks(ctx context.Context, tx *sql.Tx) error {
	// Create NarutoCard first (referenced by NarutoDeckEntry)
	if err := repository.DB.AutoMigrate(&models.NarutoCard{}); err != nil {
		return err
	}
	// Then NarutoDecks and NarutoDeckEntry (which references both)
	return repository.DB.AutoMigrate(&models.NarutoDecks{}, &models.NarutoDeckEntry{})
}

func downAddNarutoCardsAndDecks(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return repository.DB.Migrator().DropTable(&models.NarutoCard{}, &models.NarutoDeckEntry{}, &models.NarutoDecks{})
}
