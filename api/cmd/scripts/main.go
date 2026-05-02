package main

import (
	"log"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/MarcoPasta/deckmaster/api/internal/scripts"
)

func main() {
	repository.ConnectToDatabase(internal.GetDBString())

	if err := scripts.SetPokemonCardCodes(); err != nil {
		log.Fatal(err)
	}
}
