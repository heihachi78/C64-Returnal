package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) updateOverlayInput() (bool, error) {
	if g.session.GameOver {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			return g.selectGameOverOption(g.gameOverOptionAt(float64(x), float64(y)))
		}
		return false, nil
	}

	if g.session.ChestRewardActive {
		g.suppressModalHeldMovementKeys(ebiten.IsKeyPressed)
		if inpututil.IsKeyJustPressed(chestRewardAdvanceKey()) {
			return g.advanceChestReward(), nil
		}
		return false, nil
	}

	if !g.session.LevelUpChoiceActive {
		if isKillAllAndGrantExperienceJustPressed() {
			return g.handleKillAllAndGrantExperienceKeyDown(), nil
		}
		return false, nil
	}

	g.suppressModalHeldMovementKeys(ebiten.IsKeyPressed)
	if inpututil.IsKeyJustPressed(levelUpRedrawKey()) {
		g.redrawLevelUpOptions()
		return true, nil
	}

	for i, key := range levelUpOptionKeys() {
		if inpututil.IsKeyJustPressed(key) && i < len(g.session.ActiveLevelUpOptions) {
			return g.selectLevelUpOptionAt(i), nil
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.redrawRectContains(float64(x), float64(y)) {
			g.redrawLevelUpOptions()
			return true, nil
		}
		if idx := g.levelUpOptionAt(float64(x), float64(y)); idx >= 0 && idx < len(g.session.ActiveLevelUpOptions) {
			return g.selectLevelUpOptionAt(idx), nil
		}
	}
	return false, nil
}
func (g *Game) selectGameOverOption(option string) (bool, error) {
	switch option {
	case "restart":
		g.restartGame(ebiten.IsKeyPressed)
		return true, nil
	case "exit":
		return true, ebiten.Termination
	default:
		return false, nil
	}
}
func (g *Game) restartGame(isPressed func(ebiten.Key) bool) {
	g.reset()
	g.suppressHeldMovementKeys(isPressed)
}
func (g *Game) advanceChestReward() bool {
	if !g.session.ChestRewardActive {
		return false
	}
	g.session.ChestRewardActive = false
	g.session.ActiveChestRewardItems = nil
	g.session.ChestRewardOverlayTimer = 0
	g.presentNextLevelUpChoiceIfNeeded()
	return true
}
func (g *Game) selectLevelUpOptionAt(index int) bool {
	if !g.session.LevelUpChoiceActive || index < 0 || index >= len(g.session.ActiveLevelUpOptions) {
		return false
	}
	g.applyLevelUpOption(g.session.ActiveLevelUpOptions[index])
	return true
}
func levelUpOptionKeys() []ebiten.Key {
	return []ebiten.Key{ebiten.KeyQ, ebiten.KeyA, ebiten.KeyC, ebiten.KeyX}
}
func levelUpRedrawKey() ebiten.Key {
	return ebiten.KeyR
}
func chestRewardAdvanceKey() ebiten.Key {
	return ebiten.KeyQ
}
func isKillAllAndGrantExperienceKey(key ebiten.Key) bool {
	return key == ebiten.KeyDigit1 || key == ebiten.KeyNumpad1
}
func isKillAllAndGrantExperienceJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit1) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad1)
}
