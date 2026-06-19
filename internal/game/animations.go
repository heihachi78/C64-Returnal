package game

import (
	"math"
	"slices"
)

func (g *Game) updateEffects(dt float64) {
	for i := len(g.effects) - 1; i >= 0; i-- {
		g.effects[i].TTL -= dt
		if g.effects[i].TTL <= 0 {
			g.effects = slices.Delete(g.effects, i, i+1)
		}
	}
}
func (g *Game) stopPlayerAnimation() {
	g.player.Moving = false
	g.player.AnimTimer = 0
	g.player.AnimFrame = 0
}
func (g *Game) updatePlayerWalkAnimation(dt float64) {
	if !g.player.Moving {
		return
	}
	g.player.AnimTimer += dt
	if g.player.AnimTimer >= g.tuning.PlayerAnimationFrameTime {
		g.player.AnimTimer = math.Mod(g.player.AnimTimer, g.tuning.PlayerAnimationFrameTime)
		g.player.AnimFrame = (g.player.AnimFrame + 1) % 2
	}
}
func (g *Game) updatePlayerHitFlash(dt float64) {
	if g.player.HitFlash > 0 {
		g.player.HitFlash = math.Max(0, g.player.HitFlash-dt)
	}
}
func (g *Game) updatePausedAnimations(dt float64) {
	if g.session.LevelUpChoiceActive {
		g.session.LevelUpOverlayTimer += dt
		g.session.LevelUpTitleScaleTimer += dt
		g.session.LevelUpOptionFadeTimer += dt
		g.session.LevelUpRedrawFadeTimer += dt
		g.session.LevelUpRedrawCoinFadeTimer += dt
		if g.session.LevelUpRedrawStatusTimer > 0 {
			g.session.LevelUpRedrawStatusTimer = math.Max(0, g.session.LevelUpRedrawStatusTimer-dt)
		}
	}
	if g.session.ChestRewardActive {
		g.session.ChestRewardOverlayTimer += dt
	}
	if g.session.GameOver {
		g.session.GameOverOverlayTimer += dt
		g.player.DeathTimer = math.Min(playerDeathRotationDuration, g.player.DeathTimer+dt)
		progress := g.player.DeathTimer / playerDeathRotationDuration
		g.player.DeathRotation = -math.Pi / 2 * progress
		g.updateGameOverWorldActions(dt)
	}
}
func (g *Game) updateNewlyPresentedOverlayActions(dt float64) {
	if g.session.LevelUpChoiceActive {
		g.session.LevelUpOverlayTimer += dt
		g.session.LevelUpTitleScaleTimer += dt
		g.session.LevelUpOptionFadeTimer += dt
		g.session.LevelUpRedrawFadeTimer += dt
		g.session.LevelUpRedrawCoinFadeTimer += dt
		if g.session.LevelUpRedrawStatusTimer > 0 {
			g.session.LevelUpRedrawStatusTimer = math.Max(0, g.session.LevelUpRedrawStatusTimer-dt)
		}
	}
	if g.session.ChestRewardActive {
		g.session.ChestRewardOverlayTimer += dt
	}
	if g.session.GameOver {
		g.session.GameOverOverlayTimer += dt
		g.player.DeathTimer = math.Min(playerDeathRotationDuration, g.player.DeathTimer+dt)
		progress := g.player.DeathTimer / playerDeathRotationDuration
		g.player.DeathRotation = -math.Pi / 2 * progress
		g.updateCoins(dt)
		g.updateSkeletonHitFlashes(dt)
	}
}
func (g *Game) updateGameOverWorldActions(dt float64) {
	g.updateCoins(dt)
	g.updateSkeletonHitFlashes(dt)
}
func (g *Game) updateSkeletonHitFlashes(dt float64) {
	for i := range g.skeleton {
		if g.skeleton[i].HitFlash > 0 {
			g.skeleton[i].HitFlash = math.Max(0, g.skeleton[i].HitFlash-dt)
		}
	}
}
