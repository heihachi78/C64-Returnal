package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Assets struct {
	Grass     []*ebiten.Image
	Mage      []*ebiten.Image
	Skeleton  []*ebiten.Image
	Fireball  []*ebiten.Image
	Lightning *ebiten.Image
	Orb       []*ebiten.Image
	Beam      *ebiten.Image
	Meteor    []*ebiten.Image
	Life      *ebiten.Image
	Coin      []*ebiten.Image
	Chest     map[ChestTier]*ebiten.Image
}

var generatedAssetPixelCapture func(*image.RGBA)

func NewAssets(tileSize int) *Assets {
	return &Assets{
		Grass:     makeGrassTextures(tileSize),
		Mage:      []*ebiten.Image{makeMageTexture(0), makeMageTexture(1)},
		Skeleton:  []*ebiten.Image{makeSkeletonTexture(0), makeSkeletonTexture(1)},
		Fireball:  []*ebiten.Image{makeFireballTexture(0), makeFireballTexture(1)},
		Lightning: makeLightningTexture(),
		Orb:       []*ebiten.Image{makeOrbTexture(0), makeOrbTexture(1)},
		Beam:      makeBeamTexture(),
		Meteor:    []*ebiten.Image{makeMeteorTexture(0), makeMeteorTexture(1)},
		Life:      makeLifeTexture(),
		Coin:      []*ebiten.Image{makeCoinTexture(0), makeCoinTexture(1), makeCoinTexture(2), makeCoinTexture(3)},
		Chest: map[ChestTier]*ebiten.Image{
			ChestBronze: makeChestTexture(ChestBronze),
			ChestSilver: makeChestTexture(ChestSilver),
			ChestGold:   makeChestTexture(ChestGold),
		},
	}
}

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

func fill(img *image.RGBA, clr color.RGBA) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, clr)
		}
	}
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
