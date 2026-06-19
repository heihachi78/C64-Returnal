package game

import (
	"image/color"
	"math"
	"math/rand"
	"os"
	"runtime"
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
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
	if !isKillAllAndGrantExperienceKey(ebiten.KeyDigit1) || !isKillAllAndGrantExperienceKey(ebiten.KeyNumpad1) {
		t.Fatalf("kill-all keys did not include main and numpad 1")
	}
	if isKillAllAndGrantExperienceKey(ebiten.KeyDigit2) {
		t.Fatalf("kill-all accepted digit 2, want only original 1 bindings")
	}
}

func TestInitialWindowSizeMatchesOriginalSwiftApp(t *testing.T) {
	if ScreenWidth != 800 || ScreenHeight != 600 {
		t.Fatalf("initial screen size = %dx%d, want 800x600", ScreenWidth, ScreenHeight)
	}
}

func TestHUDFontSizesMatchOriginalSwiftLabels(t *testing.T) {
	tests := []struct {
		name string
		got  float64
		want float64
	}{
		{name: "status", got: statusFontSize, want: 16},
		{name: "combat", got: combatFontSize, want: 14},
		{name: "level title", got: levelUpTitleFontSize, want: 40},
		{name: "level option", got: levelUpOptionFontSize, want: 22},
		{name: "level key", got: levelUpKeyFontSize, want: 18},
		{name: "chest title", got: chestTitleFontSize, want: 34},
		{name: "chest item", got: chestItemFontSize, want: 20},
		{name: "chest continue", got: chestContinueFontSize, want: 18},
		{name: "game over title", got: gameOverTitleFontSize, want: 40},
		{name: "game over option", got: gameOverOptionFontSize, want: 22},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Fatalf("%s font size = %v, want %v", tt.name, tt.got, tt.want)
		}
	}

	if got := fontFaceForSize(levelUpTitleFontSize).Metrics().Height.Ceil(); got < 36 {
		t.Fatalf("level title font metric height = %d, want at least 36", got)
	}
}

func TestHUDFontPrefersOriginalMenloBoldWhenAvailable(t *testing.T) {
	if !fontNameMatches("Menlo Bold", []string{"menlo", "bold"}) {
		t.Fatal("fontNameMatches rejected Menlo Bold")
	}
	if fontNameMatches("Menlo Regular", []string{"menlo", "bold"}) {
		t.Fatal("fontNameMatches accepted Menlo Regular")
	}

	const menloPath = "/System/Library/Fonts/Menlo.ttc"
	if _, err := os.Stat(menloPath); err != nil {
		t.Skip("Menlo.ttc is not available on this platform")
	}

	font, name := loadSystemFontByFullName([]string{menloPath}, []string{"menlo", "bold"})
	if font == nil {
		t.Fatalf("could not load Menlo Bold from %s", menloPath)
	}
	if !fontNameMatches(name, []string{"menlo", "bold"}) {
		t.Fatalf("loaded HUD font name = %q, want Menlo Bold", name)
	}
}

func TestHUDIntervalFormattingMatchesOriginalStatusPanel(t *testing.T) {
	if got := formattedSeconds(3); got != "3.0" {
		t.Fatalf("formatted seconds for whole interval = %q, want 3.0", got)
	}
	if got := formattedSeconds(0.91); got != "0.91" {
		t.Fatalf("formatted seconds below one = %q, want 0.91", got)
	}
	if got := formattedSeconds(0.955); got != "0.95" {
		t.Fatalf("formatted seconds rounding below one = %q, want Swift-style 0.95", got)
	}
}

func TestGrassTintBlendFactorMatchesOriginalField(t *testing.T) {
	if math.Abs(grassTintBlendFactor-0.22) > 0.0001 {
		t.Fatalf("grassTintBlendFactor = %v, want 0.22", grassTintBlendFactor)
	}
}

func TestGrassGridMatchesOriginalInfiniteFieldLayout(t *testing.T) {
	startColumn, startRow, columns, rows := grassGrid(ScreenWidth, ScreenHeight, 64, Vec2{})
	if columns != 17 || rows != 14 {
		t.Fatalf("grass grid size = %dx%d, want 17x14", columns, rows)
	}
	if startColumn != -8 || startRow != -7 {
		t.Fatalf("grass grid start = (%d,%d), want (-8,-7)", startColumn, startRow)
	}

	startColumn, startRow, columns, rows = grassGrid(ScreenWidth, ScreenHeight, 64, Vec2{X: 130, Y: -130})
	if columns != 17 || rows != 14 {
		t.Fatalf("shifted grass grid size = %dx%d, want 17x14", columns, rows)
	}
	if startColumn != -6 || startRow != -10 {
		t.Fatalf("shifted grass grid start = (%d,%d), want (-6,-10)", startColumn, startRow)
	}
}

func TestGrassHashMatchesOriginalIntMinGuard(t *testing.T) {
	g := New()
	minInt := -int(^uint(0)>>1) - 1
	if got := g.grassHash(0, 0, minInt); got != 0 {
		t.Fatalf("grassHash with Int.min-equivalent salt = %d, want 0", got)
	}
}

func TestSpriteCullingMatchesOriginalNonVisibleNodeIntent(t *testing.T) {
	if spriteBoundsVisible(800, 600, -20, 300, 30, 42, 0) {
		t.Fatal("fully offscreen sprite was visible, want culled")
	}
	if !spriteBoundsVisible(800, 600, -14, 300, 30, 42, 0) {
		t.Fatal("edge-overlapping sprite was culled, want visible")
	}
	if !spriteBoundsVisible(800, 600, 808, 300, 24, 24, math.Pi/4) {
		t.Fatal("rotated edge-overlapping sprite was culled, want conservative visible")
	}
	if spriteBoundsVisible(800, 600, 900, 300, 24, 24, math.Pi/4) {
		t.Fatal("far offscreen rotated sprite was visible, want culled")
	}
}

func TestWorldRenderLayerOrderMatchesOriginalZPositions(t *testing.T) {
	want := []float64{-20, 8.5, 8.75, 9, 9.5, 10, 11, 12, 13, 14}
	got := worldRenderLayerOrder()
	if !slices.Equal(got, want) {
		t.Fatalf("world render layer order = %v, want SpriteKit z order %v", got, want)
	}
	for i := 1; i < len(got); i++ {
		if got[i] <= got[i-1] {
			t.Fatalf("world render layer order is not strictly ascending at %d: %v", i, got)
		}
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

func TestDebugKillAllKeyDownDoesNotConsumeUpdateFrameLikeOriginal(t *testing.T) {
	g := New()
	g.skeleton = nil
	g.spatial.Rebuild(g.skeleton)
	level := g.session.Progression.Level
	experience := g.session.Progression.Experience

	if g.handleKillAllAndGrantExperienceKeyDown() {
		t.Fatal("kill-all keyDown consumed frame with no skeletons, want no-op")
	}
	if g.session.Progression.Level != level || g.session.Progression.Experience != experience || g.session.LevelUpChoiceActive {
		t.Fatalf("no-op kill-all changed state: level=%d xp=%d active=%v", g.session.Progression.Level, g.session.Progression.Experience, g.session.LevelUpChoiceActive)
	}

	g.skeleton = []Skeleton{{ID: 101, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)
	if g.handleKillAllAndGrantExperienceKeyDown() {
		t.Fatal("kill-all keyDown consumed frame after clearing skeletons, want scene update to observe resulting state")
	}
	if len(g.skeleton) != 0 || !g.session.LevelUpChoiceActive {
		t.Fatalf("kill-all state skeletons=%d levelUp=%v, want cleared skeletons and modal", len(g.skeleton), g.session.LevelUpChoiceActive)
	}
	g.updatePausedAnimations(1.0 / float64(TargetTPS))
	if math.Abs(g.session.LevelUpOverlayTimer-1.0/float64(TargetTPS)) > 0.0001 {
		t.Fatalf("level-up overlay timer after non-consumed keyDown = %v, want one frame", g.session.LevelUpOverlayTimer)
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

func TestParallelSkeletonUpdateReportsLaunchedWorkerCount(t *testing.T) {
	oldMaxProcs := runtime.GOMAXPROCS(6)
	defer runtime.GOMAXPROCS(oldMaxProcs)

	g := New()
	g.tuning.ParallelSkeletonUpdateThreshold = 1
	g.skeleton = make([]Skeleton, 10)
	for i := range g.skeleton {
		g.skeleton[i] = Skeleton{ID: i + 1, Pos: Vec2{X: float64(i + 1), Y: 0}, HP: 1, Facing: 1}
	}

	g.updateSkeletons(1.0 / TargetTPS)

	if g.lastParallelJobs != 5 {
		t.Fatalf("parallel skeleton jobs = %d, want 5 launched chunks for 10 skeletons over 6 procs", g.lastParallelJobs)
	}
}

func TestParallelSkeletonUpdateMatchesOriginalSequentialMovement(t *testing.T) {
	oldMaxProcs := runtime.GOMAXPROCS(4)
	defer runtime.GOMAXPROCS(oldMaxProcs)

	sequential := New()
	parallel := New()
	sequential.player.Pos = Vec2{X: 25, Y: -12}
	parallel.player.Pos = sequential.player.Pos
	sequential.tuning.ParallelSkeletonUpdateThreshold = 1_000
	parallel.tuning.ParallelSkeletonUpdateThreshold = 1
	sequential.skeletonAnimTimer = sequential.tuning.SkeletonAnimationFrameTime - 0.001
	parallel.skeletonAnimTimer = sequential.skeletonAnimTimer

	skeletons := make([]Skeleton, 32)
	for i := range skeletons {
		skeletons[i] = Skeleton{
			ID:        i + 1,
			Pos:       Vec2{X: 90 + float64(i%8)*31, Y: -160 + float64(i/8)*47},
			HP:        1,
			Reward:    1,
			Facing:    -1,
			AnimFrame: 0,
		}
	}
	sequential.skeleton = slices.Clone(skeletons)
	parallel.skeleton = slices.Clone(skeletons)

	dt := 1.0 / float64(TargetTPS)
	sequential.updateSkeletons(dt)
	parallel.updateSkeletons(dt)

	if parallel.lastParallelJobs <= 1 {
		t.Fatalf("parallel skeleton jobs = %d, want concurrent path", parallel.lastParallelJobs)
	}
	if len(parallel.skeleton) != len(sequential.skeleton) {
		t.Fatalf("parallel skeleton count = %d, want %d", len(parallel.skeleton), len(sequential.skeleton))
	}
	for i := range sequential.skeleton {
		got := parallel.skeleton[i]
		want := sequential.skeleton[i]
		if math.Abs(got.Pos.X-want.Pos.X) > 0.0001 || math.Abs(got.Pos.Y-want.Pos.Y) > 0.0001 ||
			got.Facing != want.Facing || got.AnimFrame != want.AnimFrame {
			t.Fatalf("parallel skeleton[%d] = %+v, want sequential %+v", i, got, want)
		}
		if idx := parallel.spatial.FirstNear(got.Pos, 0.001, parallel.skeleton, func(candidate int) bool {
			return parallel.skeleton[candidate].ID == got.ID
		}); idx < 0 {
			t.Fatalf("parallel skeleton %d missing from rebuilt spatial index", got.ID)
		}
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

func TestWorldRotationConvertsToEbitenScreenSpaceLikeOriginalVisuals(t *testing.T) {
	swiftClockwiseDeathRotation := -math.Pi / 2
	if got, want := worldRotationToScreen(swiftClockwiseDeathRotation), math.Pi/2; math.Abs(got-want) > 0.0001 {
		t.Fatalf("screen death rotation = %v, want %v", got, want)
	}

	swiftUpwardProjectileRotation := math.Pi / 2
	if got, want := worldRotationToScreen(swiftUpwardProjectileRotation), -math.Pi/2; math.Abs(got-want) > 0.0001 {
		t.Fatalf("screen projectile rotation = %v, want %v", got, want)
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

func TestWorldEffectsPauseWhileLevelUpModalIsActive(t *testing.T) {
	g := New()
	g.effects = append(g.effects, Effect{Kind: EffectBeam, TTL: 0.001, MaxTTL: 0.001})
	g.queueLevelUpChoices(1)

	if !g.session.LevelUpChoiceActive {
		t.Fatal("level-up modal is inactive, want active")
	}
	if err := g.Update(); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if len(g.effects) != 1 {
		t.Fatalf("effects expired while world should be paused, count = %d", len(g.effects))
	}
}

func TestPresentingLevelUpStopsPlayerAnimationWithoutClearingDirectionLikeOriginal(t *testing.T) {
	g := New()
	g.player.Moving = true
	g.player.MoveDir = Vec2{X: 1}
	g.player.AnimTimer = 0.1
	g.player.AnimFrame = 1

	g.queueLevelUpChoices(1)

	if g.player.Moving || g.player.AnimTimer != 0 || g.player.AnimFrame != 0 {
		t.Fatalf("player animation state = %+v, want stopped at frame 0", g.player)
	}
	if g.player.MoveDir != (Vec2{X: 1}) {
		t.Fatalf("player move direction = %+v, want preserved like Swift currentPlayerMovementDirection", g.player.MoveDir)
	}
}

func TestPlayerWalkAnimationAdvancesInWorldActionPhaseLikeOriginal(t *testing.T) {
	g := New()
	g.player.Moving = true
	g.player.AnimTimer = g.tuning.PlayerAnimationFrameTime - 0.001

	g.updatePlayerWalkAnimation(0.002)

	if g.player.AnimFrame != 1 {
		t.Fatalf("player animation frame = %d, want 1", g.player.AnimFrame)
	}
	if math.Abs(g.player.AnimTimer-0.001) > 0.0001 {
		t.Fatalf("player animation timer = %v, want carried remainder 0.001", g.player.AnimTimer)
	}
}

func TestPlayerWalkAnimationDoesNotAdvanceWhenOverlayPausesWorldBeforeActionPhase(t *testing.T) {
	g := New()
	g.player.Moving = true
	g.player.AnimTimer = g.tuning.PlayerAnimationFrameTime - 0.001
	g.session.LevelUpChoiceActive = true

	g.updateNewlyPresentedOverlayActions(0.002)

	if g.player.AnimFrame != 0 {
		t.Fatalf("player animation frame = %d, want unchanged", g.player.AnimFrame)
	}
	wantTimer := g.tuning.PlayerAnimationFrameTime - 0.001
	if math.Abs(g.player.AnimTimer-wantTimer) > 0.0001 {
		t.Fatalf("player animation timer = %v, want unchanged %v", g.player.AnimTimer, wantTimer)
	}
	if g.session.LevelUpOverlayTimer == 0 {
		t.Fatal("level-up overlay timer did not advance, want overlay actions running")
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

func TestShowingChestRewardStopsPlayerAnimationWithoutClearingDirectionLikeOriginal(t *testing.T) {
	g := New()
	g.player.Moving = true
	g.player.MoveDir = Vec2{X: 1}
	g.player.AnimTimer = 0.1
	g.player.AnimFrame = 1

	g.applyChestReward(ChestBronze)

	if g.player.Moving || g.player.AnimTimer != 0 || g.player.AnimFrame != 0 {
		t.Fatalf("player animation state = %+v, want stopped at frame 0", g.player)
	}
	if g.player.MoveDir != (Vec2{X: 1}) {
		t.Fatalf("player move direction = %+v, want preserved like Swift currentPlayerMovementDirection", g.player.MoveDir)
	}
}

func TestSyncOrbitalOrbCountImmediatelyAlignsNewOrbs(t *testing.T) {
	g := New()
	g.player.Pos = Vec2{X: 12, Y: -7}
	g.session.Progression.ApplyLevelUpOption(LearnOrb)

	g.syncOrbitalOrbCount()

	if len(g.orbs) != 1 {
		t.Fatalf("orb count = %d, want 1", len(g.orbs))
	}
	want := Vec2{X: g.player.Pos.X + g.tuning.OrbitalOrbRadius, Y: g.player.Pos.Y}
	if math.Abs(g.orbs[0].Pos.X-want.X) > 0.0001 || math.Abs(g.orbs[0].Pos.Y-want.Y) > 0.0001 {
		t.Fatalf("orb position = %+v, want %+v immediately after sync", g.orbs[0].Pos, want)
	}
}

func TestOrbitalOrbAnimationResetsWhenNoOrbsAreActive(t *testing.T) {
	g := New()
	g.orbs = []OrbitalOrb{{Active: false, AnimFrame: 1}}
	g.orbAnimTimer = g.tuning.OrbitalAnimationFrameTime - 0.001

	g.updateOrbAnimation(0.01)

	if g.orbAnimTimer != 0 {
		t.Fatalf("orb animation timer = %v, want 0", g.orbAnimTimer)
	}
	if g.orbs[0].AnimFrame != 1 {
		t.Fatalf("inactive orb frame = %d, want unchanged 1", g.orbs[0].AnimFrame)
	}
}

func TestInactiveOrbAnimationFrameStaysAtFreshTextureLikeOriginal(t *testing.T) {
	g := New()
	g.orbs = []OrbitalOrb{
		{Active: true, AnimFrame: 0},
		{Active: false, AnimFrame: 0},
	}
	g.orbAnimTimer = g.tuning.OrbitalAnimationFrameTime - 0.001

	g.updateOrbAnimation(0.01)

	if g.orbs[0].AnimFrame != 1 {
		t.Fatalf("active orb frame = %d, want 1", g.orbs[0].AnimFrame)
	}
	if g.orbs[1].AnimFrame != 0 {
		t.Fatalf("inactive orb frame = %d, want fresh texture frame 0", g.orbs[1].AnimFrame)
	}
}

func TestRespawnedOrbReturnsWithFreshTextureFrame(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnOrb)
	g.orbs = []OrbitalOrb{{Active: false, MissingOrbitProgress: math.Pi*2 - 0.01, AnimFrame: 1}}

	g.updateOrbitalOrbs(0.1)

	if !g.orbs[0].Active {
		t.Fatal("orb did not respawn after completing a missing orbit")
	}
	if g.orbs[0].AnimFrame != 0 {
		t.Fatalf("respawned orb frame = %d, want fresh texture frame 0", g.orbs[0].AnimFrame)
	}
}

func TestFamilyAnimationsResyncFreshSpritesOnNextTickLikeOriginal(t *testing.T) {
	g := New()

	g.skeletonAnimFrame = 1
	g.skeleton = []Skeleton{{ID: 101, AnimFrame: 1}, {ID: 202, AnimFrame: 0}}
	g.skeletonAnimTimer = g.tuning.SkeletonAnimationFrameTime - 0.001
	g.updateSkeletonAnimation(0.01)
	if g.skeleton[0].AnimFrame != 0 || g.skeleton[1].AnimFrame != 0 {
		t.Fatalf("skeleton frames = %d,%d; want shared frame 0", g.skeleton[0].AnimFrame, g.skeleton[1].AnimFrame)
	}

	g.fireAnimFrame = 1
	g.fireball = []Fireball{{AnimFrame: 1}, {AnimFrame: 0}}
	g.fireAnimTimer = g.tuning.FireballAnimationFrameTime - 0.001
	g.updateFireballAnimation(0.01)
	if g.fireball[0].AnimFrame != 0 || g.fireball[1].AnimFrame != 0 {
		t.Fatalf("fireball frames = %d,%d; want shared frame 0", g.fireball[0].AnimFrame, g.fireball[1].AnimFrame)
	}

	g.orbAnimFrame = 1
	g.orbs = []OrbitalOrb{{Active: true, AnimFrame: 1}, {Active: true, AnimFrame: 0}}
	g.orbAnimTimer = g.tuning.OrbitalAnimationFrameTime - 0.001
	g.updateOrbAnimation(0.01)
	if g.orbs[0].AnimFrame != 0 || g.orbs[1].AnimFrame != 0 {
		t.Fatalf("orb frames = %d,%d; want shared frame 0", g.orbs[0].AnimFrame, g.orbs[1].AnimFrame)
	}

	g.meteorAnimFrame = 1
	g.meteors = []MeteorProjectile{{AnimFrame: 1}, {AnimFrame: 0}}
	g.meteorAnimTimer = g.tuning.MeteorAnimationFrameTime - 0.001
	g.updateMeteorAnimation(0.01)
	if g.meteors[0].AnimFrame != 0 || g.meteors[1].AnimFrame != 0 {
		t.Fatalf("meteor frames = %d,%d; want shared frame 0", g.meteors[0].AnimFrame, g.meteors[1].AnimFrame)
	}
}

func TestLevelUpRedrawOptionComparisonIgnoresOrder(t *testing.T) {
	first := []LevelUpOption{FireRate, ExtraFireball, LearnBeam}
	second := []LevelUpOption{LearnBeam, FireRate, ExtraFireball}

	if !sameOptionSet(first, second) {
		t.Fatal("sameOptionSet returned false for the same options in a different order")
	}
}

func TestLevelUpRedrawRetryRerunsFullOriginalOptionSelection(t *testing.T) {
	const seed = 9
	first := New()
	first.rng = rand.New(rand.NewSource(seed))
	first.tuning.HalveHordeChanceNumerator = 100
	first.tuning.HalveHordeChanceDenominator = 100
	first.tuning.ExtraOptionChanceNumerator = 100
	first.tuning.ExtraOptionChanceDenominator = 100
	first.skeleton = []Skeleton{{ID: 1, HP: 1}}
	initial := first.randomLevelUpOptions(nil)

	retry := New()
	retry.rng = rand.New(rand.NewSource(seed))
	retry.tuning = first.tuning
	retry.skeleton = []Skeleton{{ID: 1, HP: 1}}
	redrawn := retry.randomLevelUpOptions(initial)

	if len(redrawn) != 3 {
		t.Fatalf("redrawn option count = %d, want Swift retry to rerun extra-option chance and keep 3", len(redrawn))
	}
	hasHalve := false
	for _, option := range redrawn {
		if option == HalveSkeletons {
			hasHalve = true
			break
		}
	}
	if !hasHalve {
		t.Fatalf("redrawn options = %v, want retry to rerun halve-horde chance", redrawn)
	}
}

func TestModalFadeTargetsMatchOriginalPanelAndContentAlpha(t *testing.T) {
	g := New()
	if got := g.modalPanelAlpha(modalFadeDuration); got != c64Panel.A {
		t.Fatalf("panel alpha at fade end = %d, want %d", got, c64Panel.A)
	}
	if got := g.modalContentAlpha(modalFadeDuration); got != 255 {
		t.Fatalf("content alpha at fade end = %d, want 255", got)
	}
	if got := g.modalPanelAlpha(modalFadeDuration / 2); got != 79 {
		t.Fatalf("halfway panel alpha = %d, want 79", got)
	}
	if got := g.modalContentAlpha(modalFadeDuration / 2); got != 128 {
		t.Fatalf("halfway content alpha = %d, want 128", got)
	}
}

func TestEffectFadeAlphaRoundsLikeOriginalActions(t *testing.T) {
	if got := effectFadeAlpha(0.09, 0.18); got != 128 {
		t.Fatalf("halfway effect alpha = %d, want 128", got)
	}
	if got := effectFadeAlpha(0, 0.18); got != 0 {
		t.Fatalf("expired effect alpha = %d, want 0", got)
	}
	if got := effectFadeAlpha(0.18, 0.18); got != 255 {
		t.Fatalf("fresh effect alpha = %d, want 255", got)
	}
}

func TestLightningHitEffectAlphaMatchesOriginalTargetDuplicate(t *testing.T) {
	if got := lightningHitEffectAlpha(0.18, 0.18); got != 217 {
		t.Fatalf("fresh lightning hit alpha = %d, want 217", got)
	}
	if got := lightningHitEffectAlpha(0.09, 0.18); got != 108 {
		t.Fatalf("halfway lightning hit alpha = %d, want 108", got)
	}
	if got := lightningHitEffectAlpha(0, 0.18); got != 0 {
		t.Fatalf("expired lightning hit alpha = %d, want 0", got)
	}
}

func TestWorldEffectsAgeAfterActiveGameplayCreatesThemLikeOriginalActions(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 60}, HP: 2, Reward: 1}}
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.session.Casts.Beam = g.session.Progression.BeamCastInterval()

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if len(g.effects) != 1 || g.effects[0].Kind != EffectBeam {
		t.Fatalf("effects = %+v, want one beam effect", g.effects)
	}
	want := g.tuning.BeamEffectDuration - 1.0/float64(TargetTPS)
	if math.Abs(g.effects[0].TTL-want) > 0.0001 {
		t.Fatalf("new beam effect TTL = %v, want %v after same-frame action aging", g.effects[0].TTL, want)
	}
}

func TestWorldEffectsDoNotAgeWhenChestPausesWorldBeforeActionPhase(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.effects = []Effect{{Kind: EffectBeam, TTL: 1, MaxTTL: 1}}
	g.chests = []Chest{{Pos: g.player.Pos, Tier: ChestBronze}}

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.ChestRewardActive {
		t.Fatal("chest reward is inactive, want modal to pause world")
	}
	if got := g.effects[0].TTL; got != 1 {
		t.Fatalf("effect TTL after same-frame chest pause = %v, want 1", got)
	}
}

func TestChestOverlayActionsStartOnSameFrameChestIsCollectedLikeOriginal(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.chests = []Chest{{Pos: g.player.Pos, Tier: ChestBronze}}

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.ChestRewardActive {
		t.Fatal("chest reward is inactive, want active")
	}
	want := 1.0 / float64(TargetTPS)
	if math.Abs(g.session.ChestRewardOverlayTimer-want) > 0.0001 {
		t.Fatalf("chest overlay timer = %v, want %v", g.session.ChestRewardOverlayTimer, want)
	}
}

func TestLevelUpOverlayActionsStartOnSameFrameKillPresentsChoiceLikeOriginal(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.player.Facing = 1
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 60}, HP: 1, Reward: 1}}
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.session.Casts.Beam = g.session.Progression.BeamCastInterval()

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.LevelUpChoiceActive {
		t.Fatal("level-up choice is inactive, want active")
	}
	want := 1.0 / float64(TargetTPS)
	if math.Abs(g.session.LevelUpOverlayTimer-want) > 0.0001 ||
		math.Abs(g.session.LevelUpTitleScaleTimer-want) > 0.0001 ||
		math.Abs(g.session.LevelUpOptionFadeTimer-want) > 0.0001 {
		t.Fatalf(
			"level-up timers overlay=%v title=%v option=%v, want %v",
			g.session.LevelUpOverlayTimer,
			g.session.LevelUpTitleScaleTimer,
			g.session.LevelUpOptionFadeTimer,
			want,
		)
	}
}

func TestGameOverOverlayActionsStartOnSameFrameCollisionKillsPlayerLikeOriginal(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.session.PlayerLives = 1
	g.coins = []Coin{{Pos: Vec2{X: 100}, Amount: 1}}
	g.skeleton = []Skeleton{{ID: 101, Pos: g.player.Pos, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.GameOver {
		t.Fatal("game over is false, want true")
	}
	want := 1.0 / float64(TargetTPS)
	if math.Abs(g.session.GameOverOverlayTimer-want) > 0.0001 {
		t.Fatalf("game-over overlay timer = %v, want %v", g.session.GameOverOverlayTimer, want)
	}
	wantRotation := -math.Pi / 2 * (want / playerDeathRotationDuration)
	if math.Abs(g.player.DeathRotation-wantRotation) > 0.0001 {
		t.Fatalf("death rotation = %v, want %v", g.player.DeathRotation, wantRotation)
	}
	if math.Abs(g.coins[0].Phase-want) > 0.0001 {
		t.Fatalf("coin phase on game-over frame = %v, want exactly one active-frame advance %v", g.coins[0].Phase, want)
	}
}

func TestPanelCornerRadiusMatchesOriginalHUD(t *testing.T) {
	if math.Abs(panelCornerRadius-6) > 0.0001 {
		t.Fatalf("panelCornerRadius = %v, want 6", panelCornerRadius)
	}
}

func TestModalTitleScaleMatchesOriginalOverlayAction(t *testing.T) {
	if got := modalTitleScale(0); math.Abs(got-0.75) > 0.0001 {
		t.Fatalf("initial title scale = %v, want 0.75", got)
	}
	if got := modalTitleScale(modalFadeDuration / 2); math.Abs(got-0.875) > 0.0001 {
		t.Fatalf("halfway title scale = %v, want 0.875", got)
	}
	if got := modalTitleScale(modalFadeDuration); math.Abs(got-1) > 0.0001 {
		t.Fatalf("final title scale = %v, want 1", got)
	}
}

func TestScaledTextUsesOriginalLabelScaleInsteadOfScaledFontSize(t *testing.T) {
	baseFace := fontFaceForSize(levelUpTitleFontSize)
	scaledFace := fontFaceForSize(levelUpTitleFontSize * modalTitleScale(0))
	text := "LEVEL 3"
	baseWidth := font.MeasureString(baseFace, text).Ceil()
	scaledWidth := font.MeasureString(scaledFace, text).Ceil()

	layout := spriteKitScaledTextLayout(baseFace, text, true)

	if layout.Width != baseWidth+8 {
		t.Fatalf("scaled text backing width = %d, want base label width %d plus padding", layout.Width, baseWidth)
	}
	if layout.Width == scaledWidth+8 {
		t.Fatalf("scaled text backing width matched scaled font width %d; want SpriteKit node scaling from base font", scaledWidth)
	}
	if math.Abs(layout.AnchorX-(4+float64(baseWidth)/2)) > 0.0001 {
		t.Fatalf("centered scaled text anchor X = %v, want label center", layout.AnchorX)
	}
	leftLayout := spriteKitScaledTextLayout(baseFace, text, false)
	if leftLayout.AnchorX != 4 {
		t.Fatalf("left scaled text anchor X = %v, want left label origin plus padding", leftLayout.AnchorX)
	}

	g := New()
	first := g.scaledTextImage(text, levelUpTitleFontSize, true)
	second := g.scaledTextImage(text, levelUpTitleFontSize, true)
	left := g.scaledTextImage(text, levelUpTitleFontSize, false)
	if first.Image == nil || second.Image == nil {
		t.Fatal("scaled text cache returned nil image")
	}
	if first.Image != second.Image {
		t.Fatal("scaled text cache did not reuse the base label image")
	}
	if first.Image == left.Image {
		t.Fatal("scaled text cache reused centered label image for left-aligned label")
	}
}

func TestLevelUpOptionAlphaTracksOptionFadeSeparatelyFromPanel(t *testing.T) {
	if got := levelUpOptionContentAlpha(255, 0); got != 0 {
		t.Fatalf("initial option alpha = %d, want 0", got)
	}
	if got := levelUpOptionContentAlpha(255, modalFadeDuration/2); got != 128 {
		t.Fatalf("halfway option alpha = %d, want 128", got)
	}
	if got := levelUpOptionContentAlpha(90, modalFadeDuration); got != 90 {
		t.Fatalf("option alpha during panel fade = %d, want 90", got)
	}
}

func TestSuccessfulRedrawRestartsOriginalLevelUpOptionPresentation(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.Progression.Level = 4
	g.session.CollectedCoins = 4
	g.session.PendingLevelUpLevels = []int{4}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraFireball}
	g.session.LevelUpOverlayTimer = modalFadeDuration
	g.session.LevelUpTitleScaleTimer = modalFadeDuration
	g.session.LevelUpOptionFadeTimer = modalFadeDuration
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration

	g.redrawLevelUpOptions()

	if g.session.LevelUpOverlayTimer != modalFadeDuration {
		t.Fatalf("overlay timer after redraw = %v, want %v", g.session.LevelUpOverlayTimer, modalFadeDuration)
	}
	if g.session.LevelUpTitleScaleTimer != 0 {
		t.Fatalf("title scale timer after redraw = %v, want 0", g.session.LevelUpTitleScaleTimer)
	}
	if g.session.LevelUpOptionFadeTimer != 0 {
		t.Fatalf("option fade timer after redraw = %v, want 0", g.session.LevelUpOptionFadeTimer)
	}
	if g.session.LevelUpRedrawFadeTimer != redrawStatusFadeDuration {
		t.Fatalf("redraw fade timer after redraw = %v, want preserved %v", g.session.LevelUpRedrawFadeTimer, redrawStatusFadeDuration)
	}
	if g.session.LevelUpRedrawCoinFadeTimer != 0 {
		t.Fatalf("redraw coin fade timer after unaffordable redrawn presentation = %v, want 0", g.session.LevelUpRedrawCoinFadeTimer)
	}
	if math.Abs(g.session.LevelUpRedrawStatusTimer-redrawFailurePulseDuration) > 0.0001 {
		t.Fatalf("redraw status timer after unaffordable redrawn presentation = %v, want %v", g.session.LevelUpRedrawStatusTimer, redrawFailurePulseDuration)
	}
	if g.session.CollectedCoins != 0 {
		t.Fatalf("coins after redraw = %d, want 0", g.session.CollectedCoins)
	}
}

func TestSuccessfulAffordableRedrawKeepsRedrawStatusVisibleLikeOriginal(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.Progression.Level = 4
	g.session.CollectedCoins = 9
	g.session.PendingLevelUpLevels = []int{4}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraFireball}
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration

	g.redrawLevelUpOptions()

	if g.session.CollectedCoins != 5 {
		t.Fatalf("coins after redraw = %d, want 5", g.session.CollectedCoins)
	}
	if g.session.LevelUpRedrawStatusTimer != 0 {
		t.Fatalf("redraw status timer after still-affordable redraw = %v, want 0", g.session.LevelUpRedrawStatusTimer)
	}
	if g.session.LevelUpRedrawFadeTimer != redrawStatusFadeDuration {
		t.Fatalf("redraw fade timer after still-affordable redraw = %v, want preserved %v", g.session.LevelUpRedrawFadeTimer, redrawStatusFadeDuration)
	}
	if g.session.LevelUpRedrawCoinFadeTimer != redrawStatusFadeDuration {
		t.Fatalf("redraw coin fade timer after still-affordable redraw = %v, want immediate full alpha marker %v", g.session.LevelUpRedrawCoinFadeTimer, redrawStatusFadeDuration)
	}
}

func TestPresentingLevelUpPulsesUnaffordableRedrawLikeOriginal(t *testing.T) {
	g := New()
	g.session.Progression.Level = 4
	g.session.CollectedCoins = 0
	g.session.PendingLevelUpLevels = []int{4}

	g.presentNextLevelUpChoiceIfNeeded()

	if math.Abs(g.session.LevelUpRedrawStatusTimer-redrawFailurePulseDuration) > 0.0001 {
		t.Fatalf("redraw status timer after unaffordable presentation = %v, want %v", g.session.LevelUpRedrawStatusTimer, redrawFailurePulseDuration)
	}
	if g.session.LevelUpRedrawFadeTimer != 0 {
		t.Fatalf("redraw fade timer after presentation = %v, want 0", g.session.LevelUpRedrawFadeTimer)
	}
	if g.session.LevelUpRedrawCoinFadeTimer != 0 {
		t.Fatalf("redraw coin fade timer after presentation = %v, want 0", g.session.LevelUpRedrawCoinFadeTimer)
	}
}

func TestPresentingLevelUpDoesNotPulseAffordableRedrawLikeOriginal(t *testing.T) {
	g := New()
	g.session.Progression.Level = 4
	g.session.CollectedCoins = 4
	g.session.PendingLevelUpLevels = []int{4}

	g.presentNextLevelUpChoiceIfNeeded()

	if g.session.LevelUpRedrawStatusTimer != 0 {
		t.Fatalf("redraw status timer after affordable presentation = %v, want 0", g.session.LevelUpRedrawStatusTimer)
	}
	if g.session.LevelUpRedrawFadeTimer != 0 {
		t.Fatalf("redraw fade timer after presentation = %v, want 0", g.session.LevelUpRedrawFadeTimer)
	}
	if g.session.LevelUpRedrawCoinFadeTimer != redrawStatusFadeDuration {
		t.Fatalf("redraw coin fade timer after affordable presentation = %v, want immediate full alpha marker %v", g.session.LevelUpRedrawCoinFadeTimer, redrawStatusFadeDuration)
	}
}

func TestRedrawRequiresQueuedLevelLikeOriginal(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.Progression.Level = 4
	g.session.CollectedCoins = 4
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraFireball}
	g.session.LevelUpTitleScaleTimer = modalFadeDuration
	g.session.LevelUpOptionFadeTimer = modalFadeDuration
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration

	g.redrawLevelUpOptions()

	if g.session.CollectedCoins != 4 {
		t.Fatalf("coins after redraw without queued level = %d, want 4", g.session.CollectedCoins)
	}
	if !slices.Equal(g.session.ActiveLevelUpOptions, []LevelUpOption{FireRate, ExtraFireball}) {
		t.Fatalf("options after redraw without queued level = %v, want unchanged", g.session.ActiveLevelUpOptions)
	}
	if g.session.LevelUpTitleScaleTimer != modalFadeDuration || g.session.LevelUpOptionFadeTimer != modalFadeDuration || g.session.LevelUpRedrawFadeTimer != redrawStatusFadeDuration || g.session.LevelUpRedrawCoinFadeTimer != redrawStatusFadeDuration {
		t.Fatalf(
			"timers after redraw without queued level = title %.2f options %.2f redraw %.2f redrawCoin %.2f; want unchanged",
			g.session.LevelUpTitleScaleTimer,
			g.session.LevelUpOptionFadeTimer,
			g.session.LevelUpRedrawFadeTimer,
			g.session.LevelUpRedrawCoinFadeTimer,
		)
	}
}

func TestApplyingFinalLevelUpChoiceHidesOriginalPresentationState(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.PendingLevelUpLevels = []int{2}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate}
	g.session.LevelUpRedrawStatusTimer = redrawFailurePulseDuration
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpOverlayTimer = modalFadeDuration
	g.session.LevelUpTitleScaleTimer = modalFadeDuration
	g.session.LevelUpOptionFadeTimer = modalFadeDuration

	g.applyLevelUpOption(FireRate)

	if g.session.LevelUpChoiceActive {
		t.Fatal("level-up choice remained active after applying final queued choice")
	}
	if len(g.session.ActiveLevelUpOptions) != 0 {
		t.Fatalf("active options after final choice = %v, want none", g.session.ActiveLevelUpOptions)
	}
	if g.session.LevelUpRedrawStatusTimer != 0 || g.session.LevelUpRedrawFadeTimer != 0 || g.session.LevelUpRedrawCoinFadeTimer != 0 || g.session.LevelUpOverlayTimer != 0 || g.session.LevelUpTitleScaleTimer != 0 || g.session.LevelUpOptionFadeTimer != 0 {
		t.Fatalf(
			"level-up timers after final choice = redraw %.2f redrawFade %.2f redrawCoin %.2f overlay %.2f title %.2f options %.2f; want all zero",
			g.session.LevelUpRedrawStatusTimer,
			g.session.LevelUpRedrawFadeTimer,
			g.session.LevelUpRedrawCoinFadeTimer,
			g.session.LevelUpOverlayTimer,
			g.session.LevelUpTitleScaleTimer,
			g.session.LevelUpOptionFadeTimer,
		)
	}
}

func TestApplyingQueuedLevelUpChoiceStartsNextPresentationFresh(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.PendingLevelUpLevels = []int{2, 3}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate}
	g.session.LevelUpOverlayTimer = modalFadeDuration
	g.session.LevelUpTitleScaleTimer = modalFadeDuration
	g.session.LevelUpOptionFadeTimer = modalFadeDuration
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration

	g.applyLevelUpOption(FireRate)

	if !g.session.LevelUpChoiceActive {
		t.Fatal("next queued level-up was not presented")
	}
	if g.session.CurrentLevelUpPresentation != 3 {
		t.Fatalf("current level-up presentation = %d, want 3", g.session.CurrentLevelUpPresentation)
	}
	if g.session.LevelUpOverlayTimer != 0 || g.session.LevelUpTitleScaleTimer != 0 || g.session.LevelUpOptionFadeTimer != 0 || g.session.LevelUpRedrawFadeTimer != 0 || g.session.LevelUpRedrawCoinFadeTimer != 0 {
		t.Fatalf(
			"next level-up timers = overlay %.2f title %.2f options %.2f redraw %.2f redrawCoin %.2f; want all zero",
			g.session.LevelUpOverlayTimer,
			g.session.LevelUpTitleScaleTimer,
			g.session.LevelUpOptionFadeTimer,
			g.session.LevelUpRedrawFadeTimer,
			g.session.LevelUpRedrawCoinFadeTimer,
		)
	}
}

func TestHalveSkeletonsUsesTemporaryShuffledTargetsLikeOriginal(t *testing.T) {
	const seed int64 = 19
	initial := []Skeleton{
		{ID: 101},
		{ID: 202},
		{ID: 303},
		{ID: 404},
		{ID: 505},
		{ID: 606},
	}
	targetIDs := []int{101, 202, 303, 404, 505, 606}
	expectedIDs := []int{101, 202, 303, 404, 505, 606}
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(targetIDs), func(i, j int) { targetIDs[i], targetIDs[j] = targetIDs[j], targetIDs[i] })
	for _, targetID := range targetIDs[:len(initial)/2] {
		for i, id := range expectedIDs {
			if id == targetID {
				last := len(expectedIDs) - 1
				expectedIDs[i] = expectedIDs[last]
				expectedIDs = expectedIDs[:last]
				break
			}
		}
	}

	g := New()
	g.rng = rand.New(rand.NewSource(seed))
	g.skeleton = append([]Skeleton(nil), initial...)
	g.halveSkeletons()

	if len(g.skeleton) != len(expectedIDs) {
		t.Fatalf("remaining skeleton count = %d, want %d", len(g.skeleton), len(expectedIDs))
	}
	gotIDs := make([]int, len(g.skeleton))
	for i, skeleton := range g.skeleton {
		gotIDs[i] = skeleton.ID
	}
	for i, wantID := range expectedIDs {
		if gotIDs[i] != wantID {
			t.Fatalf("remaining skeleton IDs = %v, want %v", gotIDs, expectedIDs)
		}
	}
}

func TestHalveSkeletonsDoesNotCreditWeaponKillsLikeOriginal(t *testing.T) {
	g := New()
	g.rng = rand.New(rand.NewSource(4))
	g.skeleton = []Skeleton{
		{ID: 101, HP: 1, Reward: 1},
		{ID: 202, HP: 1, Reward: 1},
		{ID: 303, HP: 1, Reward: 1},
		{ID: 404, HP: 1, Reward: 1},
	}

	g.halveSkeletons()

	if g.session.Kills.TotalSkeletons != 2 {
		t.Fatalf("total kills after halve horde = %d, want 2", g.session.Kills.TotalSkeletons)
	}
	if g.session.Kills.Fireball != 0 || g.session.Kills.Lightning != 0 || g.session.Kills.OrbitalOrb != 0 || g.session.Kills.Beam != 0 || g.session.Kills.Meteor != 0 {
		t.Fatalf("weapon kills after halve horde = %+v, want no weapon attribution", g.session.Kills)
	}
}

func TestMilestoneSkeletonSpawnDefersSpatialRefreshLikeOriginal(t *testing.T) {
	g := New()
	g.skeleton = g.skeleton[:0]
	g.spatial.Rebuild(g.skeleton)
	g.session.Kills.TotalSkeletons = 1
	g.tuning.RedKillInterval = 1
	g.tuning.PurpleKillInterval = 0

	g.spawnMilestoneSkeletonsIfNeeded()

	if len(g.skeleton) != 1 || g.skeleton[0].Kind != SkeletonRed {
		t.Fatalf("milestone skeletons = %+v, want one red skeleton", g.skeleton)
	}
	pos := g.skeleton[0].Pos
	if idx := g.spatial.FirstNear(pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != -1 {
		t.Fatalf("milestone skeleton was visible in spatial index before rebuild at index %d", idx)
	}
	g.spatial.Rebuild(g.skeleton)
	if idx := g.spatial.FirstNear(pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != 0 {
		t.Fatalf("milestone skeleton index after rebuild = %d, want 0", idx)
	}
}

func TestTimedSkeletonSpawnDefersSpatialRefreshLikeOriginal(t *testing.T) {
	g := New()
	g.skeleton = g.skeleton[:0]
	g.spatial.Rebuild(g.skeleton)
	g.session.Casts.SkeletonSpawn = g.session.Progression.SkeletonSpawnInterval()

	g.updateSkeletonSpawning(0)

	if len(g.skeleton) != 1 {
		t.Fatalf("timed skeleton count = %d, want 1", len(g.skeleton))
	}
	pos := g.skeleton[0].Pos
	if idx := g.spatial.FirstNear(pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != -1 {
		t.Fatalf("timed skeleton was visible in spatial index before rebuild at index %d", idx)
	}
	g.updateSkeletons(0)
	if idx := g.spatial.FirstNear(pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != 0 {
		t.Fatalf("timed skeleton index after skeleton update = %d, want 0", idx)
	}
}

func TestDestroyTriggeredMilestoneSpawnDefersSpatialRefreshLikeOriginal(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{}, HP: 1, Reward: 0}}
	g.spatial.Rebuild(g.skeleton)
	g.session.Kills.TotalSkeletons = 0
	g.session.NextChestMilestone = 1_000_000
	g.tuning.RedKillInterval = 1
	g.tuning.PurpleKillInterval = 0

	g.destroySkeleton(0, AttackFireball)

	if len(g.skeleton) != 1 || g.skeleton[0].Kind != SkeletonRed {
		t.Fatalf("post-kill skeletons = %+v, want one deferred red milestone skeleton", g.skeleton)
	}
	pos := g.skeleton[0].Pos
	if idx := g.spatial.FirstNear(pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != -1 {
		t.Fatalf("post-kill milestone skeleton was visible in spatial index before rebuild at index %d", idx)
	}
	g.updateSkeletons(0)
	if idx := g.spatial.FirstNear(pos, g.tuning.SkeletonHitDistance, g.skeleton, func(int) bool { return true }); idx != 0 {
		t.Fatalf("post-kill milestone skeleton index after skeleton update = %d, want 0", idx)
	}
}

func TestGameOverLayoutMatchesOriginalHUDRect(t *testing.T) {
	x, y, w, h := gameOverPanelRect(800, 600)
	if x != 90 || y != 210 || w != 620 || h != 190 {
		t.Fatalf("game over panel rect = (%v,%v,%v,%v), want (90,210,620,190)", x, y, w, h)
	}
	if y+gameOverTitleOffsetY != 258 {
		t.Fatalf("game over title y = %v, want 258", y+gameOverTitleOffsetY)
	}
	if y+gameOverRestartOffsetY != 322 || y+gameOverExitOffsetY != 370 {
		t.Fatalf("game over option y values = %v,%v; want 322,370", y+gameOverRestartOffsetY, y+gameOverExitOffsetY)
	}

	g := New()
	if got := g.gameOverOptionAt(400, 322); got != "restart" {
		t.Fatalf("game over option at restart center = %q, want restart", got)
	}
	if got := g.gameOverOptionAt(400, 370); got != "exit" {
		t.Fatalf("game over option at exit center = %q, want exit", got)
	}

	face := fontFaceForSize(gameOverOptionFontSize)
	restartOutsideX := 400 + float64(font.MeasureString(face, "RESTART").Ceil())/2 + 27
	if got := g.gameOverOptionAt(restartOutsideX, 322); got != "" {
		t.Fatalf("game over option outside restart label hit area = %q, want none like Swift", got)
	}
}

func TestHUDStatusPanelLayoutMatchesOriginalHUDRect(t *testing.T) {
	x, y, w, h := topStatusPanelRect(3)
	if x != 8 || y != 9 || w != 210 || h != 104 {
		t.Fatalf("three-life top panel rect = (%v,%v,%v,%v), want (8,9,210,104)", x, y, w, h)
	}
	x, y, w, h = topStatusPanelRect(13)
	if x != 8 || y != 9 || w != 210 || h != 120 {
		t.Fatalf("thirteen-life top panel rect = (%v,%v,%v,%v), want (8,9,210,120)", x, y, w, h)
	}

	x, y = lifeIconScreenPosition(0)
	if x != 25 || y != 98 {
		t.Fatalf("first life icon position = (%v,%v), want (25,98)", x, y)
	}
	x, y = lifeIconScreenPosition(11)
	if x != 201 || y != 98 {
		t.Fatalf("twelfth life icon position = (%v,%v), want (201,98)", x, y)
	}
	x, y = lifeIconScreenPosition(12)
	if x != 25 || y != 114 {
		t.Fatalf("thirteenth life icon position = (%v,%v), want (25,114)", x, y)
	}

	x, y, w, h = combatStatusPanelRect(600)
	if x != 8 || y != 267 || w != 176 || h != 330 {
		t.Fatalf("combat panel rect = (%v,%v,%v,%v), want (8,267,176,330)", x, y, w, h)
	}
}

func TestChestRewardLayoutMatchesOriginalHUDRect(t *testing.T) {
	x, y, w, h := chestOverlayPanelRect(800, 600, 1)
	if x != 90 || y != 196 || w != 620 || h != 208 {
		t.Fatalf("one-item chest panel rect = (%v,%v,%v,%v), want (90,196,620,208)", x, y, w, h)
	}
	if chestOverlayTitleY(y) != 250 {
		t.Fatalf("one-item chest title y = %v, want 250", chestOverlayTitleY(y))
	}
	if chestOverlayContinueY(y, h) != 368 {
		t.Fatalf("one-item chest continue y = %v, want 368", chestOverlayContinueY(y, h))
	}
	if chestRewardItemY(y, h, 1, 0) != 312 {
		t.Fatalf("one-item chest item y = %v, want 312", chestRewardItemY(y, h, 1, 0))
	}

	x, y, w, h = chestOverlayPanelRect(800, 600, 2)
	if x != 90 || y != 179 || w != 620 || h != 242 {
		t.Fatalf("two-item chest panel rect = (%v,%v,%v,%v), want (90,179,620,242)", x, y, w, h)
	}
	if chestOverlayTitleY(y) != 233 {
		t.Fatalf("two-item chest title y = %v, want 233", chestOverlayTitleY(y))
	}
	if chestOverlayContinueY(y, h) != 385 {
		t.Fatalf("two-item chest continue y = %v, want 385", chestOverlayContinueY(y, h))
	}
	if chestRewardItemY(y, h, 2, 0) != 297 || chestRewardItemY(y, h, 2, 1) != 327 {
		t.Fatalf("two-item chest item y values = %v,%v; want 297,327", chestRewardItemY(y, h, 2, 0), chestRewardItemY(y, h, 2, 1))
	}
}

func TestLevelUpMouseHitAreasMatchOriginalLabelAndIconOnly(t *testing.T) {
	g := New()
	g.screenW = 800
	g.screenH = 600
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate}

	if got := g.levelUpOptionAt(250, 304); got != -1 {
		t.Fatalf("option at key hint center = %d, want no selection like Swift", got)
	}
	if got := g.levelUpOptionAt(284, 304); got != 0 {
		t.Fatalf("option at icon center = %d, want 0", got)
	}
	if got := g.levelUpOptionAt(322, 304); got != 0 {
		t.Fatalf("option at label start = %d, want 0", got)
	}
}

func TestLevelUpOptionsAreCappedToOriginalVisibleFour(t *testing.T) {
	options := visibleLevelUpOptions([]LevelUpOption{
		FireRate,
		ExtraFireball,
		ExtraLife,
		LearnLightning,
		LearnOrb,
	})

	if len(options) != 4 {
		t.Fatalf("visible option count = %d, want 4", len(options))
	}
	if options[0] != FireRate || options[3] != LearnLightning {
		t.Fatalf("visible options = %v, want first four preserved", options)
	}
}

func TestRedrawMouseHitAreasMatchOriginalNodesOnly(t *testing.T) {
	g := New()
	g.screenW = 800
	g.screenH = 600
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraFireball}
	g.session.CollectedCoins = 123456789
	g.session.Progression.Level = 999

	label := "REDRAW 999  COINS 123456789"
	face := fontFaceForSize(levelUpOptionFontSize)
	labelTailX := 322 + float64(font.MeasureString(face, label).Ceil()) + 20
	if !g.redrawRectContains(labelTailX, 422) {
		t.Fatal("redraw hit at long label tail = false, want true like Swift")
	}
	if !g.redrawRectContains(250, 422) {
		t.Fatal("redraw hit at key label center = false, want true")
	}
	if !g.redrawRectContains(284, 422) {
		t.Fatal("redraw hit at coin center = false, want true")
	}
	if !g.redrawRectContains(322, 422) {
		t.Fatal("redraw hit at label start = false, want true")
	}
}

func TestRedrawFailurePulseMatchesOriginalOverlayAction(t *testing.T) {
	if math.Abs(redrawFailurePulseDuration-0.16) > 0.0001 {
		t.Fatalf("redrawFailurePulseDuration = %v, want 0.16", redrawFailurePulseDuration)
	}
	if got := redrawPulseScale(redrawFailurePulseDuration); math.Abs(got-1) > 0.0001 {
		t.Fatalf("initial redraw pulse scale = %v, want 1", got)
	}
	if got := redrawPulseScale(redrawFailurePulseDuration - 0.08); math.Abs(got-1.08) > 0.0001 {
		t.Fatalf("mid redraw pulse scale = %v, want 1.08", got)
	}
	if got := redrawPulseScale(0); math.Abs(got-1) > 0.0001 {
		t.Fatalf("final redraw pulse scale = %v, want 1", got)
	}
}

func TestRedrawStatusFadeMatchesOriginalOverlayAction(t *testing.T) {
	if math.Abs(redrawStatusFadeDuration-0.14) > 0.0001 {
		t.Fatalf("redrawStatusFadeDuration = %v, want 0.14", redrawStatusFadeDuration)
	}
	if got := levelUpRedrawContentAlpha(0); got != 0 {
		t.Fatalf("initial redraw text alpha = %d, want 0", got)
	}
	if got := levelUpRedrawContentAlpha(redrawStatusFadeDuration / 2); got != 128 {
		t.Fatalf("halfway redraw text alpha = %d, want 128", got)
	}
	if got := levelUpRedrawContentAlpha(redrawStatusFadeDuration); got != 255 {
		t.Fatalf("final redraw text alpha = %d, want 255 independent of modal fade", got)
	}
}

func TestRedrawCoinAlphaMatchesOriginalAffordabilityAction(t *testing.T) {
	if got := redrawCoinAlpha(true, 0); got != 255 {
		t.Fatalf("initial affordable redraw coin alpha = %d, want 255", got)
	}
	if got := redrawCoinAlpha(false, 0); got != 115 {
		t.Fatalf("initial unaffordable redraw coin alpha = %d, want 115", got)
	}
	if got := redrawCoinAlpha(false, redrawStatusFadeDuration); got != 255 {
		t.Fatalf("final unaffordable redraw coin alpha = %d, want 255", got)
	}
}

func TestRedrawFailurePulseAdvancesWhileLevelUpPausesWorld(t *testing.T) {
	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.Progression.Level = 4
	g.session.CollectedCoins = 0
	g.session.PendingLevelUpLevels = []int{4}
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration

	g.redrawLevelUpOptions()
	if math.Abs(g.session.LevelUpRedrawStatusTimer-redrawFailurePulseDuration) > 0.0001 {
		t.Fatalf("redraw timer after failed redraw = %v, want %v", g.session.LevelUpRedrawStatusTimer, redrawFailurePulseDuration)
	}
	if g.session.LevelUpRedrawFadeTimer != redrawStatusFadeDuration {
		t.Fatalf("redraw fade timer after failed redraw = %v, want unchanged %v", g.session.LevelUpRedrawFadeTimer, redrawStatusFadeDuration)
	}
	if g.session.LevelUpRedrawCoinFadeTimer != 0 {
		t.Fatalf("redraw coin fade timer after failed redraw = %v, want 0", g.session.LevelUpRedrawCoinFadeTimer)
	}
	g.updatePausedAnimations(0.08)
	if math.Abs(g.session.LevelUpRedrawStatusTimer-0.08) > 0.0001 {
		t.Fatalf("redraw timer after paused animations = %v, want 0.08", g.session.LevelUpRedrawStatusTimer)
	}
	if math.Abs(g.session.LevelUpRedrawCoinFadeTimer-0.08) > 0.0001 {
		t.Fatalf("redraw coin fade timer after paused animations = %v, want 0.08", g.session.LevelUpRedrawCoinFadeTimer)
	}
	g.updatePausedAnimations(0.08)
	if g.session.LevelUpRedrawStatusTimer != 0 {
		t.Fatalf("redraw timer after pulse = %v, want 0", g.session.LevelUpRedrawStatusTimer)
	}
}

func TestSkeletonTintBlendFactorsMatchOriginalSwiftValues(t *testing.T) {
	tests := []struct {
		kind       SkeletonKind
		wantColor  [3]uint8
		wantFactor float64
	}{
		{kind: SkeletonRegular, wantColor: [3]uint8{255, 255, 255}, wantFactor: 0},
		{kind: SkeletonRed, wantColor: [3]uint8{242, 13, 10}, wantFactor: 0.68},
		{kind: SkeletonPurple, wantColor: [3]uint8{148, 31, 242}, wantFactor: 0.72},
		{kind: SkeletonBlack, wantColor: [3]uint8{5, 5, 5}, wantFactor: 0.86},
	}

	for _, tt := range tests {
		color, factor := skeletonTintBlend(tt.kind)
		if color.R != tt.wantColor[0] || color.G != tt.wantColor[1] || color.B != tt.wantColor[2] {
			t.Fatalf("kind %v color = (%d,%d,%d), want %v", tt.kind, color.R, color.G, color.B, tt.wantColor)
		}
		if math.Abs(factor-tt.wantFactor) > 0.0001 {
			t.Fatalf("kind %v factor = %v, want %v", tt.kind, factor, tt.wantFactor)
		}
	}
}

func TestPlayerSpritePresentationMatchesOriginalDeathTint(t *testing.T) {
	presentation := playerSpritePresentation(Player{HitFlash: playerHitFlashDuration / 2}, false)
	if presentation.Tint.R != 255 || presentation.Tint.G != 255 || presentation.Tint.B != 255 {
		t.Fatalf("active hit-flash tint rgb = %+v, want white with alpha-only flash", presentation.Tint)
	}
	if presentation.BlendFactor != 0 || presentation.Rotation != 0 {
		t.Fatalf("active hit-flash presentation = %+v, want no color blend or rotation", presentation)
	}

	presentation = playerSpritePresentation(Player{HitFlash: playerHitFlashDuration, DeathRotation: -math.Pi / 2}, true)
	if presentation.Tint != (color.RGBA{217, 13, 20, 115}) {
		t.Fatalf("death tint = %+v, want SpriteKit red color with 0.45 alpha", presentation.Tint)
	}
	if math.Abs(presentation.BlendFactor-0.65) > 0.0001 {
		t.Fatalf("death blend factor = %v, want 0.65", presentation.BlendFactor)
	}
	if math.Abs(presentation.Rotation-math.Pi/2) > 0.0001 {
		t.Fatalf("death screen rotation = %v, want pi/2", presentation.Rotation)
	}
}

func TestSkeletonSpritePresentationPreservesTintBlendDuringHitFlash(t *testing.T) {
	presentation := skeletonSpritePresentation(Skeleton{Kind: SkeletonPurple, HitFlash: skeletonDamageFlashDuration})
	if presentation.Tint.R != 148 || presentation.Tint.G != 31 || presentation.Tint.B != 242 || presentation.Tint.A != 255 {
		t.Fatalf("fresh purple hit-flash tint = %+v, want purple tint with full alpha", presentation.Tint)
	}
	if math.Abs(presentation.BlendFactor-0.72) > 0.0001 {
		t.Fatalf("purple blend factor = %v, want 0.72", presentation.BlendFactor)
	}

	presentation = skeletonSpritePresentation(Skeleton{Kind: SkeletonPurple, HitFlash: skeletonDamageFlashDuration - 0.06})
	if presentation.Tint.R != 148 || presentation.Tint.G != 31 || presentation.Tint.B != 242 || presentation.Tint.A != 89 {
		t.Fatalf("dim purple hit-flash tint = %+v, want purple tint with 0.35 alpha", presentation.Tint)
	}
	if math.Abs(presentation.BlendFactor-0.72) > 0.0001 {
		t.Fatalf("dim purple blend factor = %v, want 0.72", presentation.BlendFactor)
	}
}

func TestSkeletonDamageFlashDurationMatchesOriginalActionSequence(t *testing.T) {
	if math.Abs(skeletonDamageFlashDuration-0.24) > 0.0001 {
		t.Fatalf("skeletonDamageFlashDuration = %v, want 0.24", skeletonDamageFlashDuration)
	}
}

func TestSkeletonHitFlashAdvancesAsWorldActionEvenWhenSkeletonOverlapsPlayer(t *testing.T) {
	g := New()
	g.player.Pos = Vec2{}
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{}, HitFlash: skeletonDamageFlashDuration}}

	g.updateSkeletonHitFlashes(0.06)

	if math.Abs(g.skeleton[0].HitFlash-(skeletonDamageFlashDuration-0.06)) > 0.0001 {
		t.Fatalf("hit flash = %v, want %v", g.skeleton[0].HitFlash, skeletonDamageFlashDuration-0.06)
	}
}

func TestSkeletonHitFlashAgesAfterSameFrameDamageLikeOriginalActions(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 60}, HP: 2, Reward: 1}}
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.session.Casts.Beam = g.session.Progression.BeamCastInterval()

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if len(g.skeleton) != 1 {
		t.Fatalf("skeleton count = %d, want 1 damaged skeleton", len(g.skeleton))
	}
	want := skeletonDamageFlashDuration - 1.0/float64(TargetTPS)
	if math.Abs(g.skeleton[0].HitFlash-want) > 0.0001 {
		t.Fatalf("same-frame hit flash = %v, want %v", g.skeleton[0].HitFlash, want)
	}
}

func TestSkeletonHitFlashDoesNotAgeWhenChestPausesWorldBeforeActionPhase(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.skeleton = []Skeleton{{ID: 101, HitFlash: skeletonDamageFlashDuration}}
	g.chests = []Chest{{Pos: g.player.Pos, Tier: ChestBronze}}

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.ChestRewardActive {
		t.Fatal("chest reward is inactive, want modal to pause world")
	}
	if got := g.skeleton[0].HitFlash; got != skeletonDamageFlashDuration {
		t.Fatalf("hit flash after same-frame chest pause = %v, want %v", got, skeletonDamageFlashDuration)
	}
}

func TestPlayerHitFlashAgesAfterSameFrameDamageLikeOriginalAction(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.player.Pos = Vec2{}
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{}, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if got, want := g.session.PlayerLives, g.tuning.InitialPlayerLives-1; got != want {
		t.Fatalf("player lives = %d, want %d", got, want)
	}
	if got, want := g.session.PlayerHitInvulnerability, g.tuning.PlayerHitInvulnerability; got != want {
		t.Fatalf("invulnerability = %v, want %v", got, want)
	}
	wantFlash := playerHitFlashDuration - 1.0/float64(TargetTPS)
	if math.Abs(g.player.HitFlash-wantFlash) > 0.0001 {
		t.Fatalf("player hit flash = %v, want %v after same-frame action aging", g.player.HitFlash, wantFlash)
	}
}

func TestPlayerHitFlashDoesNotAgeWhenChestPausesWorldBeforeActionPhase(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.player.HitFlash = playerHitFlashDuration
	g.chests = []Chest{{Pos: g.player.Pos, Tier: ChestBronze}}

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.ChestRewardActive {
		t.Fatal("chest reward is inactive, want modal to pause world")
	}
	if got := g.player.HitFlash; got != playerHitFlashDuration {
		t.Fatalf("player hit flash after same-frame chest pause = %v, want %v", got, playerHitFlashDuration)
	}
}

func TestHitFlashAlphaCurvesMatchOriginalActionSequences(t *testing.T) {
	if math.Abs(playerHitFlashDuration-0.96) > 0.0001 {
		t.Fatalf("playerHitFlashDuration = %v, want 0.96", playerHitFlashDuration)
	}

	tests := []struct {
		name     string
		elapsed  float64
		total    float64
		fadeDown float64
		fadeUp   float64
		want     uint8
	}{
		{name: "player start", elapsed: 0, total: playerHitFlashDuration, fadeDown: 0.08, fadeUp: 0.08, want: 255},
		{name: "player dim", elapsed: 0.08, total: playerHitFlashDuration, fadeDown: 0.08, fadeUp: 0.08, want: 89},
		{name: "player after repeat", elapsed: playerHitFlashDuration, total: playerHitFlashDuration, fadeDown: 0.08, fadeUp: 0.08, want: 255},
		{name: "skeleton dim", elapsed: 0.06, total: skeletonDamageFlashDuration, fadeDown: 0.06, fadeUp: 0.06, want: 89},
	}

	for _, tt := range tests {
		if got := flashActionAlpha(tt.elapsed, tt.total, tt.fadeDown, tt.fadeUp); got != tt.want {
			t.Fatalf("%s alpha = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestMeteorImpactEffectMatchesOriginalActionSequence(t *testing.T) {
	if math.Abs(meteorImpactEffectDuration-0.32) > 0.0001 {
		t.Fatalf("meteorImpactEffectDuration = %v, want 0.32", meteorImpactEffectDuration)
	}

	effect := Effect{Kind: EffectMeteorImpact, TTL: meteorImpactEffectDuration, MaxTTL: meteorImpactEffectDuration}
	if got := meteorImpactPresentation(effect); math.Abs(got.Scale-0.25) > 0.0001 || math.Abs(got.Alpha-1) > 0.0001 {
		t.Fatalf("initial meteor presentation = %+v, want scale 0.25 alpha 1", got)
	}
	effect.TTL = meteorImpactEffectDuration - 0.08
	if got := meteorImpactPresentation(effect); math.Abs(got.Scale-1) > 0.0001 || math.Abs(got.Alpha-1) > 0.0001 {
		t.Fatalf("grown meteor presentation = %+v, want scale 1 alpha 1", got)
	}
	effect.TTL = 0
	if got := meteorImpactPresentation(effect); math.Abs(got.Scale-1.25) > 0.0001 || math.Abs(got.Alpha) > 0.0001 {
		t.Fatalf("final meteor presentation = %+v, want scale 1.25 alpha 0", got)
	}
}

func TestMeteorImpactRenderScalesWholeEffectNodeLikeOriginal(t *testing.T) {
	effect := Effect{Kind: EffectMeteorImpact, Radius: 48, TTL: meteorImpactEffectDuration, MaxTTL: meteorImpactEffectDuration}
	style := meteorImpactRenderStyle(effect)
	if math.Abs(style.Radius-12) > 0.0001 || math.Abs(style.CoreRadius-4.2) > 0.0001 {
		t.Fatalf("initial meteor radii = %+v, want radius 12 core 4.2", style)
	}
	if math.Abs(style.GlowWidth-0.75) > 0.0001 || math.Abs(style.StrokeWidth-0.5) > 0.0001 {
		t.Fatalf("initial meteor stroke widths = %+v, want glow 0.75 stroke 0.5", style)
	}

	effect.TTL = meteorImpactEffectDuration - 0.08
	style = meteorImpactRenderStyle(effect)
	if math.Abs(style.Radius-48) > 0.0001 || math.Abs(style.CoreRadius-16.8) > 0.0001 {
		t.Fatalf("grown meteor radii = %+v, want radius 48 core 16.8", style)
	}
	if math.Abs(style.GlowWidth-3) > 0.0001 || math.Abs(style.StrokeWidth-2) > 0.0001 {
		t.Fatalf("grown meteor stroke widths = %+v, want glow 3 stroke 2", style)
	}
}

func TestMeteorCastingUsesOriginalPerMeteorSpawnInterval(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnMeteor)
	g.session.Progression.ApplyLevelUpOption(ExtraMeteor)
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{}, HP: 1}}

	spawnInterval := g.session.Progression.MeteorCastInterval() / float64(g.session.Progression.MeteorCount())
	g.updateMeteorCasting(spawnInterval*2 + 0.03)

	if len(g.meteors) != 2 {
		t.Fatalf("meteor cast count = %d, want 2", len(g.meteors))
	}
	if math.Abs(g.session.Casts.Meteor-0.03) > 0.0001 {
		t.Fatalf("meteor cast remainder = %v, want 0.03", g.session.Casts.Meteor)
	}
}

func TestMeteorCastingResetsTimerWhenNoSkeletonsLikeOriginal(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnMeteor)
	g.skeleton = nil
	g.session.Casts.Meteor = g.session.Progression.MeteorCastInterval()

	g.updateMeteorCasting(0.25)

	if g.session.Casts.Meteor != 0 {
		t.Fatalf("meteor cast timer with no skeletons = %v, want 0", g.session.Casts.Meteor)
	}
	if len(g.meteors) != 0 || len(g.effects) != 0 {
		t.Fatalf("meteor state with no skeletons = meteors %v effects %v, want no cast or impact", g.meteors, g.effects)
	}
}

func TestMeteorSpawnHeightStaysReasonableForLowImpacts(t *testing.T) {
	g := New()
	g.player.Pos = Vec2{X: 10, Y: 100}
	impact := Vec2{X: 20, Y: g.player.Pos.Y - 360}

	spawn := g.meteorSpawnPosition(impact)

	wantY := g.player.Pos.Y + g.tuning.MeteorFallHeight
	if math.Abs(spawn.Y-wantY) > 0.0001 {
		t.Fatalf("low-impact meteor spawn Y = %v, want %v", spawn.Y, wantY)
	}
	if spawn.Y-impact.Y <= g.tuning.MeteorFallHeight {
		t.Fatalf("low-impact meteor vertical fall = %v, want more than baseline fall height %v", spawn.Y-impact.Y, g.tuning.MeteorFallHeight)
	}
	if math.Abs(spawn.X-impact.X) > g.tuning.MeteorFallDrift {
		t.Fatalf("meteor horizontal drift = %v, want within %v", spawn.X-impact.X, g.tuning.MeteorFallDrift)
	}
}

func TestMeteorSpawnHeightPreservesBaselineForHighImpacts(t *testing.T) {
	g := New()
	g.player.Pos = Vec2{Y: 100}
	impact := Vec2{Y: g.player.Pos.Y + 80}

	spawn := g.meteorSpawnPosition(impact)

	wantY := impact.Y + g.tuning.MeteorFallHeight
	if math.Abs(spawn.Y-wantY) > 0.0001 {
		t.Fatalf("high-impact meteor spawn Y = %v, want %v", spawn.Y, wantY)
	}
}

func TestMeteorProjectileFallsLinearlyToImpactLikeOriginal(t *testing.T) {
	g := New()
	start := Vec2{X: 10, Y: 110}
	impact := Vec2{X: 50, Y: 30}
	g.meteors = []MeteorProjectile{{Pos: start, Start: start, Impact: impact}}

	g.updateMeteors(g.tuning.MeteorFallDuration / 2)

	if len(g.meteors) != 1 {
		t.Fatalf("meteor count after half fall = %d, want 1", len(g.meteors))
	}
	if math.Abs(g.meteors[0].Pos.X-30) > 0.0001 || math.Abs(g.meteors[0].Pos.Y-70) > 0.0001 {
		t.Fatalf("meteor midpoint = %+v, want {X:30 Y:70}", g.meteors[0].Pos)
	}

	g.updateMeteors(g.tuning.MeteorFallDuration / 2)

	if len(g.meteors) != 0 {
		t.Fatalf("meteor count after full fall = %d, want 0", len(g.meteors))
	}
	if len(g.effects) != 1 || g.effects[0].Kind != EffectMeteorImpact || g.effects[0].Pos != impact {
		t.Fatalf("meteor impact effect = %+v, want one effect at %+v", g.effects, impact)
	}
}

func TestMeteorProjectileDoesNotImpactBeforeOriginalFallDuration(t *testing.T) {
	g := New()
	start := Vec2{X: 10, Y: 110}
	impact := Vec2{X: 50, Y: 30}
	g.meteors = []MeteorProjectile{{Pos: start, Start: start, Impact: impact}}

	g.updateMeteors(g.tuning.MeteorFallDuration - 0.0001)

	if len(g.meteors) != 1 {
		t.Fatalf("meteor count just before impact = %d, want 1", len(g.meteors))
	}
	if len(g.effects) != 0 {
		t.Fatalf("effects just before impact = %+v, want none", g.effects)
	}
}

func TestMeteorImpactRadiusUsesOriginalInclusiveCircle(t *testing.T) {
	g := New()
	g.tuning.MeteorImpactRadius = 48
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 48}, HP: 1},
		{ID: 202, Pos: Vec2{X: 48.0001}, HP: 1},
		{ID: 303, Pos: Vec2{Y: -48}, HP: 1},
	}
	g.spatial.Rebuild(g.skeleton)

	targets := g.meteorImpactTargetIDs(Vec2{})

	if len(targets) != 2 || !slices.Contains(targets, 101) || !slices.Contains(targets, 303) {
		t.Fatalf("meteor radius targets = %v, want exactly boundary IDs 101 and 303", targets)
	}
	if slices.Contains(targets, 202) {
		t.Fatalf("meteor radius targets = %v, included outside-radius ID 202", targets)
	}
}

func TestMeteorImpactTargetsPreserveSpatialCandidateOrderByIdentity(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 303, Pos: Vec2{X: 100, Y: 0}, HP: 1},
		{ID: 101, Pos: Vec2{X: -100, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: 0, Y: 0}, HP: 1},
	}
	g.tuning.MeteorImpactRadius = 140
	g.spatial.Rebuild(g.skeleton)

	targets := g.meteorImpactTargetIDs(Vec2{})

	if len(targets) != 3 {
		t.Fatalf("meteor target count = %d, want 3", len(targets))
	}
	if targets[0] != 101 || targets[1] != 202 || targets[2] != 303 {
		t.Fatalf("meteor targets = %v, want spatial candidate IDs [101 202 303]", targets)
	}
}

func TestMeteorImpactNoOpsAfterGameOverLikeOriginal(t *testing.T) {
	g := New()
	g.session.GameOver = true
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{}, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)

	g.impactMeteor(Vec2{})

	if len(g.effects) != 0 {
		t.Fatalf("effects after game-over meteor impact = %+v, want none", g.effects)
	}
	if len(g.skeleton) != 1 || g.skeleton[0].HP != 1 {
		t.Fatalf("skeletons after game-over meteor impact = %+v, want unchanged", g.skeleton)
	}
	if g.session.Kills.Meteor != 0 || g.session.Kills.TotalSkeletons != 0 {
		t.Fatalf("kills after game-over meteor impact = %+v, want zero", g.session.Kills)
	}
}

func TestCoinAnimationUsesOriginalLinearActionCurves(t *testing.T) {
	if got := coinFloatOffset(0); math.Abs(got) > 0.0001 {
		t.Fatalf("initial coin offset = %v, want 0", got)
	}
	if got := coinFloatOffset(0.42); math.Abs(got-5) > 0.0001 {
		t.Fatalf("half-period coin offset = %v, want 5", got)
	}
	if got := coinFloatOffset(0.84); math.Abs(got) > 0.0001 {
		t.Fatalf("full-period coin offset = %v, want 0", got)
	}
	if got := coinShimmerAlpha(0); got != 255 {
		t.Fatalf("initial coin alpha = %d, want 255", got)
	}
	if got := coinShimmerAlpha(0.28); got != 184 {
		t.Fatalf("half-period coin alpha = %d, want 184", got)
	}
	if got := coinShimmerAlpha(0.56); got != 255 {
		t.Fatalf("full-period coin alpha = %d, want 255", got)
	}
}

func TestCoinAnimationsAdvanceInWorldActionPhaseAfterPickupChecksLikeOriginal(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.coins = []Coin{
		{Pos: Vec2{X: 100}, Amount: 1},
		{Pos: g.player.Pos, Amount: 3},
	}

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if len(g.coins) != 1 {
		t.Fatalf("coin count after pickup = %d, want 1", len(g.coins))
	}
	if math.Abs(g.coins[0].Phase-1.0/float64(TargetTPS)) > 0.0001 {
		t.Fatalf("remaining coin phase = %v, want one frame advanced", g.coins[0].Phase)
	}
	if g.session.CollectedCoins != 3 {
		t.Fatalf("collected coins = %d, want 3", g.session.CollectedCoins)
	}
}

func TestCoinAnimationsDoNotAdvanceWhenChestPausesWorldBeforeActionPhase(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.coins = []Coin{{Pos: Vec2{X: 100}, Amount: 1}}
	g.chests = []Chest{{Pos: g.player.Pos, Tier: ChestBronze}}

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.ChestRewardActive {
		t.Fatal("chest reward is inactive, want modal to pause world")
	}
	if got := g.coins[0].Phase; got != 0 {
		t.Fatalf("coin phase after same-frame chest pause = %v, want 0", got)
	}
}

func TestCoinAnimationsAdvanceWhenSameFrameGameOverLeavesWorldUnpausedLikeOriginal(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.session.PlayerLives = 1
	g.coins = []Coin{{Pos: Vec2{X: 100}, Amount: 1}}
	g.skeleton = []Skeleton{{ID: 101, Pos: g.player.Pos, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)

	if err := g.Update(); err != nil {
		t.Fatal(err)
	}

	if !g.session.GameOver {
		t.Fatal("game over is false, want true")
	}
	want := 1.0 / float64(TargetTPS)
	if math.Abs(g.coins[0].Phase-want) > 0.0001 {
		t.Fatalf("coin phase after same-frame game over = %v, want %v", g.coins[0].Phase, want)
	}
}

func TestLockedCombatRowsKeepKillCountsLikeOriginalHUD(t *testing.T) {
	first, second, kills := combatRowLabels("x0", "3.0s", "KILLS 0", false)
	if first != "LOCKED" || second != "--" || kills != "KILLS 0" {
		t.Fatalf("locked row labels = %q, %q, %q; want LOCKED, --, KILLS 0", first, second, kills)
	}
	row := combatRowPresentation("x0", "3.0s", "KILLS 0", false)
	if row.Tint != (color.RGBA{255, 255, 255, 255}) {
		t.Fatalf("locked row icon tint = %+v, want unchanged white", row.Tint)
	}
	if row.TextColor != c64Text {
		t.Fatalf("locked row text color = %+v, want primary text color %+v", row.TextColor, c64Text)
	}
}

func TestCombatRowTextOffsetsMatchOriginalHUDLayout(t *testing.T) {
	if combatRowFirstOffsetY != -10 {
		t.Fatalf("first row offset = %v, want -10", combatRowFirstOffsetY)
	}
	if combatRowSecondOffsetY != 6 {
		t.Fatalf("second row offset = %v, want 6", combatRowSecondOffsetY)
	}
	if combatRowKillsOffsetY != 20 {
		t.Fatalf("kills row offset = %v, want 20", combatRowKillsOffsetY)
	}
}

func TestChestRewardItemsPreserveOptionForIconRendering(t *testing.T) {
	g := New()
	g.rng = rand.New(rand.NewSource(1))
	g.session.Progression.ApplyLevelUpOption(LearnBeam)

	items := g.chestRewardItemsForSkill(SkillBeam)
	if len(items) != 2 {
		t.Fatalf("beam reward item count = %d, want 2", len(items))
	}
	if items[0].Option != BeamRate || items[1].Option != BeamKillCount {
		t.Fatalf("beam reward options = %v, %v; want BeamRate, BeamKillCount", items[0].Option, items[1].Option)
	}
	if items[1].Title != "+2 BEAM KILL" {
		t.Fatalf("beam kill reward title = %q, want %q", items[1].Title, "+2 BEAM KILL")
	}
}

func TestChestRewardItemUsesBeamBonusOnlyForBeamSkillLikeOriginal(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.session.Progression.ApplyLevelUpOption(BeamKillCount)

	if got, want := g.chestRewardItemForSkill(BeamKillCount, SkillBeam).Title, "+3 BEAM KILL"; got != want {
		t.Fatalf("beam-context title = %q, want %q", got, want)
	}
	if got, want := g.chestRewardItemForSkill(BeamKillCount, SkillFireball).Title, "+1 BEAM KILL"; got != want {
		t.Fatalf("non-beam-context title = %q, want %q", got, want)
	}
}

func TestLightningBoltPathIsJaggedAndStable(t *testing.T) {
	g := New()
	g.rng = rand.New(rand.NewSource(3))
	points := g.lightningBoltPoints(Vec2{}, Vec2{X: 180, Y: 0})

	if len(points) < 4 {
		t.Fatalf("point count = %d, want a multi-segment bolt", len(points))
	}
	if points[0] != (Vec2{}) || points[len(points)-1] != (Vec2{X: 180, Y: 0}) {
		t.Fatalf("bolt endpoints = %v, %v; want origin and target", points[0], points[len(points)-1])
	}
	hasOffset := false
	for _, point := range points[1 : len(points)-1] {
		if point.Y != 0 {
			hasOffset = true
			break
		}
	}
	if !hasOffset {
		t.Fatal("bolt path is straight; want jagged intermediate offsets")
	}
}

func TestLightningEffectUsesSeparateOuterAndInnerBoltPathsLikeOriginal(t *testing.T) {
	g := New()
	g.rng = rand.New(rand.NewSource(9))
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 180, Y: 0}, HP: 1, Reward: 1}}

	g.castLightning()

	var effect Effect
	found := false
	for _, candidate := range g.effects {
		if candidate.Kind == EffectLightning {
			effect = candidate
			found = true
			break
		}
	}
	if !found {
		t.Fatal("lightning cast did not create a bolt effect")
	}
	if len(effect.Points) != len(effect.InnerPoints) {
		t.Fatalf("outer/inner point counts = %d/%d, want equal independent paths", len(effect.Points), len(effect.InnerPoints))
	}
	if len(effect.Points) < 2 {
		t.Fatalf("lightning point count = %d, want at least endpoints", len(effect.Points))
	}
	last := len(effect.Points) - 1
	if effect.Points[0] != effect.InnerPoints[0] || effect.Points[last] != effect.InnerPoints[last] {
		t.Fatalf("outer/inner endpoints = %v/%v and %v/%v, want shared start/end", effect.Points[0], effect.InnerPoints[0], effect.Points[last], effect.InnerPoints[last])
	}
	if slices.Equal(effect.Points, effect.InnerPoints) {
		t.Fatalf("inner lightning path reused outer path; Swift creates separate random bolt shapes")
	}
}

func TestBeamDirectionUsesCurrentMovementDirectionEvenWhenAnimationStoppedLikeOriginal(t *testing.T) {
	g := New()
	g.player.Facing = 1
	g.player.Moving = false
	g.player.MoveDir = Vec2{X: 0, Y: 1}

	if got := g.playerBeamDirection(); got != (Vec2{X: 0, Y: 1}) {
		t.Fatalf("beam direction = %+v, want preserved current movement direction", got)
	}

	g.player.MoveDir = Vec2{}
	if got := g.playerBeamDirection(); got != (Vec2{X: 1, Y: 0}) {
		t.Fatalf("beam direction without movement direction = %+v, want facing direction", got)
	}
}

func TestBeamTargetsUseStableSkeletonIdentifiers(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: 20, Y: 0}, HP: 1},
	}

	targets := g.beamTargets(Vec2{X: 1, Y: 0}, 100, 4, 2)
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want 2", len(targets))
	}
	if targets[0] != 101 || targets[1] != 202 {
		t.Fatalf("targets = %v, want stable IDs [101 202]", targets)
	}
}

func TestClosestSkeletonSelectionPreservesOriginalTieOrder(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: -10, Y: 0}, HP: 1},
		{ID: 303, Pos: Vec2{X: 0, Y: 20}, HP: 1},
	}

	targets := g.closestSkeletons(Vec2{}, nil, 2)
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want 2", len(targets))
	}
	if g.skeleton[targets[0]].ID != 101 || g.skeleton[targets[1]].ID != 202 {
		t.Fatalf("tie targets = %v (%d,%d), want first two skeletons in source order", targets, g.skeleton[targets[0]].ID, g.skeleton[targets[1]].ID)
	}
}

func TestFireballVolleySkipsAlreadyTargetedSkeletonsLikeOriginal(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: 20, Y: 0}, HP: 1},
	}
	g.fireball = []Fireball{{TargetID: 101}}

	g.spawnFireballs()

	if len(g.fireball) != 2 {
		t.Fatalf("fireball count = %d, want existing fireball plus one new volley fireball", len(g.fireball))
	}
	if g.fireball[1].TargetID != 202 {
		t.Fatalf("new fireball target = %d, want untargeted skeleton 202", g.fireball[1].TargetID)
	}
}

func TestBeamTargetSelectionPreservesOriginalTieOrder(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 1}, HP: 1},
		{ID: 202, Pos: Vec2{X: 10, Y: -1}, HP: 1},
		{ID: 303, Pos: Vec2{X: 20, Y: 0}, HP: 1},
	}

	targets := g.beamTargets(Vec2{X: 1, Y: 0}, 100, 4, 2)
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want 2", len(targets))
	}
	if targets[0] != 101 || targets[1] != 202 {
		t.Fatalf("beam tie targets = %v, want stable IDs [101 202]", targets)
	}
}

func TestLightningTargetsUseStableSkeletonIdentifiers(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 30, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: 60, Y: 0}, HP: 1},
	}

	targets := g.chainLightningTargets()
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want 2", len(targets))
	}
	if targets[0].targetID != 101 || targets[1].targetID != 202 {
		t.Fatalf("targets = %+v, want stable IDs 101 then 202", targets)
	}
}

func TestLightningTargetsSkipFireballReservedSkeletonsLikeOriginal(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: 20, Y: 0}, HP: 1},
		{ID: 303, Pos: Vec2{X: 30, Y: 0}, HP: 1},
	}
	g.fireball = []Fireball{{TargetID: 101}}

	targets := g.chainLightningTargets()
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want two non-reserved targets", len(targets))
	}
	if targets[0].targetID != 202 || targets[1].targetID != 303 {
		t.Fatalf("lightning targets = %+v, want non-fireball targets 202 then 303", targets)
	}
}

func TestLightningRemainingTargetsPreserveOriginalTieOrder(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: 20, Y: 10}, HP: 1},
		{ID: 303, Pos: Vec2{X: 20, Y: -10}, HP: 1},
	}

	targets := g.chainLightningTargets()
	if len(targets) != 3 {
		t.Fatalf("target count = %d, want 3", len(targets))
	}
	if targets[0].targetID != 101 || targets[1].targetID != 202 || targets[2].targetID != 303 {
		t.Fatalf("chain lightning targets = %+v, want IDs 101, 202, 303", targets)
	}
}

func TestLightningCastCreatesTargetHitDuplicateEffect(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 30, Y: 0}, HP: 2, Reward: 1, Facing: -1, AnimFrame: 1},
	}

	g.castLightning()

	hasBolt := false
	hasHit := false
	for _, effect := range g.effects {
		if effect.Kind == EffectLightning {
			hasBolt = true
		}
		if effect.Kind == EffectLightningHit {
			hasHit = true
			if effect.Pos != (Vec2{X: 30, Y: 0}) || effect.Frame != 1 || effect.Facing != -1 {
				t.Fatalf("hit effect = %+v, want copied target position/frame/facing", effect)
			}
		}
	}
	if !hasBolt || !hasHit {
		t.Fatalf("effects missing bolt or hit duplicate: hasBolt=%v hasHit=%v effects=%+v", hasBolt, hasHit, g.effects)
	}
}

func TestProjectileRemovalPreservesOriginalOrder(t *testing.T) {
	g := New()
	g.fireball = []Fireball{
		{Pos: Vec2{X: 1}},
		{Pos: Vec2{X: 2}},
		{Pos: Vec2{X: 3}},
	}

	g.removeFireball(0)

	if len(g.fireball) != 2 || g.fireball[0].Pos.X != 2 || g.fireball[1].Pos.X != 3 {
		t.Fatalf("fireball order after removal = %+v, want X positions 2 then 3", g.fireball)
	}
}
