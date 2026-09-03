package utils

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
)

// TODO: add errors
func CreatePokemonCardMarketList(cards []models.PokemonCard) (string, error) {

	var list string

	sort.Slice(cards, func(i, j int) bool {
		return cardNumber(cards[i].CardNumber) < cardNumber(cards[j].CardNumber)
	})

	log.Println("Sorted cards: ", len(cards))

	for i := 0; i < len(cards); {
		card := cards[i]

		if card.CardType == "Basis-Energie" || card.CardType == "Basic Energy" ||
			(strings.Contains(card.Name, "Basic") && strings.Contains(card.Name, "Energy")) {
			i++
			continue
		}

		// Count how many consecutive entries share the same CardNumber
		count := 1
		for i+count < len(cards) && cards[i+count].CardNumber == card.CardNumber {
			count++
		}

		entry := ""
		if card.Ability.Name != "" {
			entry = fmt.Sprintf("%d %s %s", count, card.Name, card.Ability.Name)
		} else {
			entry = fmt.Sprintf("%d %s", count, card.Name)
		}
		for _, attack := range card.Attacks {
			if attack.Name != "Tera" && attack.Name != "Pokémon ex rule" {
				entry += fmt.Sprintf(" %s", attack.Name)
			}
		}
		entry += "\n"
		list += entry

		i += count
	}
	log.Printf("Generated list: \n%s\n", list)
	return list, nil
}

// cardNumber extracts the leading integer from a card number string like "10/162 Common".
// Returns -1 for strings that don't start with a number so they sort last.
func cardNumber(s string) int {
	part := strings.SplitN(s, "/", 2)[0]
	n, err := strconv.Atoi(strings.TrimSpace(part))
	if err != nil {
		return -1
	}
	return n
}
