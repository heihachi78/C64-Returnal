package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

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
func formattedSeconds(value float64) string {
	if value >= 1 {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.2f", value)
}
