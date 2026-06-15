package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddTagsNaruto, downAddTagsNaruto)
}

func upAddTagsNaruto(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return repository.DB.Migrator().AddColumn(&models.NarutoCard{}, "tags")
}

func downAddTagsNaruto(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return repository.DB.Migrator().DropColumn(&models.NarutoCard{}, "tags")
}
