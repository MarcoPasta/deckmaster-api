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
	host, port, user, dbName, password = internal.GetDBString()
	dbString                           = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbName)
)

const (
	// Hier bleiben Treiber und Connection String hartcodiert
	driver     = "postgres"
	migrations = "./migrations/"
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	repository.ConnectToDatabase(internal.GetDBString())

	db, err := goose.OpenDBWithDriver(driver, dbString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
		fmt.Println("closing database connection")
	}()

	if err := goose.SetDialect(driver); err != nil {
		log.Fatalf("Fehler beim Setzen des Dialekts: %v", err)
	}

	arguments := []string{}
	if len(os.Args) > 1 {
		arguments = append(arguments, args[1:]...)
	} else {
		log.Fatal("Bitte gib einen Befehl an (z.B. ./goose up)")
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, migrations, arguments...); err != nil {
		log.Fatalf("Goose Fehler: %v", err)
	}
}
