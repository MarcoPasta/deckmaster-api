package repository

import (
	"fmt"
	"log"
	"os"

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
