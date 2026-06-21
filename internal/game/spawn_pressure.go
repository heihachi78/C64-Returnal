package game

import (
	"math"
	"slices"
)

func initialSkeletonHPPerSecond(t Tuning, p Progression) float64 {
	return capSkeletonHPPerSecond(t.InitialSkeletonHPPerSecond, p.MageRawDPS())
}

func (g *Game) updateSkeletonSpawning(dt float64) {
	rate := g.SkeletonHPPerSecond()
	if rate <= 0 {
		g.session.Casts.SkeletonSpawn = 0
		g.dynamicSpawnQueue = g.dynamicSpawnQueue[:0]
		return
	}
	g.session.Casts.SkeletonSpawn += dt * rate
	if len(g.dynamicSpawnQueue) == 0 {
		g.planDynamicSkeletonSpawns(rate)
	}
	g.spawnQueuedDynamicSkeletons()
}

func (g *Game) planDynamicSkeletonSpawns(rate float64) {
	planBudget := max(float64(SkeletonRegular.HitPoints(g.tuning)), rate*dynamicSpawnBatchInterval)
	g.dynamicSpawnQueue = appendDynamicSkeletonSpawnPlanEntries(g.dynamicSpawnQueue[:0], g.tuning, planBudget)
}

type dynamicSpawnPlanEntry struct {
	Kind  SkeletonKind
	Count int
}

func dynamicSkeletonSpawnPlanEntries(t Tuning, budget float64) []dynamicSpawnPlanEntry {
	return appendDynamicSkeletonSpawnPlanEntries(nil, t, budget)
}

func appendDynamicSkeletonSpawnPlanEntries(plan []dynamicSpawnPlanEntry, t Tuning, budget float64) []dynamicSpawnPlanEntry {
	for _, kind := range dynamicSkeletonSpawnOrder(t) {
		hp := kind.HitPoints(t)
		if kind == SkeletonRegular {
			count := int(budget / float64(hp))
			if count > 0 {
				plan = append(plan, dynamicSpawnPlanEntry{Kind: kind, Count: count})
			}
			break
		}
		if budget >= float64(hp) {
			budget -= float64(hp)
			plan = append(plan, dynamicSpawnPlanEntry{Kind: kind, Count: 1})
		}
	}
	return plan
}

func countDynamicSkeletonSpawnPlan(t Tuning, budget float64, kind SkeletonKind) int {
	count := 0
	for _, entry := range dynamicSkeletonSpawnPlanEntries(t, budget) {
		if entry.Kind == kind {
			count += entry.Count
		}
	}
	return count
}

func (g *Game) spawnQueuedDynamicSkeletons() {
	spawned := 0
	for len(g.dynamicSpawnQueue) > 0 {
		if !g.canSpawnDynamicSkeleton(spawned) {
			return
		}
		entry := &g.dynamicSpawnQueue[0]
		kind := entry.Kind
		hp := kind.HitPoints(g.tuning)
		if g.session.Casts.SkeletonSpawn < float64(hp) {
			return
		}
		g.session.Casts.SkeletonSpawn -= float64(hp)
		entry.Count--
		if entry.Count <= 0 {
			g.dynamicSpawnQueue = g.dynamicSpawnQueue[1:]
		}
		g.spawnSkeleton(kind)
		spawned++
	}
}

func (g *Game) canSpawnDynamicSkeleton(spawnedThisTick int) bool {
	maxActive := g.tuning.MaxActiveSkeletons
	if maxActive > 0 && len(g.skeleton) >= maxActive {
		return false
	}
	maxPerTick := g.tuning.MaxSkeletonSpawnsPerTick
	if maxPerTick > 0 && spawnedThisTick >= maxPerTick {
		return false
	}
	return true
}

func dynamicSkeletonSpawnOrder(t Tuning) []SkeletonKind {
	order := []SkeletonKind{SkeletonBlue, SkeletonBlack, SkeletonPurple, SkeletonRed, SkeletonRegular}
	slices.SortStableFunc(order, func(a, b SkeletonKind) int {
		return b.HitPoints(t) - a.HitPoints(t)
	})
	return order
}

func (g *Game) SkeletonHPPerSecond() float64 {
	return max(0, g.skeletonHPPerSecond)
}

func (g *Game) SkeletonSpawnInterval() float64 {
	rate := g.SkeletonHPPerSecond()
	if rate <= 0 {
		return math.Inf(1)
	}
	return float64(SkeletonRegular.HitPoints(g.tuning)) / rate
}

func (g *Game) queueDynamicSpawnPressureForLevelUp(count int) {
	if count <= 0 {
		return
	}
	if g.maxActualDPS > 0 {
		g.pendingSpawnPressureActual = max(g.pendingSpawnPressureActual, g.maxActualDPS)
	}
	g.pendingSpawnPressureLevels += count
	g.maxActualDPS = 0
	g.actualDamage = g.actualDamage[:0]
	g.actualDamageWindowTotal = 0
}

func (g *Game) applyPendingDynamicSpawnPressure() {
	if g.pendingSpawnPressureLevels <= 0 {
		g.capDynamicSpawnPressure()
		return
	}
	rawDPS := g.session.Progression.MageRawDPS()
	headroom := max(0, rawDPS-g.pendingSpawnPressureActual)
	increase := max(0, g.tuning.DynamicSpawnPressureFactor) * headroom
	g.skeletonHPPerSecond = increaseSkeletonHPPerSecond(g.SkeletonHPPerSecond(), increase, rawDPS)
	g.pendingSpawnPressureLevels--
	if g.pendingSpawnPressureLevels == 0 {
		g.pendingSpawnPressureActual = 0
	}
}

func (g *Game) capDynamicSpawnPressure() {
	g.skeletonHPPerSecond = g.SkeletonHPPerSecond()
}

func capSkeletonHPPerSecond(desired, rawDPS float64) float64 {
	desired = max(0, desired)
	if rawDPS <= 0 {
		return 0
	}
	return min(desired, rawDPS)
}

func increaseSkeletonHPPerSecond(current, increase, rawDPS float64) float64 {
	current = max(0, current)
	increase = max(0, increase)
	if rawDPS <= current {
		return current
	}
	return min(current+increase, rawDPS)
}
