package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/models"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDbInitDigimon, downDbInitDigimon)
}

func upDbInitDigimon(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return models.DbMigrator.DropTable(&models.DigimonCard{})
}

func downDbInitDigimon(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return models.DbMigrator.DropTable(&models.DigimonCard{})
}
