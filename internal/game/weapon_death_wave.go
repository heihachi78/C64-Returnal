package game

import "slices"

func (g *Game) updateDeathWaveCasting(dt float64) {
	if !g.session.Progression.DeathWaveUnlocked {
		g.session.Casts.DeathWave = 0
		return
	}
	g.session.Casts.DeathWave += dt
	for g.session.Casts.DeathWave >= g.tuning.DeathWaveInterval {
		g.session.Casts.DeathWave -= g.tuning.DeathWaveInterval
		g.castDeathWave()
	}
}

func (g *Game) castDeathWave() {
	maxRadius := max(g.tuning.DeathWaveMaxRadius, deathWaveVisibleRadius(g.screenW, g.screenH, g.tuning.SkeletonSpawnMargin))
	g.deathWaves = append(g.deathWaves, DeathWave{
		Origin:    g.player.Pos,
		MaxRadius: maxRadius,
	})
}

func deathWaveVisibleRadius(screenW, screenH int, spawnMargin float64) float64 {
	return Vec2{X: float64(screenW) / 2, Y: float64(screenH) / 2}.Len() + spawnMargin
}

func (g *Game) updateDeathWaves(dt float64) {
	for i := len(g.deathWaves) - 1; i >= 0; i-- {
		g.deathWaves[i].PreviousRadius = g.deathWaves[i].Radius
		g.deathWaves[i].Radius += g.tuning.DeathWaveSpeed * dt
		g.applyDeathWaveDamage(&g.deathWaves[i])
		if g.deathWaves[i].Radius > g.deathWaves[i].MaxRadius+g.tuning.DeathWaveWidth {
			g.deathWaves = slices.Delete(g.deathWaves, i, i+1)
		}
	}
}

func (g *Game) applyDeathWaveDamage(wave *DeathWave) {
	if wave == nil {
		return
	}
	for i := range g.skeleton {
		if g.skeleton[i].Kind == SkeletonRegular || g.skeleton[i].HP <= 1 || deathWaveAlreadyHit(wave.HitIDs, g.skeleton[i].ID) {
			continue
		}
		if !deathWaveTouchesSkeleton(*wave, g.skeleton[i], g.tuning) {
			continue
		}
		damage := g.skeleton[i].HP / 2
		g.recordActualDamage(damage)
		g.skeleton[i].HP -= damage
		g.skeleton[i].HitFlash = skeletonDamageFlashDuration
		wave.HitIDs = append(wave.HitIDs, g.skeleton[i].ID)
	}
}

func deathWaveAlreadyHit(hitIDs []int, id int) bool {
	for _, hitID := range hitIDs {
		if hitID == id {
			return true
		}
	}
	return false
}

func deathWaveTouchesSkeleton(wave DeathWave, skeleton Skeleton, tuning Tuning) bool {
	distance := skeleton.Pos.Sub(wave.Origin).Len()
	bodyRadius := skeletonBodyRadius(tuning, skeleton.Kind)
	inner := max(0, wave.PreviousRadius-tuning.DeathWaveWidth/2-bodyRadius)
	outer := wave.Radius + tuning.DeathWaveWidth/2 + bodyRadius
	return distance >= inner && distance <= outer
}
