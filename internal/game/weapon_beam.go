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
	end := g.player.Pos.Add(direction.Mul(length))
	g.effects = append(g.effects, Effect{Kind: EffectBeam, Start: g.player.Pos, End: end, TTL: g.tuning.BeamEffectDuration, MaxTTL: g.tuning.BeamEffectDuration})

	targets := g.beamTargets(direction, length, g.tuning.BeamHitWidth, g.session.Progression.BeamKillCount())
	remainingDamage := g.session.Progression.BeamKillCount()
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
	g.queueLevelUpChoices(levelUps)
}
func (g *Game) playerBeamDirection() Vec2 {
	if g.player.MoveDir != (Vec2{}) {
		return g.player.MoveDir.Normalized()
	}
	return Vec2{X: g.player.Facing, Y: 0}
}
func (g *Game) beamTargets(direction Vec2, length, hitWidth float64, limit int) []int {
	type hit struct {
		id       int
		progress float64
	}
	hits := make([]hit, 0, limit)
	for i := range g.skeleton {
		target := g.skeleton[i].Pos.Sub(g.player.Pos)
		progress := target.X*direction.X + target.Y*direction.Y
		if progress < 0 || progress > length {
			continue
		}
		closest := g.player.Pos.Add(direction.Mul(progress))
		if DistanceSq(closest, g.skeleton[i].Pos) > hitWidth*hitWidth {
			continue
		}
		if len(hits) < limit {
			hits = append(hits, hit{g.skeleton[i].ID, progress})
			for j := len(hits) - 1; j > 0 && hits[j].progress < hits[j-1].progress; j-- {
				hits[j], hits[j-1] = hits[j-1], hits[j]
			}
		} else if progress < hits[len(hits)-1].progress {
			hits[len(hits)-1] = hit{g.skeleton[i].ID, progress}
			for j := len(hits) - 1; j > 0 && hits[j].progress < hits[j-1].progress; j-- {
				hits[j], hits[j-1] = hits[j-1], hits[j]
			}
		}
	}
	result := make([]int, len(hits))
	for i, h := range hits {
		result[i] = h.id
	}
	return result
}
