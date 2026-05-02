package scripts

import (
	"fmt"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
)

func SetPokemonCardCodes() error {
	setToCode := map[string]string{
		"Optimale Ordnung":           "por",
		"Erhabene Helden":            "asc",
		"Fatale Flammen":             "pfl",
		"Mega-Entwicklung":           "meg",
		"Weiße Flammen":              "wht",
		"Schwarze Blitze":            "blk",
		"Ewige Rivalen":              "dri",
		"Reisegefährten":             "jtg",
		"Prismatische Entwicklungen": "pre",
		"Nebel der Sagen":            "sfa",
		"Stellarkrone":               "scr",
		"Stürmische Funken":          "ssp",
		"Maskerade im Zwielicht":     "twm",
		"Gewalten der Zeit":          "tef",
	}

	for set, code := range setToCode {
		result := repository.DB.Model(&models.PokemonCard{}).
			Where("set = ?", set).
			Update("code", code)

		if result.Error != nil {
			return fmt.Errorf("failed to update code for set %q: %w", set, result.Error)
		}

		fmt.Printf("Updated %d cards in set %q with code %q\n", result.RowsAffected, set, code)
	}

	return nil
}
