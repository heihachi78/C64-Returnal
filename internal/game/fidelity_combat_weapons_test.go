package game

import (
	"math"
	"math/rand"
	"runtime"
	"slices"
	"testing"
)

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

func TestSkeletonKindsMoveAtConfiguredSpeedMultipliers(t *testing.T) {
	g := New()
	g.player.Pos = Vec2{X: 1000}
	g.tuning.SkeletonSpeed = 100
	g.skeleton = []Skeleton{
		{ID: 1, Kind: SkeletonRegular},
		{ID: 2, Kind: SkeletonRed},
		{ID: 3, Kind: SkeletonPurple},
		{ID: 4, Kind: SkeletonBlack},
		{ID: 5, Kind: SkeletonBlue},
	}

	g.updateSkeletons(1)

	tests := []struct {
		kind SkeletonKind
		want float64
	}{
		{kind: SkeletonRegular, want: 100},
		{kind: SkeletonRed, want: 99},
		{kind: SkeletonPurple, want: 97},
		{kind: SkeletonBlack, want: 94},
		{kind: SkeletonBlue, want: 90},
	}
	for i, tt := range tests {
		if got := g.skeleton[i].Pos.X; math.Abs(got-tt.want) > 0.0001 {
			t.Fatalf("kind %v moved %v, want %v", tt.kind, got, tt.want)
		}
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

func TestTimedSkeletonSpawnDefersSpatialRefreshLikeOriginal(t *testing.T) {
	g := New()
	g.skeleton = g.skeleton[:0]
	g.spatial.Rebuild(g.skeleton)
	redHP := SkeletonRed.HitPoints(g.tuning)
	g.dynamicSpawnQueue = []dynamicSpawnPlanEntry{{Kind: SkeletonRed, Count: 1}}
	g.session.Casts.SkeletonSpawn = float64(redHP)

	g.updateSkeletonSpawning(0)

	if got, want := len(g.skeleton), 1; got != want {
		t.Fatalf("timed skeleton count = %d, want %d", got, want)
	}
	if g.skeleton[0].Kind != SkeletonRed {
		t.Fatalf("timed skeleton kind = %v, want red", g.skeleton[0].Kind)
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

func TestDestroyDoesNotTriggerMilestoneEnemySpawn(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{}, HP: 1, Reward: 0}}
	g.spatial.Rebuild(g.skeleton)
	g.session.Kills.TotalSkeletons = 0
	g.session.NextChestMilestone = 1_000_000

	g.destroySkeleton(0, AttackFireball)

	if len(g.skeleton) != 0 {
		t.Fatalf("post-kill skeletons = %+v, want none", g.skeleton)
	}
}

func TestDestroySkeletonDefersSpatialRebuildUntilNextQuery(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{}, HP: 1, Reward: 0},
		{ID: 202, Pos: Vec2{X: 200}, HP: 1, Reward: 0},
	}
	g.rebuildSkeletonSpatialIndex()

	g.destroySkeleton(0, AttackNone)

	if !g.skeletonSpatialDirty {
		t.Fatal("skeleton spatial dirty = false, want true after removal")
	}
	if got := g.firstSkeletonHitByPoint(Vec2{X: 200}, g.tuning.SkeletonHitDistance); got != 0 {
		t.Fatalf("hit query after dirty removal = %d, want remaining skeleton index 0", got)
	}
	if g.skeletonSpatialDirty {
		t.Fatal("skeleton spatial dirty = true, want false after query rebuild")
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

func TestDeathWaveCastingCreatesWaveEveryThirtySeconds(t *testing.T) {
	g := New()
	g.session.Progression.DeathWaveUnlocked = true

	g.updateDeathWaveCasting(g.tuning.DeathWaveInterval*2 + 0.25)

	if len(g.deathWaves) != 2 {
		t.Fatalf("death wave count = %d, want 2", len(g.deathWaves))
	}
	if math.Abs(g.session.Casts.DeathWave-0.25) > 0.0001 {
		t.Fatalf("death wave cast remainder = %v, want 0.25", g.session.Casts.DeathWave)
	}
}

func TestDeathWaveRemovesHalfHPFromTouchedNonWhiteSkeletons(t *testing.T) {
	g := New()
	g.tuning.DeathWaveWidth = 20
	g.skeleton = []Skeleton{
		{ID: 101, Kind: SkeletonRed, Pos: Vec2{X: 48}, HP: 3},
		{ID: 202, Kind: SkeletonRegular, Pos: Vec2{X: 52}, HP: 3},
		{ID: 303, Kind: SkeletonBlack, Pos: Vec2{X: 180}, HP: 29},
	}
	wave := DeathWave{Origin: Vec2{}, PreviousRadius: 0, Radius: 60}

	g.applyDeathWaveDamage(&wave)

	if g.skeleton[0].HP != 2 {
		t.Fatalf("red skeleton HP = %d, want 2", g.skeleton[0].HP)
	}
	if g.skeleton[1].HP != 3 {
		t.Fatalf("regular white skeleton HP = %d, want immune unchanged 3", g.skeleton[1].HP)
	}
	if g.skeleton[2].HP != 29 {
		t.Fatalf("outside black skeleton HP = %d, want 29", g.skeleton[2].HP)
	}
	if got, want := g.actualDamageWindowTotal, 1; got != want {
		t.Fatalf("recorded death wave damage = %d, want %d", got, want)
	}
	if !slices.Equal(wave.HitIDs, []int{101}) {
		t.Fatalf("death wave hit IDs = %v, want [101]", wave.HitIDs)
	}
}

func TestDeathWaveLeavesMinimumOneHP(t *testing.T) {
	g := New()
	g.tuning.DeathWaveWidth = 20
	g.skeleton = []Skeleton{{ID: 101, Kind: SkeletonRed, Pos: Vec2{X: 48}, HP: 2}}
	wave := DeathWave{Origin: Vec2{}, PreviousRadius: 0, Radius: 60}

	g.applyDeathWaveDamage(&wave)

	if g.skeleton[0].HP != 1 {
		t.Fatalf("red skeleton HP = %d, want minimum 1", g.skeleton[0].HP)
	}
	if got, want := g.actualDamageWindowTotal, 1; got != want {
		t.Fatalf("recorded death wave damage = %d, want %d", got, want)
	}
}

func TestDeathWaveDoesNotHitSameSkeletonTwice(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{{ID: 101, Kind: SkeletonRed, Pos: Vec2{X: 24}, HP: 5}}
	wave := DeathWave{Origin: Vec2{}, PreviousRadius: 0, Radius: 40}

	g.applyDeathWaveDamage(&wave)
	g.skeleton[0].HP = 5
	wave.PreviousRadius = wave.Radius
	wave.Radius = 80
	g.applyDeathWaveDamage(&wave)

	if g.skeleton[0].HP != 5 {
		t.Fatalf("same skeleton was hit twice; HP = %d, want reset value 5", g.skeleton[0].HP)
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

func TestMeteorImpactRadiusUsesInclusiveSkeletonBodyOverlap(t *testing.T) {
	g := New()
	g.tuning.MeteorImpactRadius = 48
	g.tuning.SkeletonHitDistance = 24
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 72}, HP: 1},
		{ID: 202, Pos: Vec2{X: 72.0001}, HP: 1},
		{ID: 303, Pos: Vec2{Y: -72}, HP: 1},
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

func TestBlueMonsterHitboxScalesWithLargerSprite(t *testing.T) {
	g := New()
	g.tuning.SkeletonHitDistance = 24

	if got, want := skeletonBodyRadius(g.tuning, SkeletonBlue), skeletonBodyRadius(g.tuning, SkeletonRegular)*skeletonSpriteScale(SkeletonBlue); got != want {
		t.Fatalf("blue body radius = %v, want %v", got, want)
	}

	g.skeleton = []Skeleton{{ID: 101, Kind: SkeletonRegular, Pos: Vec2{X: 72}, HP: 1}}
	g.spatial.Rebuild(g.skeleton)
	if idx := g.firstSkeletonHitByPoint(Vec2{}, g.tuning.SkeletonHitDistance); idx >= 0 {
		t.Fatalf("regular skeleton hit at blue-only distance index %d, want miss", idx)
	}

	g.skeleton = []Skeleton{{ID: 202, Kind: SkeletonBlue, Pos: Vec2{X: 72}, HP: 1}}
	g.spatial.Rebuild(g.skeleton)
	if idx := g.firstSkeletonHitByPoint(Vec2{}, g.tuning.SkeletonHitDistance); idx != 0 {
		t.Fatalf("blue skeleton hit at larger boundary index %d, want 0", idx)
	}
}

func TestBlueMonsterLargerHitboxAppliesToWeaponHits(t *testing.T) {
	g := New()
	g.tuning.SkeletonHitDistance = 24
	g.tuning.FireballHitDistance = 20
	g.tuning.BeamHitWidth = 18
	g.tuning.MeteorImpactRadius = 48

	g.skeleton = []Skeleton{{ID: 101, Kind: SkeletonBlue, Pos: Vec2{X: 68}, HP: 1}}
	g.spatial.Rebuild(g.skeleton)
	if idx := g.firstSkeletonHitBySegment(Vec2{Y: -10}, Vec2{Y: 10}, g.tuning.FireballHitDistance); idx != 0 {
		t.Fatalf("fireball segment blue hit index = %d, want 0", idx)
	}

	g.skeleton = []Skeleton{{ID: 202, Kind: SkeletonBlue, Pos: Vec2{X: 100, Y: 66}, HP: 1}}
	g.spatial.Rebuild(g.skeleton)
	targets := g.beamTargets(Vec2{X: 1}, 140, g.tuning.BeamHitWidth, 1)
	if len(targets) != 1 || targets[0] != 202 {
		t.Fatalf("beam blue targets = %v, want [202]", targets)
	}

	g.skeleton = []Skeleton{
		{ID: 303, Kind: SkeletonBlue, Pos: Vec2{X: 120}, HP: 1},
		{ID: 404, Kind: SkeletonBlue, Pos: Vec2{X: 120.0001}, HP: 1},
	}
	g.spatial.Rebuild(g.skeleton)
	targets = g.meteorImpactTargetIDs(Vec2{})
	if len(targets) != 1 || targets[0] != 303 {
		t.Fatalf("meteor blue targets = %v, want [303]", targets)
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

func TestMeteorStopsDamagingWhenKillQueuesLevelUp(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: -10}, HP: 1, Reward: 1},
		{ID: 202, Pos: Vec2{X: 10}, HP: 1, Reward: 1},
	}
	g.spatial.Rebuild(g.skeleton)

	g.impactMeteor(Vec2{})

	if got, want := g.session.Kills.Meteor, 1; got != want {
		t.Fatalf("meteor kills = %d, want %d", got, want)
	}
	if got, want := len(g.skeleton), 1; got != want {
		t.Fatalf("remaining skeletons = %d, want %d", got, want)
	}
	if !g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) != 1 {
		t.Fatalf("level-up state active=%v pending=%v, want one queued level-up", g.session.LevelUpChoiceActive, g.session.PendingLevelUpLevels)
	}
}

func TestOrbitalOrbsStopDamagingWhenKillQueuesLevelUp(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnOrb)
	g.session.Progression.ApplyLevelUpOption(ExtraOrb)
	g.orbs = []OrbitalOrb{
		{Pos: Vec2{X: -10}, Active: true},
		{Pos: Vec2{X: 10}, Active: true},
	}
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: -10}, HP: 1, Reward: 1},
		{ID: 202, Pos: Vec2{X: 10}, HP: 1, Reward: 1},
	}
	g.spatial.Rebuild(g.skeleton)

	g.checkOrbitalOrbCollisions()

	if got, want := g.session.Kills.OrbitalOrb, 1; got != want {
		t.Fatalf("orbital kills = %d, want %d", got, want)
	}
	if got, want := len(g.skeleton), 1; got != want {
		t.Fatalf("remaining skeletons = %d, want %d", got, want)
	}
	if !g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) != 1 {
		t.Fatalf("level-up state active=%v pending=%v, want one queued level-up", g.session.LevelUpChoiceActive, g.session.PendingLevelUpLevels)
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
		t.Fatalf("inner lightning path reused outer path; want separate random bolt shapes")
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

func TestBeamVisualStopsAtLastTargetCoveredByDamageBudget(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 60}, HP: 1},
		{ID: 202, Pos: Vec2{X: 120}, HP: 1},
	}

	g.castBeam()

	if len(g.effects) != 1 || g.effects[0].Kind != EffectBeam {
		t.Fatalf("effects = %+v, want one beam effect", g.effects)
	}
	if math.Abs(g.effects[0].End.X-60) > 0.0001 || math.Abs(g.effects[0].End.Y) > 0.0001 {
		t.Fatalf("beam visual end = %+v, want first target at x=60", g.effects[0].End)
	}
	if g.session.Kills.Beam != 1 {
		t.Fatalf("beam kills = %d, want gameplay kill unchanged at 1", g.session.Kills.Beam)
	}
}

func TestBeamVisualStopsAtPartiallyDamagedTargetWhenBudgetRunsOut(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.session.Progression.ApplyLevelUpOption(BeamKillCount)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 50}, HP: 1},
		{ID: 202, Pos: Vec2{X: 90}, HP: 5},
		{ID: 303, Pos: Vec2{X: 130}, HP: 1},
	}

	g.castBeam()

	if len(g.effects) != 1 || g.effects[0].Kind != EffectBeam {
		t.Fatalf("effects = %+v, want one beam effect", g.effects)
	}
	if math.Abs(g.effects[0].End.X-90) > 0.0001 || math.Abs(g.effects[0].End.Y) > 0.0001 {
		t.Fatalf("beam visual end = %+v, want partially damaged target at x=90", g.effects[0].End)
	}
	if g.session.Kills.Beam != 1 {
		t.Fatalf("beam kills = %d, want only first target killed", g.session.Kills.Beam)
	}
	idx := g.skeletonIndexByID(202)
	if idx < 0 {
		t.Fatal("partially damaged target was destroyed, want it alive")
	}
	if g.skeleton[idx].HP != 3 {
		t.Fatalf("partially damaged target HP = %d, want 3", g.skeleton[idx].HP)
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

func TestFireballVolleySkipsLightningReservedSkeletons(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 0}, HP: 2},
		{ID: 202, Pos: Vec2{X: 20, Y: 0}, HP: 2},
	}

	g.castLightning()
	g.spawnFireballs()

	if len(g.fireball) != 1 {
		t.Fatalf("fireball count = %d, want one fireball", len(g.fireball))
	}
	if g.fireball[0].TargetID != 202 {
		t.Fatalf("fireball target = %d, want non-lightning target 202", g.fireball[0].TargetID)
	}
}

func TestLightningStopsDamagingWhenKillQueuesLevelUp(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10}, HP: 1, Reward: 1},
		{ID: 202, Pos: Vec2{X: 20}, HP: 1, Reward: 1},
	}

	g.castLightning()

	if got, want := g.session.Kills.Lightning, 1; got != want {
		t.Fatalf("lightning kills = %d, want %d", got, want)
	}
	if got, want := len(g.skeleton), 1; got != want {
		t.Fatalf("remaining skeletons = %d, want %d", got, want)
	}
	if !g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) != 1 {
		t.Fatalf("level-up state active=%v pending=%v, want one queued level-up", g.session.LevelUpChoiceActive, g.session.PendingLevelUpLevels)
	}
}

func TestSameTickLightningAndFireballDoNotShareTargets(t *testing.T) {
	g := New()
	g.hasUpdated = true
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Casts.Lightning = g.session.Progression.LightningCastInterval()
	g.session.Casts.Fireball = g.session.Progression.FireballCastInterval()
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 100, Y: 0}, HP: 2, Facing: 1},
		{ID: 202, Pos: Vec2{X: 200, Y: 0}, HP: 2, Facing: 1},
	}
	g.spatial.Rebuild(g.skeleton)

	if err := g.Update(); err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if len(g.fireball) != 1 {
		t.Fatalf("fireball count = %d, want one fireball", len(g.fireball))
	}
	if g.fireball[0].TargetID == 101 {
		t.Fatalf("fireball targeted lightning-struck skeleton %d", g.fireball[0].TargetID)
	}
	if g.fireball[0].TargetID != 202 {
		t.Fatalf("fireball target = %d, want 202", g.fireball[0].TargetID)
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

func TestBeamStopsDamagingWhenKillQueuesLevelUp(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10}, HP: 1, Reward: 1},
		{ID: 202, Pos: Vec2{X: 20}, HP: 1, Reward: 1},
	}

	g.castBeam()

	if got, want := g.session.Kills.Beam, 1; got != want {
		t.Fatalf("beam kills = %d, want %d", got, want)
	}
	if got, want := len(g.skeleton), 1; got != want {
		t.Fatalf("remaining skeletons = %d, want %d", got, want)
	}
	if !g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) != 1 {
		t.Fatalf("level-up state active=%v pending=%v, want one queued level-up", g.session.LevelUpChoiceActive, g.session.PendingLevelUpLevels)
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

func TestLightningTargetsNearestFireballReservedSkeletonWithMoreThanOneHP(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 10, Y: 0}, HP: 2},
		{ID: 202, Pos: Vec2{X: 20, Y: 0}, HP: 1},
		{ID: 303, Pos: Vec2{X: 30, Y: 0}, HP: 1},
	}
	g.fireball = []Fireball{{TargetID: 101}}

	targets := g.chainLightningTargets()
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want nearest two targets", len(targets))
	}
	if targets[0].targetID != 101 || targets[1].targetID != 202 {
		t.Fatalf("lightning targets = %+v, want nearest targets 101 then 202", targets)
	}
}

func TestLightningSkipsFireballReservedSkeletonWithOneHP(t *testing.T) {
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
		t.Fatalf("target count = %d, want two non-finishing-fireball targets", len(targets))
	}
	if targets[0].targetID != 202 || targets[1].targetID != 303 {
		t.Fatalf("lightning targets = %+v, want non-finishing-fireball targets 202 then 303", targets)
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

func TestLightningChainsFromPreviouslyStruckEnemy(t *testing.T) {
	g := New()
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: 20, Y: 0}, HP: 1},
		{ID: 202, Pos: Vec2{X: 35, Y: 0}, HP: 1},
		{ID: 303, Pos: Vec2{X: 0, Y: 25}, HP: 1},
	}

	targets := g.chainLightningTargets()
	if len(targets) != 2 {
		t.Fatalf("target count = %d, want 2", len(targets))
	}
	if targets[0].targetID != 101 || targets[1].targetID != 202 {
		t.Fatalf("chain lightning targets = %+v, want 101 then next closest from 101: 202", targets)
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

func TestFireballHitCreatesImpactEffect(t *testing.T) {
	g := New()
	g.fireball = []Fireball{{Pos: Vec2{}, TargetID: 101}}
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 10}, HP: 2, Reward: 1}}

	g.updateHomingFireball(0, 0, 0)

	if len(g.effects) != 1 || g.effects[0].Kind != EffectFireballImpact {
		t.Fatalf("fireball impact effects = %+v, want one fireball impact", g.effects)
	}
	if g.effects[0].Pos != (Vec2{X: 10}) || g.effects[0].TTL != fireballImpactEffectDuration || g.effects[0].MaxTTL != fireballImpactEffectDuration {
		t.Fatalf("fireball impact effect = %+v, want target position and duration", g.effects[0])
	}
}

func TestUntargetedFireballHitCreatesImpactEffect(t *testing.T) {
	g := New()
	g.fireball = []Fireball{{Pos: Vec2{}, Velocity: Vec2{X: 1}}}
	g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 5}, HP: 2, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)

	g.updateUntargetedFireball(0, 0.1)

	if len(g.effects) != 1 || g.effects[0].Kind != EffectFireballImpact {
		t.Fatalf("untargeted fireball impact effects = %+v, want one fireball impact", g.effects)
	}
	if g.effects[0].Pos != (Vec2{X: 5}) {
		t.Fatalf("untargeted fireball impact position = %+v, want target position", g.effects[0].Pos)
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
