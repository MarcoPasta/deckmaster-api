package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MarcoPasta/deckmaster/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := "host=db port=5432 user=postgres dbname=deckmaster password=example"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.DigimonCard{})
	if err != nil {
		log.Fatal(err)
	}

	card := models.DigimonCard{
		Name:  "Agumon",
		Color: "Red",
	}

	ctx := context.Background()
	err = gorm.G[models.DigimonCard](db).Create(ctx, &card)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "deckmaster: OK")
	})

	log.Printf("listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
