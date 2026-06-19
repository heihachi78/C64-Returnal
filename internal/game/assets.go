package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
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
