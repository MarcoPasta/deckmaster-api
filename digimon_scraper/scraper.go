// developed by marcopasta

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type DvCost struct {
	colors []string
	cost   int8
	lvl    string
}

// TODO: add image link
type Card struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	CardType  string   `json:"cardType"`
	Color     []string `json:"color"`
	Cost      int8     `json:"cost"`
	Dp        int16    `json:"dp"`
	DvCost1   DvCost   `json:"dvCost1"`
	DvCost2   DvCost   `json:"dvCost2"`
	Form      string   `json:"form"`
	Attribute string   `json:"attribute"`
	Species   []string `json:"species"` // type is a primitve, so we have to use species here
	CardText  []string `json:"text"`
	Set       string   `json:"set"`
}

// Take string value and parse it to int8 value
func toInt8(s string) int8 {
	parsed, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	return int8(parsed)
}

// Take string value and parse it to int16 value
func toInt16(s string) int16 {
	parsed, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	return int16(parsed)
}

// Takes the concat. string and marhalls it into a DvCost object
func splitByDvCosts(dvCosts string) DvCost {
	re := regexp.MustCompile(`[A-Z][a-z]*`)
	var obj DvCost
	fields := strings.Fields(dvCosts)
	obj.colors = re.FindAllString(fields[0], -1)
	if len(fields) != 0 {
		obj.cost = toInt8(fields[1])
	}
	obj.lvl = fields[3]
	return obj
}

// AI GENERATED
// normalizeSpaces trims and replaces any whitespace run with a single space.

func normalizeSpaces(s string) string {
	var reSpaces = regexp.MustCompile(`\s+`)
	return strings.TrimSpace(reSpaces.ReplaceAllString(s, " "))
}

// AI GENERATED
// preserves bracket-only headers (e.g. "[Effect]") as items
// and joins subsequent non-header lines into one paragraph.
func cleanCardTextGrouped(lines []string) []string {
	var reHeader = regexp.MustCompile(`^\[[^\]]+\]$`)

	out := make([]string, 0, len(lines))
	var buf string
	flushBuf := func() {
		if buf != "" {
			out = append(out, buf)
			buf = ""
		}
	}

	for _, l := range lines {
		t := normalizeSpaces(l)
		if t == "" {
			continue
		}
		if reHeader.MatchString(t) {
			flushBuf()
			out = append(out, t) // header on its own line
			continue
		}
		if buf == "" {
			buf = t
		} else {
			buf += " " + t
		}
	}
	flushBuf()
	return out
}

func onHTMLScrape(e *colly.HTMLElement) Card {
	card := Card{}
	card.Id = e.ChildText(".cardNo")
	card.Name = e.ChildText(".cardTitle")
	card.CardType = e.ChildText(".cardType")

	e.ForEach(".cardInfoBox", func(i int, h *colly.HTMLElement) {

		switch h.ChildText(".cardInfoTit") {
		case "Color":
			card.Color = strings.Fields(h.ChildText(".cardInfoData"))
		case "Cost":
			card.Cost = toInt8(h.ChildText(".cardInfoData"))
		case "DP":
			card.Dp = toInt16(h.ChildText(".cardInfoData"))
		case "Digivolve Cost 1": // TODO: repair this as it is not working and also fix the digivolution costs
			dvCosts := h.ChildText(".cardInfoData.cardColor.cardColorCost")
			card.DvCost1 = splitByDvCosts(dvCosts)
		case "Form":
			card.Form = h.ChildText(".cardInfoData")
		case "Attribute":
			card.Attribute = h.ChildText(".cardInfoData")
		case "Type":
			// we have to case for "Type" here as it is the used keyword in the html
			card.Species = strings.Fields(h.ChildText(".cardInfoData"))
		case "Notes":
			// same goes for Set, the html keyword used here is "Notes"
			card.Set = strings.Join(strings.Fields(h.ChildText(".cardInfoData"))[1:4], " ")
		}
	})

	card.CardText = cleanCardTextGrouped(strings.Split(e.ChildText(".cardInfoBoxSmall"), "\n"))

	return card
}

func main() {

	// collector object
	collector := colly.NewCollector()

	/* scraping details:
	url: https://world.digimoncard.com/cards/?search=true&category=522122
	Each set or starter/advanced deck contains a category which indicates the cards being placed inside
	example: the new st-22 has the category 522122, while st-21 has category 522121 and st-20 has category 5122120.
	Opening a card results in adding a "#card-X" value to the categorie.
	All informations relevant to a card is inside a <div> which has the class "cardInfoCol"
	*/

	var scrapedCards []Card

	// notify when sending the request
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Sending request to website")
	})

	// define the scraping logic
	collector.OnHTML(".cardDetailColInner", func(e *colly.HTMLElement) {
		// TODO: add multiple set categories to iterate over and scrape each set
		scrapedCards = append(scrapedCards, onHTMLScrape(e))
	})

	collector.OnScraped(func(r *colly.Response) {

		fmt.Println("after scraping")

		jsonBytes, err := json.Marshal(scrapedCards)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create("./files/bt-24.json")
		if err != nil {
			log.Fatal(err)
		}
		file.Write(jsonBytes)
	})

	// visit page & start scraping
	collector.Visit("https://world.digimoncard.com/cards/?search=true&category=522033")

	println("done!")

}
