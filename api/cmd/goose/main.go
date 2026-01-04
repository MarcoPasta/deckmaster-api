package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/MarcoPasta/deckmaster/api/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

const (
	// Hier bleiben Treiber und Connection String hartcodiert
	driver = "postgres"
	// dbString   = "user=postgres password=example dbname=deckmaster sslmode=disable"
	dbString   = "postgres://postgres:example@db:5432/deckmaster?sslmode=disable"
	migrations = "../migrations/"
)

func main() {

	// 1. Verbindung zur Datenbank herstellen
	db, err := sql.Open(driver, dbString)
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der DB: %v", err)
	}
	defer db.Close()

	// 2. Das Migrations-Verzeichnis festlegen
	// Goose wird in diesem Ordner nach .sql Dateien suchen
	if err := goose.SetDialect(driver); err != nil {
		log.Fatalf("Fehler beim Setzen des Dialekts: %v", err)
	}

	// 3. Argumente aus dem Terminal auslesen
	// Wir nehmen alles ab dem ersten Argument (z.B. "up", "down", "status")
	arguments := []string{}
	if len(os.Args) > 1 {
		arguments = os.Args[1:]
	} else {
		log.Fatal("Bitte gib einen Befehl an (z.B. ./goose up)")
	}

	command := arguments[0]

	// 4. Goose Befehl ausführen
	// Wir übergeben die DB, den Pfad zu den Dateien und den Befehl
	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, migrations, arguments...); err != nil {
		log.Fatalf("Goose Fehler: %v", err)
	}
}
