package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"slices"
	"testing"
)

func TestFourthLevelUpKeyMatchesOriginalXBinding(t *testing.T) {
	keys := levelUpOptionKeys()
	if got := keys[3]; got != ebiten.KeyX {
		t.Fatalf("fourth level-up key = %v, want %v", got, ebiten.KeyX)
	}
}

func TestModalAndDebugKeysMatchOriginalInputBindings(t *testing.T) {
	keys := levelUpOptionKeys()
	wantLevelUp := []ebiten.Key{ebiten.KeyQ, ebiten.KeyA, ebiten.KeyC, ebiten.KeyX}
	if !slices.Equal(keys, wantLevelUp) {
		t.Fatalf("level-up keys = %v, want %v", keys, wantLevelUp)
	}
	if got := levelUpRedrawKey(); got != ebiten.KeyR {
		t.Fatalf("redraw key = %v, want %v", got, ebiten.KeyR)
	}
	if got := chestRewardAdvanceKey(); got != ebiten.KeyQ {
		t.Fatalf("chest reward advance key = %v, want %v", got, ebiten.KeyQ)
	}
	if !isJumpToLevel100DebugKey(ebiten.KeyDigit0) || !isJumpToLevel100DebugKey(ebiten.KeyNumpad0) {
		t.Fatalf("level-100 debug keys did not include main and numpad 0")
	}
	if isJumpToLevel100DebugKey(ebiten.KeyDigit2) {
		t.Fatalf("level-100 debug accepted digit 2, want only 0 bindings")
	}
}

func TestFirstUpdateEstablishesTimeWithoutAdvancingSimulationLikeOriginal(t *testing.T) {
	g := New()
	g.session.Casts.SkeletonSpawn = 0

	if err := g.Update(); err != nil {
		t.Fatalf("first update returned error: %v", err)
	}

	if g.totalTime != 0 {
		t.Fatalf("totalTime after first update = %v, want 0", g.totalTime)
	}
	if g.session.Casts.SkeletonSpawn != 0 {
		t.Fatalf("skeleton spawn timer after first update = %v, want 0", g.session.Casts.SkeletonSpawn)
	}
}

func TestOverlayEventActionsConsumeTheCurrentUpdateFrameLikeOriginal(t *testing.T) {
	g := New()
	g.session.GameOver = true
	g.totalTime = 10
	g.hasUpdated = true
	g.session.Casts.SkeletonSpawn = 5

	consumed, err := g.selectGameOverOption("restart")
	if err != nil {
		t.Fatalf("restart returned error: %v", err)
	}
	if !consumed {
		t.Fatal("restart consumed frame = false, want true")
	}
	if g.session.GameOver || g.totalTime != 0 || g.hasUpdated || g.session.Casts.SkeletonSpawn != 0 {
		t.Fatalf("restart state = gameOver %v totalTime %v hasUpdated %v spawnTimer %v; want fresh reset", g.session.GameOver, g.totalTime, g.hasUpdated, g.session.Casts.SkeletonSpawn)
	}

	g.session.ChestRewardActive = true
	g.session.ActiveChestRewardItems = []ChestRewardDisplayItem{{Option: FireRate, Title: "FASTER FIRE"}}
	if !g.advanceChestReward() {
		t.Fatal("chest advance consumed frame = false, want true")
	}
	if g.session.ChestRewardActive || len(g.session.ActiveChestRewardItems) != 0 {
		t.Fatalf("chest reward state active=%v items=%v, want closed and cleared", g.session.ChestRewardActive, g.session.ActiveChestRewardItems)
	}

	g.session.LevelUpChoiceActive = true
	g.session.PendingLevelUpLevels = []int{2}
	g.session.ActiveLevelUpOptions = []LevelUpOption{ExtraFireball}
	if !g.selectLevelUpOptionAt(0) {
		t.Fatal("level-up selection consumed frame = false, want true")
	}
	if g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) != 0 || g.session.Progression.SimultaneousFireball != 2 {
		t.Fatalf("level-up state active=%v pending=%v fireballs=%d; want applied and closed", g.session.LevelUpChoiceActive, g.session.PendingLevelUpLevels, g.session.Progression.SimultaneousFireball)
	}
}

func TestDebugJumpToLevel100QueuesProgressionChoicesAndGold(t *testing.T) {
	g := New()
	g.player.Pos = Vec2{X: 25, Y: -40}
	g.session.Progression.GainExperienceToLevel(3)
	g.session.LevelUpChoiceActive = true
	g.session.ChestRewardActive = true
	g.session.PendingLevelUpLevels = []int{2, 3}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate}
	g.session.ActiveChestRewardItems = []ChestRewardDisplayItem{{Option: FireRate, Title: "FASTER FIRE"}}
	g.chests = []Chest{{Tier: ChestBronze, Pos: Vec2{X: 999, Y: 999}}}

	if g.handleJumpToLevel100DebugKeyDown() {
		t.Fatal("level-100 debug consumed frame, want scene update to continue")
	}

	if got := g.session.Progression.Level; got != debugLevelJumpTarget {
		t.Fatalf("level after debug jump = %d, want %d", got, debugLevelJumpTarget)
	}
	if got := g.session.Progression.Experience; got != 0 {
		t.Fatalf("experience remainder after debug jump = %d, want 0", got)
	}
	if got, want := g.session.Progression.NextExperience, ExperienceRequirement(debugLevelJumpTarget); got != want {
		t.Fatalf("next experience after debug jump = %d, want %d", got, want)
	}
	if got := g.session.CollectedCoins; got != debugLevelJumpCoins {
		t.Fatalf("coins after debug jump = %d, want %d", got, debugLevelJumpCoins)
	}
	if got, want := len(g.chests), debugLevelJumpGoldChests+1; got != want {
		t.Fatalf("chests after debug jump = %d, want %d", got, want)
	}
	if g.chests[0].Tier != ChestBronze {
		t.Fatalf("existing chest tier after debug jump = %v, want preserved bronze", g.chests[0].Tier)
	}
	for i, chest := range g.chests[1:] {
		if chest.Tier != ChestGold {
			t.Fatalf("debug chest %d tier = %v, want gold", i, chest.Tier)
		}
		if got := math.Sqrt(DistanceSq(chest.Pos, g.player.Pos)); math.Abs(got-debugLevelJumpChestRadius) > 0.0001 {
			t.Fatalf("debug chest %d radius = %v, want %v", i, got, debugLevelJumpChestRadius)
		}
	}
	if g.session.Progression.LightningUnlocked || g.session.Progression.OrbitalOrbUnlocked || g.session.Progression.BeamUnlocked || g.session.Progression.MeteorUnlocked {
		t.Fatal("debug jump unlocked skills before the player chose them")
	}
	if g.session.Progression.SimultaneousFireball != 1 || g.session.PlayerLives != g.tuning.InitialPlayerLives {
		t.Fatalf("debug jump applied upgrades before choice: fireballs=%d lives=%d", g.session.Progression.SimultaneousFireball, g.session.PlayerLives)
	}
	if !g.session.LevelUpChoiceActive || g.session.ChestRewardActive {
		t.Fatalf("debug jump modal state active levelUp=%v chest=%v, want level-up only", g.session.LevelUpChoiceActive, g.session.ChestRewardActive)
	}
	if got, want := len(g.session.PendingLevelUpLevels), debugLevelJumpTarget-1; got != want {
		t.Fatalf("pending choices after debug jump = %d, want %d", got, want)
	}
	if got := g.session.CurrentLevelUpPresentation; got != 2 {
		t.Fatalf("current level-up presentation after debug jump = %d, want existing pending level 2", got)
	}
	if len(g.session.ActiveLevelUpOptions) == 0 {
		t.Fatal("debug jump did not present level-up options")
	}
	if len(g.session.ActiveChestRewardItems) != 0 {
		t.Fatalf("chest reward items after debug jump = %v, want cleared", g.session.ActiveChestRewardItems)
	}
}

func TestRestartClearsSpatialIndexUntilFirstSkeletonUpdateLikeOriginal(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 40}, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)

	g.reset()

	if len(g.skeleton) != 1 {
		t.Fatalf("skeleton count after reset = %d, want fresh initial skeleton", len(g.skeleton))
	}
	if idx := g.spatial.FirstNear(g.skeleton[0].Pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != -1 {
		t.Fatalf("fresh skeleton visible in spatial index before first update at index %d, want no candidate", idx)
	}
	g.updateSkeletons(0)
	if idx := g.spatial.FirstNear(g.skeleton[0].Pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != 0 {
		t.Fatalf("fresh skeleton index after first skeleton update = %d, want 0", idx)
	}
}

func TestRestartSuppressesHeldMovementKeysLikeOriginalPressedKeyReset(t *testing.T) {
	g := New()
	g.session.GameOver = true

	held := map[ebiten.Key]bool{ebiten.KeyArrowRight: true}
	g.restartGame(func(key ebiten.Key) bool { return held[key] })

	if !g.suppressedMovement[ebiten.KeyArrowRight] {
		t.Fatalf("suppressed movement after restart = %v, want held right key suppressed", g.suppressedMovement)
	}
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return held[key] }); got != (Vec2{}) {
		t.Fatalf("movement while restart-held key remains down = %+v, want zero", got)
	}

	held[ebiten.KeyArrowRight] = false
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return held[key] }); got != (Vec2{}) {
		t.Fatalf("movement while releasing restart-held key = %+v, want zero", got)
	}
	held[ebiten.KeyArrowRight] = true
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return held[key] }); got != (Vec2{X: 1}) {
		t.Fatalf("movement after re-pressing restart-held key = %+v, want right", got)
	}
}

func TestTriggerGameOverMatchesOriginalSessionCleanup(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.ChestRewardActive = true
	g.session.PlayerHitInvulnerability = 0.5
	g.session.PendingLevelUpLevels = []int{2, 3}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate}
	g.session.LevelUpRedrawStatusTimer = redrawFailurePulseDuration
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpOverlayTimer = modalFadeDuration
	g.session.LevelUpTitleScaleTimer = modalFadeDuration
	g.session.LevelUpOptionFadeTimer = modalFadeDuration
	g.session.ActiveChestRewardItems = []ChestRewardDisplayItem{{Option: FireRate, Title: "FASTER FIRE"}}
	g.session.ChestRewardOverlayTimer = modalFadeDuration
	g.player.HitFlash = playerHitFlashDuration / 2
	g.suppressedMovement[ebiten.KeyArrowLeft] = true
	g.effects = append(g.effects, Effect{Kind: EffectBeam, TTL: 1, MaxTTL: 1})
	g.meteors = append(g.meteors, MeteorProjectile{})

	g.triggerGameOver()

	if !g.session.GameOver {
		t.Fatal("GameOver is false, want true")
	}
	if g.session.LevelUpChoiceActive || g.session.ChestRewardActive {
		t.Fatal("modal overlays remained active after game over")
	}
	if g.session.PlayerHitInvulnerability != 0 {
		t.Fatalf("invulnerability = %v, want 0", g.session.PlayerHitInvulnerability)
	}
	if g.player.HitFlash != 0 {
		t.Fatalf("player hit flash = %v, want 0", g.player.HitFlash)
	}
	if len(g.session.PendingLevelUpLevels) != 0 || len(g.session.ActiveLevelUpOptions) != 0 || len(g.session.ActiveChestRewardItems) != 0 {
		t.Fatal("queued/modal state was not cleared")
	}
	if len(g.suppressedMovement) != 0 {
		t.Fatalf("suppressed movement keys = %v, want empty", g.suppressedMovement)
	}
	if g.session.LevelUpRedrawStatusTimer != 0 {
		t.Fatalf("redraw status timer = %v, want 0", g.session.LevelUpRedrawStatusTimer)
	}
	if g.session.LevelUpRedrawFadeTimer != 0 || g.session.LevelUpRedrawCoinFadeTimer != 0 || g.session.LevelUpOverlayTimer != 0 || g.session.LevelUpTitleScaleTimer != 0 || g.session.LevelUpOptionFadeTimer != 0 || g.session.ChestRewardOverlayTimer != 0 {
		t.Fatalf(
			"overlay timers after game over = redraw %.2f redrawCoin %.2f level %.2f title %.2f options %.2f chest %.2f; want all zero",
			g.session.LevelUpRedrawFadeTimer,
			g.session.LevelUpRedrawCoinFadeTimer,
			g.session.LevelUpOverlayTimer,
			g.session.LevelUpTitleScaleTimer,
			g.session.LevelUpOptionFadeTimer,
			g.session.ChestRewardOverlayTimer,
		)
	}
	if len(g.effects) != 0 || len(g.meteors) != 0 {
		t.Fatal("temporary effects or meteors were not cleared")
	}
	if g.player.DeathRotation != 0 {
		t.Fatalf("initial death rotation = %v, want 0 before animation advances", g.player.DeathRotation)
	}
	g.updatePausedAnimations(playerDeathRotationDuration / 2)
	if math.Abs(g.player.DeathRotation-(-math.Pi/4)) > 0.0001 {
		t.Fatalf("halfway death rotation = %v, want -pi/4", g.player.DeathRotation)
	}
	g.updatePausedAnimations(playerDeathRotationDuration)
	if math.Abs(g.player.DeathRotation-(-math.Pi/2)) > 0.0001 {
		t.Fatalf("final death rotation = %v, want -pi/2", g.player.DeathRotation)
	}
}

func TestGameOverContinuesUnpausedWorldActionsLikeOriginal(t *testing.T) {
	g := New()
	g.session.GameOver = true
	g.coins = []Coin{{Phase: 0.25}}
	g.skeleton = []Skeleton{{ID: 101, HP: 1, HitFlash: skeletonDamageFlashDuration}}
	g.fireball = []Fireball{{AnimFrame: 1}}
	g.fireAnimTimer = g.tuning.FireballAnimationFrameTime - 0.001

	g.updatePausedAnimations(0.06)

	if math.Abs(g.coins[0].Phase-0.31) > 0.0001 {
		t.Fatalf("coin phase during game over = %v, want 0.31", g.coins[0].Phase)
	}
	if math.Abs(g.skeleton[0].HitFlash-(skeletonDamageFlashDuration-0.06)) > 0.0001 {
		t.Fatalf("skeleton hit flash during game over = %v, want %v", g.skeleton[0].HitFlash, skeletonDamageFlashDuration-0.06)
	}
	if g.fireball[0].AnimFrame != 1 || math.Abs(g.fireAnimTimer-(g.tuning.FireballAnimationFrameTime-0.001)) > 0.0001 {
		t.Fatalf("fireball animation during game over = frame %d timer %v, want unchanged", g.fireball[0].AnimFrame, g.fireAnimTimer)
	}
}

func TestSuppressedMovementKeyMustBeReleasedBeforeMovingAgain(t *testing.T) {
	g := New()
	g.suppressedMovement[ebiten.KeyArrowRight] = true

	pressed := map[ebiten.Key]bool{ebiten.KeyArrowRight: true}
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return pressed[key] }); got != (Vec2{}) {
		t.Fatalf("movement while suppressed key is held = %+v, want zero", got)
	}

	pressed[ebiten.KeyArrowRight] = false
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return pressed[key] }); got != (Vec2{}) {
		t.Fatalf("movement while releasing suppressed key = %+v, want zero", got)
	}
	if g.suppressedMovement[ebiten.KeyArrowRight] {
		t.Fatal("suppressed movement key remained suppressed after release")
	}

	pressed[ebiten.KeyArrowRight] = true
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return pressed[key] }); got != (Vec2{X: 1}) {
		t.Fatalf("movement after re-press = %+v, want right", got)
	}
}

func TestOverlayPresentationSuppressesHeldMovementKeysLikeOriginalPressedKeyClear(t *testing.T) {
	g := New()

	g.suppressHeldMovementKeys(func(key ebiten.Key) bool {
		return key == ebiten.KeyArrowLeft || key == ebiten.KeyArrowUp
	})

	if !g.suppressedMovement[ebiten.KeyArrowLeft] || !g.suppressedMovement[ebiten.KeyArrowUp] {
		t.Fatalf("suppressed movement keys = %v, want left and up", g.suppressedMovement)
	}
	if g.suppressedMovement[ebiten.KeyArrowRight] || g.suppressedMovement[ebiten.KeyArrowDown] {
		t.Fatalf("unexpected suppressed movement keys = %v", g.suppressedMovement)
	}
}

func TestOverlaySuppressesMovementKeysPressedWhileModalIsActiveLikeOriginal(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	held := map[ebiten.Key]bool{ebiten.KeyArrowRight: true}

	g.suppressModalHeldMovementKeys(func(key ebiten.Key) bool { return held[key] })

	if !g.suppressedMovement[ebiten.KeyArrowRight] {
		t.Fatalf("suppressed movement keys = %v, want right key captured during modal", g.suppressedMovement)
	}

	g.session.LevelUpChoiceActive = false
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return held[key] }); got != (Vec2{}) {
		t.Fatalf("movement while modal-pressed key remains down = %+v, want zero", got)
	}
	held[ebiten.KeyArrowRight] = false
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return held[key] }); got != (Vec2{}) {
		t.Fatalf("movement while releasing modal-pressed key = %+v, want zero", got)
	}
	held[ebiten.KeyArrowRight] = true
	if got := g.playerMovementVector(func(key ebiten.Key) bool { return held[key] }); got != (Vec2{X: 1}) {
		t.Fatalf("movement after re-pressing modal-pressed key = %+v, want right", got)
	}
}
