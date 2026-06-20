package game

import (
	"math"
	"slices"
)

func (g *Game) updateFireballCasting(dt float64) {
	if len(g.skeleton) == 0 {
		g.session.Casts.Fireball = 0
		return
	}
	g.session.Casts.Fireball += dt
	interval := g.session.Progression.FireballCastInterval()
	for g.session.Casts.Fireball >= interval {
		g.session.Casts.Fireball -= interval
		g.spawnFireballs()
	}
}
func (g *Game) spawnFireballs() {
	reserved := map[int]bool{}
	for _, fire := range g.fireball {
		if fire.TargetID != 0 {
			reserved[fire.TargetID] = true
		}
	}
	for id := range g.lightningTargetReservations {
		reserved[id] = true
	}
	for _, idx := range g.closestSkeletons(g.player.Pos, reserved, g.session.Progression.SimultaneousFireball) {
		target := g.skeleton[idx]
		g.fireball = append(g.fireball, Fireball{
			Pos:      g.player.Pos,
			TargetID: target.ID,
			Velocity: target.Pos.Sub(g.player.Pos).Normalized(),
		})
	}
}
func (g *Game) closestSkeletons(pos Vec2, excluded map[int]bool, limit int) []int {
	if limit <= 0 {
		return nil
	}
	type selected struct {
		index int
		dist  float64
	}
	picks := make([]selected, 0, limit)
	for i := range g.skeleton {
		if excluded[g.skeleton[i].ID] {
			continue
		}
		dist := DistanceSq(pos, g.skeleton[i].Pos)
		if len(picks) < limit {
			picks = append(picks, selected{i, dist})
			for j := len(picks) - 1; j > 0 && picks[j].dist < picks[j-1].dist; j-- {
				picks[j], picks[j-1] = picks[j-1], picks[j]
			}
		} else if dist < picks[len(picks)-1].dist {
			picks[len(picks)-1] = selected{i, dist}
			for j := len(picks) - 1; j > 0 && picks[j].dist < picks[j-1].dist; j-- {
				picks[j], picks[j-1] = picks[j-1], picks[j]
			}
		}
	}
	result := make([]int, len(picks))
	for i, pick := range picks {
		result[i] = pick.index
	}
	return result
}
func (g *Game) updateFireballs(dt float64) {
	g.updateFireballAnimation(dt)
	for i := len(g.fireball) - 1; i >= 0; i-- {
		fire := &g.fireball[i]
		targetIndex := g.skeletonIndexByID(fire.TargetID)
		if targetIndex >= 0 {
			g.updateHomingFireball(i, targetIndex, dt)
		} else {
			fire.TargetID = 0
			g.updateUntargetedFireball(i, dt)
		}
		if g.session.LevelUpChoiceActive {
			return
		}
	}
}
func (g *Game) updateHomingFireball(i, targetIndex int, dt float64) {
	fire := &g.fireball[i]
	toTarget := g.skeleton[targetIndex].Pos.Sub(fire.Pos)
	distanceSq := toTarget.LenSq()
	travel := g.tuning.FireballSpeed * dt
	hitDistance := g.skeletonCollisionRadius(g.tuning.FireballHitDistance, g.skeleton[targetIndex].Kind) + travel
	if distanceSq == 0 || distanceSq <= hitDistance*hitDistance {
		g.damageSkeleton(targetIndex, 1, AttackFireball, true)
		g.removeFireball(i)
		return
	}
	fire.Velocity = toTarget.Normalized()
	fire.Pos = fire.Pos.Add(fire.Velocity.Mul(travel))
}
func (g *Game) updateUntargetedFireball(i int, dt float64) {
	fire := &g.fireball[i]
	start := fire.Pos
	fire.TimeWithoutTarget += dt
	fire.Pos = fire.Pos.Add(fire.Velocity.Mul(g.tuning.FireballSpeed * dt))
	if idx := g.firstSkeletonHitBySegment(start, fire.Pos, g.tuning.FireballHitDistance); idx >= 0 {
		g.damageSkeleton(idx, 1, AttackFireball, true)
		g.removeFireball(i)
		return
	}
	if fire.TimeWithoutTarget >= g.tuning.FireballUntargetedLifetime {
		g.removeFireball(i)
	}
}
func (g *Game) firstSkeletonHitBySegment(start, end Vec2, radius float64) int {
	delta := end.Sub(start)
	lengthSq := delta.LenSq()
	searchRadius := g.maxSkeletonCollisionRadius(radius)
	minPos := Vec2{X: math.Min(start.X, end.X) - searchRadius, Y: math.Min(start.Y, end.Y) - searchRadius}
	maxPos := Vec2{X: math.Max(start.X, end.X) + searchRadius, Y: math.Max(start.Y, end.Y) + searchRadius}
	bestIndex := -1
	bestProgress := math.Inf(1)
	g.spatial.ForEachRect(minPos, maxPos, func(i int) bool {
		progress := 0.0
		if lengthSq > 0 {
			progress = Clamp(g.skeleton[i].Pos.Sub(start).X*delta.X/lengthSq+g.skeleton[i].Pos.Sub(start).Y*delta.Y/lengthSq, 0, 1)
		}
		closest := start.Add(delta.Mul(progress))
		hitRadius := g.skeletonCollisionRadius(radius, g.skeleton[i].Kind)
		hitRadiusSq := hitRadius * hitRadius
		if DistanceSq(closest, g.skeleton[i].Pos) <= hitRadiusSq && progress < bestProgress {
			bestIndex = i
			bestProgress = progress
		}
		return true
	})
	return bestIndex
}
func (g *Game) removeFireball(i int) {
	g.fireball = slices.Delete(g.fireball, i, i+1)
}
func (g *Game) updateFireballAnimation(dt float64) {
	if len(g.fireball) == 0 {
		g.fireAnimTimer = 0
		return
	}
	g.fireAnimTimer += dt
	if g.fireAnimTimer >= g.tuning.FireballAnimationFrameTime {
		g.fireAnimTimer = math.Mod(g.fireAnimTimer, g.tuning.FireballAnimationFrameTime)
		g.fireAnimFrame = (g.fireAnimFrame + 1) % 2
		for i := range g.fireball {
			g.fireball[i].AnimFrame = g.fireAnimFrame
		}
	}
}
