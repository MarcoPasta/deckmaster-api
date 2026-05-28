package utils

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ImportDeck struct {
	Count  int
	Name   string
	Code   string
	Number string
}

func ParseDecklist(ctx *gin.Context, decklist string) []ImportDeck {

	var imported []ImportDeck
	lines := strings.Split(decklist, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if matched, _ := regexp.MatchString(`^(Pokémon|Pokemon|Trainer|Energy):\s*\d+$`, line); matched {
			continue
		}

		match := regexp.MustCompile(`^(\d+)\s+(.+?)\s+([A-Z]{2,})\s+(\d+)$`).FindStringSubmatch(line)
		matchEnergy := regexp.MustCompile(`^(\d+)\s+(.+?)\s+(\d+)$`).FindStringSubmatch(line)

		var count int
		var name string
		var code string
		var number string

		if match != nil {
			count, _ = strconv.Atoi(match[1])
			name = match[2]
			code = match[3]
			number = match[4]
		} else if matchEnergy != nil {
			count, _ = strconv.Atoi(matchEnergy[1])
			name = matchEnergy[2]
			number = matchEnergy[3]
			code = "mee"
		} else {
			// log.Printf(`Line did not match expected format: "%s"`, line)
			continue
		}

		if strings.HasSuffix(name, "Energy") {
			log.Printf("Found energy card line: \"%s\"", name)
			match := regexp.MustCompile(`^Basic\s+{(.+?)}.+|$`).FindStringSubmatch(name)  // matches only basic energies
			matchName := regexp.MustCompile(`^.+\s+{(.+?)}.+|$`).FindStringSubmatch(name) // matches with energies that do not start with basic
			if match != nil || matchName != nil {
				log.Println("Match Name: ", matchName)
				code, number = findEnergyCard(match, code, number) // pass original code & number incase it erases
				name = renameEnergyCard(matchName, name)           // pass original name incase it erases the name
				log.Println("After renaming Card: ", name, code, number)
			}
		}

		// searches for pattern that begins with a number, followed by n literals until "PR-" get hit
		// ex: 1 Teal Mask Ogerpon ex PR-166
		matchPromo := regexp.MustCompile(`^^(.+?\s+)+(PR-.+)$`).FindStringSubmatch(name)

		if matchPromo != nil {
			name = strings.Replace(name, matchPromo[2], "", 1) // remove count
		}

		log.Println("NEW NAME AFTER FINDIN PROMO: ", name)

		imported = append(imported, ImportDeck{count, name, code, number})
	}
	return imported
}

func findEnergyCard(match []string, code string, number string) (string, string) {
	switch match[1] {
	case "G":
		return "mee", "1"
	case "R":
		return "mee", "2"
	case "W":
		return "mee", "3"
	case "L":
		return "mee", "4"
	case "P":
		return "mee", "5"
	case "F":
		return "mee", "6"
	case "D":
		return "mee", "7"
	case "M":
		return "mee", "8"
	default:
		return code, number
	}
}

func handlePromoCards() {

}

func renameEnergyCard(match []string, name string) string {

	switch match[1] {
	case "G":
		return strings.ReplaceAll(match[0], "{G}", "Grass")
	case "R":
		return strings.ReplaceAll(match[0], "{R}", "Fire")
	case "W":
		return strings.ReplaceAll(match[0], "{W}", "Water")
	case "L":
		return strings.ReplaceAll(match[0], "{L}", "Lightning")
	case "P":
		return strings.ReplaceAll(match[0], "{P}", "Psychic")
	case "F":
		return strings.ReplaceAll(match[0], "{F}", "Fighting")
	case "D":
		return strings.ReplaceAll(match[0], "{D}", "Darkness")
	case "M":
		return strings.ReplaceAll(match[0], "{M}", "Metal")
	default:
		return name
	}
}
