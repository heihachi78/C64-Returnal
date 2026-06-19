package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

func (g *Game) drawGrass(screen *ebiten.Image) {
	tile := g.tuning.TileSize
	startColumn, startRow, columns, rows := grassGrid(g.screenW, g.screenH, tile, g.player.Pos)

	for rowOffset := 0; rowOffset < rows; rowOffset++ {
		for columnOffset := 0; columnOffset < columns; columnOffset++ {
			column := startColumn + columnOffset
			row := startRow + rowOffset
			world := Vec2{
				X: float64(column)*tile + tile/2,
				Y: float64(row)*tile + tile/2,
			}
			sx, sy := g.worldToScreen(world)
			img := g.assets.Grass[g.grassHash(column, row, 17)%len(g.assets.Grass)]
			flip := g.grassHash(column, row, 31)%2 == 1
			tint := g.grassTint(column, row)
			g.drawSpriteRotatedBlend(screen, img, sx, sy, tile, tile, 0, flip, tint, grassTintBlendFactor)
		}
	}
}
func grassGrid(screenW, screenH int, tile float64, center Vec2) (startColumn, startRow, columns, rows int) {
	columns = max(6, int(math.Ceil(float64(screenW)/tile))+4)
	rows = max(6, int(math.Ceil(float64(screenH)/tile))+4)
	centerColumn := int(math.Floor(center.X / tile))
	centerRow := int(math.Floor(center.Y / tile))
	startColumn = centerColumn - columns/2
	startRow = centerRow - rows/2
	return startColumn, startRow, columns, rows
}
func (g *Game) drawPlayer(screen *ebiten.Image) {
	x, y := g.worldToScreen(g.player.Pos)
	presentation := playerSpritePresentation(g.player, g.session.GameOver)
	g.drawSpriteRotatedBlend(screen, g.assets.Mage[g.player.AnimFrame%len(g.assets.Mage)], x, y, 32, 44, presentation.Rotation, g.player.Facing < 0, presentation.Tint, presentation.BlendFactor)
}
func (g *Game) drawSkeleton(screen *ebiten.Image, skeleton Skeleton) {
	x, y := g.worldToScreen(skeleton.Pos)
	presentation := skeletonSpritePresentation(skeleton)
	w, h := skeletonSpriteSize(skeleton.Kind)
	g.drawSpriteRotatedBlend(screen, g.assets.Skeleton[skeleton.AnimFrame%len(g.assets.Skeleton)], x, y, w, h, 0, skeleton.Facing < 0, presentation.Tint, presentation.BlendFactor)
}

type spritePresentation struct {
	Tint        color.RGBA
	BlendFactor float64
	Rotation    float64
}

func playerSpritePresentation(player Player, gameOver bool) spritePresentation {
	presentation := spritePresentation{Tint: color.RGBA{255, 255, 255, 255}}
	if player.HitFlash > 0 {
		elapsed := playerHitFlashDuration - player.HitFlash
		presentation.Tint.A = flashActionAlpha(elapsed, playerHitFlashDuration, 0.08, 0.08)
	}
	if gameOver {
		presentation.Tint = color.RGBA{217, 13, 20, 115}
		presentation.BlendFactor = 0.65
		presentation.Rotation = worldRotationToScreen(player.DeathRotation)
	}
	return presentation
}
func skeletonSpritePresentation(skeleton Skeleton) spritePresentation {
	tint, blendFactor := skeletonTintBlend(skeleton.Kind)
	if skeleton.HitFlash > 0 {
		elapsed := skeletonDamageFlashDuration - skeleton.HitFlash
		tint.A = flashActionAlpha(elapsed, skeletonDamageFlashDuration, 0.06, 0.06)
	}
	return spritePresentation{Tint: tint, BlendFactor: blendFactor}
}
func skeletonSpriteSize(kind SkeletonKind) (float64, float64) {
	if kind == SkeletonBlue {
		return 90, 126
	}
	return 30, 42
}
func (g *Game) drawFireball(screen *ebiten.Image, fire Fireball) {
	x, y := g.worldToScreen(fire.Pos)
	angle := worldRotationToScreen(math.Atan2(fire.Velocity.Y, fire.Velocity.X))
	g.drawSpriteRotated(screen, g.assets.Fireball[fire.AnimFrame%len(g.assets.Fireball)], x, y, 18, 18, angle, false, color.RGBA{255, 255, 255, 255})
}
func (g *Game) drawOrb(screen *ebiten.Image, orb OrbitalOrb) {
	x, y := g.worldToScreen(orb.Pos)
	g.drawSprite(screen, g.assets.Orb[orb.AnimFrame%len(g.assets.Orb)], x, y, 20, 20, false, color.RGBA{255, 255, 255, 255})
}
func (g *Game) drawMeteor(screen *ebiten.Image, meteor MeteorProjectile) {
	x, y := g.worldToScreen(meteor.Pos)
	direction := meteor.Impact.Sub(meteor.Start)
	angle := worldRotationToScreen(math.Atan2(direction.Y, direction.X))
	g.drawSpriteRotated(screen, g.assets.Meteor[meteor.AnimFrame%len(g.assets.Meteor)], x, y, 24, 24, angle, false, color.RGBA{255, 255, 255, 255})
}
func (g *Game) drawCoin(screen *ebiten.Image, coin Coin) {
	x, y := g.worldToScreen(coin.Pos)
	y += coinFloatOffset(coin.Phase)
	frame := int(coin.Phase/g.tuning.CoinAnimationFrameTime) % len(g.assets.Coin)
	alpha := coinShimmerAlpha(coin.Phase)
	g.drawSprite(screen, g.assets.Coin[frame], x, y, 28, 28, false, color.RGBA{255, 255, 255, alpha})
}
func (g *Game) drawChest(screen *ebiten.Image, chest Chest) {
	x, y := g.worldToScreen(chest.Pos)
	g.drawSprite(screen, g.assets.Chest[chest.Tier], x, y, 32, 28, false, color.RGBA{255, 255, 255, 255})
}
func (g *Game) drawPickupIndicators(screen *ebiten.Image) {
	for _, chest := range g.chests {
		x, y := g.worldToScreen(chest.Pos)
		if spriteBoundsVisible(g.screenW, g.screenH, x, y, chestSpriteWidth, chestSpriteHeight, 0) {
			continue
		}
		ix, iy := edgeIndicatorPosition(g.screenW, g.screenH, x, y, pickupIndicatorEdgeInset)
		g.drawPickupIndicatorBacking(screen, ix, iy)
		g.drawSpriteScreen(screen, g.assets.Chest[chest.Tier], ix, iy, pickupIndicatorChestWidth, pickupIndicatorChestHeight, false, color.RGBA{255, 255, 255, 255})
	}
	for _, coin := range g.coins {
		x, y := g.worldToScreen(coin.Pos)
		y += coinFloatOffset(coin.Phase)
		if spriteBoundsVisible(g.screenW, g.screenH, x, y, coinSpriteSize, coinSpriteSize, 0) {
			continue
		}
		ix, iy := edgeIndicatorPosition(g.screenW, g.screenH, x, y, pickupIndicatorEdgeInset)
		frame := int(coin.Phase/g.tuning.CoinAnimationFrameTime) % len(g.assets.Coin)
		alpha := coinShimmerAlpha(coin.Phase)
		g.drawPickupIndicatorBacking(screen, ix, iy)
		g.drawSpriteScreen(screen, g.assets.Coin[frame], ix, iy, pickupIndicatorCoinSize, pickupIndicatorCoinSize, false, color.RGBA{255, 255, 255, alpha})
	}
}
func (g *Game) drawPickupIndicatorBacking(screen *ebiten.Image, x, y float64) {
	panelColor := c64Panel
	panelColor.A = 196
	size := float64(pickupIndicatorBackingSize)
	drawFilledRoundedRect(screen, x-size/2, y-size/2, size, size, 5, panelColor)
}
func edgeIndicatorPosition(screenW, screenH int, targetX, targetY, inset float64) (float64, float64) {
	width := float64(screenW)
	height := float64(screenH)
	centerX := width / 2
	centerY := height / 2
	left := math.Min(inset, centerX)
	right := math.Max(left, width-left)
	top := math.Min(inset, centerY)
	bottom := math.Max(top, height-top)
	dx := targetX - centerX
	dy := targetY - centerY
	if dx == 0 && dy == 0 {
		return centerX, centerY
	}

	scale := math.Inf(1)
	if dx > 0 {
		scale = math.Min(scale, (right-centerX)/dx)
	} else if dx < 0 {
		scale = math.Min(scale, (left-centerX)/dx)
	}
	if dy > 0 {
		scale = math.Min(scale, (bottom-centerY)/dy)
	} else if dy < 0 {
		scale = math.Min(scale, (top-centerY)/dy)
	}
	if math.IsInf(scale, 1) || scale < 0 {
		scale = 0
	}
	return centerX + dx*scale, centerY + dy*scale
}
func (g *Game) worldToScreen(pos Vec2) (float64, float64) {
	return float64(g.screenW)/2 + pos.X - g.player.Pos.X, float64(g.screenH)/2 - (pos.Y - g.player.Pos.Y)
}
func worldRotationToScreen(angle float64) float64 {
	return -angle
}
func (g *Game) grassHash(column, row, salt int) int {
	hash := (column * 73856093) ^ (row * 19349663) ^ (salt * 83492791)
	if hash < 0 {
		negated := -hash
		if negated < 0 {
			return 0
		}
		hash = negated
	}
	return hash
}
func (g *Game) grassTint(column, row int) color.RGBA {
	palette := []color.RGBA{
		rgb(0.27, 0.55, 0.21),
		rgb(0.20, 0.47, 0.18),
		rgb(0.34, 0.62, 0.24),
		rgb(0.16, 0.40, 0.17),
		rgb(0.42, 0.67, 0.25),
	}
	return palette[g.grassHash(column, row, 0)%len(palette)]
}
