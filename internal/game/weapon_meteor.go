package game

import (
	"math"
	"slices"
)

func (g *Game) updateMeteorCasting(dt float64) {
	if !g.session.Progression.MeteorUnlocked || g.session.Progression.MeteorCount() <= 0 || len(g.skeleton) == 0 {
		g.session.Casts.Meteor = 0
		return
	}
	g.session.Casts.Meteor += dt
	interval := g.session.Progression.MeteorCastInterval() / float64(g.session.Progression.MeteorCount())
	for g.session.Casts.Meteor >= interval {
		g.session.Casts.Meteor -= interval
		g.castMeteor()
	}
}
func (g *Game) castMeteor() {
	angle := g.randRange(0, math.Pi*2)
	distance := math.Sqrt(g.randRange(0, 1)) * g.tuning.OrbitalOrbRadius * g.tuning.MeteorTargetMultiplier
	impact := Vec2{X: g.player.Pos.X + math.Cos(angle)*distance, Y: g.player.Pos.Y + math.Sin(angle)*distance}
	start := g.meteorSpawnPosition(impact)
	g.meteors = append(g.meteors, MeteorProjectile{Pos: start, Start: start, Impact: impact})
}
func (g *Game) meteorSpawnPosition(impact Vec2) Vec2 {
	return Vec2{
		X: impact.X + g.randRange(-g.tuning.MeteorFallDrift, g.tuning.MeteorFallDrift),
		Y: math.Max(impact.Y+g.tuning.MeteorFallHeight, g.player.Pos.Y+g.tuning.MeteorFallHeight),
	}
}
func (g *Game) updateMeteors(dt float64) {
	g.updateMeteorAnimation(dt)
	for i := len(g.meteors) - 1; i >= 0; i-- {
		g.meteors[i].Age += dt
		progress := Clamp(g.meteors[i].Age/g.tuning.MeteorFallDuration, 0, 1)
		g.meteors[i].Pos = g.meteors[i].Start.Add(g.meteors[i].Impact.Sub(g.meteors[i].Start).Mul(progress))
		if progress >= 1 {
			impact := g.meteors[i].Impact
			g.meteors = slices.Delete(g.meteors, i, i+1)
			g.impactMeteor(impact)
			if g.session.LevelUpChoiceActive {
				return
			}
		}
	}
}
func (g *Game) impactMeteor(pos Vec2) {
	if g.session.GameOver {
		return
	}
	g.effects = append(g.effects, Effect{Kind: EffectMeteorImpact, Pos: pos, Radius: g.tuning.MeteorImpactRadius, TTL: meteorImpactEffectDuration, MaxTTL: meteorImpactEffectDuration})
	targets := g.meteorImpactTargetIDs(pos)
	levelUps := 0
	for _, id := range targets {
		if idx := g.skeletonIndexByID(id); idx >= 0 {
			levelUps += g.damageSkeleton(idx, 1, AttackMeteor, false)
			if levelUps > 0 {
				break
			}
		}
	}
	g.queueLevelUpChoices(levelUps)
}
func (g *Game) meteorImpactTargetIDs(pos Vec2) []int {
	g.ensureSkeletonSpatialIndex()
	targets := []int{}
	damageRadius := g.meteorImpactDamageRadius()
	g.spatial.ForEachNear(pos, damageRadius, g.skeleton, func(i int) bool {
		radius := g.tuning.MeteorImpactRadius + skeletonBodyRadius(g.tuning, g.skeleton[i].Kind)
		radiusSq := radius * radius
		if DistanceSq(pos, g.skeleton[i].Pos) <= radiusSq {
			targets = append(targets, g.skeleton[i].ID)
		}
		return true
	})
	return targets
}
func (g *Game) meteorImpactDamageRadius() float64 {
	return g.tuning.MeteorImpactRadius + skeletonBodyRadius(g.tuning, SkeletonBlue)
}
func (g *Game) updateMeteorAnimation(dt float64) {
	if len(g.meteors) == 0 {
		g.meteorAnimTimer = 0
		return
	}
	g.meteorAnimTimer += dt
	if g.meteorAnimTimer >= g.tuning.MeteorAnimationFrameTime {
		g.meteorAnimTimer = math.Mod(g.meteorAnimTimer, g.tuning.MeteorAnimationFrameTime)
		g.meteorAnimFrame = (g.meteorAnimFrame + 1) % 2
		for i := range g.meteors {
			g.meteors[i].AnimFrame = g.meteorAnimFrame
		}
	}
}
