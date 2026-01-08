package models

import (
	"fmt"
	"log"

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
