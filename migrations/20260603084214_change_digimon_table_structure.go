package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upChangeDigimonTableStructure, downChangeDigimonTableStructure)
}

func upChangeDigimonTableStructure(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	err := repository.DB.Migrator().DropTable(&models.DigimonCard{})
	if err != nil {
		return err
	}
	return repository.DB.Migrator().CreateTable(&models.DigimonCard{})
}

func downChangeDigimonTableStructure(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.

	return repository.DB.Migrator().DropTable(&models.DigimonCard{})
}
