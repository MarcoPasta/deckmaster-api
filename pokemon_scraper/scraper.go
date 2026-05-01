// developed by marcopasta

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Card struct {
	CardType    string   `json:"cardType"`
	Phase       string   `json:"phase"`
	Name        string   `json:"name"`
	Kp          int16    `json:"kp"`
	EnergyType  string   `json:"color"`
	Ability     string   `json:"ability"`
	Attacks     []Attack `json:"attacks"`
	Weakness    string   `json:"weakness"`
	Resistance  string   `json:"resistance"`
	RetreatCost int8     `json:"retreatCost"`
	CardNumber  string   `json:"cardNumber"`
	Set         string   `json:"set"`
}

type Attack struct {
	Cost   []string `json:"cost"`
	Name   string   `json:"name"`
	Effect string   `json:"effect"`
	Damage string   `json:"damage"`
}

type Set struct {
	Set      string `json:"set"`
	Uri      string `json:"uri"`
	Category string `json:"category"`
	Amount   int    `json:"amount"`
}

// Take string value and parse it to int16 value
func toInt16(s string) int16 {
	parsed, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	return int16(parsed)
}

// scrapeCard uses chromedp to extract card data from the rendered DOM
func scrapeCard(ctx context.Context, browserCtx context.Context) (Card, error) {
	var card Card

	var name, phase, hpRaw, energyType, cardNumber string
	var weaknessTitle, resistanceTitle string
	var retreatCount int
	var attacks []Attack

	// use Evaluate with optional chaining so missing elements (e.g. Trainer cards) don't error out
	err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector(".color-block h1")?.textContent?.trim() || ""`, &name),
		chromedp.Evaluate(`document.querySelector(".card-type h2")?.textContent?.trim() || ""`, &phase),
		chromedp.Evaluate(`document.querySelector(".card-hp")?.textContent?.trim() || ""`, &hpRaw),
		chromedp.Evaluate(`document.querySelector(".card-basic-info .energy")?.getAttribute("title") || ""`, &energyType),
		chromedp.Evaluate(`document.querySelector(".stats-footer span")?.textContent?.trim() || ""`, &cardNumber),
	)
	if err != nil {
		return card, fmt.Errorf("scraping basic info: %w", err)
	}

	card.Name = strings.TrimSpace(name)
	card.Phase = strings.TrimSpace(phase)
	card.EnergyType = strings.TrimSpace(energyType)
	card.CardNumber = strings.TrimSpace(cardNumber)

	if strings.Contains(card.Phase, "Pokémon") || strings.Contains(card.Phase, "Mega-Entwicklungs-ex") {
		card.CardType = "Pokemon"
	} else if strings.Contains(card.Phase, "Item") || strings.Contains(card.Phase, "Unterstützer") || strings.Contains(card.Phase, "Stadion") || strings.Contains(card.Phase, "Ausrüstung") {
		card.CardType = "Trainer"
		card.Phase = "" // phase is not relevant for Trainer cards, so clear it out
	}

	// parse HP: strip "KP" prefix
	hpStr := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(hpRaw), "KP"))
	if hpStr != "" {
		card.Kp = toInt16(hpStr)
	}

	// weakness: get title from first li in the first .stat block
	chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector(".pokemon-stats .stat:first-child .card-energies li")?.getAttribute("title") || ""`, &weaknessTitle),
	)
	card.Weakness = strings.TrimSpace(weaknessTitle)

	// resistance: get title from li in the second .stat block
	chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector(".pokemon-stats .stat:nth-child(2) .card-energies li")?.getAttribute("title") || ""`, &resistanceTitle),
	)
	card.Resistance = strings.TrimSpace(resistanceTitle)

	// retreat cost = count of li elements inside .stat.last .card-energies
	chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelectorAll(".pokemon-stats .stat.last .card-energies li").length`, &retreatCount),
	)
	card.RetreatCost = int8(retreatCount)

	// set name from stats-footer
	var setName string
	chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector(".stats-footer h3 a")?.textContent?.trim() || ""`, &setName),
	)
	card.Set = strings.TrimSpace(setName)

	// attacks: extract each .ability block
	var attackCount int
	chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelectorAll(".pokemon-abilities .ability").length`, &attackCount),
	)

	for i := range attackCount {
		var atk Attack
		prefix := fmt.Sprintf(".pokemon-abilities .ability:nth-child(%d)", i+1)

		var atkName, atkDamage, atkEffect string
		chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(
				`document.querySelector("%s h4.label")?.textContent?.trim() || ""`, prefix), &atkName),
		)
		chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(
				`document.querySelector("%s span.right.plus")?.textContent?.trim() || ""`, prefix), &atkDamage),
		)
		chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(
				`document.querySelector("%s pre")?.textContent?.trim() || ""`, prefix), &atkEffect),
		)

		// attack cost: get title attrs from each li in the attack's ul
		var costTitles []string
		chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(
				`Array.from(document.querySelectorAll("%s ul.left li")).map(li => li.getAttribute("title") || "")`,
				prefix,
			), &costTitles),
		)

		atk.Name = strings.TrimSpace(atkName)
		atk.Damage = strings.TrimSpace(atkDamage)
		atk.Effect = strings.TrimSpace(atkEffect)
		atk.Cost = costTitles
		attacks = append(attacks, atk)
	}
	card.Attacks = attacks

	fmt.Printf("CARD DETAILS: %+v\n", card)
	return card, nil
}

func checkPathExist(path string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Path %v does not exist, creating path..\n", path)
			err := os.MkdirAll(path, 0750)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else { // path does exist
		return
	}
}

// downloadImage downloads an image from the given URL and saves it to the set directory.
// The file is named after the set with the card index.
func downloadImage(imgUrl string, set string, index int) error {
	path := os.Getenv("POKEMON_IMAGES_DIR") + set + "/"
	checkPathExist(path)

	fileName := set + "_" + strconv.Itoa(index)

	response, err := http.Get(imgUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("got status %d for %s", response.StatusCode, imgUrl)
	}

	filePath := path + fileName + ".png"
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}

	fmt.Printf("Saved image %v successfully!\n", fileName)
	return nil
}

// user agents to rotate through on each browser restart
var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/126.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.5; rv:126.0) Gecko/20100101 Firefox/126.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
}

var uaIndex int

/* Random Strings to supplement User Agent creation*/
// func RandomString() string {
// 	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	b := make([]byte, rand.Intn(10)+10)
// 	for i := range b {
// 		b[i] = letterBytes[rand.Intn(len(letterBytes))]
// 	}
// 	return string(b)
// }

// newBrowserContext creates a chromedp context with anti-detection flags
func newBrowserContext() (context.Context, context.CancelFunc) {
	ua := userAgents[uaIndex%len(userAgents)]
	uaIndex++

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium"),
		chromedp.Flag("headless", "new"),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent(ua),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, ctxCancel := chromedp.NewContext(allocCtx)

	cancel := func() {
		ctxCancel()
		allocCancel()
	}
	fmt.Printf("Using User-Agent: %s\n", ua)
	return ctx, cancel
}

func main() {
	// CLI usage: go run scraper.go [cards|images]
	mode := "cards"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	if mode != "cards" && mode != "images" {
		log.Fatalf("Usage: %s [cards|images]", os.Args[0])
	}
	fmt.Printf("Running in mode: %s\n", mode)

	baseUrlDE := "https://www.pokemon.com/de/pokemon-sammelkartenspiel/pokemon-karten/series/"
	// baseUrlEN := "https://www.pokemon.com/us/pokemon-tcg/pokemon-cards/series/"
	imgURLDE := "https://assets.pokemon.com/static-assets/content-assets/cms2-de-de/img/cards/web/"

	/* scraping details:
	url: https://www.pokemon.com/de/pokemon-sammelkartenspiel/pokemon-karten/series/EX0?/Y/ --> Extension & Number followed by the card number, Example: /me02/1/
	Each set or starter/advanced deck contains a category which indicates the cards being placed inside
	*/

	// inject the set to be scraped from a collection
	var sets []Set
	data, err := os.ReadFile("./collection.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &sets)
	if err != nil {
		log.Fatal(err)
	}

	// check path before writing to insure path and env exists before scraping cards for a loooong time
	writePath := os.Getenv("POKEMON_SCRAPER_DIR")
	checkPathExist(writePath)

	// create a single browser context -- first one started inside the loop
	for _, set := range sets {
		var scrapedCards []Card

		// browser restart interval to prevent resource exhaustion and rotate user agent /////
		const restartEvery = 30
		var browserCtx context.Context
		var browserCancel context.CancelFunc
		//////////////////////////////////////////////////////////////////////////////////////////

		for i := 1; i <= set.Amount; i++ {
			fmt.Printf("##### CARD %v/%v #####\n", i, set.Amount)

			// add a random delay between requests (1-3s) to avoid rate limiting
			if i > 1 {
				delay := time.Duration(1000+rand.Intn(2000)) * time.Millisecond
				time.Sleep(delay)
			}
			////////////////////////////////////////////////////////////////////////

			switch mode {
			case "cards":
				// start or restart browser every N cards
				if browserCtx == nil || ((i-1)%restartEvery == 0 && i > 1) {
					if browserCancel != nil {
						browserCancel()
					}
					browserCtx, browserCancel = newBrowserContext()
					fmt.Printf("Browser (re)started at card %d\n", i)
				}

				// define the curling url
				currUrl := baseUrlDE + set.Uri + fmt.Sprintf("%d/", i)
				fmt.Printf("Visiting url: %v\n", currUrl)
				fmt.Printf("Start scraping set: %v\n", set.Set)

				// create a new tab for each card to avoid context cancellation issues
				tabCtx, tabCancel := chromedp.NewContext(browserCtx)
				timeoutCtx, timeoutCancel := context.WithTimeout(tabCtx, 20*time.Second)

				// navigate and wait for the card content to render
				err := chromedp.Run(timeoutCtx,
					chromedp.Navigate(currUrl),
					chromedp.WaitVisible(".card-detail", chromedp.ByQuery),
				)
				if err != nil {
					timeoutCancel()
					tabCancel()
					fmt.Printf("Failed to load page for set %v card %d: %v\n", set.Set, i, context.Cause(timeoutCtx))
					continue
				}

				card, err := scrapeCard(timeoutCtx, browserCtx)
				timeoutCancel()
				tabCancel()
				if err != nil {
					fmt.Printf("Failed to scrape card for set %v card %d: %v\n", set.Set, i, err)
					continue
				}
				scrapedCards = append(scrapedCards, card)

			case "images":
				// directly download the image file via URL
				upperCaseUri := strings.ToUpper(set.Uri)
				curlImgUrl := imgURLDE + upperCaseUri + strings.Trim(upperCaseUri, "/") + "_DE_" + fmt.Sprintf("%d", i) + ".png"
				fmt.Printf("Downloading image: %v\n", curlImgUrl)

				err := downloadImage(curlImgUrl, set.Set, i)
				if err != nil {
					fmt.Printf("Failed to download image for set %v card %d: %v\n", set.Set, i, err)
					continue
				}
			}
		}

		if browserCancel != nil {
			browserCancel()
			browserCtx = nil
			browserCancel = nil
		}

		if mode == "cards" {
			// write scraped data to JSON
			fmt.Printf("Scraping set %v finished, writing data...\n", set.Set)
			jsonBytes, err := json.Marshal(scrapedCards)
			if err != nil {
				log.Fatal(err)
			}

			// writing files after scraping
			file, err := os.Create(writePath + set.Set + ".json")
			if err != nil {
				log.Fatal(err)
			}
			file.Write(jsonBytes)
			file.Close()
		}

		fmt.Printf("Finished scraping: %v\n", set.Set)
	}

	println("done!")
}
