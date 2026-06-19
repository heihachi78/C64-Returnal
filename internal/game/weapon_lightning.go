package game

import (
	"math"
	"slices"
)

func (g *Game) updateLightningCasting(dt float64) {
	if !g.session.Progression.LightningUnlocked || len(g.skeleton) == 0 {
		g.session.Casts.Lightning = 0
		return
	}
	g.session.Casts.Lightning += dt
	interval := g.session.Progression.LightningCastInterval()
	for g.session.Casts.Lightning >= interval {
		g.session.Casts.Lightning -= interval
		g.castLightning()
		if g.session.LevelUpChoiceActive {
			return
		}
	}
}
func (g *Game) castLightning() {
	strikes := g.chainLightningTargets()
	g.reserveLightningTargets(strikes)
	levelUps := g.applyLightningStrikes(strikes)
	g.queueLevelUpChoices(levelUps)
}
func (g *Game) reserveLightningTargets(strikes []lightningStrikeTarget) {
	if len(strikes) == 0 {
		return
	}
	if g.lightningTargetReservations == nil {
		g.lightningTargetReservations = map[int]bool{}
	}
	for _, strike := range strikes {
		g.lightningTargetReservations[strike.targetID] = true
	}
}
func (g *Game) applyLightningStrikes(strikes []lightningStrikeTarget) int {
	levelUps := 0
	start := g.player.Pos
	for _, strike := range strikes {
		idx := g.skeletonIndexByID(strike.targetID)
		if idx < 0 {
			continue
		}
		end := strike.end
		g.effects = append(g.effects, Effect{
			Kind:        EffectLightning,
			Start:       start,
			End:         end,
			Points:      g.lightningBoltPoints(start, end),
			InnerPoints: g.lightningBoltPoints(start, end),
			TTL:         g.tuning.LightningEffectDuration,
			MaxTTL:      g.tuning.LightningEffectDuration,
		})
		g.effects = append(g.effects, Effect{
			Kind:   EffectLightningHit,
			Pos:    end,
			Frame:  g.skeleton[idx].AnimFrame,
			Facing: g.skeleton[idx].Facing,
			TTL:    g.tuning.LightningEffectDuration,
			MaxTTL: g.tuning.LightningEffectDuration,
		})
		levelUps += g.damageSkeleton(idx, 1, AttackLightning, false)
		start = end
	}
	return levelUps
}

type lightningStrikeTarget struct {
	targetID int
	end      Vec2
}

func (g *Game) lightningBoltPoints(start, end Vec2) []Vec2 {
	delta := end.Sub(start)
	distance := math.Max(1, delta.Len())
	normal := Vec2{X: -delta.Y / distance, Y: delta.X / distance}
	segmentCount := max(3, min(9, int(distance/30)))
	points := make([]Vec2, 0, segmentCount+1)
	points = append(points, start)
	for segment := 1; segment < segmentCount; segment++ {
		progress := float64(segment) / float64(segmentCount)
		base := start.Add(delta.Mul(progress))
		points = append(points, base.Add(normal.Mul(g.randRange(-8, 8))))
	}
	points = append(points, end)
	return points
}
func (g *Game) chainLightningTargets() []lightningStrikeTarget {
	count := g.session.Progression.LightningStrikeCount()
	if count <= 0 {
		return nil
	}
	reserved := map[int]bool{}
	for _, fire := range g.fireball {
		if fire.TargetID != 0 {
			reserved[fire.TargetID] = true
		}
	}
	remaining := make([]int, 0, len(g.skeleton))
	for i := range g.skeleton {
		if !reserved[g.skeleton[i].ID] {
			remaining = append(remaining, i)
		}
	}
	targets := make([]lightningStrikeTarget, 0, count)
	origin := g.player.Pos
	for len(targets) < count && len(remaining) > 0 {
		best := 0
		bestDist := DistanceSq(origin, g.skeleton[remaining[0]].Pos)
		for i := 1; i < len(remaining); i++ {
			dist := DistanceSq(origin, g.skeleton[remaining[i]].Pos)
			if dist < bestDist {
				best = i
				bestDist = dist
			}
		}
		idx := remaining[best]
		targets = append(targets, lightningStrikeTarget{targetID: g.skeleton[idx].ID, end: g.skeleton[idx].Pos})
		origin = g.skeleton[idx].Pos
		remaining = slices.Delete(remaining, best, best+1)
	}
	return targets
}
