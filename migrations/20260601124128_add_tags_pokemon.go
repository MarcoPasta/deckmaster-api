package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddTagsPokemon, downAddTagsPokemon)
}

func upAddTagsPokemon(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return repository.DB.Migrator().AutoMigrate(&models.PokemonCard{})
}

func downAddTagsPokemon(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return repository.DB.Migrator().DropColumn(&models.PokemonCard{}, "Tags")
}
