package game

import "math/rand"

func (g *Game) skeletonIndexByID(id int) int {
	if id == 0 {
		return -1
	}
	for i := range g.skeleton {
		if g.skeleton[i].ID == id {
			return i
		}
	}
	return -1
}
func (g *Game) randRange(minValue, maxValue float64) float64 {
	return minValue + g.rng.Float64()*(maxValue-minValue)
}
func chance(rng *rand.Rand, numerator, denominator int) bool {
	return numerator > 0 && denominator > 0 && rng.Intn(denominator)+1 <= numerator
}
func sameOptionSet(a, b []LevelUpOption) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[LevelUpOption]int, len(a))
	for _, option := range a {
		seen[option]++
	}
	for _, option := range b {
		seen[option]--
		if seen[option] < 0 {
			return false
		}
	}
	return true
}
