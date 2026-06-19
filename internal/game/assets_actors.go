package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func makeMageTexture(variant int) *ebiten.Image {
	return makePixelImage(16, 22, func(img *image.RGBA) {
		leftFootX, rightFootX, staffX, crystalX, crystalY := 4, 9, 13, 12, 17
		if variant != 0 {
			leftFootX, rightFootX, staffX, crystalX, crystalY = 3, 10, 14, 13, 16
		}
		px(img, 5, 20, 6, 1, rgb(0.13, 0.07, 0.24))
		px(img, 6, 19, 4, 1, rgb(0.38, 0.18, 0.65))
		px(img, 4, 17, 8, 2, rgb(0.27, 0.12, 0.49))
		px(img, 3, 16, 10, 1, rgb(0.13, 0.07, 0.24))
		px(img, 4, 11, 8, 6, rgb(0.32, 0.13, 0.56))
		px(img, 6, 12, 4, 4, rgb(0.88, 0.65, 0.47))
		px(img, 7, 14, 1, 1, rgb(0.12, 0.08, 0.08))
		px(img, 10, 14, 1, 1, rgb(0.12, 0.08, 0.08))
		px(img, 4, 4, 8, 8, rgb(0.14, 0.24, 0.66))
		px(img, 3, 3, 10, 2, rgb(0.09, 0.15, 0.39))
		px(img, 5, 5, 2, 6, rgb(0.22, 0.35, 0.83))
		px(img, 9, 5, 2, 6, rgb(0.09, 0.15, 0.39))
		if variant == 0 {
			px(img, 2, 6, 3, 3, rgb(0.32, 0.13, 0.56))
			px(img, 11, 7, 2, 2, rgb(0.32, 0.13, 0.56))
		} else {
			px(img, 3, 7, 2, 2, rgb(0.32, 0.13, 0.56))
			px(img, 11, 5, 3, 3, rgb(0.32, 0.13, 0.56))
		}
		px(img, leftFootX, 1, 3, 2, rgb(0.08, 0.07, 0.16))
		px(img, rightFootX, 1, 3, 2, rgb(0.08, 0.07, 0.16))
		px(img, staffX, 4, 1, 13, rgb(0.37, 0.20, 0.09))
		px(img, crystalX, crystalY, 3, 3, rgb(0.32, 0.86, 0.95))
		px(img, crystalX+1, crystalY+1, 1, 1, color.RGBA{255, 255, 255, 255})
	})
}
func makeSkeletonTexture(variant int) *ebiten.Image {
	return makePixelImage(16, 22, func(img *image.RGBA) {
		bone := rgb(0.82, 0.84, 0.76)
		bright := rgb(0.94, 0.96, 0.86)
		shadow := rgb(0.48, 0.51, 0.46)
		dark := rgb(0.06, 0.07, 0.08)
		rust := rgb(0.42, 0.28, 0.13)
		px(img, 5, 14, 7, 6, bone)
		px(img, 6, 18, 5, 2, bright)
		px(img, 5, 14, 1, 4, shadow)
		px(img, 7, 17, 1, 1, dark)
		px(img, 10, 17, 1, 1, dark)
		px(img, 9, 15, 1, 1, dark)
		px(img, 7, 11, 3, 3, bone)
		px(img, 5, 8, 7, 4, bone)
		px(img, 6, 9, 5, 1, dark)
		px(img, 5, 7, 1, 4, shadow)
		px(img, 11, 7, 1, 4, shadow)
		if variant == 0 {
			px(img, 3, 8, 2, 1, bone)
			px(img, 2, 5, 1, 4, bone)
			px(img, 12, 8, 2, 1, bone)
			px(img, 13, 5, 1, 4, bone)
			px(img, 6, 4, 2, 4, bone)
			px(img, 10, 4, 2, 4, bone)
			px(img, 5, 2, 3, 2, bone)
			px(img, 10, 2, 3, 2, bone)
			px(img, 5, 1, 4, 1, shadow)
			px(img, 10, 1, 4, 1, shadow)
		} else {
			px(img, 3, 7, 2, 1, bone)
			px(img, 2, 4, 1, 4, bone)
			px(img, 12, 9, 2, 1, bone)
			px(img, 14, 6, 1, 4, bone)
			px(img, 5, 4, 2, 4, bone)
			px(img, 11, 4, 2, 4, bone)
			px(img, 4, 2, 3, 2, bone)
			px(img, 11, 2, 3, 2, bone)
			px(img, 4, 1, 4, 1, shadow)
			px(img, 11, 1, 4, 1, shadow)
		}
		px(img, 13, 7, 1, 8, rust)
		px(img, 12, 14, 3, 1, rust)
	})
}
