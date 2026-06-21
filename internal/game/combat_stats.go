package game

func (g *Game) recordActualDamage(amount int) {
	if amount <= 0 {
		return
	}
	g.actualDamage = append(g.actualDamage, actualDamageSample{Time: g.totalTime, Amount: amount})
	g.actualDamageWindowTotal += amount
	g.pruneActualDamageSamples()
	g.maxActualDPS = max(g.maxActualDPS, g.ActualDPS())
}

func (g *Game) ActualDPS() float64 {
	g.pruneActualDamageSamples()
	return float64(g.actualDamageWindowTotal) / actualDPSWindow
}

func (g *Game) pruneActualDamageSamples() {
	cutoff := g.totalTime - actualDPSWindow
	first := 0
	for first < len(g.actualDamage) && g.actualDamage[first].Time < cutoff {
		g.actualDamageWindowTotal -= g.actualDamage[first].Amount
		first++
	}
	g.actualDamage = g.actualDamage[first:]
	if g.actualDamageWindowTotal < 0 {
		g.actualDamageWindowTotal = 0
	}
}
