package controller

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "image/jpeg"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"github.com/MarcoPasta/deckmaster/api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	_ "golang.org/x/image/webp"
	xdraw "golang.org/x/image/draw"
)

//go:embed ../assets/pattern-bg-dark.svg
var patternBgDarkSVG []byte

//go:embed ../assets/pattern-bg-naruto.svg
var patternBgNarutoSVG []byte

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

func GetPokemonDeckImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}

	deck, err := repository.GetPokemonDeckById(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}

	groups := groupPokemonCards(deck.Entries)
	img := renderDeckImage(groups, "pokemon")

	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", `attachment; filename="`+deck.Name+`.png"`)
	png.Encode(ctx.Writer, img)
}

func GetNarutoDeckImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid deck id"})
		return
	}

	deck, err := repository.GetNarutoDeckById(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "deck not found"})
		return
	}

	groups := groupNarutoCards(deck.Entries)
	img := renderDeckImage(groups, "naruto")

	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", `attachment; filename="`+deck.Name+`.png"`)
	png.Encode(ctx.Writer, img)
}

func groupPokemonCards(entries []models.PokemonDeckEntry) []cardGroup {
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

func groupNarutoCards(entries []models.NarutoDeckEntry) []cardGroup {
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

func renderDeckImage(groups []cardGroup, gameType string) *image.RGBA {
	maxStack := 1
	for _, g := range groups {
		if g.count > maxStack {
			maxStack = g.count
		}
	}
	extraOffset := (maxStack - 1) * stackOffset
	rows := (len(groups) + cols - 1) / cols
	canvasWidth := padding*2 + cols*cardWidth + (cols-1)*gap + extraOffset
	canvasHeight := padding*2 + rows*(cardHeight+extraOffset) + (rows-1)*gap

	img := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	bgImg := rasterizeSVGBackground(gameType, canvasWidth, canvasHeight)
	if bgImg != nil {
		draw.Draw(img, img.Bounds(), bgImg, image.Point{}, draw.Src)
	} else {
		bg := color.RGBA{0x16, 0x17, 0x1d, 0xff}
		draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	}

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
			drawScaledHighQuality(img, cardImg, x, y, cardWidth, cardHeight)
		}

		if g.count > 1 {
			badgeX := baseX + (g.count-1)*stackOffset + cardWidth - 16
			badgeY := baseY + (g.count-1)*stackOffset + cardHeight - 16
			drawBadge(img, badgeX, badgeY, g.count)
		}
	}

	return img
}

func rasterizeSVGBackground(gameType string, canvasWidth, canvasHeight int) *image.RGBA {
	var svgData []byte
	switch gameType {
	case "pokemon":
		svgData = patternBgDarkSVG
	case "naruto":
		svgData = patternBgNarutoSVG
	default:
		return nil
	}

	icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData))
	if err != nil {
		return nil
	}

	svgW := int(icon.ViewBox.W)
	svgH := int(icon.ViewBox.H)
	if svgW == 0 || svgH == 0 {
		return nil
	}

	icon.SetTarget(0, 0, float64(svgW), float64(svgH))
	tile := image.NewRGBA(image.Rect(0, 0, svgW, svgH))
	scanner := rasterx.NewScannerGV(svgW, svgH, tile, tile.Bounds())
	dasher := rasterx.NewDasher(svgW, svgH, scanner)
	icon.Draw(dasher, 1.0)

	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	for y := 0; y < canvasHeight; y += svgH {
		for x := 0; x < canvasWidth; x += svgW {
			draw.Draw(canvas, image.Rect(x, y, x+svgW, y+svgH), tile, image.Point{}, draw.Src)
		}
	}
	return canvas
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

func drawScaled(dst *image.RGBA, src image.Image, x, y, w, h int) {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	for dy := 0; dy < h; dy++ {
		srcY := srcBounds.Min.Y + dy*srcH/h
		for dx := 0; dx < w; dx++ {
			srcX := srcBounds.Min.X + dx*srcW/w
			dst.Set(x+dx, y+dy, src.At(srcX, srcY))
		}
	}
}

func drawBadge(img *image.RGBA, x, y, count int) {
	text := strconv.Itoa(count)
	badgeWidth := 22
	if len(text) > 1 {
		badgeWidth = 28
	}
	badgeHeight := 22

	bx := x - badgeWidth
	by := y - badgeHeight

	purple := color.RGBA{0x7c, 0x4d, 0xff, 0xff}
	for dy := 0; dy < badgeHeight; dy++ {
		for dx := 0; dx < badgeWidth; dx++ {
			px := bx + dx
			py := by + dy
			if px >= 0 && py >= 0 && px < img.Bounds().Dx() && py < img.Bounds().Dy() {
				cx := dx - badgeWidth/2
				cy := dy - badgeHeight/2
				rx := float64(badgeWidth) / 2
				ry := float64(badgeHeight) / 2
				if float64(cx*cx)/(rx*rx)+float64(cy*cy)/(ry*ry) <= 1.0 {
					img.Set(px, py, purple)
				}
			}
		}
	}

	drawDigit(img, bx+badgeWidth/2, by+badgeHeight/2, text)
}

func drawDigit(img *image.RGBA, cx, cy int, text string) {
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}

	digits := map[byte][]string{
		'1': {
			"  #  ",
			" ##  ",
			"  #  ",
			"  #  ",
			" ### ",
		},
		'2': {
			" ### ",
			"#   #",
			"  ## ",
			" #   ",
			"#####",
		},
		'3': {
			" ### ",
			"#   #",
			"  ## ",
			"#   #",
			" ### ",
		},
		'4': {
			"#  # ",
			"#  # ",
			"#####",
			"   # ",
			"   # ",
		},
		'5': {
			"#####",
			"#    ",
			"#### ",
			"    #",
			"#### ",
		},
		'6': {
			" ### ",
			"#    ",
			"#### ",
			"#   #",
			" ### ",
		},
		'7': {
			"#####",
			"   # ",
			"  #  ",
			" #   ",
			" #   ",
		},
		'8': {
			" ### ",
			"#   #",
			" ### ",
			"#   #",
			" ### ",
		},
		'9': {
			" ### ",
			"#   #",
			" ####",
			"    #",
			" ### ",
		},
		'0': {
			" ### ",
			"#   #",
			"#   #",
			"#   #",
			" ### ",
		},
	}

	totalWidth := len(text)*5 + (len(text)-1)*1
	startX := cx - totalWidth/2

	for i, ch := range []byte(text) {
		glyph, ok := digits[ch]
		if !ok {
			continue
		}
		ox := startX + i*6
		oy := cy - 2
		for row, line := range glyph {
			for col, c := range line {
				if c == '#' {
					img.Set(ox+col, oy+row, white)
				}
			}
		}
	}
}
