package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	migrations = "./migrations/"
)

func main() {

	// 1. Pfad der Executable ermitteln (z.B. /home/user/project/build/goose)
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	fmt.Println(exePath)

	// 2. Das Verzeichnis der Executable bestimmen (/home/user/project/build)
	exeDir := filepath.Dir(exePath)
	fmt.Println(exeDir)

	// 3. Den Pfad relativ zur Binary berechnen
	// Join kümmert sich um die korrekten Trennzeichen (/ oder \)
	migrationPath := filepath.Join(exeDir, "..", "migrations")

	// Zur Kontrolle ausgeben
	fmt.Println("Suche (clean) Migrations in:", filepath.Clean(migrationPath))
	fmt.Println("Suche Migrations in:", migrationPath)

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
