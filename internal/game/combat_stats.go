package game

func (g *Game) recordActualDamage(amount int) {
	if amount <= 0 {
		return
	}
	g.actualDamage = append(g.actualDamage, actualDamageSample{Time: g.totalTime, Amount: amount})
	g.pruneActualDamageSamples()
}

func (g *Game) ActualDPS() float64 {
	if actualDPSWindow <= 0 {
		return 0
	}
	cutoff := g.totalTime - actualDPSWindow
	damage := 0
	for _, sample := range g.actualDamage {
		if sample.Time >= cutoff {
			damage += sample.Amount
		}
	}
	return float64(damage) / actualDPSWindow
}

func (g *Game) pruneActualDamageSamples() {
	cutoff := g.totalTime - actualDPSWindow
	first := 0
	for first < len(g.actualDamage) && g.actualDamage[first].Time < cutoff {
		first++
	}
	g.actualDamage = g.actualDamage[first:]
}
