package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"gorm.io/gorm"
)

// mock data
func initDigimonCards() []models.DigimonCard {
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

func initPokemonCards() []models.PokemonCard {
	return []models.PokemonCard{
		{
			CardType:   "Pokemon",
			Phase:      "Basis-Pokémon",
			Name:       "Webarak",
			Kp:         60,
			EnergyType: "Pflanze",
			Ability:    "",
			Attacks: []models.Attack{
				{
					Cost:   []string{"Pflanze"},
					Name:   "Klebriger Faden",
					Effect: "Während des nächsten Zuges deines Gegners kann sich das Verteidigende Pokémon nicht zurückziehen.",
					Damage: "10",
				},
			},
			Weakness:    "Feuer",
			Resistance:  "",
			RetreatCost: 1,
			CardNumber:  "1/88 Common",
			Set:         "Optimale Ordnung",
		},
		{
			CardType:   "Pokemon",
			Phase:      "Phase-1-Pokémon",
			Name:       "Ariados",
			Kp:         110,
			EnergyType: "Pflanze",
			Ability:    "",
			Attacks: []models.Attack{
				{
					Cost:   []string{"Pflanze"},
					Name:   "Giftring",
					Effect: "Das Aktive Pokémon deines Gegners ist jetzt vergiftet. Während des nächsten Zuges deines Gegners kann sich jenes Pokémon nicht zurückziehen.",
					Damage: "50",
				},
			},
			Weakness:    "Feuer",
			Resistance:  "",
			RetreatCost: 1,
			CardNumber:  "2/88 Common",
			Set:         "Optimale Ordnung",
		},
	}
}

func TestReadDirectory(t *testing.T) {
	logInfo("Test: readDirectory")
	// Arrange
	var tests = []struct {
		mode string
		dir  string
	}{
		{"pokemon", "POKEMON_SEEDER_DIR"},
		{"digimon", "DIGIMON_SEEDER_DIR"},
	}

	for i, tt := range tests {
		fmt.Printf("seeder_test.go::TestReadDirectory() -> Running test %v for mode: %v\n", i, tt.mode)

		path := os.Getenv(tt.dir)
		if path == "" {
			t.Errorf("seeder_test.go::TestReadDirectory() -> %s environment variable is not set", tt.dir)
		}
		files, err := readDirectory(path)
		if err != nil {
			t.Fatalf("seeder_test.go::readDirectory -> Error reading directory: %v", err)
		}
		if len(files) == 0 {
			t.Errorf("seeder_test.go::readDirectory -> No files found in directory: %v", path)
		}
		clearTests()
	}
}

func TestUnmarshalData(t *testing.T) {
	logInfo("Test: unmarshalData")
	// Arrange
	var tests = []struct {
		mode string
		dir  string
		file string
	}{
		{"pokemon", "pokemon", "OptimaleOrdnung.json"},
		{"digimon", "digimon", "BT24.json"},
	}

	for i, tt := range tests {
		fmt.Printf("seeder_test.go::TestUnmarshalData() -> Running test %v for mode: %v\n", i, tt.mode)
		dig, pok, err := unmarshalData(tt.dir, tt.file, tt.mode)
		if err != nil {
			t.Fatalf("seeder_test.go::TestUnmarshalData() in file %v: -> %v", tt.file, err)
		}
		switch tt.mode {
		case "pokemon":
			want := initPokemonCards()
			if !reflect.DeepEqual(want[0].Name, pok[0].Name) {
				t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want[0].Name, pok[0].Name)
			} else {
				logSuccess("seeder_test.go::TestUnmarshalData() -> Test successful")
			}

			if !reflect.DeepEqual(want[1].Attacks, pok[1].Attacks) {
				t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want[1].Attacks, pok[1].Attacks)
			} else {
				logSuccess("seeder_test.go::TestUnmarshalData() -> Test successful")
			}
		case "digimon":
			want := initDigimonCards()
			if !reflect.DeepEqual(want[0].Id, dig[0].Id) {
				t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want[0].Id, dig[0].Id)
			} else {
				logSuccess("seeder_test.go::TestUnmarshalData() -> Test successful")
			}

			if !reflect.DeepEqual(want[1].Text, dig[1].Text) {
				t.Errorf("Data not identical!\nExpected: %+v\nReceived: %+v", want[1].Text, dig[1].Text)
			} else {
				logSuccess("seeder_test.go::TestUnmarshalData() -> Test successful")
			}
		}
		clearTests()
	}
}

func TestWriteToDB(t *testing.T) {
	logInfo("Test: writeToDB")
	// Arrange
	var tests = []struct {
		mode string
		dir  string
		file string
	}{
		{"pokemon", "pokemon", "OptimaleOrdnung.json"},
		{"digimon", "digimon", "BT24.json"},
	}

	ctx := context.Background()

	for i, tt := range tests {
		fmt.Printf("seeder_test.go::TestWriteToDB() -> Running test %v for mode: %v\n", i, tt.mode)
		dig, pok, err := unmarshalData(tt.dir, tt.file, tt.mode)
		// fmt.Println("DIGIMON: ", dig)
		// fmt.Println("POKEMON: ", pok)
		if err != nil {
			t.Fatalf("seeder_test.go::TestWriteToDB() -> Error unmarshaling data from: %v: %v", tt.file, err)
		}
		var writtenCards []string

		// Act
		switch tt.mode {
		case "pokemon":
			writeToDB(db, pok, &writtenCards)

		case "digimon":
			writeToDB(db, dig, &writtenCards)
		}

		// Assert
		switch tt.mode {
		case "pokemon":
			have, err := gorm.G[models.PokemonCard](db).Where("set = ?", "Optimale Ordnung").Find(ctx)
			if err != nil {
				t.Fatalf("seeder_test.go::TestWriteToDB() -> Error reading from database: %v", err)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
			if pok[0].Name != have[0].Name {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", pok[0].Name, have[0].Name)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}

			if !reflect.DeepEqual(pok[0].EnergyType, have[0].EnergyType) {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", pok[0].EnergyType, have[0].EnergyType)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
			if !reflect.DeepEqual(pok[0].Attacks, have[0].Attacks) {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", pok[0].Attacks, have[0].Attacks)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
		case "digimon":
			have, err := gorm.G[models.DigimonCard](db).Where("set = ?", "TIME STRANGER [BT-24]").Find(ctx)
			if err != nil {
				t.Fatalf("seeder_test.go::TestWriteToDB() -> Error reading from database: %v", err)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
			if dig[0].Id != have[0].Id {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", dig[0].Id, have[0].Id)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
			if !reflect.DeepEqual(dig[0].Id, have[0].Id) {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", dig[0].Id, have[0].Id)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
			if !reflect.DeepEqual(dig[0].Name, have[0].Name) {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", dig[0].Name, have[0].Name)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
			if !reflect.DeepEqual(dig[0].Color, have[0].Color) {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", dig[0].Color, have[0].Color)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
			if !reflect.DeepEqual(dig[0].Text, have[0].Text) {
				t.Errorf("seeder_test.go::TestWriteToDB() -> Expected return value: %v\nReceived value: %v\n", dig[0].Text, have[0].Text)
			} else {
				logSuccess("seeder_test.go::TestWriteToDB() -> Test successful")
			}
		}
		clearTests()
	}
}

func TestRun(t *testing.T) {
	logInfo("Test: run")
	ctx := context.Background()

	// Arrange
	var tests = []struct {
		mode string
		dir  string
	}{
		{"pokemon", "pokemon"},
		{"digimon", "digimon"},
	}

	for i, tt := range tests {
		fmt.Printf("seeder_test.go::TestRun() -> Running test %v for mode: %v\n", i, tt.mode)
		// Act
		err := run(db, tt.dir, tt.mode)
		if err != nil {
			t.Errorf("seeder_test.go::TestRun() -> Error running seeder: %v", err)
		} else {
			logSuccess("seeder_test.go::TestRun() -> Test successful")
		}

		switch tt.mode {
		case "pokemon":
			have, err := gorm.G[models.PokemonCard](db).Where("set = ?", "Optimale Ordnung").Find(ctx)
			if err != nil {
				t.Errorf("seeder_test.go::TestRun() -> Could not retrieve from Pokemon Database: %v\n", err)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
			want := initPokemonCards()
			// Assert
			if have != nil && !reflect.DeepEqual(want[1].Name, have[1].Name) {
				t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[1].Name, have[1].Name)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
			if have != nil && !reflect.DeepEqual(want[1].Attacks, have[1].Attacks) {
				t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[1].Attacks, have[1].Attacks)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
			if have != nil && !reflect.DeepEqual(want[0].EnergyType, have[0].EnergyType) {
				t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[0].EnergyType, have[0].EnergyType)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}

		case "digimon":
			have, err := gorm.G[models.DigimonCard](db).Where("set = ?", "TIME STRANGER [BT-24]").Find(ctx)
			if err != nil {
				t.Errorf("seeder_test.go::TestRun() -> Could not retrieve from Digimon Database: %v\n", err)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
			want := initDigimonCards()
			// Assert
			if have != nil && want[0].Id != have[0].Id {
				t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[0].Id, have[0].Id)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
			if have != nil && !reflect.DeepEqual(want[1].Text, have[1].Text) {
				t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[1].Text, have[1].Text)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
			if have != nil && !reflect.DeepEqual(want[0].Color, have[0].Color) {
				t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[0].Color, have[0].Color)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
			if have != nil && !reflect.DeepEqual(want[1].Name, have[1].Name) {
				t.Errorf("seeder_test.go::TestRun() -> Expected return value: %v\nReceived value: %v\n", want[1].Name, have[1].Name)
			} else {
				logSuccess("seeder_test.go::TestRun() -> Test successful")
			}
		}
		clearTests()
	}
}

func clearTests() {
	fmt.Println("clearing testsuite")
	db.Unscoped().Where("1 = 1").Delete(&models.DigimonCard{})
	db.Unscoped().Where("1 = 1").Delete(&models.PokemonCard{})
}

func logSuccess(message string) {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), message, string(colorReset))
}

func logInfo(message string) {
	colorReset := "\033[0m"
	colorBlue := "\033[34m"
	fmt.Println(string(colorBlue), message, string(colorReset))
}

var db *gorm.DB

func TestMain(m *testing.M) {
	// connect to test database
	repository.ConnectToDatabase(internal.GetTestDBString())
	db = repository.DB
	// run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}
