package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func countIdDuplicates(cards []string, id string) int {
	counter := 0
	for _, card := range cards {
		if card == id {
			counter++
		}
	}
	return counter
}

func writeToDB(db *gorm.DB, cards []models.DigimonCard) {
	ctx := context.Background()
	result := gorm.WithResult()

	var writtenCards []string
	for _, card := range cards {
		fmt.Println("Current Card: ", card.Id)
		if counter := countIdDuplicates(writtenCards, card.Id); counter == 0 {
			err := gorm.G[models.DigimonCard](db, result).Create(ctx, &card)
			if err != nil {
				log.Fatalf("seeder.go::writeToDB -> Error writing card %v to DB: %v", card.Id, err)
			}
			writtenCards = append(writtenCards, card.Id)
			fmt.Printf("Currently written Cards:\n %+v\n", writtenCards)
		} else {
			fmt.Printf("Duplicate found at: %v; creating alt art...\n", card.Id)
			altArt := models.DigimonCardAltArt{CardID: card.Id, Image: fmt.Sprintf("%v_P%v", card.Id, counter)}
			fmt.Println("Alt Art Id->", altArt.Image)
			err := gorm.G[models.DigimonCardAltArt](db, result).Create(ctx, &altArt)
			if err != nil {
				log.Fatalf("seeder.go::writeToDB -> Error writing card %v to DB: %v", card.Id, err)
			}
			writtenCards = append(writtenCards, card.Id)
		}
		fmt.Println("Card written: ", card.Id)
	}

	defer 
}

func readDirectory(path string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, fmt.Sprintf("%v/%v", path, entry.Name()))
		}
	}
	defer entries.Close()
	return files, nil
}

func main() {
	models.ConnectToDatabase(internal.GetDBString())
	// db := models.DB
	// cards, err := unmarshalData("./BT24.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// writeToDB(db, cards)
}
