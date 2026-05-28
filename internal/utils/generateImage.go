package utils

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
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

func getFilesDir() string {
	if dir := os.Getenv("FILES_DIR"); dir != "" {
		return dir
	}
	if _, err := os.Stat("/files"); err == nil {
		return "/files"
	}
	return "files"
}

const (
	cardWidth   = 240
	cardHeight  = 336
	cols        = 6
	gap         = 32
	stackOffset = 12
	padding     = 40
)

type cardGroup struct {
	imagePath string
	count     int
}

func GroupPokemonCards(entries []models.PokemonDeckEntry) []cardGroup {
	type seen struct {
		index int
		count int
	}
	orderMap := make(map[uint]*seen)
	var groups []cardGroup
	base := getFilesDir()

	for _, e := range entries {
		if s, ok := orderMap[e.Card.ID]; ok {
			s.count++
			groups[s.index].count = s.count
		} else {
			number := strings.Split(e.Card.CardNumber, "/")[0]
			code := strings.ToLower(e.Card.Code)
			path := filepath.Join(base, "pokemon", "images", code, code+"_"+number+".png")
			orderMap[e.Card.ID] = &seen{index: len(groups), count: 1}
			groups = append(groups, cardGroup{imagePath: path, count: 1})
		}
	}
	return groups
}

func GroupNarutoCards(entries []models.NarutoDeckEntry) []cardGroup {
	type seen struct {
		index int
		count int
	}
	orderMap := make(map[uint]*seen)
	var groups []cardGroup
	base := getFilesDir()

	for _, e := range entries {
		if s, ok := orderMap[e.Card.ID]; ok {
			s.count++
			groups[s.index].count = s.count
		} else {
			path := filepath.Join(base, "naruto", "images", e.Card.CardNumber+"_"+e.Card.Rarity+".webp")
			orderMap[e.Card.ID] = &seen{index: len(groups), count: 1}
			groups = append(groups, cardGroup{imagePath: path, count: 1})
		}
	}
	return groups
}

func RenderDeckImage(groups []cardGroup, gameType string) *image.RGBA {
	maxStack := 1
	for _, g := range groups {
		if g.count > maxStack {
			maxStack = g.count
		}
	}

	// offset calc
	extraOffset := (maxStack - 1) * stackOffset
	rows := (len(groups) + cols - 1) / cols
	canvasWidth := padding*2 + cols*cardWidth + (cols-1)*gap + extraOffset
	canvasHeight := padding*2 + rows*(cardHeight+extraOffset) + (rows-1)*gap

	// log.Println("Values of Canvas: ", canvasWidth, canvasHeight)

	img := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	bgImg, err := loadBackgroundImage(gameType)
	if err != nil {
		log.Fatalf("Error loading background image: %v\n", err)
	}

	// draw background for entire canvas first

	drawBackground(img, bgImg, 0, 0, canvasWidth, canvasHeight)

	for i, g := range groups {
		cardImg := loadImageFile(g.imagePath)
		if cardImg == nil {
			continue
		}

		col := i % cols
		row := i / cols
		baseX := padding + col*(cardWidth+gap)
		baseY := padding + row*(cardHeight+extraOffset+gap)

		for n := 0; n < g.count; n++ {
			x := baseX + n*stackOffset
			y := baseY + n*stackOffset
			drawScaled(img, bgImg, cardImg, x, y, cardWidth, cardHeight)
		}

		// if g.count > 1 {
		// 	badgeX := baseX + (g.count-1)*stackOffset + cardWidth - 16
		// 	badgeY := baseY + (g.count-1)*stackOffset + cardHeight - 16
		// 	drawBadge(img, badgeX, badgeY, g.count)
		// }
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

func loadImageFile(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return img
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

	// then draw card on top
	for dy := 0; dy < h; dy++ {
		srcY := srcBounds.Min.Y + dy*srcH/h
		for dx := 0; dx < w; dx++ {
			srcX := srcBounds.Min.X + dx*srcW/w
			dst.Set(x+dx, y+dy, src.At(srcX, srcY))
		}
	}
}

// func drawBadge(img *image.RGBA, x, y, count int) {
// 	text := strconv.Itoa(count)
// 	badgeWidth := 44
// 	if len(text) > 1 {
// 		badgeWidth = 52
// 	}
// 	badgeHeight := 44

// 	bx := x - badgeWidth
// 	by := y - badgeHeight

// 	var badgeFile []byte

// 	switch count {
// 	case 0:
// 		badgeFile = digit0PNG
// 	case 1:
// 		badgeFile = digit1PNG
// 	case 2:
// 		badgeFile = digit2PNG
// 	case 3:
// 		badgeFile = digit3PNG
// 	case 4:
// 		badgeFile = digit4PNG
// 	case 5:
// 		badgeFile = digit5PNG
// 	case 6:
// 		badgeFile = digit6PNG
// 	case 7:
// 		badgeFile = digit7PNG
// 	case 8:
// 		badgeFile = digit8PNG
// 	case 9:
// 		badgeFile = digit9PNG
// 	}

// 	badge, _, err := image.Decode(bytes.NewReader(badgeFile))
// 	if err != nil {
// 		log.Println("Error decoding badge image: ", err)
// 		return
// 	}

// 	for dy := 0; dy < badgeHeight; dy++ {
// 		for dx := 0; dx < badgeWidth; dx++ {
// 			c := badge.At(badge.Bounds().Min.X+dx, badge.Bounds().Min.Y+dy)
// 			if _, _, _, a := c.RGBA(); a > 0 {
// 				img.Set(bx+dx, by+dy, c)
// 			}
// 		}
// 	}

// }
