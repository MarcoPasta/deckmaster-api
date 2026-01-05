package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/models"
	"gorm.io/gorm"
)

func unmarshalData(file string) ([]models.DigimonCard, error) {
	rawFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var jsonData []models.DigimonCard

	err = json.Unmarshal(rawFile, &jsonData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		// log.Printf("JSON: %#v", jsonData)
		return jsonData, err
	}

}

func writeToDB(db *gorm.DB, cards []models.DigimonCard) {

	// create db instance
	db.Create(&cards)
}

func main() {
	models.ConnectToDatabase(internal.GetDBString())
	db := models.DB
	cards, err := unmarshalData("./BT24.json")
	if err != nil {
		log.Fatal(err)
	}
	writeToDB(db, cards)
}
