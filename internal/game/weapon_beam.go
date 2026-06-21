package game

import "math"

func (g *Game) updateBeamCasting(dt float64) {
	if !g.session.Progression.BeamUnlocked || len(g.skeleton) == 0 {
		g.session.Casts.Beam = 0
		return
	}
	g.session.Casts.Beam += dt
	interval := g.session.Progression.BeamCastInterval()
	for g.session.Casts.Beam >= interval {
		g.session.Casts.Beam -= interval
		g.castBeam()
		if g.session.LevelUpChoiceActive {
			return
		}
	}
}
func (g *Game) castBeam() {
	direction := g.playerBeamDirection()
	length := math.Max(700, math.Hypot(float64(g.screenW), float64(g.screenH))/2+g.tuning.SkeletonSpawnMargin)
	damageBudget := g.session.Progression.BeamKillCount()
	targets := g.beamTargets(direction, length, g.tuning.BeamHitWidth, damageBudget)
	end := g.beamVisualEnd(direction, length, targets, damageBudget)
	g.effects = append(g.effects, Effect{Kind: EffectBeam, Start: g.player.Pos, End: end, TTL: g.tuning.BeamEffectDuration, MaxTTL: g.tuning.BeamEffectDuration})

	levelUps := g.applyBeamDamage(targets, damageBudget)
	g.queueLevelUpChoices(levelUps)
}

func (g *Game) beamVisualEnd(direction Vec2, length float64, targets []int, damageBudget int) Vec2 {
	remainingDamage := damageBudget
	visualLength := 0.0
	for _, id := range targets {
		idx := g.skeletonIndexByID(id)
		if idx < 0 || remainingDamage <= 0 {
			break
		}
		damage := min(remainingDamage, g.skeleton[idx].HP)
		if damage <= 0 {
			continue
		}
		target := g.skeleton[idx].Pos.Sub(g.player.Pos)
		progress := target.X*direction.X + target.Y*direction.Y
		visualLength = Clamp(progress, 0, length)
		remainingDamage -= damage
	}
	return g.player.Pos.Add(direction.Mul(visualLength))
}

func (g *Game) applyBeamDamage(targets []int, damageBudget int) int {
	remainingDamage := damageBudget
	levelUps := 0
	for _, id := range targets {
		idx := g.skeletonIndexByID(id)
		if idx < 0 || remainingDamage <= 0 {
			break
		}
		damage := min(remainingDamage, g.skeleton[idx].HP)
		remainingDamage -= damage
		levelUps += g.damageSkeleton(idx, damage, AttackBeam, false)
	}
	return levelUps
}
func (g *Game) playerBeamDirection() Vec2 {
	if g.player.MoveDir != (Vec2{}) {
		return g.player.MoveDir.Normalized()
	}
	return Vec2{X: g.player.Facing, Y: 0}
}
func (g *Game) beamTargets(direction Vec2, length, hitWidth float64, limit int) []int {
	if limit <= 0 {
		return nil
	}
	hits := g.beamTargetScratch[:0]
	if cap(hits) < limit {
		hits = make([]beamTargetHit, 0, limit)
	}
	for i := range g.skeleton {
		target := g.skeleton[i].Pos.Sub(g.player.Pos)
		progress := target.X*direction.X + target.Y*direction.Y
		if progress < 0 || progress > length {
			continue
		}
		closest := g.player.Pos.Add(direction.Mul(progress))
		radius := g.skeletonCollisionRadius(hitWidth, g.skeleton[i].Kind)
		if DistanceSq(closest, g.skeleton[i].Pos) > radius*radius {
			continue
		}
		if len(hits) < limit {
			hits = append(hits, beamTargetHit{g.skeleton[i].ID, progress})
			for j := len(hits) - 1; j > 0 && hits[j].progress < hits[j-1].progress; j-- {
				hits[j], hits[j-1] = hits[j-1], hits[j]
			}
		} else if progress < hits[len(hits)-1].progress {
			hits[len(hits)-1] = beamTargetHit{g.skeleton[i].ID, progress}
			for j := len(hits) - 1; j > 0 && hits[j].progress < hits[j-1].progress; j-- {
				hits[j], hits[j-1] = hits[j-1], hits[j]
			}
		}
	}
	g.beamTargetScratch = hits
	result := g.beamTargetResult[:0]
	if cap(result) < len(hits) {
		result = make([]int, 0, len(hits))
	}
	for _, h := range hits {
		result = append(result, h.id)
	}
	g.beamTargetResult = result
	return result
}

type beamTargetHit struct {
	id       int
	progress float64
}
