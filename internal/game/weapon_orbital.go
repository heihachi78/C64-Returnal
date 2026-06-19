package game

import "math"

func (g *Game) updateOrbitalOrbs(dt float64) {
	if !g.session.Progression.OrbitalOrbUnlocked {
		return
	}
	g.syncOrbitalOrbCount()
	if len(g.orbs) == 0 {
		return
	}
	angleDelta := g.session.Progression.OrbitalAngularSpeed() * dt
	g.session.OrbitalOrbAngle += angleDelta
	g.updateOrbAnimation(dt)
	for i := range g.orbs {
		if !g.orbs[i].Active {
			g.orbs[i].MissingOrbitProgress += math.Abs(angleDelta)
			if g.orbs[i].MissingOrbitProgress >= math.Pi*2 {
				g.orbs[i].Active = true
				g.orbs[i].MissingOrbitProgress = 0
				g.orbs[i].AnimFrame = 0
			}
		}
	}
	g.alignOrbitalOrbs()
	g.checkOrbitalOrbCollisions()
}
func (g *Game) syncOrbitalOrbCount() {
	target := g.session.Progression.OrbitalOrbCount()
	for len(g.orbs) < target {
		g.orbs = append(g.orbs, OrbitalOrb{Active: true})
	}
	for len(g.orbs) > target {
		g.orbs = g.orbs[:len(g.orbs)-1]
	}
	g.alignOrbitalOrbs()
}
func (g *Game) alignOrbitalOrbs() {
	if len(g.orbs) == 0 {
		return
	}
	spacing := math.Pi * 2 / float64(len(g.orbs))
	for i := range g.orbs {
		angle := g.session.OrbitalOrbAngle + spacing*float64(i)
		g.orbs[i].Pos = Vec2{
			X: g.player.Pos.X + math.Cos(angle)*g.tuning.OrbitalOrbRadius,
			Y: g.player.Pos.Y + math.Sin(angle)*g.tuning.OrbitalOrbRadius,
		}
	}
}
func (g *Game) checkOrbitalOrbCollisions() {
	levelUps := 0
	for i := range g.orbs {
		if !g.orbs[i].Active {
			continue
		}
		idx := g.spatial.FirstNear(g.orbs[i].Pos, g.tuning.OrbitalHitDistance, g.skeleton, func(int) bool { return true })
		if idx < 0 {
			continue
		}
		g.orbs[i].Active = false
		g.orbs[i].MissingOrbitProgress = 0
		g.orbs[i].AnimFrame = 0
		levelUps += g.damageSkeleton(idx, 1, AttackOrbitalOrb, false)
	}
	g.queueLevelUpChoices(levelUps)
}
func (g *Game) updateOrbAnimation(dt float64) {
	hasActive := false
	for i := range g.orbs {
		if g.orbs[i].Active {
			hasActive = true
			break
		}
	}
	if !hasActive {
		g.orbAnimTimer = 0
		return
	}
	g.orbAnimTimer += dt
	if g.orbAnimTimer >= g.tuning.OrbitalAnimationFrameTime {
		g.orbAnimTimer = math.Mod(g.orbAnimTimer, g.tuning.OrbitalAnimationFrameTime)
		g.orbAnimFrame = (g.orbAnimFrame + 1) % 2
		for i := range g.orbs {
			if g.orbs[i].Active {
				g.orbs[i].AnimFrame = g.orbAnimFrame
			}
		}
	}
}
