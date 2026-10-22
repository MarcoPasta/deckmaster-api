package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcoPasta/deckmaster-api/internal"
	"github.com/MarcoPasta/deckmaster-api/internal/repository"
	_ "github.com/MarcoPasta/deckmaster-api/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	// dir                                     = flags.String("dir", ".", "directory with migration files")
	host, port, user, dbName, password      = internal.GetDBString()
	tHost, tPort, tUser, tDbName, tPassword = internal.GetTestDBString()
	devDbString                             = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbName)
	testDbString                            = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", tUser, tPassword, tHost, tPort, tDbName)
)

const (
	// Hier bleiben Treiber und Connection String hartcodiert
	driver     = "postgres"
	migrations = "./migrations/"
)

func main() {
	// ###################################################################################
	/// wip
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	dbArg, command := args[0], args[1]
	if dbArg == "dev" {
		dbArg = devDbString
		repository.ConnectToDatabase(internal.GetDBString())
		fmt.Println("DB String: ", dbArg)
	} else {
		dbArg = testDbString
		repository.ConnectToDatabase(internal.GetTestDBString())
		fmt.Println("DB String: ", dbArg)
	}

	// ###################################################################################

	// 1. Verbindung zur Datenbank herstellen
	db, err := goose.OpenDBWithDriver(driver, dbArg)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
		fmt.Println("closing database connection")
	}()

	// 2. Das Migrations-Verzeichnis festlegen
	// Goose wird in diesem Ordner nach .sql Dateien suchen
	if err := goose.SetDialect(driver); err != nil {
		log.Fatalf("Fehler beim Setzen des Dialekts: %v", err)
	}

	// // 3. Argumente aus dem Terminal auslesen
	// // Wir nehmen alles ab dem ersten Argument (z.B. "up", "down", "status")
	arguments := []string{}
	if len(os.Args) > 2 {
		arguments = os.Args[2:]
	} else {
		log.Fatal("Bitte gib einen Befehl an (z.B. ./goose up)")
	}

	// 4. Goose Befehl ausführen
	// Wir übergeben die DB, den Pfad zu den Dateien und den Befehl
	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, migrations, arguments...); err != nil {
		log.Fatalf("Goose Fehler: %v", err)
	}
}
