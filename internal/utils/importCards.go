package utils

import (
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/gin-gonic/gin"
)

func DecodeCards(ctx *gin.Context, file multipart.FileHeader, game interface{}) ([]interface{}, error) {

	switch game.(type) {
	case models.PokemonCard:
		log.Println("Importing Pokemon cards from file: ", file.Filename)
		var cards []models.PokemonCard
		fileContent, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %s, error: %v\n", file.Filename, err)
			return nil, err
		}
		defer fileContent.Close()

		err = json.NewDecoder(fileContent).Decode(&cards)
		if err != nil {
			return nil, err
		}

		var result []interface{}
		for _, card := range cards {
			result = append(result, card)
		}

		return result, nil

	case models.DigimonCard:
		log.Println("Importing Digimon cards from file: ", file.Filename)
		// Implement the logic to read the file and import Digimon cards

		var cards []models.DigimonCard
		// var writtenCards []string
		fileContent, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %s, error: %v\n", file.Filename, err)
			return nil, err
		}
		defer fileContent.Close()

		err = json.NewDecoder(fileContent).Decode(&cards)
		if err != nil {
			return nil, err
		}

		var result []interface{}
		for _, card := range cards {
			result = append(result, card)
		}

		return result, nil
		// for _, card := range cards {
		// 	/* Keep incase there are ID suplicates constraints that will stop the data import */
		// 	// if containsIdDuplicate(&writtenCards, card, "Digimon") {
		// 	// 	log.Printf("Duplicate card ID found: %s, skipping...\n", card.Id)
		// 	// 	continue
		// 	// }
		// 	err := gorm.G[models.DigimonCard](repository.DB).Create(ctx, &card)
		// 	if err != nil {
		// 		log.Printf("Error saving Digimon card to database: %s, error: %v\n", card.Id, err)
		// 		return err
		// 	}
		// }

	case models.NarutoCard:
		log.Println("Importing Naruto cards from file: ", file.Filename)
		var cards []models.NarutoCard
		fileContent, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %s, error: %v\n", file.Filename, err)
			return nil, err
		}
		defer fileContent.Close()

		err = json.NewDecoder(fileContent).Decode(&cards)
		if err != nil {
			return nil, err
		}

		var result []interface{}
		for _, card := range cards {
			result = append(result, card)
		}
		return result, nil
		// for _, card := range cards {
		// 	err := gorm.G[models.NarutoCard](repository.DB).Create(ctx, &card)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

	default:
		log.Printf("Unsupported game type for file: %s\n", file.Filename)
		return nil, errors.New("Unsupported game type")
	}
}
