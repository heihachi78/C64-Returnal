package game

import (
	"math"
	"runtime"
	"sync"
)

func (g *Game) updateSkeletons(dt float64) {
	g.lastParallelJobs = 1

	if len(g.skeleton) < g.tuning.ParallelSkeletonUpdateThreshold {
		g.updateSkeletonRange(0, len(g.skeleton), dt, g.player.Pos)
		g.updateSkeletonAnimation(dt)
		g.spatial.Rebuild(g.skeleton)
		return
	}

	jobs := min(runtime.GOMAXPROCS(0), len(g.skeleton))
	chunk := (len(g.skeleton) + jobs - 1) / jobs
	playerPos := g.player.Pos
	launchedJobs := 0
	var wg sync.WaitGroup
	for job := 0; job < jobs; job++ {
		start := job * chunk
		end := min(len(g.skeleton), start+chunk)
		if start >= end {
			continue
		}
		launchedJobs++
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			g.updateSkeletonRange(start, end, dt, playerPos)
		}(start, end)
	}
	wg.Wait()
	g.lastParallelJobs = launchedJobs
	g.updateSkeletonAnimation(dt)
	g.spatial.Rebuild(g.skeleton)
}
func (g *Game) updateSkeletonRange(start, end int, dt float64, playerPos Vec2) {
	for i := start; i < end; i++ {
		toPlayer := playerPos.Sub(g.skeleton[i].Pos)
		if toPlayer.X == 0 && toPlayer.Y == 0 {
			continue
		}
		move := toPlayer.Normalized()
		g.skeleton[i].Pos = g.skeleton[i].Pos.Add(move.Mul(g.tuning.SkeletonSpeed * dt))
		if move.X < 0 {
			g.skeleton[i].Facing = -1
		} else if move.X > 0 {
			g.skeleton[i].Facing = 1
		}
	}
}
func (g *Game) updateSkeletonAnimation(dt float64) {
	if len(g.skeleton) == 0 {
		g.skeletonAnimTimer = 0
		return
	}
	g.skeletonAnimTimer += dt
	if g.skeletonAnimTimer < g.tuning.SkeletonAnimationFrameTime {
		return
	}
	g.skeletonAnimTimer = math.Mod(g.skeletonAnimTimer, g.tuning.SkeletonAnimationFrameTime)
	g.skeletonAnimFrame = (g.skeletonAnimFrame + 1) % 2
	for i := range g.skeleton {
		g.skeleton[i].AnimFrame = g.skeletonAnimFrame
	}
}
func (g *Game) updateSkeletonSpawning(dt float64) {
	g.session.Casts.SkeletonSpawn += dt
	for {
		interval := g.session.Progression.SkeletonSpawnInterval()
		if g.session.Casts.SkeletonSpawn < interval {
			break
		}
		g.session.Casts.SkeletonSpawn -= interval
		g.spawnTimedSkeleton()
	}
}
func (g *Game) spawnTimedSkeleton() {
	if g.shouldSpawnBlueMonster() {
		g.spawnBlueMonster()
		return
	}
	g.spawnSkeleton(g.timedSkeletonSpawnKind())
}
func (g *Game) timedSkeletonSpawnKind() SkeletonKind {
	if g.session.Progression.Level >= g.tuning.BlackOnlyLevel {
		return SkeletonBlack
	}
	if g.session.Progression.Level >= g.tuning.PurpleOnlyLevel {
		return SkeletonPurple
	}
	if g.session.Progression.Level >= g.tuning.RedOnlyLevel {
		return SkeletonRed
	}
	return SkeletonRegular
}
func (g *Game) spawnSkeleton(kind SkeletonKind) {
	g.addSkeleton(kind)
}
func (g *Game) shouldSpawnBlueMonster() bool {
	return g.session.Progression.Level >= g.tuning.BlueMonsterMinimumLevel &&
		g.tuning.BlueMonsterMinimumEnemies > 0 &&
		len(g.skeleton) > g.tuning.BlueMonsterMinimumEnemies
}
func (g *Game) spawnBlueMonster() {
	blue := g.addSkeleton(SkeletonBlue)
	g.session.Progression.SlowSkeletonSpawnRate()
	g.cullHalfEnemiesForBlueMonster(blue.ID)
}
func (g *Game) addSkeleton(kind SkeletonKind) Skeleton {
	s := Skeleton{
		ID:     g.nextID,
		Pos:    g.skeletonSpawnPosition(),
		Kind:   kind,
		HP:     kind.HitPoints(g.tuning),
		Reward: kind.ExperienceReward(),
		Facing: 1,
	}
	g.nextID++
	g.skeleton = append(g.skeleton, s)
	return s
}
func (g *Game) cullHalfEnemiesForBlueMonster(blueID int) {
	targetIDs := make([]int, 0, len(g.skeleton))
	for _, skeleton := range g.skeleton {
		if skeleton.ID != blueID {
			targetIDs = append(targetIDs, skeleton.ID)
		}
	}
	cullCount := len(targetIDs) / 2
	if cullCount <= 0 {
		return
	}
	g.rng.Shuffle(len(targetIDs), func(i, j int) { targetIDs[i], targetIDs[j] = targetIDs[j], targetIDs[i] })
	cull := make(map[int]bool, cullCount)
	for _, id := range targetIDs[:cullCount] {
		cull[id] = true
	}
	remaining := g.skeleton[:0]
	for _, skeleton := range g.skeleton {
		if !cull[skeleton.ID] {
			remaining = append(remaining, skeleton)
		}
	}
	g.skeleton = remaining
	g.spatial.Rebuild(g.skeleton)
}
func (g *Game) skeletonSpawnPosition() Vec2 {
	halfW := float64(g.screenW) / 2
	halfH := float64(g.screenH) / 2
	spawnDistance := math.Hypot(halfW, halfH) + g.tuning.SkeletonSpawnMargin
	target := Vec2{}

	switch g.rng.Intn(4) {
	case 0:
		target = Vec2{X: g.player.Pos.X - halfW, Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
	case 1:
		target = Vec2{X: g.player.Pos.X + halfW, Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
	case 2:
		target = Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y - halfH}
	default:
		target = Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y + halfH}
	}

	return g.player.Pos.Add(target.Sub(g.player.Pos).Normalized().Mul(spawnDistance))
}
func (g *Game) checkSkeletonCollisions() {
	if g.session.PlayerHitInvulnerability > 0 {
		return
	}
	idx := g.spatial.FirstNear(g.player.Pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true })
	if idx >= 0 {
		g.damagePlayer()
	}
}
func (g *Game) damageSkeleton(index, amount int, attack AttackKind, queueLevelUp bool) int {
	if index < 0 || index >= len(g.skeleton) || amount <= 0 {
		return 0
	}
	if g.skeleton[index].HP > amount {
		g.skeleton[index].HP -= amount
		g.skeleton[index].HitFlash = skeletonDamageFlashDuration
		return 0
	}
	levelUps := g.destroySkeleton(index, attack)
	if queueLevelUp {
		g.queueLevelUpChoices(levelUps)
	}
	return levelUps
}
func (g *Game) destroySkeleton(index int, attack AttackKind) int {
	if index < 0 || index >= len(g.skeleton) {
		return 0
	}
	kind := g.skeleton[index].Kind
	reward := g.skeleton[index].Reward
	last := len(g.skeleton) - 1
	g.skeleton[index] = g.skeleton[last]
	g.skeleton = g.skeleton[:last]
	g.spatial.Rebuild(g.skeleton)
	g.session.Kills.TotalSkeletons++
	g.spawnBlackSkeletonIfNeeded(kind)
	g.spawnMilestoneSkeletonsIfNeeded()
	g.spawnChestsForMilestones()
	g.session.RegisterAttackKill(attack)
	levelUps := g.session.Progression.GainExperience(reward)
	return levelUps
}
func (g *Game) spawnBlackSkeletonIfNeeded(kind SkeletonKind) {
	if kind != SkeletonPurple || g.tuning.BlackPurpleKillInterval <= 0 {
		return
	}
	g.session.Kills.PurpleSkeletons++
	if g.session.Kills.PurpleSkeletons%g.tuning.BlackPurpleKillInterval == 0 {
		g.addSkeleton(SkeletonBlack)
	}
}
func (g *Game) spawnMilestoneSkeletonsIfNeeded() {
	if g.session.Progression.Level >= g.tuning.RedOnlyLevel {
		return
	}
	if g.tuning.RedKillInterval > 0 && g.session.Kills.TotalSkeletons%g.tuning.RedKillInterval == 0 {
		g.addSkeleton(SkeletonRed)
	}
	if g.tuning.PurpleKillInterval > 0 && g.session.Kills.TotalSkeletons%g.tuning.PurpleKillInterval == 0 {
		g.addSkeleton(SkeletonPurple)
	}
}
func (g *Game) spawnChestsForMilestones() {
	for g.session.Kills.TotalSkeletons >= g.session.NextChestMilestone {
		if tier, ok := chestTier(g.tuning, g.session.NextChestMilestone, g.session.Progression.Level); ok {
			g.spawnChest(tier)
		}
		g.session.NextChestMilestone += g.tuning.BronzeKillInterval
	}
}
func chestTier(t Tuning, milestone, playerLevel int) (ChestTier, bool) {
	if milestone%t.GoldKillInterval == 0 {
		return ChestGold, true
	}
	if milestone%t.SilverKillInterval == 0 {
		return ChestSilver, playerLevel <= t.SilverMaximumLevel
	}
	return ChestBronze, playerLevel <= t.BronzeMaximumLevel
}
func (g *Game) killAllEnemiesAndGrantExperience() bool {
	if len(g.skeleton) == 0 {
		return false
	}
	reward := 0
	for _, skeleton := range g.skeleton {
		reward += skeleton.Reward
	}
	g.skeleton = g.skeleton[:0]
	g.spatial.Rebuild(g.skeleton)
	levelUps := g.session.Progression.GainExperience(reward)
	g.queueLevelUpChoices(levelUps)
	return true
}
func (g *Game) handleKillAllAndGrantExperienceKeyDown() bool {
	g.killAllEnemiesAndGrantExperience()
	return false
}
