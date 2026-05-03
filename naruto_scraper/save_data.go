//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type RawCard struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rarity      string `json:"rarity"`
	Chakra      string `json:"chakra"`
	Power       string `json:"power"`
	Keyword     string `json:"keyword"`
	Group       string `json:"group"`
	Effect      string `json:"effect"`
	ImgSrc      string `json:"imgSrc"`
}

type NarutoCard struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Cost        int      `json:"cost"`
	Power       int      `json:"power"`
	Keyword     string   `json:"keyword"`
	Group       string   `json:"group"`
	Effect      string   `json:"effect"`
	EffectTypes []string `json:"effect_types"`
	Rarity      string   `json:"rarity"`
	CardNumber  string   `json:"card_number"`
	Set         string   `json:"set"`
}

func main() {
	dir := os.Getenv("NARUTO_SCRAPER_DIR")
	if dir == "" {
		dir = "../files/naruto/json/"
	}
	os.MkdirAll(dir, 0750)
	imgDir := filepath.Join(dir, "images")
	os.MkdirAll(imgDir, 0750)

	rawData, err := os.ReadFile("/tmp/naruto_raw.json")
	if err != nil {
		log.Fatal(err)
	}

	var rawCards []RawCard
	if err := json.Unmarshal(rawData, &rawCards); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processing %d cards...\n", len(rawCards))

	var cards []NarutoCard
	for i, raw := range rawCards {
		card := NarutoCard{
			Name:        raw.Name,
			Description: raw.Description,
			Rarity:      raw.Rarity,
			Keyword:     raw.Keyword,
			Group:       raw.Group,
			Effect:      raw.Effect,
			EffectTypes: extractEffectTypes(raw.Effect),
			Cost:        parseIntSafe(raw.Chakra),
			Power:       parseIntSafe(raw.Power),
			CardNumber:  extractCardNumber(raw.ImgSrc),
			Set:         "1",
		}
		cards = append(cards, card)

		// Download image
		if raw.ImgSrc != "" {
			downloadImage(raw.ImgSrc, imgDir, card.CardNumber, card.Rarity)
		}

		fmt.Printf("[%d/%d] %s (%s) [%s]\n", i+1, len(rawCards), card.Name, card.CardNumber, card.Rarity)
	}

	outputPath := filepath.Join(dir, "Set1.json")
	data, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(outputPath, data, 0644)
	fmt.Printf("\nDone! Saved %d cards to %s\n", len(cards), outputPath)
}

func extractCardNumber(imgSrc string) string {
	parts := strings.Split(imgSrc, "/")
	if len(parts) == 0 {
		return ""
	}
	filename := parts[len(parts)-1]
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	re := regexp.MustCompile(`^(\d+)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1]
	}
	re2 := regexp.MustCompile(`^([A-Za-z]+\d*)`)
	matches2 := re2.FindStringSubmatch(filename)
	if len(matches2) > 1 {
		return matches2[1]
	}
	return filename
}

func extractEffectTypes(effect string) []string {
	if effect == "" {
		return nil
	}
	var types []string
	seen := make(map[string]bool)
	re := regexp.MustCompile("(MAIN|UPGRADE|PASSIVE|AMBUSH|SCORE)\\s*[↯⧗:]")
	matches := re.FindAllStringSubmatch(effect, -1)
	for _, m := range matches {
		if !seen[m[1]] {
			types = append(types, m[1])
			seen[m[1]] = true
		}
	}
	return types
}

func parseIntSafe(s string) int {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" {
		return 0
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

func downloadImage(url string, dir string, cardNumber string, rarity string) {
	ext := ".webp"
	if strings.Contains(url, ".png") {
		ext = ".png"
	}
	filename := fmt.Sprintf("%s_%s%s", cardNumber, rarity, ext)
	outputPath := filepath.Join(dir, filename)

	// Skip if already downloaded
	if _, err := os.Stat(outputPath); err == nil {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download %s: %d", url, resp.StatusCode)
		return
	}
	file, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Failed to create %s: %v", outputPath, err)
		return
	}
	defer file.Close()
	io.Copy(file, resp.Body)
}
