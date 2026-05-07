package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upPokemonAddSetCode, downPokemonAddSetCode)
}

func upPokemonAddSetCode(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return repository.DbMigrator.AddColumn(&models.PokemonCard{}, "code")
}

func downPokemonAddSetCode(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return repository.DbMigrator.DropColumn(&models.PokemonCard{}, "code")
}
