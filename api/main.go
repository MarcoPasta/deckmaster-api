package main

import (
	"net/http"

	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/MarcoPasta/deckmaster/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	dsn := "host=db port=5432 user=postgres dbname=deckmaster password=example"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.DigimonDecks{}, &models.DigimonCard{})
	if err != nil {
		log.Fatal(err)
	}

	card := models.DigimonCard{
		Name:  "Agumon",
		Color: "Red",
	}

	deck := models.DigimonDecks{
		Name:  "Mein Deck",
		Cards: []models.DigimonCard{card},
	}

	ctx := context.Background()
	err = gorm.G[models.DigimonCard](db).Create(ctx, &card)
	err = gorm.G[models.DigimonDecks](db).Create(ctx, &deck)

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
