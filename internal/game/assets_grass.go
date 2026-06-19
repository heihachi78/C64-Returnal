package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

func makeGrassTextures(size int) []*ebiten.Image {
	textures := make([]*ebiten.Image, 8)
	for variant := range textures {
		textures[variant] = makePixelImage(size, size, func(img *image.RGBA) {
			fill(img, grassColor(0.22, 0.47, 0.17, variant, 0))
			drawGrassGroundFlecks(img, size, variant)
			drawGrassDetails(img, size, variant)
		})
	}
	return textures
}
func drawGrassGroundFlecks(img *image.RGBA, size, variant int) {
	flecks := [][5]int{
		{4, 7, 3, 1, 0}, {13, 23, 2, 1, 1}, {27, 9, 4, 1, 2}, {39, 18, 2, 1, 3},
		{52, 5, 3, 1, 4}, {7, 39, 2, 1, 5}, {20, 53, 3, 1, 6}, {34, 41, 4, 1, 7},
		{49, 55, 2, 1, 8}, {58, 29, 3, 1, 9}, {2, 58, 2, 1, 10}, {45, 34, 3, 1, 11},
	}
	for index, fleck := range flecks {
		clr := grassColor(0.30, 0.56, 0.20, variant, index)
		if fleck[4]%2 != 0 {
			clr = grassColor(0.16, 0.36, 0.14, variant, index)
		}
		px(img, wrappedPixel(fleck[0]+grassOffset(variant, index, 0), size), wrappedPixel(fleck[1]+grassOffset(variant, index, 1), size), fleck[2], fleck[3], clr)
	}
}
func drawGrassDetails(img *image.RGBA, size, variant int) {
	dark := grassColor(0.10, 0.29, 0.11, variant, 17)
	mid := grassColor(0.28, 0.61, 0.20, variant, 19)
	light := grassColor(0.50, 0.82, 0.31, variant, 23)
	root := grassColor(0.08, 0.24, 0.10, variant, 29)
	tufts := [][3]int{{7, 9, 8}, {22, 5, 11}, {43, 8, 9}, {56, 17, 10}, {13, 29, 12}, {33, 25, 8}, {49, 36, 13}, {5, 49, 9}, {25, 51, 10}, {39, 55, 7}, {59, 53, 11}}
	for index, tuft := range tufts {
		baseX := wrappedPixel(tuft[0]+grassOffset(variant, index, 4), size)
		baseY := wrappedPixel(tuft[1]+grassOffset(variant, index, 5), size)
		height := max(5, tuft[2]+grassOffset(variant, index, 6)/2)
		drawGrassTuft(img, baseX, baseY, height, dark, mid, light, root)
	}
	blades := [][4]int{{4, 21, 6, 1}, {17, 42, 5, -1}, {30, 15, 7, 0}, {41, 44, 6, 1}, {54, 28, 5, -1}, {61, 6, 6, 0}, {9, 60, 4, 1}, {31, 36, 5, -1}}
	for index, blade := range blades {
		x := wrappedPixel(blade[0]+grassOffset(variant, index, 7), size)
		y := wrappedPixel(blade[1]+grassOffset(variant, index, 8), size)
		height := max(3, blade[2]+grassOffset(variant, index, 9)/2)
		drawGrassBlade(img, x, y, height, blade[3], light, 1)
	}
}
func drawGrassTuft(img *image.RGBA, baseX, baseY, height int, dark, mid, light, root color.RGBA) {
	px(img, baseX-2, baseY, 6, 2, root)
	px(img, baseX-1, baseY+1, 4, 1, dark)
	drawGrassBlade(img, baseX-2, baseY+1, height-2, -1, dark, 1)
	drawGrassBlade(img, baseX, baseY+1, height, 0, mid, 2)
	drawGrassBlade(img, baseX+2, baseY+1, height-1, 1, mid, 1)
	drawGrassBlade(img, baseX+1, baseY+2, max(3, height-4), -1, light, 1)
	px(img, baseX, baseY+height+1, 1, 1, light)
}
func drawGrassBlade(img *image.RGBA, baseX, baseY, height, lean int, clr color.RGBA, baseWidth int) {
	for step := 0; step < height; step++ {
		leanOffset := 0
		if lean != 0 {
			leanOffset = step * lean / 3
		}
		width := 1
		if step < 2 {
			width = baseWidth
		}
		px(img, baseX+leanOffset, baseY+step, width, 1, clr)
	}
}
func grassColor(red, green, blue float64, variant, shift int) color.RGBA {
	amount := float64(((variant*17+shift*11)%9)-4) * 0.012
	return rgb(red+amount*0.7, green+amount, blue+amount*0.45)
}
func grassOffset(variant, index, axis int) int {
	return (variant*37+index*19+axis*13)%11 - 5
}
func wrappedPixel(value, size int) int {
	return ((value % size) + size) % size
}
