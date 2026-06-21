package game

func (g *Game) recordActualDamage(amount int) {
	if amount <= 0 {
		return
	}
	g.actualDamageLevelTotal += amount
}

func (g *Game) ActualDPS() float64 {
	if g.actualDamageLevelTotal <= 0 {
		return 0
	}
	elapsed := g.totalTime - g.actualDamageLevelStartTime - g.actualDamageLevelPausedTime
	elapsed = max(elapsed, 1.0)
	return float64(g.actualDamageLevelTotal) / elapsed
}

func (g *Game) pauseActualDamageLevelStats(dt float64) {
	if dt <= 0 {
		return
	}
	g.actualDamageLevelPausedTime += dt
}

func (g *Game) resetActualDamageLevelStats() {
	g.actualDamageLevelTotal = 0
	g.actualDamageLevelStartTime = g.totalTime
	g.actualDamageLevelPausedTime = 0
}
