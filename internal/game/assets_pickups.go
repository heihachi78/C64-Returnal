package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

func makeLifeTexture() *ebiten.Image {
	return makePixelImage(12, 12, func(img *image.RGBA) {
		red := rgb(0.86, 0.05, 0.10)
		bright := rgb(1, 0.18, 0.22)
		highlight := rgb(1, 0.62, 0.62)
		shadow := rgb(0.45, 0.02, 0.06)
		px(img, 2, 7, 3, 3, red)
		px(img, 7, 7, 3, 3, red)
		px(img, 1, 5, 10, 3, bright)
		px(img, 2, 3, 8, 2, red)
		px(img, 4, 1, 4, 2, shadow)
		px(img, 3, 8, 2, 1, highlight)
	})
}
func makeChestTexture(tier ChestTier) *ebiten.Image {
	return makePixelImage(16, 14, func(img *image.RGBA) {
		base, light, dark := chestColors(tier)
		outline := rgb(0.09, 0.05, 0.03)
		darkStrap := rgb(0.14, 0.08, 0.04)
		lock := rgb(1.0, 0.86, 0.28)
		px(img, 3, 11, 10, 1, outline)
		px(img, 2, 9, 12, 2, dark)
		px(img, 3, 10, 10, 1, light)
		px(img, 1, 3, 14, 6, outline)
		px(img, 2, 4, 12, 5, base)
		px(img, 2, 7, 12, 1, light)
		px(img, 2, 4, 12, 1, dark)
		px(img, 1, 2, 14, 1, outline)
		px(img, 7, 3, 2, 7, darkStrap)
		px(img, 6, 5, 4, 3, outline)
		px(img, 7, 5, 2, 2, lock)
		px(img, 3, 8, 3, 1, light)
	})
}
func chestColors(tier ChestTier) (color.RGBA, color.RGBA, color.RGBA) {
	switch tier {
	case ChestSilver:
		return rgb(0.58, 0.63, 0.68), rgb(0.88, 0.93, 0.96), rgb(0.34, 0.38, 0.43)
	case ChestGold:
		return rgb(0.86, 0.58, 0.08), rgb(1.0, 0.86, 0.25), rgb(0.48, 0.30, 0.03)
	default:
		return rgb(0.55, 0.30, 0.13), rgb(0.86, 0.52, 0.23), rgb(0.32, 0.16, 0.07)
	}
}
func makeCoinTexture(variant int) *ebiten.Image {
	return makePixelImage(12, 12, func(img *image.RGBA) {
		outline := rgb(0.42, 0.24, 0.03)
		dark := rgb(0.72, 0.43, 0.04)
		gold := rgb(1, 0.73, 0.08)
		yellow := rgb(1, 0.92, 0.24)
		white := rgb(1, 0.98, 0.72)
		switch variant {
		case 0:
			px(img, 3, 1, 6, 1, outline)
			px(img, 2, 2, 8, 1, outline)
			px(img, 1, 3, 10, 6, outline)
			px(img, 2, 9, 8, 1, outline)
			px(img, 3, 10, 6, 1, outline)
			px(img, 2, 3, 8, 6, gold)
			px(img, 3, 2, 6, 8, gold)
			px(img, 3, 7, 5, 2, yellow)
			px(img, 4, 4, 3, 2, dark)
			px(img, 4, 8, 2, 1, white)
		case 1, 3:
			px(img, 4, 1, 4, 1, outline)
			px(img, 3, 2, 6, 1, outline)
			px(img, 3, 3, 6, 6, outline)
			px(img, 3, 9, 6, 1, outline)
			px(img, 4, 10, 4, 1, outline)
			px(img, 4, 2, 4, 8, gold)
			px(img, 5, 3, 2, 6, yellow)
			px(img, 4, 4, 1, 4, dark)
			px(img, 5, 8, 1, 1, white)
		default:
			px(img, 5, 1, 2, 1, outline)
			px(img, 4, 2, 4, 1, outline)
			px(img, 4, 3, 4, 6, outline)
			px(img, 4, 9, 4, 1, outline)
			px(img, 5, 10, 2, 1, outline)
			px(img, 5, 2, 2, 8, gold)
			px(img, 6, 3, 1, 6, yellow)
			px(img, 5, 4, 1, 4, dark)
			px(img, 6, 8, 1, 1, white)
		}
	})
}
