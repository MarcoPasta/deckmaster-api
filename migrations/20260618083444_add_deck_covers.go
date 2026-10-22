package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddDeckCovers, downAddDeckCovers)
}

func upAddDeckCovers(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	repository.DB.Migrator().AutoMigrate(&models.PokemonDecks{})
	repository.DB.Migrator().AutoMigrate(&models.NarutoDecks{})
	repository.DB.Migrator().AutoMigrate(&models.DigimonDecks{})
	return nil
}

func downAddDeckCovers(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	repository.DB.Migrator().DropColumn(&models.PokemonDecks{}, "cover_card")
	repository.DB.Migrator().DropColumn(&models.NarutoDecks{}, "cover_card")
	repository.DB.Migrator().DropColumn(&models.DigimonDecks{}, "cover_card")
	return nil
}
