package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DbMigrator gorm.Migrator

func ConnectToDatabase() {
	dsn := "host=db port=5432 user=postgres dbname=deckmaster password=example"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	DB = db
	DbMigrator = db.Migrator()
}
