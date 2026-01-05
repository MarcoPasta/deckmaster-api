package migrations

import (
	"context"
	"database/sql"

	"github.com/MarcoPasta/deckmaster/api/models"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTableDbwritetest, downCreateTableDbwritetest)
}

type customTable struct {
	models.DigimonCard
}

func (customTable) TableName() string {
	return "WriteTest"
}

func upCreateTableDbwritetest(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return models.DbMigrator.CreateTable(&customTable{})
}

func downCreateTableDbwritetest(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return models.DbMigrator.DropTable(&customTable{})

}
