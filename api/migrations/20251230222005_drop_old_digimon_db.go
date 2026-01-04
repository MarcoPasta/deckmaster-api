package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/models"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitDigimonDb, downInitDigimonDb)
}

func upInitDigimonDb(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return models.DbMigrator.DropTable(&models.DigimonCard{})
}

func downInitDigimonDb(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return models.DbMigrator.CreateTable(&models.DigimonCard{})
}
