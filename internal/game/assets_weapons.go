package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

func makeFireballTexture(variant int) *ebiten.Image {
	return makePixelImage(12, 12, func(img *image.RGBA) {
		red := rgb(0.92, 0.08, 0.03)
		darkRed := rgb(0.52, 0.03, 0.02)
		orange := rgb(1.00, 0.40, 0.03)
		yellow := rgb(1.00, 0.88, 0.18)
		core, flame := yellow, orange
		if variant != 0 {
			core, flame = orange, yellow
		}
		px(img, 1, 4, 2, 4, darkRed)
		px(img, 2, 2, 7, 8, red)
		px(img, 4, 3, 5, 6, flame)
		px(img, 5, 4, 3, 4, core)
		px(img, 9, 5, 2, 3, flame)
		if variant == 0 {
			px(img, 3, 9, 3, 1, yellow)
			px(img, 2, 1, 2, 1, orange)
		} else {
			px(img, 7, 9, 2, 1, yellow)
			px(img, 1, 2, 2, 1, orange)
		}
	})
}
func makeLightningTexture() *ebiten.Image {
	return makePixelImage(12, 12, func(img *image.RGBA) {
		white := color.RGBA{255, 255, 255, 255}
		blue := rgb(0.18, 0.82, 1)
		darkBlue := rgb(0.04, 0.25, 0.68)
		px(img, 6, 10, 3, 2, white)
		px(img, 5, 8, 4, 2, blue)
		px(img, 4, 6, 3, 2, white)
		px(img, 5, 4, 3, 2, blue)
		px(img, 3, 2, 3, 2, white)
		px(img, 2, 0, 2, 2, blue)
		px(img, 8, 8, 2, 2, darkBlue)
		px(img, 7, 3, 2, 2, darkBlue)
	})
}
func makeOrbTexture(variant int) *ebiten.Image {
	return makePixelImage(12, 12, func(img *image.RGBA) {
		deep := rgb(0.20, 0.04, 0.42)
		purple := rgb(0.48, 0.11, 0.78)
		bright := rgb(0.78, 0.30, 1)
		core := rgb(0.96, 0.78, 1)
		halo, inner := bright, purple
		if variant != 0 {
			halo, inner = purple, bright
		}
		px(img, 4, 1, 4, 1, halo)
		px(img, 2, 3, 8, 6, deep)
		px(img, 3, 2, 6, 8, purple)
		px(img, 4, 3, 5, 6, inner)
		px(img, 5, 4, 3, 4, core)
		if variant == 0 {
			px(img, 2, 8, 2, 2, bright)
			px(img, 8, 2, 2, 2, bright)
		} else {
			px(img, 1, 5, 2, 2, bright)
			px(img, 7, 9, 2, 2, bright)
		}
	})
}
func makeBeamTexture() *ebiten.Image {
	return makePixelImage(12, 12, func(img *image.RGBA) {
		gold := rgb(1.0, 0.72, 0.08)
		yellow := rgb(1.0, 0.94, 0.22)
		white := color.RGBA{255, 255, 255, 255}
		px(img, 1, 4, 10, 4, gold)
		px(img, 0, 5, 12, 2, yellow)
		px(img, 3, 5, 6, 2, white)
		px(img, 2, 8, 2, 1, yellow)
		px(img, 8, 3, 2, 1, yellow)
	})
}
func makeMeteorTexture(variant int) *ebiten.Image {
	return makePixelImage(12, 12, func(img *image.RGBA) {
		darkBrown := rgb(0.22, 0.12, 0.06)
		brown := rgb(0.42, 0.24, 0.11)
		warm := rgb(0.58, 0.36, 0.17)
		highlight := rgb(0.76, 0.53, 0.28)
		shadow, mid := darkBrown, brown
		if variant != 0 {
			shadow, mid = brown, warm
		}
		px(img, 4, 10, 4, 1, shadow)
		px(img, 2, 8, 8, 2, darkBrown)
		px(img, 1, 4, 10, 4, brown)
		px(img, 3, 2, 7, 2, shadow)
		px(img, 4, 5, 5, 3, mid)
		px(img, 5, 7, 3, 2, highlight)
		if variant == 0 {
			px(img, 2, 6, 2, 2, warm)
			px(img, 8, 3, 2, 1, darkBrown)
		} else {
			px(img, 8, 6, 2, 2, warm)
			px(img, 2, 3, 2, 1, darkBrown)
		}
	})
}
