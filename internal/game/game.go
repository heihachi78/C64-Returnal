package game

import (
	"math"
	"math/rand"
	"runtime"
	"slices"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	tuning Tuning
	rng    *rand.Rand
	assets *Assets

	screenW int
	screenH int
	nextID  int

	player   Player
	session  Session
	spatial  SpatialIndex
	skeleton []Skeleton
	fireball []Fireball
	orbs     []OrbitalOrb
	meteors  []MeteorProjectile
	chests   []Chest
	coins    []Coin
	effects  []Effect

	skeletonAnimTimer  float64
	skeletonAnimFrame  int
	fireAnimTimer      float64
	fireAnimFrame      int
	orbAnimTimer       float64
	orbAnimFrame       int
	meteorAnimTimer    float64
	meteorAnimFrame    int
	totalTime          float64
	hasUpdated         bool
	lastParallelJobs   int
	suppressedMovement map[ebiten.Key]bool
	scaledTextCache    map[scaledTextCacheKey]scaledTextCacheEntry
}

func New() *Game {
	tuning := DefaultTuning()
	g := &Game{
		tuning:             tuning,
		rng:                rand.New(rand.NewSource(rand.Int63())),
		assets:             NewAssets(int(tuning.TileSize)),
		screenW:            ScreenWidth,
		screenH:            ScreenHeight,
		spatial:            NewSpatialIndex(tuning.SpatialIndexCellSize),
		suppressedMovement: map[ebiten.Key]bool{},
		scaledTextCache:    map[scaledTextCacheKey]scaledTextCacheEntry{},
	}
	g.reset()
	return g
}

func (g *Game) reset() {
	g.nextID = 1
	g.player = Player{Facing: 1}
	g.session = NewSession(g.tuning)
	g.skeleton = g.skeleton[:0]
	g.fireball = g.fireball[:0]
	g.orbs = g.orbs[:0]
	g.meteors = g.meteors[:0]
	g.chests = g.chests[:0]
	g.coins = g.coins[:0]
	g.effects = g.effects[:0]
	g.spatial.Rebuild(g.skeleton)
	g.skeletonAnimTimer = 0
	g.skeletonAnimFrame = 0
	g.fireAnimTimer = 0
	g.fireAnimFrame = 0
	g.orbAnimTimer = 0
	g.orbAnimFrame = 0
	g.meteorAnimTimer = 0
	g.meteorAnimFrame = 0
	g.totalTime = 0
	g.hasUpdated = false
	g.lastParallelJobs = 0
	clear(g.suppressedMovement)
	g.spawnSkeleton(SkeletonRegular)
	g.spawnCoinForLevel(g.session.Progression.Level)
}

func (g *Game) Update() error {
	dt := 0.0
	if g.hasUpdated {
		dt = 1.0 / float64(TargetTPS)
	}
	g.hasUpdated = true
	g.totalTime += dt
	consumedFrame, err := g.updateOverlayInput()
	if err != nil {
		return err
	}
	if consumedFrame {
		return nil
	}
	g.updatePausedAnimations(dt)

	if g.session.GameOver || g.session.LevelUpChoiceActive || g.session.ChestRewardActive {
		return nil
	}

	g.updatePlayer(dt)
	g.checkCoinPickups()
	g.checkChestPickups()
	if g.session.ChestRewardActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateSkeletons(dt)
	g.updateOrbitalOrbs(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateLightningCasting(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateFireballCasting(dt)
	g.updateFireballs(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateBeamCasting(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateMeteorCasting(dt)
	g.updateMeteors(dt)
	g.updateInvulnerability(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.checkSkeletonCollisions()
	g.updateSkeletonSpawning(dt)
	if g.session.GameOver {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}
	g.updatePlayerWalkAnimation(dt)
	g.updatePlayerHitFlash(dt)
	g.updateSkeletonHitFlashes(dt)
	g.updateCoins(dt)
	g.updateEffects(dt)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenW = max(1, outsideWidth)
	g.screenH = max(1, outsideHeight)
	return g.screenW, g.screenH
}

func (g *Game) updateOverlayInput() (bool, error) {
	if g.session.GameOver {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			return g.selectGameOverOption(g.gameOverOptionAt(float64(x), float64(y)))
		}
		return false, nil
	}

	if g.session.ChestRewardActive {
		g.suppressModalHeldMovementKeys(ebiten.IsKeyPressed)
		if inpututil.IsKeyJustPressed(chestRewardAdvanceKey()) {
			return g.advanceChestReward(), nil
		}
		return false, nil
	}

	if !g.session.LevelUpChoiceActive {
		if isKillAllAndGrantExperienceJustPressed() {
			return g.handleKillAllAndGrantExperienceKeyDown(), nil
		}
		return false, nil
	}

	g.suppressModalHeldMovementKeys(ebiten.IsKeyPressed)
	if inpututil.IsKeyJustPressed(levelUpRedrawKey()) {
		g.redrawLevelUpOptions()
		return true, nil
	}

	for i, key := range levelUpOptionKeys() {
		if inpututil.IsKeyJustPressed(key) && i < len(g.session.ActiveLevelUpOptions) {
			return g.selectLevelUpOptionAt(i), nil
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.redrawRectContains(float64(x), float64(y)) {
			g.redrawLevelUpOptions()
			return true, nil
		}
		if idx := g.levelUpOptionAt(float64(x), float64(y)); idx >= 0 && idx < len(g.session.ActiveLevelUpOptions) {
			return g.selectLevelUpOptionAt(idx), nil
		}
	}
	return false, nil
}

func (g *Game) selectGameOverOption(option string) (bool, error) {
	switch option {
	case "restart":
		g.restartGame(ebiten.IsKeyPressed)
		return true, nil
	case "exit":
		return true, ebiten.Termination
	default:
		return false, nil
	}
}

func (g *Game) restartGame(isPressed func(ebiten.Key) bool) {
	g.reset()
	g.suppressHeldMovementKeys(isPressed)
}

func (g *Game) advanceChestReward() bool {
	if !g.session.ChestRewardActive {
		return false
	}
	g.session.ChestRewardActive = false
	g.session.ActiveChestRewardItems = nil
	g.session.ChestRewardOverlayTimer = 0
	g.presentNextLevelUpChoiceIfNeeded()
	return true
}

func (g *Game) selectLevelUpOptionAt(index int) bool {
	if !g.session.LevelUpChoiceActive || index < 0 || index >= len(g.session.ActiveLevelUpOptions) {
		return false
	}
	g.applyLevelUpOption(g.session.ActiveLevelUpOptions[index])
	return true
}

func levelUpOptionKeys() []ebiten.Key {
	return []ebiten.Key{ebiten.KeyQ, ebiten.KeyA, ebiten.KeyC, ebiten.KeyX}
}

func levelUpRedrawKey() ebiten.Key {
	return ebiten.KeyR
}

func chestRewardAdvanceKey() ebiten.Key {
	return ebiten.KeyQ
}

func isKillAllAndGrantExperienceKey(key ebiten.Key) bool {
	return key == ebiten.KeyDigit1 || key == ebiten.KeyNumpad1
}

func isKillAllAndGrantExperienceJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDigit1) || inpututil.IsKeyJustPressed(ebiten.KeyNumpad1)
}

func (g *Game) updatePlayer(dt float64) {
	move := g.playerMovementVector(ebiten.IsKeyPressed)

	if move.X != 0 || move.Y != 0 {
		move = move.Normalized()
		g.player.MoveDir = move
		g.player.Moving = true
		g.player.Pos = g.player.Pos.Add(move.Mul(g.tuning.PlayerSpeed * dt))
		if move.X < 0 {
			g.player.Facing = -1
		} else if move.X > 0 {
			g.player.Facing = 1
		}
		return
	}

	g.player.Moving = false
	g.player.MoveDir = Vec2{}
	g.player.AnimTimer = 0
	g.player.AnimFrame = 0
}

func (g *Game) playerMovementVector(isPressed func(ebiten.Key) bool) Vec2 {
	for _, key := range movementKeys() {
		if g.suppressedMovement[key] && !isPressed(key) {
			delete(g.suppressedMovement, key)
		}
	}

	move := Vec2{}
	if isPressed(ebiten.KeyArrowLeft) && !g.suppressedMovement[ebiten.KeyArrowLeft] {
		move.X--
	}
	if isPressed(ebiten.KeyArrowRight) && !g.suppressedMovement[ebiten.KeyArrowRight] {
		move.X++
	}
	if isPressed(ebiten.KeyArrowUp) && !g.suppressedMovement[ebiten.KeyArrowUp] {
		move.Y++
	}
	if isPressed(ebiten.KeyArrowDown) && !g.suppressedMovement[ebiten.KeyArrowDown] {
		move.Y--
	}
	return move
}

func (g *Game) suppressHeldMovementKeys(isPressed func(ebiten.Key) bool) {
	if g.suppressedMovement == nil {
		g.suppressedMovement = map[ebiten.Key]bool{}
	}
	for _, key := range movementKeys() {
		if isPressed(key) {
			g.suppressedMovement[key] = true
		}
	}
}

func (g *Game) suppressModalHeldMovementKeys(isPressed func(ebiten.Key) bool) {
	if g.session.LevelUpChoiceActive || g.session.ChestRewardActive {
		g.suppressHeldMovementKeys(isPressed)
	}
}

func movementKeys() []ebiten.Key {
	return []ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyArrowUp, ebiten.KeyArrowDown}
}

func (g *Game) updateInvulnerability(dt float64) {
	if g.session.PlayerHitInvulnerability > 0 {
		g.session.PlayerHitInvulnerability = math.Max(0, g.session.PlayerHitInvulnerability-dt)
	}
}

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
	interval := g.session.Progression.SkeletonSpawnInterval()
	for g.session.Casts.SkeletonSpawn >= interval {
		g.session.Casts.SkeletonSpawn -= interval
		g.spawnSkeleton(g.timedSkeletonSpawnKind())
	}
}

func (g *Game) timedSkeletonSpawnKind() SkeletonKind {
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

func (g *Game) addSkeleton(kind SkeletonKind) {
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

func (g *Game) damagePlayer() {
	g.session.PlayerLives = max(0, g.session.PlayerLives-1)
	if g.session.PlayerLives == 0 {
		g.triggerGameOver()
		return
	}
	g.session.PlayerHitInvulnerability = g.tuning.PlayerHitInvulnerability
	g.player.HitFlash = playerHitFlashDuration
}

func (g *Game) triggerGameOver() {
	if g.session.GameOver {
		return
	}
	g.session.GameOver = true
	g.session.LevelUpChoiceActive = false
	g.session.ChestRewardActive = false
	g.session.PlayerHitInvulnerability = 0
	g.session.PendingLevelUpLevels = g.session.PendingLevelUpLevels[:0]
	g.session.ActiveLevelUpOptions = nil
	g.session.ActiveChestRewardItems = nil
	clear(g.suppressedMovement)
	g.hideLevelUpPresentation()
	g.session.ChestRewardOverlayTimer = 0
	g.player.Moving = false
	g.player.AnimTimer = 0
	g.player.AnimFrame = 0
	g.player.HitFlash = 0
	g.player.DeathTimer = 0
	g.player.DeathRotation = 0
	g.session.GameOverOverlayTimer = 0
	g.effects = g.effects[:0]
	g.meteors = g.meteors[:0]
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
	if distanceSq == 0 || distanceSq <= (g.tuning.FireballHitDistance+travel)*(g.tuning.FireballHitDistance+travel) {
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
	minPos := Vec2{X: math.Min(start.X, end.X) - radius, Y: math.Min(start.Y, end.Y) - radius}
	maxPos := Vec2{X: math.Max(start.X, end.X) + radius, Y: math.Max(start.Y, end.Y) + radius}
	hitRadiusSq := radius * radius
	bestIndex := -1
	bestProgress := math.Inf(1)
	g.spatial.ForEachRect(minPos, maxPos, func(i int) bool {
		progress := 0.0
		if lengthSq > 0 {
			progress = Clamp(g.skeleton[i].Pos.Sub(start).X*delta.X/lengthSq+g.skeleton[i].Pos.Sub(start).Y*delta.Y/lengthSq, 0, 1)
		}
		closest := start.Add(delta.Mul(progress))
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
	g.queueLevelUpChoices(levelUps)
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
		}
	}
	g.queueLevelUpChoices(levelUps)
}

func (g *Game) meteorImpactTargetIDs(pos Vec2) []int {
	targets := []int{}
	radiusSq := g.tuning.MeteorImpactRadius * g.tuning.MeteorImpactRadius
	g.spatial.ForEachNear(pos, g.tuning.MeteorImpactRadius, g.skeleton, func(i int) bool {
		if DistanceSq(pos, g.skeleton[i].Pos) <= radiusSq {
			targets = append(targets, g.skeleton[i].ID)
		}
		return true
	})
	return targets
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

func (g *Game) spawnChest(tier ChestTier) {
	g.chests = append(g.chests, Chest{Pos: g.randomChestPosition(), Tier: tier})
}

func (g *Game) randomChestPosition() Vec2 {
	halfW := math.Max(48, float64(g.screenW)/2-g.tuning.ChestSpawnMargin)
	halfH := math.Max(48, float64(g.screenH)/2-g.tuning.ChestSpawnMargin)
	minDistSq := g.tuning.ChestPickupDistance * g.tuning.ChestPickupDistance * 4
	for range 12 {
		pos := Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
		if DistanceSq(pos, g.player.Pos) >= minDistSq {
			return pos
		}
	}
	return Vec2{X: g.player.Pos.X + halfW, Y: g.player.Pos.Y}
}

func (g *Game) checkChestPickups() {
	distSq := g.tuning.ChestPickupDistance * g.tuning.ChestPickupDistance
	for i := len(g.chests) - 1; i >= 0; i-- {
		if DistanceSq(g.chests[i].Pos, g.player.Pos) <= distSq {
			chest := g.chests[i]
			g.chests = slices.Delete(g.chests, i, i+1)
			g.applyChestReward(chest.Tier)
			return
		}
	}
}

func (g *Game) applyChestReward(tier ChestTier) {
	skills := g.session.Progression.LearnedSkills()
	if len(skills) == 0 {
		return
	}
	items := []ChestRewardDisplayItem{}
	switch tier {
	case ChestBronze:
		skill := skills[g.rng.Intn(len(skills))]
		options := skill.UpgradeOptions()
		option := options[g.rng.Intn(len(options))]
		items = append(items, g.chestRewardItemForSkill(option, skill))
		g.applyUpgradeEffect(option)
	case ChestSilver:
		skill := skills[g.rng.Intn(len(skills))]
		items = append(items, g.chestRewardItemsForSkill(skill)...)
		g.session.Progression.UpgradeAllProperties(skill)
	case ChestGold:
		g.rng.Shuffle(len(skills), func(i, j int) { skills[i], skills[j] = skills[j], skills[i] })
		for _, skill := range skills[:min(2, len(skills))] {
			items = append(items, g.chestRewardItemsForSkill(skill)...)
			g.session.Progression.UpgradeAllProperties(skill)
		}
	}
	g.syncOrbitalOrbCount()
	if len(items) > 0 {
		g.session.ChestRewardActive = true
		g.session.ActiveChestTier = tier
		g.session.ActiveChestRewardItems = items
		g.session.ChestRewardOverlayTimer = 0
		g.suppressHeldMovementKeys(ebiten.IsKeyPressed)
		g.stopPlayerAnimation()
	}
}

func (g *Game) chestRewardItemsForSkill(skill LearnedSkill) []ChestRewardDisplayItem {
	options := skill.UpgradeOptions()
	items := make([]ChestRewardDisplayItem, 0, len(options))
	for _, option := range options {
		items = append(items, g.chestRewardItemForSkill(option, skill))
	}
	return items
}

func (g *Game) chestRewardItemForSkill(option LevelUpOption, skill LearnedSkill) ChestRewardDisplayItem {
	beamKillBonus := 0
	if skill == SkillBeam {
		beamKillBonus = g.session.Progression.BeamKillUpgradeBonus()
	}
	return ChestRewardDisplayItem{
		Option: option,
		Title:  option.Title(beamKillBonus),
	}
}

func (g *Game) queueLevelUpChoices(count int) {
	if count <= 0 || g.session.GameOver {
		return
	}
	first := g.session.Progression.Level - count + 1
	for level := first; level <= g.session.Progression.Level; level++ {
		g.session.PendingLevelUpLevels = append(g.session.PendingLevelUpLevels, level)
		g.spawnCoinForLevel(level)
	}
	g.presentNextLevelUpChoiceIfNeeded()
}

func (g *Game) presentNextLevelUpChoiceIfNeeded() {
	if g.session.GameOver || g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) == 0 {
		return
	}
	g.session.LevelUpChoiceActive = true
	g.session.CurrentLevelUpPresentation = g.session.PendingLevelUpLevels[0]
	g.session.ActiveLevelUpOptions = visibleLevelUpOptions(g.randomLevelUpOptions(nil))
	g.session.LevelUpOverlayTimer = 0
	g.session.LevelUpTitleScaleTimer = 0
	g.session.LevelUpOptionFadeTimer = 0
	g.showLevelUpRedrawPresentation(true)
	g.suppressHeldMovementKeys(ebiten.IsKeyPressed)
	g.stopPlayerAnimation()
}

func (g *Game) randomLevelUpOptions(excluding []LevelUpOption) []LevelUpOption {
	selected := g.randomLevelUpOptionsCandidate()
	if len(excluding) > 0 {
		for tries := 0; tries < 8 && sameOptionSet(selected, excluding); tries++ {
			selected = g.randomLevelUpOptionsCandidate()
		}
	}
	return selected
}

func (g *Game) randomLevelUpOptionsCandidate() []LevelUpOption {
	hasSkeletons := len(g.skeleton) > 0
	count := 2
	if chance(g.rng, g.tuning.ExtraOptionChanceNumerator, g.tuning.ExtraOptionChanceDenominator) {
		count = 3
	}
	available := slices.Clone(g.session.Progression.AvailableLevelUpOptions())
	available = slices.DeleteFunc(available, func(o LevelUpOption) bool { return o == HalveSkeletons })
	g.rng.Shuffle(len(available), func(i, j int) { available[i], available[j] = available[j], available[i] })
	selected := slices.Clone(available[:min(count, len(available))])
	if hasSkeletons && len(selected) > 0 && chance(g.rng, g.tuning.HalveHordeChanceNumerator, g.tuning.HalveHordeChanceDenominator) {
		selected[g.rng.Intn(len(selected))] = HalveSkeletons
	}
	return selected
}

func (g *Game) applyLevelUpOption(option LevelUpOption) {
	g.applyUpgradeEffect(option)
	g.syncOrbitalOrbCount()
	if len(g.session.PendingLevelUpLevels) > 0 {
		g.session.PendingLevelUpLevels = g.session.PendingLevelUpLevels[1:]
	}
	g.session.LevelUpChoiceActive = false
	g.session.ActiveLevelUpOptions = nil
	g.hideLevelUpPresentation()
	g.presentNextLevelUpChoiceIfNeeded()
}

func (g *Game) applyUpgradeEffect(option LevelUpOption) {
	switch option {
	case ExtraLife:
		g.session.PlayerLives++
	case HalveSkeletons:
		g.halveSkeletons()
	default:
		g.session.Progression.ApplyLevelUpOption(option)
	}
}

func (g *Game) halveSkeletons() {
	killCount := len(g.skeleton) / 2
	if killCount <= 0 {
		return
	}
	targetIDs := make([]int, len(g.skeleton))
	for i, skeleton := range g.skeleton {
		targetIDs[i] = skeleton.ID
	}
	g.rng.Shuffle(len(targetIDs), func(i, j int) { targetIDs[i], targetIDs[j] = targetIDs[j], targetIDs[i] })
	levelUps := 0
	for _, id := range targetIDs[:killCount] {
		if idx := g.skeletonIndexByID(id); idx >= 0 {
			levelUps += g.destroySkeleton(idx, AttackNone)
		}
	}
	g.queueLevelUpChoices(levelUps)
}

func (g *Game) redrawLevelUpOptions() {
	if !g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) == 0 {
		return
	}
	cost := g.levelUpRedrawCost()
	if g.session.CollectedCoins < cost {
		g.session.LevelUpRedrawStatusTimer = redrawFailurePulseDuration
		g.session.LevelUpRedrawCoinFadeTimer = 0
		return
	}
	g.session.CollectedCoins -= cost
	previous := slices.Clone(g.session.ActiveLevelUpOptions)
	g.session.ActiveLevelUpOptions = visibleLevelUpOptions(g.randomLevelUpOptions(previous))
	g.session.LevelUpTitleScaleTimer = 0
	g.session.LevelUpOptionFadeTimer = 0
	g.showLevelUpRedrawPresentation(false)
}

func visibleLevelUpOptions(options []LevelUpOption) []LevelUpOption {
	return slices.Clone(options[:min(4, len(options))])
}

func (g *Game) showLevelUpRedrawPresentation(resetTextFade bool) {
	g.session.LevelUpRedrawStatusTimer = 0
	if g.session.CollectedCoins < g.levelUpRedrawCost() {
		g.session.LevelUpRedrawStatusTimer = redrawFailurePulseDuration
		g.session.LevelUpRedrawCoinFadeTimer = 0
	} else {
		g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration
	}
	if resetTextFade {
		g.session.LevelUpRedrawFadeTimer = 0
	}
}

func (g *Game) levelUpRedrawCost() int {
	return max(1, g.session.Progression.Level)
}

func (g *Game) hideLevelUpPresentation() {
	g.session.LevelUpRedrawStatusTimer = 0
	g.session.LevelUpRedrawFadeTimer = 0
	g.session.LevelUpRedrawCoinFadeTimer = 0
	g.session.LevelUpOverlayTimer = 0
	g.session.LevelUpTitleScaleTimer = 0
	g.session.LevelUpOptionFadeTimer = 0
}

func (g *Game) updateEffects(dt float64) {
	for i := len(g.effects) - 1; i >= 0; i-- {
		g.effects[i].TTL -= dt
		if g.effects[i].TTL <= 0 {
			g.effects = slices.Delete(g.effects, i, i+1)
		}
	}
}

func (g *Game) stopPlayerAnimation() {
	g.player.Moving = false
	g.player.AnimTimer = 0
	g.player.AnimFrame = 0
}

func (g *Game) updatePlayerWalkAnimation(dt float64) {
	if !g.player.Moving {
		return
	}
	g.player.AnimTimer += dt
	if g.player.AnimTimer >= g.tuning.PlayerAnimationFrameTime {
		g.player.AnimTimer = math.Mod(g.player.AnimTimer, g.tuning.PlayerAnimationFrameTime)
		g.player.AnimFrame = (g.player.AnimFrame + 1) % 2
	}
}

func (g *Game) updatePlayerHitFlash(dt float64) {
	if g.player.HitFlash > 0 {
		g.player.HitFlash = math.Max(0, g.player.HitFlash-dt)
	}
}

func (g *Game) updatePausedAnimations(dt float64) {
	if g.session.LevelUpChoiceActive {
		g.session.LevelUpOverlayTimer += dt
		g.session.LevelUpTitleScaleTimer += dt
		g.session.LevelUpOptionFadeTimer += dt
		g.session.LevelUpRedrawFadeTimer += dt
		g.session.LevelUpRedrawCoinFadeTimer += dt
		if g.session.LevelUpRedrawStatusTimer > 0 {
			g.session.LevelUpRedrawStatusTimer = math.Max(0, g.session.LevelUpRedrawStatusTimer-dt)
		}
	}
	if g.session.ChestRewardActive {
		g.session.ChestRewardOverlayTimer += dt
	}
	if g.session.GameOver {
		g.session.GameOverOverlayTimer += dt
		g.player.DeathTimer = math.Min(playerDeathRotationDuration, g.player.DeathTimer+dt)
		progress := g.player.DeathTimer / playerDeathRotationDuration
		g.player.DeathRotation = -math.Pi / 2 * progress
		g.updateGameOverWorldActions(dt)
	}
}

func (g *Game) updateNewlyPresentedOverlayActions(dt float64) {
	if g.session.LevelUpChoiceActive {
		g.session.LevelUpOverlayTimer += dt
		g.session.LevelUpTitleScaleTimer += dt
		g.session.LevelUpOptionFadeTimer += dt
		g.session.LevelUpRedrawFadeTimer += dt
		g.session.LevelUpRedrawCoinFadeTimer += dt
		if g.session.LevelUpRedrawStatusTimer > 0 {
			g.session.LevelUpRedrawStatusTimer = math.Max(0, g.session.LevelUpRedrawStatusTimer-dt)
		}
	}
	if g.session.ChestRewardActive {
		g.session.ChestRewardOverlayTimer += dt
	}
	if g.session.GameOver {
		g.session.GameOverOverlayTimer += dt
		g.player.DeathTimer = math.Min(playerDeathRotationDuration, g.player.DeathTimer+dt)
		progress := g.player.DeathTimer / playerDeathRotationDuration
		g.player.DeathRotation = -math.Pi / 2 * progress
		g.updateCoins(dt)
		g.updateSkeletonHitFlashes(dt)
	}
}

func (g *Game) updateGameOverWorldActions(dt float64) {
	g.updateCoins(dt)
	g.updateSkeletonHitFlashes(dt)
}

func (g *Game) updateSkeletonHitFlashes(dt float64) {
	for i := range g.skeleton {
		if g.skeleton[i].HitFlash > 0 {
			g.skeleton[i].HitFlash = math.Max(0, g.skeleton[i].HitFlash-dt)
		}
	}
}

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
