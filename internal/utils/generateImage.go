package utils

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/storage"
)

//go:embed assets/pattern-bg-dark.png
var patternBgDarkPNG []byte

//go:embed assets/pattern-bg-naruto.png
var patternBgNarutoPNG []byte

//go:embed assets/digit_0.png
var digit0PNG []byte

//go:embed assets/digit_1.png
var digit1PNG []byte

//go:embed assets/digit_2.png
var digit2PNG []byte

//go:embed assets/digit_3.png
var digit3PNG []byte

//go:embed assets/digit_4.png
var digit4PNG []byte

//go:embed assets/digit_5.png
var digit5PNG []byte

//go:embed assets/digit_6.png
var digit6PNG []byte

//go:embed assets/digit_7.png
var digit7PNG []byte

//go:embed assets/digit_8.png
var digit8PNG []byte

//go:embed assets/digit_9.png
var digit9PNG []byte

const (
	cardWidth   = 240
	cardHeight  = 336
	cols        = 6
	gap         = 32
	stackOffset = 12
	padding     = 40
)

type cardGroup struct {
	s3Key string
	count int
}

func GroupPokemonCards(entries []models.PokemonDeckEntry) []cardGroup {
	type seen struct {
		index int
		count int
	}
	orderMap := make(map[uint]*seen)
	var groups []cardGroup

	for _, e := range entries {
		if s, ok := orderMap[e.Card.ID]; ok {
			s.count++
			groups[s.index].count = s.count
		} else {
			number := strings.Split(e.Card.CardNumber, "/")[0]
			code := strings.ToLower(e.Card.Code)
			key := "cards/pokemon/" + code + "_" + number + ".png"
			orderMap[e.Card.ID] = &seen{index: len(groups), count: 1}
			groups = append(groups, cardGroup{s3Key: key, count: 1})
		}
	}
	return groups
}

func formatCardNumber(entry string) string {
	// var number int

	fmt.Println("Entry before formatting: ", entry)
	if !strings.Contains(entry, "M") {
		number, err := strconv.Atoi(strings.Split(entry, "/")[0])
		if err != nil {
			log.Println("Error formatting card number: ", err)
		}
		fmt.Println("number after formatting: ", number)

		if number < 10 {
			log.Println("Value of number: ", number)
			log.Println("Number: ", fmt.Sprintf("00%v", number))
			return fmt.Sprintf("00%v", number)
		}

		if number < 100 {
			log.Println("Value of number: ", number)
			log.Println("Number: ", fmt.Sprintf("0%v", number))
			return fmt.Sprintf("0%v", number)
		}
	}

	return entry
}

func GroupNarutoCards(entries []models.NarutoDeckEntry) []cardGroup {
	type seen struct {
		index int
		count int
	}
	orderMap := make(map[uint]*seen)
	var groups []cardGroup

	for _, e := range entries {
		if s, ok := orderMap[e.Card.ID]; ok {
			s.count++
			groups[s.index].count = s.count
		} else {
			key := "cards/naruto/" + e.Card.Set + "_" + formatCardNumber(e.Card.CardNumber) + "_" + e.Card.Rarity + ".webp"
			orderMap[e.Card.ID] = &seen{index: len(groups), count: 1}
			groups = append(groups, cardGroup{s3Key: key, count: 1})
		}
	}
	return groups
}

func RenderDeckImage(ctx context.Context, s3Client storage.Presigner, groups []cardGroup, gameType string) *image.RGBA {
	maxStack := 1
	for _, g := range groups {
		if g.count > maxStack {
			maxStack = g.count
		}
	}

	// Download all card images from S3 concurrently
	cardImages := make([]image.Image, len(groups))
	var wg sync.WaitGroup
	for i, g := range groups {
		wg.Add(1)
		go func(idx int, key string) {
			defer wg.Done()
			data, err := s3Client.DownloadObject(ctx, key)
			if err != nil {
				log.Printf("Failed to download %s: %v\n", key, err)
				return
			}
			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				log.Printf("Failed to decode %s: %v\n", key, err)
				return
			}
			cardImages[idx] = img
		}(i, g.s3Key)
	}
	wg.Wait()

	// offset calc
	extraOffset := (maxStack - 1) * stackOffset
	rows := (len(groups) + cols - 1) / cols
	canvasWidth := padding*2 + cols*cardWidth + (cols-1)*gap + extraOffset
	canvasHeight := padding*2 + rows*(cardHeight+extraOffset) + (rows-1)*gap

	img := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	bgImg, err := loadBackgroundImage(gameType)
	if err != nil {
		log.Fatalf("Error loading background image: %v\n", err)
	}

	drawBackground(img, bgImg, 0, 0, canvasWidth, canvasHeight)

	for i, g := range groups {
		cardImg := cardImages[i]
		if cardImg == nil {
			continue
		}

		col := i % cols
		row := i / cols
		baseX := padding + col*(cardWidth+gap)
		baseY := padding + row*(cardHeight+extraOffset+gap)

		for n := 0; n < g.count; n++ {
			/* This part of the code creates the card stacks */
			x := baseX // + n*stackOffset
			y := baseY + n*stackOffset
			drawScaled(img, bgImg, cardImg, x, y, cardWidth, cardHeight)
		}

		if g.count > 1 {
			drawBadge(img, baseX+cardWidth, baseY+cardHeight, g.count)
		}
	}

	return img
}

func loadBackgroundImage(gameType string) (image.Image, error) {
	var imgData []byte
	switch gameType {
	case "pokemon":
		imgData = patternBgDarkPNG
	case "naruto":
		imgData = patternBgNarutoPNG
	default:
		return nil, errors.New("unsupported game type")
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatalf("Error decoding background image: %v\n", err)
		return nil, err
	}
	return img, nil
}

func drawBackground(dst *image.RGBA, bg image.Image, x, y, w, h int) {
	bgBounds := bg.Bounds()
	bgW := bgBounds.Dx()
	bgH := bgBounds.Dy()

	for dy := 0; dy < h; dy++ {
		bgY := bgBounds.Min.Y + (y+dy)%bgH
		for dx := 0; dx < w; dx++ {
			bgX := bgBounds.Min.X + (x+dx)%bgW
			dst.Set(x+dx, y+dy, bg.At(bgX, bgY))
		}
	}
}

func drawScaled(dst *image.RGBA, bg image.Image, src image.Image, x, y, w, h int) {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	const radius = 12

	for dy := 0; dy < h; dy++ {
		srcY := srcBounds.Min.Y + dy*srcH/h
		for dx := 0; dx < w; dx++ {
			if !insideRoundedRect(dx, dy, w, h, radius) {
				continue
			}
			srcX := srcBounds.Min.X + dx*srcW/w
			c := src.At(srcX, srcY)
			if _, _, _, a := c.RGBA(); a > 0 {
				dst.Set(x+dx, y+dy, c)
			}
		}
	}
}

// insideRoundedRect returns true if (px, py) is inside a w×h rectangle
// with the given corner radius.
func insideRoundedRect(px, py, w, h, r int) bool {
	// Check each corner quadrant
	// Top-left
	if px < r && py < r {
		return circleContains(r-px, r-py, r)
	}
	// Top-right
	if px >= w-r && py < r {
		return circleContains(px-(w-r-1), r-py, r)
	}
	// Bottom-left
	if px < r && py >= h-r {
		return circleContains(r-px, py-(h-r-1), r)
	}
	// Bottom-right
	if px >= w-r && py >= h-r {
		return circleContains(px-(w-r-1), py-(h-r-1), r)
	}
	return true
}

func circleContains(dx, dy, r int) bool {
	return dx*dx+dy*dy <= r*r
}

func drawBadge(img *image.RGBA, x, y, count int) {
	text := strconv.Itoa(count)
	badgeWidth := 56
	if len(text) > 1 {
		badgeWidth = 64
	}
	badgeHeight := 56

	bx := x - badgeWidth
	by := y - badgeHeight

	var badgeFile []byte

	switch count {
	case 0:
		badgeFile = digit0PNG
	case 1:
		badgeFile = digit1PNG
	case 2:
		badgeFile = digit2PNG
	case 3:
		badgeFile = digit3PNG
	case 4:
		badgeFile = digit4PNG
	case 5:
		badgeFile = digit5PNG
	case 6:
		badgeFile = digit6PNG
	case 7:
		badgeFile = digit7PNG
	case 8:
		badgeFile = digit8PNG
	case 9:
		badgeFile = digit9PNG
	}

	badge, _, err := image.Decode(bytes.NewReader(badgeFile))
	if err != nil {
		log.Println("Error decoding badge image: ", err)
		return
	}

	srcBounds := badge.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	for dy := 0; dy < badgeHeight; dy++ {
		srcY := srcBounds.Min.Y + dy*srcH/badgeHeight
		for dx := 0; dx < badgeWidth; dx++ {
			srcX := srcBounds.Min.X + dx*srcW/badgeWidth
			c := badge.At(srcX, srcY)
			if _, _, _, a := c.RGBA(); a > 0 {
				img.Set(bx+dx, by+dy, c)
			}
		}
	}

}
