package game

import "image/color"

var (
	c64Green     = color.RGBA{46, 94, 41, 255}
	c64Panel     = color.RGBA{5, 5, 5, 158}
	c64Text      = color.RGBA{245, 237, 212, 255}
	c64Gold      = color.RGBA{255, 219, 66, 255}
	c64Orange    = color.RGBA{241, 105, 37, 255}
	c64Blue      = color.RGBA{71, 187, 255, 255}
	c64Purple    = color.RGBA{148, 31, 242, 255}
	c64Red       = color.RGBA{242, 13, 20, 255}
	c64Black     = color.RGBA{5, 5, 5, 255}
	c64White     = color.RGBA{246, 243, 231, 255}
	c64Dim       = color.RGBA{148, 138, 115, 255}
	c64PanelEdge = color.RGBA{0, 0, 0, 0}
)

const (
	statusFontSize         = 16
	combatFontSize         = 14
	levelUpTitleFontSize   = 40
	levelUpOptionFontSize  = 22
	levelUpKeyFontSize     = 18
	chestTitleFontSize     = 34
	chestItemFontSize      = 20
	chestContinueFontSize  = 18
	gameOverTitleFontSize  = 40
	gameOverOptionFontSize = 22

	coinSpriteSize             = 28
	chestSpriteWidth           = 32
	chestSpriteHeight          = 28
	pickupIndicatorBackingSize = 14
	pickupIndicatorCoinSize    = 9
	pickupIndicatorChestWidth  = 10
	pickupIndicatorChestHeight = 9
	pickupIndicatorEdgeInset   = 11

	combatRowFirstOffsetY  = -10
	combatRowSecondOffsetY = 6
	combatRowKillsOffsetY  = 20

	worldLayerGrass        = -20.0
	worldLayerGroundEffect = -19.0
	worldLayerChest        = 8.5
	worldLayerCoin         = 8.75
	worldLayerSkeleton     = 9.0
	worldLayerLightningHit = 9.5
	worldLayerPlayer       = 10.0
	worldLayerOrb          = 11.0
	worldLayerFireball     = 12.0
	worldLayerEffect       = 13.0
	worldLayerMeteor       = 14.0
)
