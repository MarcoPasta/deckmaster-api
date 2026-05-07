package repository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DbMigrator gorm.Migrator

func ConnectToDatabase(host string, port string, user string, dbName string, password string) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", host, port, user, dbName, password)
	// fmt.Printf("host=%v port=%v user=%v dbname=%v password=%v", host, port, user, dbName, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	DB = db
	DbMigrator = db.Migrator()
}

func RunMigrations() {
	log.Println("Running database migrations...")
	sqlDb, err := DB.DB()
	if err != nil {
		panic(err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	migrationsDir := os.Getenv("GOOSE_MIGRATION_DIR")
	if migrationsDir == "" {
		migrationsDir = "/migrations"
	}

	if err := goose.Up(sqlDb, migrationsDir); err != nil {
		panic(err)
	}
}

func CheckTableValue() error {
	ctx := context.Background()

	rtvPok, err := gorm.G[models.PokemonCard](DB).Find(ctx)
	log.Println("retrieved pokemon cards")
	if err != nil {
		log.Fatal(err)
	}

	if len(rtvPok) == 0 {
		fmt.Println("No Pokémon cards found in the database. Please run the seeders.")
		return fmt.Errorf("pokemon") // RunSeeder("pokemon")
	}

	rtvNaruto, err := gorm.G[models.NarutoCard](DB).Find(ctx)
	log.Println("retrieved naruto cards")
	if err != nil {
		log.Fatal(err)
	}

	if len(rtvNaruto) == 0 {
		fmt.Println("No Naruto cards found in the database. Please run the seeders.")
		return fmt.Errorf("naruto") // RunSeeder("naruto")
	}

	rtvDigimon, err := gorm.G[models.DigimonCard](DB).Find(ctx)
	log.Println("retrieved digimon cards")
	if err != nil {
		log.Fatal(err)
	}

	if len(rtvDigimon) == 0 {
		fmt.Println("No Digimon cards found in the database. Please run the seeders.")
		return fmt.Errorf("digimon") // RunSeeder("digimon")
	}

	return nil
}
