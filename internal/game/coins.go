package game

import (
	"math"
	"slices"
)

func (g *Game) spawnCoinForLevel(level int) {
	if level <= 0 || g.session.SpawnedCoinLevels[level] {
		return
	}
	g.session.SpawnedCoinLevels[level] = true
	minReward := max(1, g.tuning.CoinMinimumReward)
	maxReward := max(minReward, g.tuning.CoinMaximumReward)
	g.coins = append(g.coins, Coin{
		Pos:    g.randomCoinPosition(),
		Amount: g.rng.Intn(maxReward-minReward+1) + minReward,
		Level:  level,
	})
}
func (g *Game) randomCoinPosition() Vec2 {
	halfW := math.Max(1, float64(g.screenW)/2)
	halfH := math.Max(1, float64(g.screenH)/2)
	margin := math.Max(1, g.tuning.CoinSpawnMargin)
	switch g.rng.Intn(4) {
	case 0:
		return Vec2{X: g.player.Pos.X - halfW - margin - g.randRange(0, halfW), Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
	case 1:
		return Vec2{X: g.player.Pos.X + halfW + margin + g.randRange(0, halfW), Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
	case 2:
		return Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y - halfH - margin - g.randRange(0, halfH)}
	default:
		return Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y + halfH + margin + g.randRange(0, halfH)}
	}
}
func (g *Game) checkCoinPickups() {
	distSq := g.tuning.CoinPickupDistance * g.tuning.CoinPickupDistance
	for i := len(g.coins) - 1; i >= 0; i-- {
		if DistanceSq(g.coins[i].Pos, g.player.Pos) <= distSq {
			g.session.CollectedCoins += g.coins[i].Amount
			g.coins = slices.Delete(g.coins, i, i+1)
			return
		}
	}
}
func (g *Game) updateCoins(dt float64) {
	for i := range g.coins {
		g.coins[i].Phase += dt
	}
}
