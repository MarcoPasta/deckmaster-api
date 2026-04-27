// developed by marcopasta

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/gocolly/colly/v2"
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
}

// // Take string value and parse it to int8 value
// func toInt8(s string) int8 {
// 	parsed, err := strconv.ParseInt(s, 10, 8)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return int8(parsed)
// }

// Take string value and parse it to int16 value
func toInt16(s string) int16 {
	parsed, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	return int16(parsed)
}

// // normalizeSpaces trims and replaces any whitespace run with a single space.
// func normalizeSpaces(s string) string {
// 	var reSpaces = regexp.MustCompile(`\s+`)
// 	return strings.TrimSpace(reSpaces.ReplaceAllString(s, " "))
// }

// // preserves bracket-only headers (e.g. "[Effect]") as items
// // and joins subsequent non-header lines into one paragraph.
// func cleanGrouped(lines []string) []string {
// 	var reHeader = regexp.MustCompile(`^\[[^\]]+\]$`)

// 	out := make([]string, 0, len(lines))
// 	var buf string
// 	flushBuf := func() {
// 		if buf != "" {
// 			out = append(out, buf)
// 			buf = ""
// 		}
// 	}

// 	for _, l := range lines {
// 		t := normalizeSpaces(l)
// 		if t == "" {
// 			continue
// 		}
// 		if reHeader.MatchString(t) {
// 			flushBuf()
// 			out = append(out, t) // header on its own line
// 			continue
// 		}
// 		if buf == "" {
// 			buf = t
// 		} else {
// 			buf += " " + t
// 		}
// 	}
// 	flushBuf()
// 	return out
// }

// // if values are given that are camel cased, it searches for those values and splits them with a seperator
// func splitByCapital(s string, sep string) string {
// 	var result strings.Builder
// 	var prev rune
// 	for i, curr := range s {
// 		if i > 0 && curr >= 'A' && curr <= 'Z' && prev != ' ' {
// 			result.WriteString(sep)
// 		}
// 		result.WriteRune(curr)
// 		prev = curr
// 	}
// 	return result.String()
// }

// // searches through each css identifier containing the information for the digivolve cost. If multiple identifiers are found,
// // it loops through each and creates a new string inside
// func groupDvCosts(cost string) []string {
// 	var out []string
// 	split := strings.Split(cost, "\n")
// 	for _, s := range split {
// 		normalized := normalizeSpaces(s)
// 		transformed := splitByCapital(normalized, "/")
// 		out = append(out, transformed)
// 	}
// 	return out
// }

func onHTMLScrape(e *colly.HTMLElement) Card {
	card := Card{}
	// card.Id = e.ChildText(".cardNo")
	card.Name = e.ChildText(".color-block")
	card.Kp = toInt16(e.ChildText(".card-hp"))
	card.Phase = e.ChildText(".card-type")
	card.EnergyType = e.ChildAttr(".energy", "title")

	// e.ForEach(".cardInfoBox", func(i int, h *colly.HTMLElement) {

	// 	switch h.ChildText(".cardInfoTit") {
	// 	case "Color":
	// 		card.Color = strings.Fields(h.ChildText(".cardInfoData"))
	// 	case "Cost":
	// 		card.Cost = toInt8(h.ChildText(".cardInfoData"))
	// 	case "DP":
	// 		card.Dp = toInt16(h.ChildText(".cardInfoData"))
	// 	case "Form":
	// 		card.Form = h.ChildText(".cardInfoData")
	// 	case "Attribute":
	// 		card.Attribute = h.ChildText(".cardInfoData")
	// 	case "Type":
	// 		// we have to case for "Type" here as it is the used keyword in the html
	// 		card.Species = strings.Fields(h.ChildText(".cardInfoData"))
	// 	case "Notes":
	// 		// same goes for Set, the html keyword used here is "Notes"
	// 		card.Set = strings.Join(strings.Fields(h.ChildText(".cardInfoData"))[1:4], " ")
	// 	}
	// })
	// card.DvCost = groupDvCosts(e.ChildText(".cardInfoData.cardColor.cardColorCost"))
	// card.Text = cleanGrouped(strings.Split(e.ChildText(".cardInfoBoxSmall"), "\n"))
	// return card
	fmt.Println("CARD DETAILS: ", card)
	return Card{}
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

// func scrapeImage(url string, set string, cards []Card) {
// 	path := os.Getenv("DIGIMON_IMAGES_DIR") + set + "/"

// 	println("Start scraping images")
// 	for _, card := range cards {
// 		imgUrl := url + card.Id + ".png"
// 		response, err := http.Get(imgUrl)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer response.Body.Close()
// 		checkPathExist(path)
// 		file, err := os.Create(path + card.Id + ".png")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer file.Close()
// 		_, err = io.Copy(file, response.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("Saved image %v successfully!\n", card.Id)
// 	}

// }

// used to change user agent after each scrape to avoid being blocked by the server
func RandomString() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {

	// collector object
	collector := colly.NewCollector()

	baseUrl := "https://www.pokemon.com/de/pokemon-sammelkartenspiel/pokemon-karten/series"
	// imgUrl := "https://world.digimoncard.com/images/cardlist/card/"

	/* scraping details:
	url: https://www.pokemon.com/de/pokemon-sammelkartenspiel/pokemon-karten/series/EX0?/Y/ --> Extension & Number followed by the card number, Example: /me02/1/
	Each set or starter/advanced deck contains a category which indicates the cards being placed inside
	example: the new st-22 has the category 522122, while st-21 has category 522121 and st-20 has category 5122120.
	Opening a card results in adding a "#card-X" value to the categorie.
	All informations relevant to a card is inside a <div> which has the class "cardInfoCol"
	*/

	var scrapedCards []Card

	// category json object
	var sets []Set
	data, err := os.ReadFile("./collection.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &sets)
	if err != nil {
		log.Fatal(err)
	}

	var currSet *string

	// notify when sending the request
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Start scraping...\n", *currSet)
	})

	// define the scraping logic
	collector.OnHTML(".cardDetailColInner", func(e *colly.HTMLElement) {
		scrapedCards = append(scrapedCards, onHTMLScrape(e))
	})

	collector.OnScraped(func(r *colly.Response) {

		r.Headers.Set("User-Agent", RandomString())
		fmt.Printf("Scraping set %v finished, writing data...\n", *currSet)
		fmt.Println("Changing User-Agent...")

		jsonBytes, err := json.Marshal(scrapedCards)
		if err != nil {
			log.Fatal(err)
		}

		writePath := os.Getenv("POKEMON_JSON_DIR")
		checkPathExist(writePath) // ensure path exists, if not create it
		file, err := os.Create(writePath + *currSet + ".json")
		if err != nil {
			log.Fatal(err)
		}
		file.Write(jsonBytes)
		defer file.Close()

		// scrapeImage(imgUrl, *currSet, scrapedCards)

		println("finished scraping: ", *currSet)
		println("Clearing Array....")
		scrapedCards = nil
	})

	for _, set := range sets {
		currUrl := baseUrl + set.Uri
		currSet = &set.Set

		// visit page & start scraping
		collector.Visit(currUrl)

	}

	println("done!")

}
