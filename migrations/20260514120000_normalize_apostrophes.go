package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upNormalizeApostrophes, downNormalizeApostrophes)
}

func upNormalizeApostrophes(ctx context.Context, tx *sql.Tx) error {
	curly := "’"
	straight := "'"

	return repository.DB.Exec(
		`UPDATE pokemon_cards SET
			name    = REPLACE(name,    ?, ?),
			ability = REPLACE(ability, ?, ?),
			attacks = REPLACE(attacks, ?, ?)`,
		curly, straight,
		curly, straight,
		curly, straight,
	).Error
}

func downNormalizeApostrophes(ctx context.Context, tx *sql.Tx) error {
	return nil
}
