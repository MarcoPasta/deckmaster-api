package seeder

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

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
	// fmt.Println("found files: ", files)
	return files, nil
}

func unmarshalData(path string, file string, mode string) ([]models.DigimonCard, []models.PokemonCard, []models.NarutoCard, error) {
	rawFile, err := os.ReadFile(filepath.Join(path, file))
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
			// fmt.Println("Successfully unmarshaled Naruto data: ", narutoData)
			return nil, nil, narutoData, err
		}
	} else {
		return nil, nil, nil, fmt.Errorf("seeder.go::unmarshalData() -> Invalid mode: %s\n", mode)
	}

}

func containsIdDuplicate(writtenCards *[]string, card models.DigimonCard, mode string) bool {

	// log.Println("Checking card with ID:", card.Id)
	if slices.Contains(*writtenCards, string(card.Id)) {
		// log.Println("Duplicate ID found:", card.Id)
		return true
	} else {

		*writtenCards = append(*writtenCards, card.Id)
		return false
	}

	// switch mode {
	// case "digimon":
	// case "pokemon":
	// 	id := reflect.ValueOf(card).FieldByName("Id")
	// 	log.Println("Checking card with ID:", id.Uint())
	// 	if slices.Contains(*writtenCards, string(id.Uint())) {
	// 		log.Println("Duplicate ID found:", id)
	// 		return true
	// 	} else {
	// 		*writtenCards = append(*writtenCards, string(id.Uint()))
	// 	}
	// case "naruto":
	// 	id := reflect.ValueOf(card).FieldByName("Id")
	// 	log.Println("Checking card with ID:", id.Uint())
	// 	if slices.Contains(*writtenCards, string(id.Uint())) {
	// 		log.Println("Duplicate ID found:", id)
	// 		return true
	// 	} else {
	// 		*writtenCards = append(*writtenCards, string(id.Uint()))
	// 	}
	// }
}

func cardMapper[T TCard](cards []T) []models.DigimonCard {
	result := make([]models.DigimonCard, len(cards))
	for i, card := range cards {
		result[i] = any(card).(models.DigimonCard)
	}

	return result
}

func writeToDB[T TCard](db *gorm.DB, cards []T, writtenCards *[]string, mode string) error {
	ctx := context.Background()
	result := gorm.WithResult()
	var digCards []models.DigimonCard

	if mode == "digimon" {
		// log.Println("Type matches Digimon Type")
		digCards = cardMapper(cards)
		// log.Println(digCards)
		for _, card := range digCards {

			if !containsIdDuplicate(writtenCards, card, mode) {
				err := gorm.G[models.DigimonCard](db, result).Create(ctx, &card)
				if err != nil {
					fmt.Printf("seeder.go::writeToDB() -> Error writing card %v to DB: %v", card, err)
					return err
				}
				// log.Println("Card written to DB: ", card.Id)
			} else {
				// log.Println("inside continue")
				continue
			}
		}
	} else {
		// log.Println("Type does not match Digimon Type")
		for _, card := range cards {
			// log.Println("Current Card: ", card)
			// log.Println("Card written to DB")
			err := gorm.G[T](db, result).Create(ctx, &card)
			if err != nil {
				fmt.Printf("seeder.go::writeToDB() -> Error writing card %v to DB: %v", card, err)
				return err
			}
		}
	}
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
			err = writeToDB(db, digimonCards, &writtenCards, mode)
		} else if mode == "pokemon" {
			err = writeToDB(db, pokemonCards, &writtenCards, mode)
		} else if mode == "naruto" {
			err = writeToDB(db, narutoCards, &writtenCards, mode)
		}
		if err != nil {
			fmt.Printf("seeder.go::run() -> Returned error writing to database: %v", err)
			return err
		}
	}

	return nil
}

func Main(args ...string) error {

	var mode *string
	if len(args) > 0 {
		fmt.Println("Received arguments:", args)
		mode = flag.String(args[0], args[1], "ich schreib mal iwas rein")
		flag.Parse()
	} else {
		mode = flag.String("mode", "pokemon", "ich schreib mal iwas rein")
		flag.Parse()
	}
	var seedDir string
	if *mode == "pokemon" {
		seedDir = "POKEMON_SEEDER_DIR"
	} else if *mode == "digimon" {
		seedDir = "DIGIMON_SEEDER_DIR"
	} else if *mode == "naruto" {
		seedDir = "NARUTO_SEEDER_DIR"
		fmt.Println("Running Naruto seeder...")
	} else {
		return fmt.Errorf("seeder.go::main() -> Invalid mode: %s", *mode)
	}
	repository.ConnectToDatabase(internal.GetDBString())
	err := run(repository.DB, os.Getenv(seedDir), *mode)
	// Source - https://stackoverflow.com/a/57972717
	// Posted by VolatileCoder, modified by community. See post 'Timeline' for change history
	// Retrieved 2026-05-07, License - CC BY-SA 4.0

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) //flags are now reset

	return err
}
