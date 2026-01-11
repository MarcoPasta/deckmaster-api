package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/models"
	"gorm.io/gorm"
)

// mock data
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
			Text: []string{
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
			Text: []string{
				"[Inherited Effect]",
				"[End of Your Turn] [Once Per Turn] By paying 1 cost, this blue Digimon with the [TS] trait unsuspends.",
			},
			Set: "TIME STRANGER [BT-24]",
		},
	}
}

func TestReadDirectory(t *testing.T) {
	path := os.Getenv("SEEDER_FILE_DIR")
	if path == "" {
		t.Errorf("seeder_test.go::TestReadDirectory() -> SEEDER_FILE_DIR environment variable is not set")
	}
	files, err := readDirectory(path)
	if err != nil {
		t.Fatalf("seeder_test.go::readDirectory -> Error reading directory: %v", err)
	}
	if len(files) == 0 {
		t.Errorf("seeder_test.go::readDirectory -> No files found in directory: %v", path)
	}
}

func TestUnmarshalData(t *testing.T) {

	want := initCards()
	file := "./BT24.json"
	have, err := unmarshalData(os.Getenv("SEEDER_FILE_DIR"), file)
	if err != nil {
		t.Fatalf("seeder_test.go::TestUnmarshalData() in file %v: -> %v", file, err)
	}

	if !reflect.DeepEqual(want[0].Id, have[0].Id) {
		t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want[0].Id, have[0].Id)
	}

	if !reflect.DeepEqual(want[1].Text, have[1].Text) {
		t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want[1].Text, have[1].Text)
	}

}

func TestWriteToDB(t *testing.T) {
	file := "BT24.json"
	wantCards, err := unmarshalData(os.Getenv("SEEDER_FILE_DIR"), file)
	if err != nil {
		t.Fatalf("seeder_test.go::TestWriteToDB() -> Error unmarshaling data from: %v: %v", file, err)
	}
	var writtenCards []string
	var rtvCards []models.DigimonCard
	var rtvAltCards []models.DigimonCardAltArt

	writeToDB(db, wantCards, &writtenCards)

	res := db.Find(&rtvCards)
	if res.Error != nil {
		t.Fatalf("seeder_test.go::TestWriteToDB() -> Error reading from database: %v", res.Error)
	}
	if wantCards[0].Id != rtvCards[0].Id {
		t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", wantCards[0].Id, rtvCards[0].Id)
	}
	if !reflect.DeepEqual(wantCards[0].Id, rtvCards[0].Id) {
		t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", wantCards[0].Id, rtvCards[0].Id)
	}
	if !reflect.DeepEqual(wantCards[0].Name, rtvCards[0].Name) {
		t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", wantCards[0].Name, rtvCards[0].Name)
	}
	if !reflect.DeepEqual(wantCards[0].Color, rtvCards[0].Color) {
		t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", wantCards[0].Color, rtvCards[0].Color)
	}
	if !reflect.DeepEqual(wantCards[0].Text, rtvCards[0].Text) {
		t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", wantCards[0].Text, rtvCards[0].Text)
	}

	res = db.Find(&rtvAltCards)
	if res.Error != nil {
		t.Fatalf("seeder_test.go::TestWriteToDB() -> Error reading alt art from database: %v", res.Error)
	}
	if len(rtvAltCards) == 0 {
		t.Errorf("seeder_test.go::TestWriteToDB() -> Expected at least one alt art entry in database, found none.")
	}

	clearTests()

}

func TestRun(t *testing.T) {
	// Arrange
	path := os.Getenv("SEEDER_FILE_DIR")
	want := initCards()
	ctx := context.Background()
	if path == "" {
		t.Errorf("seeder_test.go::TestRun() -> SEEDER_FILE_DIR environment variable is not set")
	}

	// Act
	err := run(db, path)
	if err != nil {
		t.Errorf("seeder_test.go::TestRun() -> Error running seeder: %v", err)
	}

	// Assert
	have, err := gorm.G[models.DigimonCard](db).Where("set = ?", "TIME STRANGER [BT-24]").Find(ctx)

	if want[0].Id != have[0].Id {
		t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[0].Id, have[0].Id)
	}
	if !reflect.DeepEqual(want[1].Text, have[1].Text) {
		t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[1].Text, have[1].Text)
	}
	if !reflect.DeepEqual(want[0].Color, have[0].Color) {
		t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[0].Color, have[0].Color)
	}
	if !reflect.DeepEqual(want[1].Name, have[1].Name) {
		t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[1].Name, have[1].Name)
	}

	clearTests()
}

func clearTests() {
	fmt.Println("clearing test")
	db.Unscoped().Where("1 = 1").Delete(&models.DigimonCard{})
	db.Unscoped().Where("1 = 1").Delete(&models.DigimonCardAltArt{})
}

var db *gorm.DB

func TestMain(m *testing.M) {
	// connect to test database
	models.ConnectToDatabase(internal.GetTestDBString())
	db = models.DB
	// run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}
