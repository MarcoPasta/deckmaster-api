package repository_test

import (
	"os"
	"testing"

	"github.com/MarcoPasta/deckmaster/api/internal"
	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initTestDeck() models.DigimonDecks {
	return models.DigimonDecks{

		UserId: "test-user-id",
		Name:   "Test Deck",
		Cards: []models.DigimonCard{
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
		},
	}
}

func TestCreateDigimonDeck(t *testing.T) {
	var ctx gin.Context
	insertDeck := initTestDeck()

	// Act
	createdDeck, err := repository.CreateDigimonDeck(&ctx, insertDeck)
	if err != nil {
		t.Errorf("repository_test.go::TestCreateDigimonDeck() -> error creating digimon deck: %v", err)
	}

	if createdDeck.UserId != insertDeck.UserId {
		t.Errorf("repository_test.go::TestCreateDigimonDeck() -> expected UserId %s, got %s", insertDeck.UserId, createdDeck.UserId)
	}
	if createdDeck.Name != insertDeck.Name {
		t.Errorf("repository_test.go::TestCreateDigimonDeck() -> expected Name %s, got %s", insertDeck.Name, createdDeck.Name)
	}
	if len(createdDeck.Cards) != len(insertDeck.Cards) {
		t.Errorf("repository_test.go::TestCreateDigimonDeck() -> expected %d cards, got %d", len(insertDeck.Cards), len(createdDeck.Cards))
	}

	clearTests()
}

func TestCreateDigimonDeck_EmptyUserIdNoCards(t *testing.T) {
	var ctx gin.Context
	deck := models.DigimonDecks{
		UserId: "",
		Name:   "Invalid Deck",
		Cards:  []models.DigimonCard{},
	}

	createdDeck, err := repository.CreateDigimonDeck(&ctx, deck)
	if err != nil {
		t.Errorf("repository_test.go::TestCreateDigimonDeck_EmptyUserId() -> unexpected error: %v", err)
	}

	if createdDeck.UserId != "" {
		t.Errorf("repository_test.go::TestCreateDigimonDeck_EmptyUserId() -> expected empty UserId, got %s", createdDeck.UserId)
	}
	if len(createdDeck.Cards) != 0 {
		t.Errorf("repository_test.go::TestCreateDigimonDeck_EmptyUserId() -> expected 0 cards, got %d", len(createdDeck.Cards))
	}

	clearTests()
}

func clearTests() {
	db.Unscoped().Where("1 = 1").Delete(&models.DigimonDecks{})
}

var db *gorm.DB

func TestMain(m *testing.M) {
	// connect to test database
	repository.ConnectToDatabase(internal.GetTestDBString())
	db = repository.DB
	// set Gin to test mode
	gin.SetMode(gin.TestMode)
	// run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}
