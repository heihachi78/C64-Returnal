package game

import "github.com/hajimehoshi/ebiten/v2"

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
