package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MarcoPasta/deckmaster/api/models"
)

func unmarshalData() ([]models.DigimonCard, error) {
	rawFile, err := os.ReadFile("./unmarshalData.json")
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

func writeToDB(cards []models.DigimonCard) {

	// create db instance
	db := models.DB
	db.Create(&cards)
}

func main() {
	models.ConnectToDatabase()
	cards, err := unmarshalData()
	if err != nil {
		log.Fatal(err)
	}
	writeToDB(cards)
}
