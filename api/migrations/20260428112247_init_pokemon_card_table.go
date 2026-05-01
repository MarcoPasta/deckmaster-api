package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitPokemonCardTable, downInitPokemonCardTable)
}

func upInitPokemonCardTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return repository.DbMigrator.CreateTable(&models.PokemonCard{})
}

func downInitPokemonCardTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return repository.DbMigrator.DropTable(&models.PokemonCard{})
}
