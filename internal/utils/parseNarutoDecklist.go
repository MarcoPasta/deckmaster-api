package utils

import (
	"strconv"
	"strings"
)

type ImportNarutoDeck struct {
	Set    string
	Number string
	Rarity string
	Count  int
}

// Naruto decklist come in the format of: "Set-Number-Rarity--Count|Set-Number-Rarity--Count|..."
// By splitting correspondingly, we can extract the relevant information for each card and return a structured list of ImportNarutoDeck objects.
func ParseNarutoDecklist(decklist string) []ImportNarutoDeck {
	var imported []ImportNarutoDeck
	var importRun ImportNarutoDeck
	lines := strings.Split(decklist, "|")

	for _, line := range lines {
		splits := strings.Split(line, "-")
		splitsCount := strings.Split(line, "--")
		if len(splits) > 1 {
			importRun.Set = splits[0]
			importRun.Number = splits[1]
			importRun.Rarity = splits[2]
			if importRun.Rarity == "MMS" {
				importRun.Rarity = "MSS"
				num, err := strconv.Atoi(importRun.Number)
				if err == nil {
					importRun.Number = "M" + strconv.Itoa(num)
				}
			}
			importRun.Count, _ = strconv.Atoi(splitsCount[1])
		}

		imported = append(imported, importRun)
		importRun = ImportNarutoDeck{}
	}

	return imported
}
