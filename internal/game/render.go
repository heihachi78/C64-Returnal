package game

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

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

	combatRowFirstOffsetY  = -10
	combatRowSecondOffsetY = 6
	combatRowKillsOffsetY  = 20

	worldLayerGrass        = -20.0
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

var (
	hudFont       *opentype.Font
	hudFontSource string
	hudFontFaces  = map[int]font.Face{}
)

func init() {
	hudFont, hudFontSource = loadHUDFont()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(c64Green)
	g.drawGrass(screen)
	g.drawWorld(screen)
	g.drawHUD(screen)

	if g.session.LevelUpChoiceActive {
		g.drawLevelUpOverlay(screen)
	}
	if g.session.ChestRewardActive {
		g.drawChestOverlay(screen)
	}
	if g.session.GameOver {
		g.drawGameOver(screen)
	}
}

func (g *Game) drawWorld(screen *ebiten.Image) {
	for _, chest := range g.chests {
		g.drawChest(screen, chest)
	}
	for _, coin := range g.coins {
		g.drawCoin(screen, coin)
	}
	for _, skeleton := range g.skeleton {
		g.drawSkeleton(screen, skeleton)
	}
	for _, effect := range g.effects {
		if effect.Kind == EffectLightningHit {
			g.drawEffect(screen, effect)
		}
	}
	g.drawPlayer(screen)
	for _, orb := range g.orbs {
		if orb.Active {
			g.drawOrb(screen, orb)
		}
	}
	for _, fire := range g.fireball {
		g.drawFireball(screen, fire)
	}
	for _, effect := range g.effects {
		if effect.Kind != EffectLightningHit {
			g.drawEffect(screen, effect)
		}
	}
	for _, meteor := range g.meteors {
		g.drawMeteor(screen, meteor)
	}
}

func worldRenderLayerOrder() []float64 {
	return []float64{
		worldLayerGrass,
		worldLayerChest,
		worldLayerCoin,
		worldLayerSkeleton,
		worldLayerLightningHit,
		worldLayerPlayer,
		worldLayerOrb,
		worldLayerFireball,
		worldLayerEffect,
		worldLayerMeteor,
	}
}

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
	g.drawSpriteRotatedBlend(screen, g.assets.Skeleton[skeleton.AnimFrame%len(g.assets.Skeleton)], x, y, 30, 42, 0, skeleton.Facing < 0, presentation.Tint, presentation.BlendFactor)
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

func (g *Game) drawEffect(screen *ebiten.Image, effect Effect) {
	alpha := effectFadeAlpha(effect.TTL, effect.MaxTTL)
	switch effect.Kind {
	case EffectLightning:
		innerPoints := effect.InnerPoints
		if len(innerPoints) == 0 {
			innerPoints = effect.Points
		}
		g.drawBolt(screen, effect.Points, g.tuning.LightningBranchWidth+5, color.RGBA{33, 163, 255, alpha / 2})
		g.drawBolt(screen, effect.Points, g.tuning.LightningBranchWidth, color.RGBA{33, 163, 255, alpha})
		g.drawBolt(screen, innerPoints, 2, color.RGBA{255, 255, 255, alpha / 2})
		g.drawBolt(screen, innerPoints, 1, color.RGBA{255, 255, 255, alpha})
		endX, endY := g.worldToScreen(effect.End)
		g.drawSpriteScreen(screen, g.assets.Lightning, endX, endY, 24, 24, false, color.RGBA{255, 255, 255, alpha})
	case EffectLightningHit:
		x, y := g.worldToScreen(effect.Pos)
		hitAlpha := lightningHitEffectAlpha(effect.TTL, effect.MaxTTL)
		g.drawSpriteRotatedBlend(
			screen,
			g.assets.Skeleton[effect.Frame%len(g.assets.Skeleton)],
			x,
			y,
			30,
			42,
			0,
			effect.Facing < 0,
			color.RGBA{89, 219, 255, hitAlpha},
			0.8,
		)
	case EffectBeam:
		startX, startY := g.worldToScreen(effect.Start)
		endX, endY := g.worldToScreen(effect.End)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 16, color.RGBA{255, 184, 20, alpha / 4}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 9, color.RGBA{255, 184, 20, alpha}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 7, color.RGBA{255, 240, 56, alpha / 3}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 4, color.RGBA{255, 240, 56, alpha}, false)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 1, color.RGBA{255, 255, 255, alpha}, false)
	case EffectMeteorImpact:
		x, y := g.worldToScreen(effect.Pos)
		style := meteorImpactRenderStyle(effect)
		vector.DrawFilledCircle(screen, float32(x), float32(y), float32(style.Radius), color.RGBA{74, 43, 20, scaleAlpha(115, style.Alpha)}, false)
		vector.DrawFilledCircle(screen, float32(x), float32(y), float32(style.CoreRadius), color.RGBA{156, 97, 43, scaleAlpha(166, style.Alpha)}, false)
		vector.StrokeCircle(screen, float32(x), float32(y), float32(style.Radius), float32(style.GlowWidth), color.RGBA{194, 128, 61, scaleAlpha(64, style.Alpha)}, false)
		vector.StrokeCircle(screen, float32(x), float32(y), float32(style.Radius), float32(style.StrokeWidth), color.RGBA{194, 128, 61, scaleAlpha(217, style.Alpha)}, false)
	}
}

type meteorImpactStyle struct {
	Scale float64
	Alpha float64
}

type meteorImpactRenderMetrics struct {
	Radius      float64
	CoreRadius  float64
	GlowWidth   float64
	StrokeWidth float64
	Alpha       float64
}

func meteorImpactPresentation(effect Effect) meteorImpactStyle {
	if effect.MaxTTL <= 0 {
		return meteorImpactStyle{Scale: 1.25, Alpha: 0}
	}
	age := Clamp(effect.MaxTTL-effect.TTL, 0, effect.MaxTTL)
	switch {
	case age < 0.08:
		return meteorImpactStyle{Scale: 0.25 + 0.75*(age/0.08), Alpha: 1}
	case age < 0.16:
		return meteorImpactStyle{Scale: 1, Alpha: 1}
	default:
		fade := Clamp((age-0.16)/0.16, 0, 1)
		return meteorImpactStyle{Scale: 1 + 0.25*fade, Alpha: 1 - fade}
	}
}

func meteorImpactRenderStyle(effect Effect) meteorImpactRenderMetrics {
	presentation := meteorImpactPresentation(effect)
	radius := effect.Radius * presentation.Scale
	return meteorImpactRenderMetrics{
		Radius:      radius,
		CoreRadius:  radius * 0.35,
		GlowWidth:   3 * presentation.Scale,
		StrokeWidth: 2 * presentation.Scale,
		Alpha:       presentation.Alpha,
	}
}

func scaleAlpha(base uint8, alpha float64) uint8 {
	return uint8(math.Round(float64(base) * Clamp(alpha, 0, 1)))
}

func flashActionAlpha(elapsed, total, fadeDown, fadeUp float64) uint8 {
	if elapsed < 0 || elapsed >= total || fadeDown <= 0 || fadeUp <= 0 {
		return 255
	}
	cycle := fadeDown + fadeUp
	phase := math.Mod(elapsed, cycle)
	const minAlpha = 0.35
	if phase < fadeDown {
		progress := phase / fadeDown
		return uint8(math.Round(255 * (1 - (1-minAlpha)*progress)))
	}
	progress := (phase - fadeDown) / fadeUp
	return uint8(math.Round(255 * (minAlpha + (1-minAlpha)*progress)))
}

func coinFloatOffset(phase float64) float64 {
	return 5 * linearPingPong(phase, 0.42)
}

func coinShimmerAlpha(phase float64) uint8 {
	alpha := 1 - 0.28*linearPingPong(phase, 0.28)
	return uint8(math.Round(255 * alpha))
}

func linearPingPong(phase, halfPeriod float64) float64 {
	if halfPeriod <= 0 {
		return 0
	}
	t := math.Mod(phase, halfPeriod*2)
	if t < 0 {
		t += halfPeriod * 2
	}
	if t <= halfPeriod {
		return t / halfPeriod
	}
	return 1 - (t-halfPeriod)/halfPeriod
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	left := 18.0
	top := 24.0
	topX, topY, topW, topH := topStatusPanelRect(g.session.PlayerLives)
	g.panel(screen, topX, topY, topW, topH)
	g.drawTextSize(screen, fmt.Sprintf("LV %d", g.session.Progression.Level), left, top, statusFontSize, c64Text)
	g.drawTextSize(screen, fmt.Sprintf("XP %d/%d", g.session.Progression.Experience, g.session.Progression.NextExperience), left, top+24, statusFontSize, c64Text)
	g.drawSpriteScreen(screen, g.assets.Coin[0], left+8, top+48, 16, 16, false, color.RGBA{255, 255, 255, 255})
	g.drawTextSize(screen, fmt.Sprintf("%d", max(0, g.session.CollectedCoins)), left+24, top+48, statusFontSize, c64Text)
	for i := 0; i < g.session.PlayerLives; i++ {
		x, y := lifeIconScreenPosition(i)
		g.drawSpriteScreen(screen, g.assets.Life, x, y, 14, 14, false, color.RGBA{255, 255, 255, 255})
	}

	bottomX, bottomY, bottomW, bottomH := combatStatusPanelRect(g.screenH)
	g.panel(screen, bottomX, bottomY, bottomW, bottomH)
	bottomPanelTop := bottomY
	g.drawCombatRow(screen, g.assets.Fireball[0], left+9, bottomPanelTop+295, fmt.Sprintf("x%d", g.session.Progression.SimultaneousFireball), fmt.Sprintf("%ss", formattedSeconds(g.session.Progression.FireballCastInterval())), fmt.Sprintf("KILLS %d", g.session.Kills.Fireball), true)
	g.drawCombatRow(screen, g.assets.Lightning, left+9, bottomPanelTop+241, fmt.Sprintf("x%d", g.session.Progression.LightningStrikeCount()), fmt.Sprintf("%ss", formattedSeconds(g.session.Progression.LightningCastInterval())), fmt.Sprintf("KILLS %d", g.session.Kills.Lightning), g.session.Progression.LightningUnlocked)
	g.drawCombatRow(screen, g.assets.Orb[0], left+9, bottomPanelTop+187, fmt.Sprintf("x%d", g.session.Progression.OrbitalOrbCount()), fmt.Sprintf("%.1fr/s", g.session.Progression.OrbitalAngularSpeed()), fmt.Sprintf("KILLS %d", g.session.Kills.OrbitalOrb), g.session.Progression.OrbitalOrbUnlocked)
	g.drawCombatRow(screen, g.assets.Beam, left+9, bottomPanelTop+133, fmt.Sprintf("x%d", g.session.Progression.BeamKillCount()), fmt.Sprintf("%ss", formattedSeconds(g.session.Progression.BeamCastInterval())), fmt.Sprintf("KILLS %d", g.session.Kills.Beam), g.session.Progression.BeamUnlocked)
	g.drawCombatRow(screen, g.assets.Meteor[0], left+9, bottomPanelTop+79, fmt.Sprintf("x%d", g.session.Progression.MeteorCount()), fmt.Sprintf("%ss", formattedSeconds(g.session.Progression.MeteorCastInterval())), fmt.Sprintf("KILLS %d", g.session.Kills.Meteor), g.session.Progression.MeteorUnlocked)
	g.drawSpriteScreen(screen, g.assets.Skeleton[0], left+9, bottomPanelTop+25, 16, 22, false, color.RGBA{255, 255, 255, 255})
	g.drawTextSize(screen, fmt.Sprintf("x%d", len(g.skeleton)), left+28, bottomPanelTop+15, combatFontSize, c64Text)
	g.drawTextSize(screen, fmt.Sprintf("%ss", formattedSeconds(g.session.Progression.SkeletonSpawnInterval())), left+28, bottomPanelTop+31, combatFontSize, c64Text)
}

func topStatusPanelRect(lives int) (x, y, w, h float64) {
	lifeRows := max(1, (max(1, lives)+11)/12)
	return 8, 9, 210, 104 + float64(lifeRows-1)*16
}

func lifeIconScreenPosition(index int) (x, y float64) {
	column := index % 12
	row := index / 12
	return 25 + float64(column)*16, 98 + float64(row)*16
}

func combatStatusPanelRect(screenH int) (x, y, w, h float64) {
	return 8, float64(screenH) - 333, 176, 330
}

func (g *Game) drawLevelUpOverlay(screen *ebiten.Image) {
	w := math.Min(math.Max(360, float64(g.screenW)-48), 620)
	h := 250.0 + math.Max(0, float64(max(2, len(g.session.ActiveLevelUpOptions))-2))*52
	x := float64(g.screenW)/2 - w/2
	y := float64(g.screenH)/2 - 90
	panelAlpha := g.modalPanelAlpha(g.session.LevelUpOverlayTimer)
	contentAlpha := g.modalContentAlpha(g.session.LevelUpOverlayTimer)
	optionAlpha := levelUpOptionContentAlpha(contentAlpha, g.session.LevelUpOptionFadeTimer)
	redrawAlpha := levelUpRedrawContentAlpha(g.session.LevelUpRedrawFadeTimer)
	g.panelWithAlpha(screen, x, y, w, h, panelAlpha)
	g.drawCenteredTextSizeScaled(screen, fmt.Sprintf("LEVEL %d", g.session.CurrentLevelUpPresentation), x+w/2, y+34, levelUpTitleFontSize, modalTitleScale(g.session.LevelUpTitleScaleTimer), withAlpha(c64Text, contentAlpha))
	keys := []string{"[Q]", "[A]", "[C]", "[X]"}
	for i, option := range g.session.ActiveLevelUpOptions {
		oy := y + 94 + float64(i)*52
		g.drawCenteredTextSize(screen, keys[i], x+w/2-150, oy, levelUpKeyFontSize, withAlpha(c64Gold, optionAlpha))
		g.drawOptionIconTinted(screen, option, x+w/2-116, oy, color.RGBA{255, 255, 255, optionAlpha})
		g.drawTextSize(screen, option.Title(g.session.Progression.BeamKillUpgradeBonus()), x+w/2-78, oy, levelUpOptionFontSize, withAlpha(c64Text, optionAlpha))
	}
	redrawY := y + 94 + float64(max(2, len(g.session.ActiveLevelUpOptions)))*52 + 14
	redrawColor := c64Gold
	canRedraw := g.session.CollectedCoins >= g.levelUpRedrawCost()
	if !canRedraw {
		redrawColor = c64Dim
	}
	g.drawCenteredTextSize(screen, "[R]", x+w/2-150, redrawY, levelUpKeyFontSize, withAlpha(redrawColor, redrawAlpha))
	g.drawSpriteScreen(screen, g.assets.Coin[0], x+w/2-116, redrawY, 24, 24, false, color.RGBA{255, 255, 255, redrawCoinAlpha(canRedraw, g.session.LevelUpRedrawCoinFadeTimer)})
	g.drawTextSizeScaled(screen, fmt.Sprintf("REDRAW %d  COINS %d", g.levelUpRedrawCost(), max(0, g.session.CollectedCoins)), x+w/2-78, redrawY, levelUpOptionFontSize, redrawPulseScale(g.session.LevelUpRedrawStatusTimer), withAlpha(redrawColor, redrawAlpha))
}

func (g *Game) drawChestOverlay(screen *ebiten.Image) {
	itemCount := len(g.session.ActiveChestRewardItems)
	x, y, w, h := chestOverlayPanelRect(g.screenW, g.screenH, itemCount)
	panelAlpha := g.modalPanelAlpha(g.session.ChestRewardOverlayTimer)
	contentAlpha := g.modalContentAlpha(g.session.ChestRewardOverlayTimer)
	g.panelWithAlpha(screen, x, y, w, h, panelAlpha)
	g.drawCenteredTextSizeScaled(screen, chestRewardTitle(g.session.ActiveChestTier), x+w/2, chestOverlayTitleY(y), chestTitleFontSize, modalTitleScale(g.session.ChestRewardOverlayTimer), withAlpha(c64Text, contentAlpha))
	for i, item := range g.session.ActiveChestRewardItems {
		itemY := chestRewardItemY(y, h, itemCount, i)
		g.drawOptionIconTinted(screen, item.Option, x+w/2-104, itemY, color.RGBA{255, 255, 255, contentAlpha})
		g.drawTextSize(screen, item.Title, x+w/2-66, itemY, chestItemFontSize, withAlpha(c64Text, contentAlpha))
	}
	g.drawCenteredTextSize(screen, "[Q] CONTINUE", x+w/2, chestOverlayContinueY(y, h), chestContinueFontSize, withAlpha(c64Gold, contentAlpha))
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	x, y, w, h := gameOverPanelRect(g.screenW, g.screenH)
	panelAlpha := g.modalPanelAlpha(g.session.GameOverOverlayTimer)
	contentAlpha := g.modalContentAlpha(g.session.GameOverOverlayTimer)
	g.panelWithAlpha(screen, x, y, w, h, panelAlpha)
	g.drawCenteredTextSizeScaled(screen, fmt.Sprintf("YOU DIED AT LEVEL %d", g.session.Progression.Level), x+w/2, y+gameOverTitleOffsetY, gameOverTitleFontSize, modalTitleScale(g.session.GameOverOverlayTimer), withAlpha(c64Red, contentAlpha))
	g.drawCenteredTextSize(screen, "RESTART", x+w/2, y+gameOverRestartOffsetY, gameOverOptionFontSize, withAlpha(c64Red, contentAlpha))
	g.drawCenteredTextSize(screen, "EXIT", x+w/2, y+gameOverExitOffsetY, gameOverOptionFontSize, withAlpha(c64Red, contentAlpha))
}

func (g *Game) worldToScreen(pos Vec2) (float64, float64) {
	return float64(g.screenW)/2 + pos.X - g.player.Pos.X, float64(g.screenH)/2 - (pos.Y - g.player.Pos.Y)
}

func worldRotationToScreen(angle float64) float64 {
	return -angle
}

func chestOverlayPanelRect(screenW, screenH, itemCount int) (x, y, w, h float64) {
	w = math.Min(math.Max(360, float64(screenW)-48), 620)
	h = math.Min(float64(screenH)-64, 174+float64(max(1, itemCount))*34)
	x = float64(screenW)/2 - w/2
	y = float64(screenH)/2 - h/2
	return x, y, w, h
}

func chestRewardTitle(tier ChestTier) string {
	return tier.Title() + " CHEST"
}

func chestOverlayTitleY(panelY float64) float64 {
	return panelY + 54
}

func chestOverlayContinueY(panelY, panelHeight float64) float64 {
	return panelY + panelHeight - 36
}

func chestRewardItemY(panelY, panelHeight float64, itemCount, index int) float64 {
	totalHeight := float64(max(0, itemCount-1)) * 30
	return panelY + panelHeight/2 - totalHeight/2 + float64(index)*30 + 12
}

func (g *Game) drawSprite(screen, img *ebiten.Image, x, y, w, h float64, flipX bool, tint color.RGBA) {
	g.drawSpriteScreen(screen, img, x, y, w, h, flipX, tint)
}

func (g *Game) drawSpriteScreen(screen, img *ebiten.Image, x, y, w, h float64, flipX bool, tint color.RGBA) {
	g.drawSpriteRotated(screen, img, x, y, w, h, 0, flipX, tint)
}

func (g *Game) drawSpriteRotated(screen, img *ebiten.Image, x, y, w, h, rotation float64, flipX bool, tint color.RGBA) {
	g.drawSpriteRotatedBlend(screen, img, x, y, w, h, rotation, flipX, tint, 0)
}

func (g *Game) drawSpriteRotatedBlend(screen, img *ebiten.Image, x, y, w, h, rotation float64, flipX bool, tint color.RGBA, blendFactor float64) {
	if !spriteBoundsVisible(g.screenW, g.screenH, x, y, w, h, rotation) {
		return
	}
	bounds := img.Bounds()
	scaleX := w / float64(bounds.Dx())
	scaleY := h / float64(bounds.Dy())
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
	if flipX {
		op.GeoM.Scale(-scaleX, scaleY)
	} else {
		op.GeoM.Scale(scaleX, scaleY)
	}
	if rotation != 0 {
		op.GeoM.Rotate(rotation)
	}
	if blendFactor > 0 {
		op.ColorM.Scale(1-blendFactor, 1-blendFactor, 1-blendFactor, 1)
		op.ColorM.Translate(float64(tint.R)/255*blendFactor, float64(tint.G)/255*blendFactor, float64(tint.B)/255*blendFactor, 0)
		if tint.A != 255 {
			op.ColorScale.ScaleAlpha(float32(tint.A) / 255)
		}
	} else if tint != (color.RGBA{255, 255, 255, 255}) {
		op.ColorScale.ScaleWithColor(tint)
	}
	op.GeoM.Translate(x, y)
	screen.DrawImage(img, op)
}

func spriteBoundsVisible(screenW, screenH int, x, y, w, h, rotation float64) bool {
	if screenW <= 0 || screenH <= 0 || w <= 0 || h <= 0 {
		return false
	}
	halfW := w / 2
	halfH := h / 2
	if rotation != 0 {
		radius := math.Hypot(w, h) / 2
		halfW = radius
		halfH = radius
	}
	return x+halfW >= 0 && x-halfW <= float64(screenW) &&
		y+halfH >= 0 && y-halfH <= float64(screenH)
}

func skeletonTintBlend(kind SkeletonKind) (color.RGBA, float64) {
	switch kind {
	case SkeletonRed:
		return color.RGBA{242, 13, 10, 255}, 0.68
	case SkeletonPurple:
		return color.RGBA{148, 31, 242, 255}, 0.72
	case SkeletonBlack:
		return color.RGBA{5, 5, 5, 255}, 0.86
	default:
		return color.RGBA{255, 255, 255, 255}, 0
	}
}

func (g *Game) drawBolt(screen *ebiten.Image, points []Vec2, width float32, clr color.RGBA) {
	if len(points) < 2 {
		return
	}
	for i := 1; i < len(points); i++ {
		startX, startY := g.worldToScreen(points[i-1])
		endX, endY := g.worldToScreen(points[i])
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), width, clr, false)
	}
}

func (g *Game) drawCombatRow(screen, icon *ebiten.Image, x, y float64, first, second, kills string, unlocked bool) {
	row := combatRowPresentation(first, second, kills, unlocked)
	g.drawSpriteScreen(screen, icon, x, y, 18, 18, false, row.Tint)
	g.drawTextSize(screen, row.First, x+19, y+combatRowFirstOffsetY, combatFontSize, row.TextColor)
	g.drawTextSize(screen, row.Second, x+19, y+combatRowSecondOffsetY, combatFontSize, row.TextColor)
	if kills != "" {
		g.drawTextSize(screen, row.Kills, x+19, y+combatRowKillsOffsetY, combatFontSize, row.TextColor)
	}
}

type combatRowStyle struct {
	First     string
	Second    string
	Kills     string
	Tint      color.RGBA
	TextColor color.RGBA
}

func combatRowPresentation(first, second, kills string, unlocked bool) combatRowStyle {
	first, second, kills = combatRowLabels(first, second, kills, unlocked)
	return combatRowStyle{
		First:     first,
		Second:    second,
		Kills:     kills,
		Tint:      color.RGBA{255, 255, 255, 255},
		TextColor: c64Text,
	}
}

func combatRowLabels(first, second, kills string, unlocked bool) (string, string, string) {
	if !unlocked {
		return "LOCKED", "--", kills
	}
	return first, second, kills
}

func (g *Game) drawOptionIcon(screen *ebiten.Image, option LevelUpOption, x, y float64) {
	g.drawOptionIconTinted(screen, option, x, y, color.RGBA{255, 255, 255, 255})
}

func (g *Game) drawOptionIconTinted(screen *ebiten.Image, option LevelUpOption, x, y float64, tint color.RGBA) {
	var img *ebiten.Image
	switch option {
	case FireRate, ExtraFireball:
		img = g.assets.Fireball[0]
	case ExtraLife:
		img = g.assets.Life
	case HalveSkeletons:
		img = g.assets.Skeleton[0]
	case LearnLightning, LightningBounce, LightningRate:
		img = g.assets.Lightning
	case LearnOrb, ExtraOrb, OrbitalSpeed:
		img = g.assets.Orb[0]
	case LearnBeam, BeamRate, BeamKillCount:
		img = g.assets.Beam
	case LearnMeteor, ExtraMeteor, MeteorRate:
		img = g.assets.Meteor[0]
	default:
		img = g.assets.Fireball[0]
	}
	g.drawSpriteScreen(screen, img, x, y, 24, 24, false, tint)
}

func (g *Game) panel(screen *ebiten.Image, x, y, w, h float64) {
	g.panelWithAlpha(screen, x, y, w, h, c64Panel.A)
}

func (g *Game) panelWithAlpha(screen *ebiten.Image, x, y, w, h float64, alpha uint8) {
	panelColor := c64Panel
	panelColor.A = alpha
	drawFilledRoundedRect(screen, x, y, w, h, panelCornerRadius, panelColor)
	if c64PanelEdge.A > 0 {
		vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1, c64PanelEdge, false)
	}
}

func drawFilledRoundedRect(screen *ebiten.Image, x, y, w, h, radius float64, clr color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	r := math.Min(radius, math.Min(w/2, h/2))
	if r <= 0 {
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), clr, false)
		return
	}
	vector.DrawFilledRect(screen, float32(x+r), float32(y), float32(w-2*r), float32(h), clr, false)
	vector.DrawFilledRect(screen, float32(x), float32(y+r), float32(r), float32(h-2*r), clr, false)
	vector.DrawFilledRect(screen, float32(x+w-r), float32(y+r), float32(r), float32(h-2*r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+r), float32(y+r), float32(r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+w-r), float32(y+r), float32(r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+r), float32(y+h-r), float32(r), clr, false)
	vector.DrawFilledCircle(screen, float32(x+w-r), float32(y+h-r), float32(r), clr, false)
}

func (g *Game) modalPanelAlpha(timer float64) uint8 {
	progress := Clamp(timer/modalFadeDuration, 0, 1)
	return uint8(math.Round(float64(c64Panel.A) * progress))
}

func (g *Game) modalContentAlpha(timer float64) uint8 {
	progress := Clamp(timer/modalFadeDuration, 0, 1)
	return uint8(math.Round(255 * progress))
}

func effectFadeAlpha(ttl, maxTTL float64) uint8 {
	if maxTTL <= 0 {
		return 0
	}
	return uint8(math.Round(255 * Clamp(ttl/maxTTL, 0, 1)))
}

func lightningHitEffectAlpha(ttl, maxTTL float64) uint8 {
	if maxTTL <= 0 {
		return 0
	}
	return uint8(math.Round(255 * 0.85 * Clamp(ttl/maxTTL, 0, 1)))
}

func modalTitleScale(timer float64) float64 {
	return 0.75 + 0.25*Clamp(timer/modalFadeDuration, 0, 1)
}

func levelUpOptionContentAlpha(contentAlpha uint8, timer float64) uint8 {
	return min(contentAlpha, uint8(math.Round(255*Clamp(timer/modalFadeDuration, 0, 1))))
}

func levelUpRedrawContentAlpha(timer float64) uint8 {
	return uint8(math.Round(255 * Clamp(timer/redrawStatusFadeDuration, 0, 1)))
}

func redrawPulseScale(timer float64) float64 {
	if timer <= 0 {
		return 1
	}
	elapsed := Clamp(redrawFailurePulseDuration-timer, 0, redrawFailurePulseDuration)
	if elapsed < 0.08 {
		return 1 + 0.08*(elapsed/0.08)
	}
	return 1.08 - 0.08*Clamp((elapsed-0.08)/0.08, 0, 1)
}

func redrawCoinAlpha(canRedraw bool, timer float64) uint8 {
	if canRedraw {
		return 255
	}
	const dimmedCoinAlpha = 0.45
	progress := Clamp(timer/redrawStatusFadeDuration, 0, 1)
	return uint8(math.Round(255 * (dimmedCoinAlpha + (1-dimmedCoinAlpha)*progress)))
}

func withAlpha(clr color.RGBA, alpha uint8) color.RGBA {
	clr.A = alpha
	return clr
}

func loadHUDFont() (*opentype.Font, string) {
	if font, name := loadSystemFontByFullName(
		[]string{
			"/System/Library/Fonts/Menlo.ttc",
			"/Library/Fonts/Menlo.ttc",
		},
		[]string{"menlo", "bold"},
	); font != nil {
		return font, name
	}
	font, _ := opentype.Parse(gomonobold.TTF)
	return font, "Go Mono Bold"
}

func loadSystemFontByFullName(paths, required []string) (*opentype.Font, string) {
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		collection, err := opentype.ParseCollection(data)
		if err != nil {
			continue
		}
		for i := 0; i < collection.NumFonts(); i++ {
			font, err := collection.Font(i)
			if err != nil {
				continue
			}
			name, err := font.Name(nil, sfnt.NameIDFull)
			if err != nil || !fontNameMatches(name, required) {
				continue
			}
			return font, name
		}
	}
	return nil, ""
}

func fontNameMatches(name string, required []string) bool {
	lower := strings.ToLower(name)
	for _, token := range required {
		if !strings.Contains(lower, strings.ToLower(token)) {
			return false
		}
	}
	return true
}

func fontFaceForSize(size float64) font.Face {
	if hudFont == nil || size <= 0 {
		return basicfont.Face7x13
	}
	key := fontSizeKey(size)
	if face := hudFontFaces[key]; face != nil {
		return face
	}
	face, err := opentype.NewFace(hudFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return basicfont.Face7x13
	}
	hudFontFaces[key] = face
	return face
}

func centeredTextBaseline(face font.Face, centerY float64) int {
	metrics := face.Metrics()
	return int(math.Round(centerY + float64((metrics.Ascent-metrics.Descent).Round())/2))
}

func (g *Game) drawTextSize(screen *ebiten.Image, s string, x, centerY, size float64, clr color.Color) {
	face := fontFaceForSize(size)
	text.Draw(screen, s, face, int(math.Round(x)), centeredTextBaseline(face, centerY), clr)
}

func (g *Game) drawTextSizeScaled(screen *ebiten.Image, s string, x, centerY, size, scale float64, clr color.Color) {
	g.drawScaledTextImage(screen, s, x, centerY, size, scale, false, clr)
}

func (g *Game) drawCenteredTextSize(screen *ebiten.Image, s string, x, centerY, size float64, clr color.Color) {
	face := fontFaceForSize(size)
	width := font.MeasureString(face, s).Ceil()
	text.Draw(screen, s, face, int(math.Round(x))-width/2, centeredTextBaseline(face, centerY), clr)
}

func (g *Game) drawCenteredTextSizeScaled(screen *ebiten.Image, s string, x, centerY, size, scale float64, clr color.Color) {
	g.drawScaledTextImage(screen, s, x, centerY, size, scale, true, clr)
}

type scaledTextLayout struct {
	Width    int
	Height   int
	AnchorX  float64
	AnchorY  float64
	Baseline int
}

type scaledTextCacheKey struct {
	Text     string
	SizeKey  int
	Centered bool
}

type scaledTextCacheEntry struct {
	Image  *ebiten.Image
	Layout scaledTextLayout
}

func fontSizeKey(size float64) int {
	return int(math.Round(size * 10))
}

func baseScaledTextLayout(face font.Face, s string, centered bool) scaledTextLayout {
	const padding = 4
	metrics := face.Metrics()
	textWidth := max(1, font.MeasureString(face, s).Ceil())
	textHeight := max(1, (metrics.Ascent + metrics.Descent).Ceil())
	anchorX := float64(padding)
	if centered {
		anchorX += float64(textWidth) / 2
	}
	return scaledTextLayout{
		Width:    textWidth + padding*2,
		Height:   textHeight + padding*2,
		AnchorX:  anchorX,
		AnchorY:  float64(padding) + float64(textHeight)/2,
		Baseline: padding + metrics.Ascent.Ceil(),
	}
}

func (g *Game) scaledTextImage(s string, size float64, centered bool) scaledTextCacheEntry {
	key := scaledTextCacheKey{Text: s, SizeKey: fontSizeKey(size), Centered: centered}
	if entry := g.scaledTextCache[key]; entry.Image != nil {
		return entry
	}

	face := fontFaceForSize(size)
	layout := baseScaledTextLayout(face, s, centered)
	img := ebiten.NewImage(layout.Width, layout.Height)
	text.Draw(img, s, face, 4, layout.Baseline, color.White)
	entry := scaledTextCacheEntry{Image: img, Layout: layout}
	g.scaledTextCache[key] = entry
	return entry
}

func (g *Game) drawScaledTextImage(screen *ebiten.Image, s string, x, centerY, size, scale float64, centered bool, clr color.Color) {
	if scale <= 0 {
		return
	}
	entry := g.scaledTextImage(s, size, centered)
	layout := entry.Layout

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x-layout.AnchorX*scale, centerY-layout.AnchorY*scale)
	op.ColorScale.ScaleWithColor(clr)
	screen.DrawImage(entry.Image, op)
}

func (g *Game) pointInRect(px, py, x, y, w, h float64) bool {
	return px >= x && px <= x+w && py >= y && py <= y+h
}

func (g *Game) levelUpOptionAt(px, py float64) int {
	w := math.Min(math.Max(360, float64(g.screenW)-48), 620)
	h := 250.0 + math.Max(0, float64(max(2, len(g.session.ActiveLevelUpOptions))-2))*52
	x := float64(g.screenW)/2 - w/2
	y := float64(g.screenH)/2 - 90
	_ = h
	for i := range g.session.ActiveLevelUpOptions {
		oy := y + 94 + float64(i)*52
		option := g.session.ActiveLevelUpOptions[i]
		if labelHitContains(px, py, x+w/2-78, oy, levelUpOptionFontSize, option.Title(g.session.Progression.BeamKillUpgradeBonus())) ||
			centeredBoxHitContains(px, py, x+w/2-116, oy, 24, 24, 14, 14) {
			return i
		}
	}
	return -1
}

func (g *Game) redrawRectContains(px, py float64) bool {
	w := math.Min(math.Max(360, float64(g.screenW)-48), 620)
	x := float64(g.screenW)/2 - w/2
	y := float64(g.screenH)/2 - 90
	redrawY := y + 94 + float64(max(2, len(g.session.ActiveLevelUpOptions)))*52 + 14
	label := fmt.Sprintf("REDRAW %d  COINS %d", g.levelUpRedrawCost(), max(0, g.session.CollectedCoins))
	return labelHitContains(px, py, x+w/2-78, redrawY, levelUpOptionFontSize, label) ||
		centeredBoxHitContains(px, py, x+w/2-116, redrawY, 24, 24, 14, 14) ||
		centeredLabelHitContains(px, py, x+w/2-150, redrawY, levelUpKeyFontSize, "[R]")
}

func labelHitContains(px, py, x, centerY, size float64, label string) bool {
	face := fontFaceForSize(size)
	width := float64(font.MeasureString(face, label).Ceil())
	height := float64(face.Metrics().Height.Ceil())
	return px >= x-26 && px <= x+width+26 && py >= centerY-height/2-12 && py <= centerY+height/2+12
}

func centeredLabelHitContains(px, py, centerX, centerY, size float64, label string) bool {
	face := fontFaceForSize(size)
	width := float64(font.MeasureString(face, label).Ceil())
	return labelHitContains(px, py, centerX-width/2, centerY, size, label)
}

func centeredBoxHitContains(px, py, centerX, centerY, width, height, insetX, insetY float64) bool {
	return px >= centerX-width/2-insetX &&
		px <= centerX+width/2+insetX &&
		py >= centerY-height/2-insetY &&
		py <= centerY+height/2+insetY
}

func (g *Game) gameOverOptionAt(px, py float64) string {
	x, y, w, _ := gameOverPanelRect(g.screenW, g.screenH)
	if centeredLabelHitContains(px, py, x+w/2, y+gameOverRestartOffsetY, gameOverOptionFontSize, "RESTART") {
		return "restart"
	}
	if centeredLabelHitContains(px, py, x+w/2, y+gameOverExitOffsetY, gameOverOptionFontSize, "EXIT") {
		return "exit"
	}
	return ""
}

const (
	gameOverTitleOffsetY   = 48.0
	gameOverRestartOffsetY = 112.0
	gameOverExitOffsetY    = 160.0
)

func gameOverPanelRect(screenW, screenH int) (x, y, w, h float64) {
	w = math.Min(math.Max(360, float64(screenW)-48), 620)
	h = 190
	x = float64(screenW)/2 - w/2
	y = float64(screenH)/2 - 90
	return x, y, w, h
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

func formattedSeconds(value float64) string {
	if value >= 1 {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.2f", value)
}
