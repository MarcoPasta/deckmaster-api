package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/MarcoPasta/deckmaster/api/models"
	"gorm.io/gorm"
)

func initCards() []models.DigimonCard {
	return []models.DigimonCard{
		{
			Id:        "BT24-001",
			Name:      "Gigimon",
			CardType:  "Digi-Egg",
			Color:     []string{"Red"},
			Cost:      0,
			Dp:        0,
			DvCost:    []string{""},
			Form:      "In-Training",
			Attribute: "",
			Species:   []string{"Lesser/LIBERATOR"},
			CardText: []string{
				"[Inherited Effect]",
				"[Your Turn] [Once Per Turn] When your opponent's security stack is removed from, you may delete 1 of their Digimon with 3000 DP or less.",
			},
			Set: "TIME STRANGER [BT-24]",
		},
		{
			Id:        "BT24-002",
			Name:      "Bukamon",
			CardType:  "Digi-Egg",
			Color:     []string{"Blue"},
			Cost:      0,
			Dp:        0,
			DvCost:    []string{""},
			Form:      "In-Training",
			Attribute: "",
			Species:   []string{"Lesser/Iliad/TS"},
			CardText: []string{
				"[Inherited Effect]",
				"[End of Your Turn] [Once Per Turn] By paying 1 cost, this blue Digimon with the [TS] trait unsuspends.",
			},
			Set: "TIME STRANGER [BT-24]",
		},
	}
}

func TestUnmarshalData(t *testing.T) {

	want := initCards()

	have, err := unmarshalData()
	if err != nil {
		t.Fatalf("seeder_test.go::unmarshalData -> %v", err)
	}

	if !reflect.DeepEqual(want, have) {
		t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want, have)
	}

}

func TestWriteToDB(t *testing.T) {
	db = models.DB
	wantCards := initCards()
	writeToDB(wantCards)
	var rtvCards []models.DigimonCard
	res := db.Find(&rtvCards)

	if res.Error != nil {
		t.Fatalf("seeder_test.go73::writeToDB -> Error reading from database: %v", res.Error)
	}

	if !reflect.DeepEqual(wantCards[0], rtvCards[0]) {
		t.Errorf("seeder_test.go80::writeToDB -> Expected return value: %v\nReceived value: %v\n", wantCards[0], rtvCards[0])

	}
}

var db *gorm.DB

func TestMain(m *testing.M) {
	// 1. Setup: Datenbank einmalig verbinden
	models.ConnectToDatabase()

	// 2. Alle Tests ausführen
	exitCode := m.Run()

	// 3. Aufräumen (optional) und beenden
	db = models.DB
	cards := initCards()
	db.Delete(&cards)
	println("Deleting inserted data...")
	os.Exit(exitCode)
}
