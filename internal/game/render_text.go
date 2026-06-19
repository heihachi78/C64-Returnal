package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"image/color"
	"math"
	"os"
	"strings"
)

var (
	hudFont       *opentype.Font
	hudFontSource string
	hudFontFaces  = map[int]font.Face{}
)

var hudFontPaths = []string{
	"/System/Library/Fonts/Menlo.ttc",
	"/Library/Fonts/Menlo.ttc",
}

var newOpenTypeFace = opentype.NewFace

func init() {
	hudFont, hudFontSource = loadHUDFont()
}
func loadHUDFont() (*opentype.Font, string) {
	if font, name := loadSystemFontByFullName(
		hudFontPaths,
		[]string{"menlo", "bold"},
	); font != nil {
		return font, name
	}
	font, _ := opentype.Parse(gomonobold.TTF)
	return font, "Go Mono Bold"
}
func loadSystemFontByFullName(paths, required []string) (*opentype.Font, string) {
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		collection, err := opentype.ParseCollection(data)
		if err != nil {
			continue
		}
		for i := 0; i < collection.NumFonts(); i++ {
			font, _ := collection.Font(i)
			name, err := font.Name(nil, sfnt.NameIDFull)
			if err != nil || !fontNameMatches(name, required) {
				continue
			}
			return font, name
		}
	}
	return nil, ""
}
func fontNameMatches(name string, required []string) bool {
	lower := strings.ToLower(name)
	for _, token := range required {
		if !strings.Contains(lower, strings.ToLower(token)) {
			return false
		}
	}
	return true
}
func fontFaceForSize(size float64) font.Face {
	if hudFont == nil || size <= 0 {
		return basicfont.Face7x13
	}
	key := fontSizeKey(size)
	if face := hudFontFaces[key]; face != nil {
		return face
	}
	face, err := newOpenTypeFace(hudFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return basicfont.Face7x13
	}
	hudFontFaces[key] = face
	return face
}
func centeredTextBaseline(face font.Face, centerY float64) int {
	metrics := face.Metrics()
	return int(math.Round(centerY + float64((metrics.Ascent-metrics.Descent).Round())/2))
}
func (g *Game) drawTextSize(screen *ebiten.Image, s string, x, centerY, size float64, clr color.Color) {
	face := fontFaceForSize(size)
	text.Draw(screen, s, face, int(math.Round(x)), centeredTextBaseline(face, centerY), clr)
}
func (g *Game) drawTextSizeScaled(screen *ebiten.Image, s string, x, centerY, size, scale float64, clr color.Color) {
	g.drawScaledTextImage(screen, s, x, centerY, size, scale, false, clr)
}
func (g *Game) drawCenteredTextSize(screen *ebiten.Image, s string, x, centerY, size float64, clr color.Color) {
	face := fontFaceForSize(size)
	width := font.MeasureString(face, s).Ceil()
	text.Draw(screen, s, face, int(math.Round(x))-width/2, centeredTextBaseline(face, centerY), clr)
}
func (g *Game) drawCenteredTextSizeScaled(screen *ebiten.Image, s string, x, centerY, size, scale float64, clr color.Color) {
	g.drawScaledTextImage(screen, s, x, centerY, size, scale, true, clr)
}

type scaledTextLayout struct {
	Width    int
	Height   int
	AnchorX  float64
	AnchorY  float64
	Baseline int
}
type scaledTextCacheKey struct {
	Text     string
	SizeKey  int
	Centered bool
}
type scaledTextCacheEntry struct {
	Image  *ebiten.Image
	Layout scaledTextLayout
}

func fontSizeKey(size float64) int {
	return int(math.Round(size * 10))
}
func baseScaledTextLayout(face font.Face, s string, centered bool) scaledTextLayout {
	const padding = 4
	metrics := face.Metrics()
	textWidth := max(1, font.MeasureString(face, s).Ceil())
	textHeight := max(1, (metrics.Ascent + metrics.Descent).Ceil())
	anchorX := float64(padding)
	if centered {
		anchorX += float64(textWidth) / 2
	}
	return scaledTextLayout{
		Width:    textWidth + padding*2,
		Height:   textHeight + padding*2,
		AnchorX:  anchorX,
		AnchorY:  float64(padding) + float64(textHeight)/2,
		Baseline: padding + metrics.Ascent.Ceil(),
	}
}
func (g *Game) scaledTextImage(s string, size float64, centered bool) scaledTextCacheEntry {
	key := scaledTextCacheKey{Text: s, SizeKey: fontSizeKey(size), Centered: centered}
	if entry := g.scaledTextCache[key]; entry.Image != nil {
		return entry
	}

	face := fontFaceForSize(size)
	layout := baseScaledTextLayout(face, s, centered)
	img := ebiten.NewImage(layout.Width, layout.Height)
	text.Draw(img, s, face, 4, layout.Baseline, color.White)
	entry := scaledTextCacheEntry{Image: img, Layout: layout}
	g.scaledTextCache[key] = entry
	return entry
}
func (g *Game) drawScaledTextImage(screen *ebiten.Image, s string, x, centerY, size, scale float64, centered bool, clr color.Color) {
	if scale <= 0 {
		return
	}
	entry := g.scaledTextImage(s, size, centered)
	layout := entry.Layout

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x-layout.AnchorX*scale, centerY-layout.AnchorY*scale)
	op.ColorScale.ScaleWithColor(clr)
	screen.DrawImage(entry.Image, op)
}
