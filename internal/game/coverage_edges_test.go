package game

import (
	"bytes"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/gofont/gomonobold"
)

func TestCoverageDeathWaveLifecycleDamageAndDrawing(t *testing.T) {
	g := New()
	g.screenW = ScreenWidth
	g.screenH = ScreenHeight
	g.session.Progression.DeathWaveUnlocked = true
	g.session.Casts.DeathWave = g.tuning.DeathWaveInterval
	g.updateDeathWaveCasting(0)
	if len(g.deathWaves) != 1 {
		t.Fatalf("death waves = %d, want 1", len(g.deathWaves))
	}
	if g.deathWaves[0].MaxRadius < deathWaveVisibleRadius(g.screenW, g.screenH, g.tuning.SkeletonSpawnMargin) {
		t.Fatalf("death wave max radius = %v, want at least visible radius", g.deathWaves[0].MaxRadius)
	}

	g.skeleton = []Skeleton{
		{ID: 1, Kind: SkeletonRegular, Pos: Vec2{X: 10}, HP: 1},
		{ID: 2, Kind: SkeletonRed, Pos: Vec2{X: 10}, HP: 1},
		{ID: 3, Kind: SkeletonPurple, Pos: Vec2{X: 10}, HP: 10},
		{ID: 4, Kind: SkeletonBlack, Pos: Vec2{X: 1000}, HP: 10},
	}
	g.deathWaves[0].Radius = 12
	g.applyDeathWaveDamage(&g.deathWaves[0])
	if got := g.skeleton[2].HP; got != 5 {
		t.Fatalf("purple skeleton hp = %d, want half damage to 5", got)
	}
	if g.skeleton[0].HP != 1 || g.skeleton[1].HP != 1 || g.skeleton[3].HP != 10 {
		t.Fatalf("death wave damaged excluded skeletons: %+v", g.skeleton)
	}
	g.applyDeathWaveDamage(&g.deathWaves[0])
	if got := g.skeleton[2].HP; got != 5 {
		t.Fatalf("already-hit skeleton hp = %d, want unchanged", got)
	}
	g.applyDeathWaveDamage(nil)

	screen := ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.drawDeathWave(screen, DeathWave{})
	g.drawDeathWave(screen, DeathWave{Origin: g.player.Pos, Radius: 1, MaxRadius: 1})
	g.drawDeathWave(screen, DeathWave{Origin: g.player.Pos, Radius: 12, MaxRadius: 100})
	g.drawDeathWave(screen, DeathWave{Origin: g.player.Pos, Radius: 32, MaxRadius: 100})
	g.drawDeathWaveSparks(screen, DeathWave{Origin: g.player.Pos, Radius: 10, MaxRadius: 100}, 255)
	g.drawWorld(screen)

	g.deathWaves = []DeathWave{{Origin: g.player.Pos, Radius: 20, MaxRadius: 30}}
	g.updateDeathWaves(1)
	if len(g.deathWaves) != 0 {
		t.Fatalf("expired death waves = %d, want 0", len(g.deathWaves))
	}

	locked := New()
	locked.session.Casts.DeathWave = 99
	locked.updateDeathWaveCasting(1)
	if locked.session.Casts.DeathWave != 0 {
		t.Fatalf("locked death wave cast timer = %v, want reset", locked.session.Casts.DeathWave)
	}
}

func TestCoverageProgressionAndSpawnPressureEdges(t *testing.T) {
	tuning := DefaultTuning()
	p := NewProgression(tuning)
	if got := p.MeteorSpawnInterval(); !math.IsInf(got, 1) {
		t.Fatalf("locked MeteorSpawnInterval = %v, want +Inf", got)
	}
	p.ApplyLevelUpOption(LearnMeteor)
	p.ApplyLevelUpOption(ExtraMeteor)
	if got, want := p.MeteorSpawnInterval(), p.MeteorCastInterval()/2; math.Abs(got-want) > 0.000001 {
		t.Fatalf("unlocked MeteorSpawnInterval = %v, want %v", got, want)
	}
	if got := p.GainExperienceToLevel(1); got != 0 {
		t.Fatalf("GainExperienceToLevel current = %d, want 0", got)
	}
	if got := p.UpgradeAllProperties(LearnedSkill(999)); got != nil {
		t.Fatalf("unknown skill upgrades = %v, want nil", got)
	}
	capped := NewProgression(tuning)
	capped.tuning.InitialFireballCast = 0.018
	capped.tuning.FireballIntervalMultiplier = 0.9
	applied := capped.UpgradeAllProperties(SkillFireball)
	if len(applied) != 1 || applied[0] != ExtraFireball {
		t.Fatalf("capped fireball upgrades = %v, want only ExtraFireball", applied)
	}
	if p.skillUnlocked(LearnedSkill(999)) {
		t.Fatal("unknown skill reported unlocked")
	}
	if got := orbitalOrbHitInterval(0); got != 0 {
		t.Fatalf("orbitalOrbHitInterval(0) = %v, want 0", got)
	}

	g := New()
	g.skeletonHPPerSecond = -10
	g.dynamicSpawnQueue = []dynamicSpawnPlanEntry{{Kind: SkeletonRegular, Count: 1}}
	g.session.Casts.SkeletonSpawn = 99
	g.updateSkeletonSpawning(1)
	if g.session.Casts.SkeletonSpawn != 0 || len(g.dynamicSpawnQueue) != 0 {
		t.Fatalf("zero spawn pressure state = casts %v queue %v", g.session.Casts.SkeletonSpawn, g.dynamicSpawnQueue)
	}
	if got := g.SkeletonSpawnInterval(); !math.IsInf(got, 1) {
		t.Fatalf("zero spawn interval = %v, want +Inf", got)
	}

	g = New()
	g.skeletonHPPerSecond = 2
	g.session.Casts.SkeletonSpawn = 0
	g.dynamicSpawnQueue = []dynamicSpawnPlanEntry{{Kind: SkeletonBlue, Count: 1}}
	beforeSkeletons := len(g.skeleton)
	g.spawnQueuedDynamicSkeletons()
	if len(g.skeleton) != beforeSkeletons || len(g.dynamicSpawnQueue) != 1 {
		t.Fatalf("underfunded spawn changed skeletons=%d queue=%v", len(g.skeleton), g.dynamicSpawnQueue)
	}
	g.session.Casts.SkeletonSpawn = 5000
	g.tuning.MaxSkeletonSpawnsPerTick = 1
	g.dynamicSpawnQueue = []dynamicSpawnPlanEntry{{Kind: SkeletonRegular, Count: 2}}
	beforeSkeletons = len(g.skeleton)
	g.spawnQueuedDynamicSkeletons()
	if len(g.skeleton) != beforeSkeletons+1 || len(g.dynamicSpawnQueue) != 1 || g.dynamicSpawnQueue[0].Count != 1 {
		t.Fatalf("per-tick limited spawn skeletons=%d queue=%v", len(g.skeleton), g.dynamicSpawnQueue)
	}
	if g.canSpawnDynamicSkeleton(1) {
		t.Fatal("canSpawnDynamicSkeleton ignored max per tick")
	}
	g.tuning.MaxActiveSkeletons = 1
	if g.canSpawnDynamicSkeleton(0) {
		t.Fatal("canSpawnDynamicSkeleton ignored max active skeletons")
	}

	g = New()
	g.queueDynamicSpawnPressureForLevelUp(0)
	if g.pendingSpawnPressureLevels != 0 {
		t.Fatalf("zero level-up queued pressure levels = %d, want 0", g.pendingSpawnPressureLevels)
	}
	g.queueDynamicSpawnPressureForLevelUp(-1)
	if g.pendingSpawnPressureLevels != 0 {
		t.Fatalf("negative level-up queued pressure levels = %d, want 0", g.pendingSpawnPressureLevels)
	}

	g = New()
	g.maxActualDPS = 8
	g.queueDynamicSpawnPressureForLevelUp(2)
	if g.pendingSpawnPressureActual != 8 || g.pendingSpawnPressureLevels != 2 || g.maxActualDPS != 0 {
		t.Fatalf("queued spawn pressure state = actual %v levels %d max %v", g.pendingSpawnPressureActual, g.pendingSpawnPressureLevels, g.maxActualDPS)
	}
	g.applyPendingDynamicSpawnPressure()
	if g.skeletonHPPerSecond <= 0 || g.pendingSpawnPressureLevels != 1 {
		t.Fatalf("actual-target spawn pressure = hp/s %v levels %d", g.skeletonHPPerSecond, g.pendingSpawnPressureLevels)
	}
	g.pendingSpawnPressureActual = 0
	g.applyPendingDynamicSpawnPressure()
	if g.pendingSpawnPressureActual != 0 || g.pendingSpawnPressureLevels != 0 {
		t.Fatalf("finished spawn pressure state = actual %v levels %d", g.pendingSpawnPressureActual, g.pendingSpawnPressureLevels)
	}

	if got := capSkeletonHPPerSecond(5, 0); got != 0 {
		t.Fatalf("capSkeletonHPPerSecond no DPS = %v, want 0", got)
	}
	if got := increaseSkeletonHPPerSecond(5, 1, 4); got != 5 {
		t.Fatalf("increase capped by current = %v, want 5", got)
	}
	if got := dynamicSpawnPressureActualTarget(10, 100, -1); got != 10 {
		t.Fatalf("negative-factor actual target = %v, want raw DPS 10", got)
	}
}

func TestCoverageRenderingAndFontEdges(t *testing.T) {
	g := New()
	g.screenW = 120
	g.screenH = 90
	g.chests = []Chest{{Pos: Vec2{X: 1000, Y: 0}, Tier: ChestGold}}
	g.coins = []Coin{{Pos: Vec2{X: -1000, Y: 0}, Phase: 0.13, Amount: 1}}
	screen := ebiten.NewImage(g.screenW, g.screenH)
	g.drawPickupIndicators(screen)
	g.drawPickupIndicatorBacking(screen, 20, 20)
	x, y := edgeIndicatorPosition(g.screenW, g.screenH, float64(g.screenW)/2, float64(g.screenH)/2, 20)
	if x != float64(g.screenW)/2 || y != float64(g.screenH)/2 {
		t.Fatalf("center edge indicator = %v,%v", x, y)
	}

	g.drawFireballImpact(screen, Effect{MaxTTL: 0})
	g.drawMeteorImpactGroundShake(screen, Effect{Kind: EffectMeteorImpact, Pos: Vec2{}, Radius: 48, TTL: 1, MaxTTL: 0})
	g.drawMeteorImpactGroundShake(screen, Effect{Kind: EffectMeteorImpact, Pos: Vec2{}, Radius: 65, TTL: 0.9, MaxTTL: 1})
	g.drawOptionIconTinted(screen, BuyDeathWaveScroll, 10, 10, color.RGBA{255, 255, 255, 255})
	if got := flashActionAlpha(-1, 1, 0.1, 0.1); got != 255 {
		t.Fatalf("negative flash alpha = %d, want 255", got)
	}
	if got := flashActionAlpha(0.15, 1, 0.1, 0.2); got >= 255 {
		t.Fatalf("fade-up flash alpha = %d, want below 255", got)
	}

	oldPaths := hudFontPaths
	oldFont := hudFont
	oldSource := hudFontSource
	oldFaces := hudFontFaces
	defer func() {
		hudFontPaths = oldPaths
		hudFont = oldFont
		hudFontSource = oldSource
		hudFontFaces = oldFaces
	}()

	fontPath := filepath.Join(t.TempDir(), "menlo.ttf")
	fontData := bytes.ReplaceAll(gomonobold.TTF, []byte("Go Mono Bold"), []byte("Menlo GoBold"))
	fontData = bytes.ReplaceAll(fontData, []byte{0, 'G', 0, 'o', 0, ' ', 0, 'M', 0, 'o', 0, 'n', 0, 'o', 0, ' ', 0, 'B', 0, 'o', 0, 'l', 0, 'd'}, []byte{0, 'M', 0, 'e', 0, 'n', 0, 'l', 0, 'o', 0, ' ', 0, 'G', 0, 'o', 0, 'B', 0, 'o', 0, 'l', 0, 'd'})
	if err := os.WriteFile(fontPath, fontData, 0o600); err != nil {
		t.Fatal(err)
	}
	hudFontPaths = []string{fontPath}
	loadedFont, name := loadHUDFont()
	if loadedFont == nil || name != "Menlo GoBold" {
		t.Fatalf("system HUD font = %v %q, want patched Menlo font", loadedFont, name)
	}
	if loadedFont, name := loadSystemFontByFullName([]string{fontPath}, []string{"not-present"}); loadedFont != nil || name != "" {
		t.Fatalf("non-matching system font = %v %q, want none", loadedFont, name)
	}
}

func TestCoverageWeaponReservationAndVisualEdges(t *testing.T) {
	g := New()
	g.fireballTargetReservations = nil
	g.session.Progression.SimultaneousFireball = 1
	g.skeleton = []Skeleton{{ID: 1, Pos: Vec2{X: 1}}}
	g.spawnFireballs()
	if g.fireballTargetReservations == nil || len(g.fireball) != 1 {
		t.Fatalf("nil reservation fireball state = reservations %v fireballs %+v", g.fireballTargetReservations, g.fireball)
	}

	g = New()
	g.lightningTargetReservations = map[int]bool{7: true}
	g.session.Progression.SimultaneousFireball = 2
	g.skeleton = []Skeleton{{ID: 7, Pos: Vec2{X: 1}}, {ID: 8, Pos: Vec2{X: 2}}, {ID: 9, Pos: Vec2{X: 3}}}
	g.spawnFireballs()
	if len(g.fireball) != 2 || g.fireball[0].TargetID == 7 || g.fireball[1].TargetID == 7 {
		t.Fatalf("reserved target fireballs = %+v, want target 7 skipped", g.fireball)
	}

	g = New()
	g.lightningTargetReservations = nil
	g.reserveLightningTargets(nil)
	if g.lightningTargetReservations != nil {
		t.Fatal("empty lightning reservation initialized map")
	}
	g.reserveLightningTargets([]lightningStrikeTarget{{targetID: 12}})
	if !g.lightningTargetReservations[12] {
		t.Fatalf("lightning reservations = %v, want target 12", g.lightningTargetReservations)
	}
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.fireballTargetReservations = nil
	g.skeleton = []Skeleton{{ID: 1, Pos: Vec2{X: 10}, HP: 1}, {ID: 2, Pos: Vec2{X: 20}, HP: 2}}
	targets := g.chainLightningTargets()
	if len(targets) != 2 || g.fireballTargetReservations == nil {
		t.Fatalf("lightning targets = %v reservations=%v, want two with initialized reservations", targets, g.fireballTargetReservations)
	}

	g = New()
	g.player.Facing = 1
	if got := g.beamTargets(Vec2{X: 1}, 100, g.tuning.BeamHitWidth, 0); got != nil {
		t.Fatalf("zero-limit beam targets = %v, want nil", got)
	}
	g.skeleton = []Skeleton{{ID: 20, Pos: Vec2{X: 25}, HP: 0}, {ID: 21, Pos: Vec2{X: 50}, HP: 3}}
	if got := g.beamVisualEnd(Vec2{X: 1}, 100, []int{20, 21}, 3); got != (Vec2{X: 50}) {
		t.Fatalf("beam visual end = %+v, want x=50 after skipping zero-HP target", got)
	}
	if got := g.beamVisualEnd(Vec2{X: 1}, 100, []int{999, 21}, 3); got != (Vec2{}) {
		t.Fatalf("beam visual end missing target = %+v, want player position", got)
	}

	g = New()
	g.session.Progression.ApplyLevelUpOption(LearnMeteor)
	if got := g.session.Progression.MeteorSpawnInterval(); got <= 0 || math.IsInf(got, 1) {
		t.Fatalf("meteor spawn interval = %v, want finite positive", got)
	}
}

func TestCoverageCombatAndTypeEdges(t *testing.T) {
	g := New()
	g.recordActualDamage(0)
	if len(g.actualDamage) != 0 {
		t.Fatalf("zero damage samples = %v, want none", g.actualDamage)
	}
	g.totalTime = actualDPSWindow + 1
	g.actualDamage = []actualDamageSample{{Time: 0, Amount: 10}}
	g.actualDamageWindowTotal = 5
	g.pruneActualDamageSamples()
	if g.actualDamageWindowTotal != 0 || len(g.actualDamage) != 0 {
		t.Fatalf("pruned damage state = total %d samples %v, want empty nonnegative", g.actualDamageWindowTotal, g.actualDamage)
	}

	tuning := DefaultTuning()
	tuning.RedHitPoints = 0
	tuning.PurpleHitPoints = -2
	tuning.BlackHitPoints = -3
	for _, kind := range []SkeletonKind{SkeletonRegular, SkeletonRed, SkeletonPurple, SkeletonBlack} {
		if got := kind.HitPoints(tuning); got != 1 {
			t.Fatalf("%v hit points = %d, want clamped 1", kind, got)
		}
	}
	if got := SkeletonBlue.ExperienceReward(); got != 75 {
		t.Fatalf("blue reward = %d, want 75", got)
	}
	if got := SkeletonKind(999).SpeedMultiplier(); got != 1 {
		t.Fatalf("unknown speed multiplier = %v, want 1", got)
	}

	g = New()
	g.spawnGoldChestsAroundPlayer(0, 10)
	g.spawnGoldChestsAroundPlayer(2, -10)
	if len(g.chests) != 2 || !reflect.DeepEqual(g.chests[0].Pos, g.player.Pos) || !reflect.DeepEqual(g.chests[1].Pos, g.player.Pos) {
		t.Fatalf("gold chests around player = %+v, want two at player position", g.chests)
	}

	New().applyChestReward(ChestTier(999))
}
