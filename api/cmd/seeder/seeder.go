package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"gorm.io/gorm"
)

type TCard interface {
	models.DigimonCard | models.PokemonCard | models.NarutoCard
}

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
	fmt.Println("found files: ", files)
	return files, nil
}

func unmarshalData(path string, file string, mode string) ([]models.DigimonCard, []models.PokemonCard, []models.NarutoCard, error) {
	rawFile, err := os.ReadFile(path + "/" + file)
	if err != nil {
		log.Fatalf("seeder.go::unmarshalData() -> Error reading file: %v\n", err)
		return nil, nil, nil, err
	}

	var digimonData []models.DigimonCard
	var pokemonData []models.PokemonCard
	var narutoData []models.NarutoCard

	if mode == "digimon" {
		err = json.Unmarshal(rawFile, &digimonData)
		if err != nil {
			log.Fatalf("seeder.go::unmarshalData() -> Error unmarshaling data: %v\n", err)
			return nil, nil, nil, err
		} else {
			return digimonData, nil, nil, err
		}
	} else if mode == "pokemon" {
		err = json.Unmarshal(rawFile, &pokemonData)
		if err != nil {
			log.Fatalf("seeder.go::unmarshalData() -> Error unmarshaling data: %v\n", err)
			return nil, nil, nil, err
		} else {
			return nil, pokemonData, nil, err
		}
	} else if mode == "naruto" {
		err = json.Unmarshal(rawFile, &narutoData)
		if err != nil {
			log.Fatalf("seeder.go::unmarshalData() -> Error unmarshaling data: %v\n", err)
			return nil, nil, nil, err
		} else {
			fmt.Println("Successfully unmarshaled Naruto data: ", narutoData)
			return nil, nil, narutoData, err
		}
	} else {
		return nil, nil, nil, fmt.Errorf("seeder.go::unmarshalData() -> Invalid mode: %s\n", mode)
	}

}

// func countIdDuplicates(cards []string, id string) int {
// 	counter := 0
// 	for _, card := range cards {
// 		if card == id {
// 			counter++
// 		}
// 	}
// 	return counter
// }

func writeToDB[T TCard](db *gorm.DB, cards []T, writtenCards *[]string) error {
	ctx := context.Background()
	result := gorm.WithResult()

	for _, card := range cards {
		err := gorm.G[T](db, result).Create(ctx, &card)
		if err != nil {
			fmt.Printf("seeder.go::writeToDB() -> Error writing card %v to DB: %v", card, err)
			return err
		}
	}
	*writtenCards = append(*writtenCards, "test")
	return nil
}

func run(db *gorm.DB, path string, mode string) error {
	files, err := readDirectory(path)
	if err != nil {
		fmt.Printf("seeder.go::run() -> Returned error reading directory: %v\n", err)
		return err
	}

	var writtenCards []string
	for _, file := range files {
		digimonCards, pokemonCards, narutoCards, err := unmarshalData(path, file, mode)
		if err != nil {
			fmt.Printf("seeder.go::run() -> Returned error unmarshaling data from file %v: %v\n", file, err)
			return err
		}
		if mode == "digimon" {
			err = writeToDB(db, digimonCards, &writtenCards)
		} else if mode == "pokemon" {
			err = writeToDB(db, pokemonCards, &writtenCards)
		} else if mode == "naruto" {
			err = writeToDB(db, narutoCards, &writtenCards)
		}
		if err != nil {
			fmt.Printf("seeder.go::run() -> Returned error writing to database: %v", err)
			return err
		}
	}

	return nil
}

func main() {
	mode := flag.String("mode", "pokemon", "ich schreib mal iwas rein")
	flag.Parse()
	var seedDir string
	if *mode == "pokemon" {
		seedDir = "POKEMON_SEEDER_DIR"
	} else if *mode == "digimon" {
		seedDir = "DIGIMON_SEEDER_DIR"
	} else if *mode == "naruto" {
		seedDir = "NARUTO_SEEDER_DIR"
		fmt.Println("Running Naruto seeder...")
	} else {
		log.Fatalf("seeder.go::main() -> Invalid mode: %s", *mode)
	}
	repository.ConnectToDatabase(internal.GetDBString())
	run(repository.DB, os.Getenv(seedDir), *mode)
}
