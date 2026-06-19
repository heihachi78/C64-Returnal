package game

import (
	"image"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestGeneratedAssetsMatchOriginalTextureDimensions(t *testing.T) {
	assets := NewAssets(64)
	tests := []struct {
		name   string
		width  int
		height int
		got    *ebiten.Image
	}{
		{name: "grass", width: 64, height: 64, got: assets.Grass[0]},
		{name: "mage", width: 16, height: 22, got: assets.Mage[0]},
		{name: "skeleton", width: 16, height: 22, got: assets.Skeleton[0]},
		{name: "fireball", width: 12, height: 12, got: assets.Fireball[0]},
		{name: "lightning", width: 12, height: 12, got: assets.Lightning},
		{name: "orb", width: 12, height: 12, got: assets.Orb[0]},
		{name: "beam", width: 12, height: 12, got: assets.Beam},
		{name: "meteor", width: 12, height: 12, got: assets.Meteor[0]},
		{name: "life", width: 12, height: 12, got: assets.Life},
		{name: "coin", width: 12, height: 12, got: assets.Coin[0]},
		{name: "chest", width: 16, height: 14, got: assets.Chest[ChestBronze]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bounds := tt.got.Bounds()
			if bounds.Dx() != tt.width || bounds.Dy() != tt.height {
				t.Fatalf("bounds = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), tt.width, tt.height)
			}
		})
	}
}

func TestWindowIconsIncludeAppIconSizes(t *testing.T) {
	icons := WindowIcons()
	wantSizes := []int{16, 32, 64, 128, 256, 512, 1024}
	if len(icons) != len(wantSizes) {
		t.Fatalf("icon count = %d, want %d", len(icons), len(wantSizes))
	}
	for i, want := range wantSizes {
		bounds := icons[i].Bounds()
		if bounds.Dx() != want || bounds.Dy() != want {
			t.Fatalf("icon %d bounds = %dx%d, want %dx%d", i, bounds.Dx(), bounds.Dy(), want, want)
		}
	}
}

func TestGeneratedAssetVariantCountsMatchOriginalFactories(t *testing.T) {
	assets := NewAssets(64)
	tests := []struct {
		name string
		got  int
		want int
	}{
		{name: "grass", got: len(assets.Grass), want: 8},
		{name: "mage", got: len(assets.Mage), want: 2},
		{name: "skeleton", got: len(assets.Skeleton), want: 2},
		{name: "fireball", got: len(assets.Fireball), want: 2},
		{name: "orb", got: len(assets.Orb), want: 2},
		{name: "meteor", got: len(assets.Meteor), want: 2},
		{name: "coin", got: len(assets.Coin), want: 4},
		{name: "chest tiers", got: len(assets.Chest), want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("variant count = %d, want %d from PixelArtFactory", tt.got, tt.want)
			}
		})
	}
}

func TestCalibratedColorsRoundLikeOriginalRasterizedTextures(t *testing.T) {
	tests := []struct {
		name string
		got  color.RGBA
		want color.RGBA
	}{
		{name: "hud primary text", got: c64Text, want: color.RGBA{245, 237, 212, 255}},
		{name: "mage face", got: rgb(0.88, 0.65, 0.47), want: color.RGBA{224, 166, 120, 255}},
		{name: "skeleton highlight", got: rgb(0.94, 0.96, 0.86), want: color.RGBA{240, 245, 219, 255}},
		{name: "orb core", got: rgb(0.96, 0.78, 1.0), want: color.RGBA{245, 199, 255, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("color = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}

func TestGeneratedAssetPixelsMatchOriginalPixelArtFactory(t *testing.T) {
	tests := []struct {
		name string
		img  func() *image.RGBA
		x    int
		y    int
		want color.RGBA
	}{
		{name: "mage transparent corner", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeMageTexture(0) }) }, x: 0, y: 0, want: color.RGBA{0, 0, 0, 0}},
		{name: "mage face", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeMageTexture(0) }) }, x: 6, y: 12, want: rgb(0.88, 0.65, 0.47)},
		{name: "skeleton bright skull", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeSkeletonTexture(0) }) }, x: 6, y: 18, want: rgb(0.94, 0.96, 0.86)},
		{name: "fireball yellow cap", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeFireballTexture(0) }) }, x: 3, y: 9, want: rgb(1.00, 0.88, 0.18)},
		{name: "lightning white top", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeLightningTexture() }) }, x: 6, y: 10, want: color.RGBA{255, 255, 255, 255}},
		{name: "orb core", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeOrbTexture(0) }) }, x: 5, y: 4, want: rgb(0.96, 0.78, 1.0)},
		{name: "beam white core", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeBeamTexture() }) }, x: 3, y: 5, want: color.RGBA{255, 255, 255, 255}},
		{name: "meteor highlight", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeMeteorTexture(0) }) }, x: 5, y: 7, want: rgb(0.76, 0.53, 0.28)},
		{name: "life shadow", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeLifeTexture() }) }, x: 4, y: 1, want: rgb(0.45, 0.02, 0.06)},
		{name: "bronze chest lock", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeChestTexture(ChestBronze) }) }, x: 7, y: 5, want: rgb(1.0, 0.86, 0.28)},
		{name: "coin shine", img: func() *image.RGBA { return capturedAssetPixels(func() { _ = makeCoinTexture(0) }) }, x: 4, y: 8, want: rgb(1.0, 0.98, 0.72)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logicalPixel(tt.img(), tt.x, tt.y); got != tt.want {
				t.Fatalf("pixel at (%d,%d) = %#v, want %#v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func capturedAssetPixels(makeAsset func()) *image.RGBA {
	var captured *image.RGBA
	previousCapture := generatedAssetPixelCapture
	generatedAssetPixelCapture = func(img *image.RGBA) {
		captured = img
	}
	defer func() {
		generatedAssetPixelCapture = previousCapture
	}()

	makeAsset()
	return captured
}

func logicalPixel(img *image.RGBA, x, y int) color.RGBA {
	bounds := img.Bounds()
	return color.RGBAModel.Convert(img.At(x, bounds.Dy()-y-1)).(color.RGBA)
}
