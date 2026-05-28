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
		UserId:  "test-user-id",
		Name:    "Test Deck",
		Entries: []models.DigimonDeckEntry{},
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
	if len(createdDeck.Entries) != len(insertDeck.Entries) {
		t.Errorf("repository_test.go::TestCreateDigimonDeck() -> expected %d entries, got %d", len(insertDeck.Entries), len(createdDeck.Entries))
	}

	clearTests()
}

func TestCreateDigimonDeck_EmptyUserIdNoCards(t *testing.T) {
	var ctx gin.Context
	deck := models.DigimonDecks{
		UserId:  "",
		Name:    "Invalid Deck",
		Entries: []models.DigimonDeckEntry{},
	}

	createdDeck, err := repository.CreateDigimonDeck(&ctx, deck)
	if err != nil {
		t.Errorf("repository_test.go::TestCreateDigimonDeck_EmptyUserId() -> unexpected error: %v", err)
	}

	if createdDeck.UserId != "" {
		t.Errorf("repository_test.go::TestCreateDigimonDeck_EmptyUserId() -> expected empty UserId, got %s", createdDeck.UserId)
	}
	if len(createdDeck.Entries) != 0 {
		t.Errorf("repository_test.go::TestCreateDigimonDeck_EmptyUserId() -> expected 0 entries, got %d", len(createdDeck.Entries))
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
