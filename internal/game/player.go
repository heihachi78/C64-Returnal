package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

func (g *Game) updatePlayer(dt float64) {
	move := g.playerMovementVector(ebitenIsKeyPressed)

	if move.X != 0 || move.Y != 0 {
		move = move.Normalized()
		g.player.MoveDir = move
		g.player.Moving = true
		g.player.Pos = g.player.Pos.Add(move.Mul(g.tuning.PlayerSpeed * dt))
		if move.X < 0 {
			g.player.Facing = -1
		} else if move.X > 0 {
			g.player.Facing = 1
		}
		return
	}

	g.player.Moving = false
	g.player.MoveDir = Vec2{}
	g.player.AnimTimer = 0
	g.player.AnimFrame = 0
}
func (g *Game) playerMovementVector(isPressed func(ebiten.Key) bool) Vec2 {
	for _, key := range movementKeys() {
		if g.suppressedMovement[key] && !isPressed(key) {
			delete(g.suppressedMovement, key)
		}
	}

	move := Vec2{}
	if isPressed(ebiten.KeyArrowLeft) && !g.suppressedMovement[ebiten.KeyArrowLeft] {
		move.X--
	}
	if isPressed(ebiten.KeyArrowRight) && !g.suppressedMovement[ebiten.KeyArrowRight] {
		move.X++
	}
	if isPressed(ebiten.KeyArrowUp) && !g.suppressedMovement[ebiten.KeyArrowUp] {
		move.Y++
	}
	if isPressed(ebiten.KeyArrowDown) && !g.suppressedMovement[ebiten.KeyArrowDown] {
		move.Y--
	}
	return move
}
func (g *Game) suppressHeldMovementKeys(isPressed func(ebiten.Key) bool) {
	if g.suppressedMovement == nil {
		g.suppressedMovement = map[ebiten.Key]bool{}
	}
	for _, key := range movementKeys() {
		if isPressed(key) {
			g.suppressedMovement[key] = true
		}
	}
}
func (g *Game) suppressModalHeldMovementKeys(isPressed func(ebiten.Key) bool) {
	if g.session.LevelUpChoiceActive || g.session.ChestRewardActive {
		g.suppressHeldMovementKeys(isPressed)
	}
}
func movementKeys() []ebiten.Key {
	return []ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyArrowUp, ebiten.KeyArrowDown}
}
func (g *Game) updateInvulnerability(dt float64) {
	if g.session.PlayerHitInvulnerability > 0 {
		g.session.PlayerHitInvulnerability = math.Max(0, g.session.PlayerHitInvulnerability-dt)
	}
}
func (g *Game) damagePlayer() {
	g.session.PlayerLives = max(0, g.session.PlayerLives-1)
	if g.session.PlayerLives == 0 {
		g.triggerGameOver()
		return
	}
	g.session.PlayerHitInvulnerability = g.tuning.PlayerHitInvulnerability
	g.player.HitFlash = playerHitFlashDuration
}
func (g *Game) triggerGameOver() {
	if g.session.GameOver {
		return
	}
	g.session.GameOver = true
	g.session.LevelUpChoiceActive = false
	g.session.ChestRewardActive = false
	g.session.PlayerHitInvulnerability = 0
	g.session.PendingLevelUpLevels = g.session.PendingLevelUpLevels[:0]
	g.session.ActiveLevelUpOptions = nil
	g.session.ActiveChestRewardItems = nil
	clear(g.suppressedMovement)
	g.hideLevelUpPresentation()
	g.session.ChestRewardOverlayTimer = 0
	g.player.Moving = false
	g.player.AnimTimer = 0
	g.player.AnimFrame = 0
	g.player.HitFlash = 0
	g.player.DeathTimer = 0
	g.player.DeathRotation = 0
	g.session.GameOverOverlayTimer = 0
	g.effects = g.effects[:0]
	g.meteors = g.meteors[:0]
}
