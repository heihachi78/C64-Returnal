package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image/color"
	"math"
)

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
	case BuyDeathWaveScroll:
		img = g.assets.Beam
	default:
		img = g.assets.Fireball[0]
	}
	g.drawSpriteScreen(screen, img, x, y, 24, 24, false, tint)
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
