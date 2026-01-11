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

func readDirectory(path string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func unmarshalData(path, file string) ([]models.DigimonCard, error) {
	rawFile, err := os.ReadFile(path + "/" + file)
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

func writeToDB(db *gorm.DB, cards []models.DigimonCard, writtenCards *[]string) error {
	ctx := context.Background()
	result := gorm.WithResult()

	for _, card := range cards {
		if counter := countIdDuplicates(*writtenCards, card.Id); counter == 0 {
			err := gorm.G[models.DigimonCard](db, result).Create(ctx, &card)
			if err != nil {
				fmt.Printf("seeder.go::writeToDB() -> Error writing card %v to DB: %v", card.Id, err)
				return err
			}
		} else {
			altArt := models.DigimonCardAltArt{CardID: card.Id, Image: fmt.Sprintf("%v_P%v", card.Id, counter)}
			err := gorm.G[models.DigimonCardAltArt](db, result).Create(ctx, &altArt)
			if err != nil {
				fmt.Printf("seeder.go::writeToDB() -> Error writing card %v to DB: %v", card.Id, err)
				return err
			}
		}
		*writtenCards = append(*writtenCards, card.Id)
	}
	return nil
}

func run(db *gorm.DB, path string) error {

	files, err := readDirectory(path)
	if err != nil {
		fmt.Printf("seeder.go::run() -> Returned error reading directory: %v", err)
		return err
	}

	var writtenCards []string
	for _, file := range files {
		cards, err := unmarshalData(path, file)
		if err != nil {
			fmt.Printf("seeder.go::run() -> Returned error unmarshaling data from file %v: %v", file, err)
			return err
		}
		err = writeToDB(db, cards, &writtenCards)
		fmt.Println("Written cards after function", writtenCards)
		if err != nil {
			fmt.Printf("seeder.go::run() -> Returned error writing to database: %v", err)
			return err
		}
	}

	return nil
}

func main() {
	models.ConnectToDatabase(internal.GetDBString())
	run(models.DB, os.Getenv("SEEDER_FILE_DIR"))
}
