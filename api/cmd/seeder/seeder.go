package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MarcoPasta/deckmaster/api/models"
)

func main() {
	models.ConnectToDatabase()

	rawFile, err := os.ReadFile("../../data/jsonDigimon/BT24.json")
	if err != nil {
		log.Fatal(err)
	}

	var jsonData *any

	err = json.Unmarshal(rawFile, jsonData)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("JSON: %#v", jsonData)

	// ctx := context.Background()
	// err = gorm.G[models.DigimonDecks](models.DB).Create(ctx, &deck)
}
