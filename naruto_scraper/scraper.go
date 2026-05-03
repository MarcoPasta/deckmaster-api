package main

import (
	"context"
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
	"time"

	"github.com/chromedp/chromedp"
)

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

const pageURL = "https://www.narutotcgmythos.com/card-gallery"

func main() {
	dir := os.Getenv("NARUTO_SCRAPER_DIR")
	if dir == "" {
		log.Fatal("NARUTO_SCRAPER_DIR environment variable not set")
	}

	if err := os.MkdirAll(dir, 0750); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	imgDir := filepath.Join(dir, "images")
	if err := os.MkdirAll(imgDir, 0750); err != nil {
		log.Fatalf("Failed to create images directory: %v", err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	// Navigate and wait for cards to load
	fmt.Println("Navigating to page...")
	err := chromedp.Run(ctx, chromedp.Navigate(pageURL))
	if err != nil {
		log.Fatalf("Failed to navigate: %v", err)
	}

	// Poll for card items to appear
	fmt.Println("Waiting for cards to render...")
	var cardCount int
	for attempts := 0; attempts < 60; attempts++ {
		err = chromedp.Run(ctx,
			chromedp.Sleep(2*time.Second),
			chromedp.Evaluate(`document.querySelectorAll('.card-item img').length`, &cardCount),
		)
		if err != nil {
			log.Fatalf("Failed to evaluate card count: %v", err)
		}
		fmt.Printf("  Attempt %d: found %d card images\n", attempts+1, cardCount)
		if cardCount > 100 {
			break
		}
	}

	if cardCount == 0 {
		// Dump page info for debugging
		var bodyLen int
		var imgCount int
		_ = chromedp.Run(ctx,
			chromedp.Evaluate(`document.body.innerHTML.length`, &bodyLen),
			chromedp.Evaluate(`document.querySelectorAll('img').length`, &imgCount),
		)
		fmt.Printf("  DEBUG: body length=%d, total imgs=%d\n", bodyLen, imgCount)
		log.Fatal("No cards found. The site's SPA may not render in headless mode in this environment.")
	}

	fmt.Printf("Found %d cards. Starting scrape...\n", cardCount)

	// Click the first card to open the modal
	fmt.Println("Opening first card dialog...")
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelectorAll('.card-item img')[0].click()`, nil),
		chromedp.Sleep(1500*time.Millisecond),
		chromedp.WaitVisible(`#cardDetailModal`, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatalf("Failed to open first card dialog: %v", err)
	}
	fmt.Println("Dialog opened!")

	// Now iterate through all cards using the next button
	var cards []NarutoCard
	for i := 0; i < cardCount; i++ {
		fmt.Printf("[%d/%d] Scraping...", i+1, cardCount)

		card, imgURL, err := extractCardData(ctx)
		if err != nil {
			log.Printf(" ERROR: %v\n", err)
			// Try to click next anyway
			_ = chromedp.Run(ctx,
				chromedp.Evaluate(`document.getElementById('nextCardBtn').click()`, nil),
				chromedp.Sleep(800*time.Millisecond),
			)
			continue
		}
		card.Set = "1"
		cards = append(cards, card)
		fmt.Printf(" %s (%s) [%s]\n", card.Name, card.CardNumber, card.Rarity)

		// Download image
		if imgURL != "" {
			downloadImage(imgURL, imgDir, card.CardNumber, card.Rarity)
		}

		// Click next (unless last card)
		if i < cardCount-1 {
			err = chromedp.Run(ctx,
				chromedp.Evaluate(`document.getElementById('nextCardBtn').click()`, nil),
				chromedp.Sleep(800*time.Millisecond),
			)
			if err != nil {
				log.Printf("  Warning: failed to click next: %v\n", err)
			}
		}
	}

	// Close modal
	_ = chromedp.Run(ctx, chromedp.Evaluate(`document.getElementById('closeDetailModalBtn').click()`, nil))

	// Save to JSON
	outputPath := filepath.Join(dir, "Set1.json")
	data, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("\nDone! Saved %d cards to %s\n", len(cards), outputPath)
}

func extractCardData(ctx context.Context) (NarutoCard, string, error) {
	var card NarutoCard

	var jsonStr string
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				const get = id => (document.getElementById(id)?.textContent || '').trim();
				const imgEl = document.getElementById('detailImage');
				const imgSrc = imgEl ? imgEl.src : '';
				const effectEl = document.getElementById('detailText');
				const effect = effectEl ? effectEl.textContent.trim() : '';

				return JSON.stringify({
					name: get('detailTitle'),
					description: get('detailSubtitle'),
					rarity: get('detailRarity'),
					chakra: get('detailChakra'),
					power: get('detailPower'),
					keyword: get('detailKeywords'),
					group: get('detailGroup'),
					effect: effect,
					imgSrc: imgSrc
				});
			})()
		`, &jsonStr),
	)
	if err != nil {
		return card, "", fmt.Errorf("failed to evaluate: %w", err)
	}

	var data struct {
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

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return card, "", fmt.Errorf("failed to parse JSON: %w (raw: %s)", err, jsonStr)
	}

	card.Name = data.Name
	card.Description = data.Description
	card.Rarity = data.Rarity
	card.Keyword = data.Keyword
	card.Group = data.Group
	card.Effect = data.Effect
	card.EffectTypes = extractEffectTypes(data.Effect)
	card.Cost = parseIntSafe(data.Chakra)
	card.Power = parseIntSafe(data.Power)
	card.CardNumber = extractCardNumber(data.ImgSrc)

	return card, data.ImgSrc, nil
}

func extractCardNumber(imgSrc string) string {
	parts := strings.Split(imgSrc, "/")
	if len(parts) == 0 {
		return ""
	}
	filename := parts[len(parts)-1]
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	// Try to extract leading number (e.g., "01" from "01_C_Hiruzen_Professor")
	re := regexp.MustCompile(`^(\d+)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1]
	}

	// For special cards like "Legendary_Naruto_" or "M1_Call-for-support"
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

	re := regexp.MustCompile(`(MAIN|UPGRADE|PASSIVE)\s*↯`)
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
	s = strings.Trim(s, `"`)
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
	if url == "" {
		return
	}

	ext := ".webp"
	if strings.Contains(url, ".png") {
		ext = ".png"
	} else if strings.Contains(url, ".jpg") || strings.Contains(url, ".jpeg") {
		ext = ".jpg"
	}

	filename := fmt.Sprintf("%s_%s%s", cardNumber, rarity, ext)
	outputPath := filepath.Join(dir, filename)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to download image %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download image %s: status %d", url, resp.StatusCode)
		return
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", outputPath, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Failed to save image %s: %v", outputPath, err)
	}
}
