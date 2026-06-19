package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"math"
)

func makePixelImage(width, height int, draw func(*image.RGBA)) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw(img)
	if generatedAssetPixelCapture != nil {
		generatedAssetPixelCapture(img)
	}
	return ebiten.NewImageFromImage(img)
}
func px(img *image.RGBA, x, y, w, h int, clr color.RGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	top := img.Bounds().Dy() - y - h
	for yy := max(0, top); yy < min(img.Bounds().Dy(), top+h); yy++ {
		for xx := max(0, x); xx < min(img.Bounds().Dx(), x+w); xx++ {
			img.SetRGBA(xx, yy, clr)
		}
	}
}
func rgb(r, g, b float64) color.RGBA {
	return color.RGBA{
		uint8(math.Round(clamp01(r) * 255)),
		uint8(math.Round(clamp01(g) * 255)),
		uint8(math.Round(clamp01(b) * 255)),
		255,
	}
}
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
func fill(img *image.RGBA, clr color.RGBA) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, clr)
		}
	}
}
