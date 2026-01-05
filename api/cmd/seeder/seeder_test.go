package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/models"
	"gorm.io/gorm"
)

func initCards() []models.DigimonCard {
	return []models.DigimonCard{
		{
			CardId:    "BT24-001",
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
			CardId:    "BT24-002",
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
	file := "./BT24.json"
	have, err := unmarshalData(file)
	if err != nil {
		t.Fatalf("seeder_test.go::unmarshalData in file %v: -> %v", file, err)
	}

	if !reflect.DeepEqual(want, have[0:2]) {
		t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want, have)
	}

}

func TestWriteToDB(t *testing.T) {
	db := models.DB
	file := "./BT24.json"
	wantCards, err := unmarshalData(file)
	if err != nil {
		t.Fatalf("seeder_test.go::writeToDB -> Error unmarshaling data from: %v: %v", file, err)
	}

	writeToDB(db.Table("WriteTest"), wantCards)
	var rtvCards []models.DigimonCard
	res := db.Find(&rtvCards)
	if res.Error != nil {
		t.Fatalf("seeder_test.go::writeToDB -> Error reading from database: %v", res.Error)
	}

	if wantCards[0].CardId == rtvCards[0].CardId {
		t.Errorf("seeder_test.go::writeToDB -> Expected return value: %v\nReceived value: %v\n", wantCards[0].CardId, rtvCards[0].CardId)
	}

	defer func() {
		db.Delete(&rtvCards)
	}()

	// if !reflect.DeepEqual(wantCards[0].CardId, rtvCards[0].CardId) {
	// 	t.Errorf("seeder_test.go80::writeToDB -> Expected return value: %v\nReceived value: %v\n", wantCards[0].CardId, rtvCards[0].CardId)
	// }
	// if !reflect.DeepEqual(wantCards[0].Name, rtvCards[0].Name) {
	// 	t.Errorf("seeder_test.go80::writeToDB -> Expected return value: %v\nReceived value: %v\n", wantCards[0].Name, rtvCards[0].Name)
	// }
	// if !reflect.DeepEqual(wantCards[0].Color, rtvCards[0].Color) {
	// 	t.Errorf("seeder_test.go80::writeToDB -> Expected return value: %v\nReceived value: %v\n", wantCards[0].Color, rtvCards[0].Color)
	// }
	// if !reflect.DeepEqual(wantCards[0].CardText, rtvCards[0].CardText) {
	// 	t.Errorf("seeder_test.go80::writeToDB -> Expected return value: %v\nReceived value: %v\n", wantCards[0].CardText, rtvCards[0].CardText)
	// }
}

var db *gorm.DB

func TestMain(m *testing.M) {
	// 1. Setup: Datenbank einmalig verbinden
	models.ConnectToDatabase(internal.GetTestDBString())

	// 2. Alle Tests ausführen
	exitCode := m.Run()

	// 3. Aufräumen (optional) und beenden
	// db = models.DB
	// cards := initCards()
	// db.Delete(&cards)
	// println("Deleting inserted data...")
	os.Exit(exitCode)
}
