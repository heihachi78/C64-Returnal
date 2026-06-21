package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ebitenIsKeyPressed                = ebiten.IsKeyPressed
	ebitenCursorPosition              = ebiten.CursorPosition
	inpututilIsKeyJustPressed         = inpututil.IsKeyJustPressed
	inpututilIsMouseButtonJustPressed = inpututil.IsMouseButtonJustPressed
)

func (g *Game) updateOverlayInput() (bool, error) {
	if g.session.GameOver {
		if inpututilIsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebitenCursorPosition()
			return g.selectGameOverOption(g.gameOverOptionAt(float64(x), float64(y)))
		}
		return false, nil
	}

	if isJumpToLevel100DebugJustPressed() {
		return g.handleJumpToLevel100DebugKeyDown(), nil
	}

	if g.session.ChestRewardActive {
		g.suppressModalHeldMovementKeys(ebitenIsKeyPressed)
		if inpututilIsKeyJustPressed(chestRewardAdvanceKey()) {
			return g.advanceChestReward(), nil
		}
		return false, nil
	}

	if !g.session.LevelUpChoiceActive {
		return false, nil
	}

	g.suppressModalHeldMovementKeys(ebitenIsKeyPressed)
	if inpututilIsKeyJustPressed(levelUpRedrawKey()) {
		g.redrawLevelUpOptions()
		return true, nil
	}

	for i, key := range levelUpOptionKeys() {
		if inpututilIsKeyJustPressed(key) && i < len(g.session.ActiveLevelUpOptions) {
			return g.selectLevelUpOptionAt(i), nil
		}
	}

	if inpututilIsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebitenCursorPosition()
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
		g.restartGame(ebitenIsKeyPressed)
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
func isJumpToLevel100DebugKey(key ebiten.Key) bool {
	return key == ebiten.KeyDigit0 || key == ebiten.KeyNumpad0
}
func isJumpToLevel100DebugJustPressed() bool {
	return inpututilIsKeyJustPressed(ebiten.KeyDigit0) || inpututilIsKeyJustPressed(ebiten.KeyNumpad0)
}
