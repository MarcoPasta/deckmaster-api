package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster-api/internal/models"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdatePokemonAbility, downUpdatePokemonAbility)
}

func upUpdatePokemonAbility(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	err := repository.DB.Migrator().DropColumn(&models.PokemonCard{}, "ability")
	if err != nil {
		return err
	}

	return repository.DB.Migrator().AutoMigrate(&models.PokemonCard{})
}

func downUpdatePokemonAbility(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	err := repository.DB.Migrator().DropColumn(&models.PokemonCard{}, "ability")
	if err != nil {
		return err
	}
	return repository.DB.Migrator().AddColumn(&models.PokemonCard{}, "ability")
}
