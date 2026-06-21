package game

func (g *Game) Update() error {
	dt := g.beginFrame()
	wasWorldPaused := g.worldPaused()
	if wasWorldPaused {
		g.pauseActualDamageLevelStats(dt)
	}

	consumedFrame, err := g.updateOverlayInput()
	if err != nil {
		return err
	}
	if consumedFrame {
		return nil
	}

	g.updatePausedAnimations(dt)
	if g.worldPaused() {
		return nil
	}

	g.updatePlayer(dt)
	if g.updatePickups() {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	if g.updateCombat(dt) {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.checkSkeletonCollisions()
	g.updateSkeletonSpawning(dt)
	if g.session.GameOver {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateWorldActions(dt)
	return nil
}

func (g *Game) beginFrame() float64 {
	dt := 0.0
	if g.hasUpdated {
		dt = 1.0 / float64(TargetTPS)
	}
	g.hasUpdated = true
	g.totalTime += dt
	clear(g.lightningTargetReservations)
	return dt
}

func (g *Game) worldPaused() bool {
	return g.session.GameOver ||
		g.session.LevelUpChoiceActive ||
		g.session.ChestRewardActive
}

func (g *Game) updatePickups() bool {
	g.checkCoinPickups()
	g.checkChestPickups()
	return g.session.ChestRewardActive
}

func (g *Game) updateCombat(dt float64) bool {
	g.updateSkeletons(dt)
	g.updateOrbitalOrbs(dt)
	if g.session.LevelUpChoiceActive {
		return true
	}

	g.updateLightningCasting(dt)
	if g.session.LevelUpChoiceActive {
		return true
	}

	g.updateFireballCasting(dt)
	g.updateFireballs(dt)
	if g.session.LevelUpChoiceActive {
		return true
	}

	g.updateBeamCasting(dt)
	if g.session.LevelUpChoiceActive {
		return true
	}

	g.updateMeteorCasting(dt)
	g.updateMeteors(dt)
	g.updateDeathWaveCasting(dt)
	g.updateDeathWaves(dt)
	g.updateInvulnerability(dt)
	return g.session.LevelUpChoiceActive
}

func (g *Game) updateWorldActions(dt float64) {
	g.updatePlayerWalkAnimation(dt)
	g.updatePlayerHitFlash(dt)
	g.updateSkeletonHitFlashes(dt)
	g.updateCoins(dt)
	g.updateEffects(dt)
}
