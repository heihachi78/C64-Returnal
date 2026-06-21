package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) drawHUD(screen *ebiten.Image) {
	left := 18.0
	g.drawStatusPanel(screen, left)
	g.drawCombatStatusPanel(screen, left)
	g.drawDPSPanel(screen)
}

func (g *Game) drawStatusPanel(screen *ebiten.Image, left float64) {
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
}

func (g *Game) drawCombatStatusPanel(screen *ebiten.Image, left float64) {
	bottomX, bottomY, bottomW, bottomH := combatStatusPanelRect(g.screenH)
	g.panel(screen, bottomX, bottomY, bottomW, bottomH)
	bottomPanelTop := bottomY
	for _, row := range g.combatHUDRows() {
		g.drawCombatRow(screen, row.icon, left+9, bottomPanelTop+row.y, row.first, row.second, row.kills, row.unlocked)
	}
	g.drawSpriteScreen(screen, g.assets.Skeleton[0], left+9, bottomPanelTop+25, 16, 22, false, color.RGBA{255, 255, 255, 255})
	g.drawTextSize(screen, fmt.Sprintf("x%d", len(g.skeleton)), left+28, bottomPanelTop+15, combatFontSize, c64Text)
	g.drawTextSize(screen, fmt.Sprintf("%ss", formattedSeconds(g.SkeletonSpawnInterval())), left+28, bottomPanelTop+31, combatFontSize, c64Text)
}

func (g *Game) drawDPSPanel(screen *ebiten.Image) {
	dpsX, dpsY, dpsW, dpsH := dpsPanelRect(g.screenW, g.screenH)
	g.panel(screen, dpsX, dpsY, dpsW, dpsH)
	hpRate, maxActual := g.dpsPanelReadouts()
	g.drawCenteredTextSize(screen, hpRate, dpsX+dpsW/2, dpsY+19, combatFontSize, c64Text)
	g.drawCenteredTextSize(screen, maxActual, dpsX+dpsW/2, dpsY+39, combatFontSize, c64Gold)
}

func (g *Game) dpsPanelReadouts() (hpRate, maxActual string) {
	return fmt.Sprintf("HP/S %.2f", g.SkeletonHPPerSecond()),
		fmt.Sprintf("DPS %.2f", g.ActualDPS())
}

type combatHUDRow struct {
	icon     *ebiten.Image
	y        float64
	first    string
	second   string
	kills    string
	unlocked bool
}

func (g *Game) combatHUDRows() []combatHUDRow {
	p := g.session.Progression
	kills := g.session.Kills
	return []combatHUDRow{
		{
			icon:     g.assets.Fireball[0],
			y:        295,
			first:    fmt.Sprintf("x%d", p.SimultaneousFireball),
			second:   fmt.Sprintf("%ss", formattedSeconds(p.FireballCastInterval())),
			kills:    fmt.Sprintf("KILLS %d", kills.Fireball),
			unlocked: true,
		},
		{
			icon:     g.assets.Lightning,
			y:        241,
			first:    fmt.Sprintf("x%d", p.LightningStrikeCount()),
			second:   fmt.Sprintf("%ss", formattedSeconds(p.LightningCastInterval())),
			kills:    fmt.Sprintf("KILLS %d", kills.Lightning),
			unlocked: p.LightningUnlocked,
		},
		{
			icon:     g.assets.Orb[0],
			y:        187,
			first:    fmt.Sprintf("x%d", p.OrbitalOrbCount()),
			second:   fmt.Sprintf("%.1fr/s", p.OrbitalAngularSpeed()),
			kills:    fmt.Sprintf("KILLS %d", kills.OrbitalOrb),
			unlocked: p.OrbitalOrbUnlocked,
		},
		{
			icon:     g.assets.Beam,
			y:        133,
			first:    fmt.Sprintf("x%d", p.BeamKillCount()),
			second:   fmt.Sprintf("%ss", formattedSeconds(p.BeamCastInterval())),
			kills:    fmt.Sprintf("KILLS %d", kills.Beam),
			unlocked: p.BeamUnlocked,
		},
		{
			icon:     g.assets.Meteor[0],
			y:        79,
			first:    fmt.Sprintf("x%d", p.MeteorCount()),
			second:   fmt.Sprintf("%ss", formattedSeconds(p.MeteorCastInterval())),
			kills:    fmt.Sprintf("KILLS %d", kills.Meteor),
			unlocked: p.MeteorUnlocked,
		},
	}
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
func dpsPanelRect(screenW, screenH int) (x, y, w, h float64) {
	w = 136
	h = 58
	return float64(screenW) - w - 8, float64(screenH) - h - 8, w, h
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
