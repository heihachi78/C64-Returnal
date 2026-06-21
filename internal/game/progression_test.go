package game

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

func TestExperienceRequirementMatchesOriginalCurve(t *testing.T) {
	tests := map[int]int{
		1: 1,
		2: 2,
		3: 4,
		4: 8,
		5: 12,
	}
	for level, want := range tests {
		if got := ExperienceRequirement(level); got != want {
			t.Fatalf("ExperienceRequirement(%d) = %d, want %d", level, got, want)
		}
	}
}

func TestLevelUpOptionTitlesMatchOriginalHUDModels(t *testing.T) {
	tests := []struct {
		option        LevelUpOption
		beamKillBonus int
		want          string
	}{
		{option: FireRate, want: "FASTER FIRE"},
		{option: ExtraFireball, want: "+1 FIREBALL"},
		{option: ExtraLife, want: "+1 LIFE"},
		{option: HalveSkeletons, want: "HALVE HORDE"},
		{option: LearnLightning, want: "LEARN BOLT"},
		{option: LightningBounce, want: "+1 CHAIN"},
		{option: LightningRate, want: "FASTER BOLT"},
		{option: LearnOrb, want: "LEARN ORB"},
		{option: ExtraOrb, want: "+1 ORB"},
		{option: OrbitalSpeed, want: "FASTER ORB"},
		{option: LearnBeam, want: "LEARN BEAM"},
		{option: BeamRate, want: "FASTER BEAM"},
		{option: BeamKillCount, want: "+1 BEAM KILL"},
		{option: BeamKillCount, beamKillBonus: 4, want: "+4 BEAM KILL"},
		{option: LearnMeteor, want: "LEARN METEOR"},
		{option: ExtraMeteor, want: "+1 METEOR"},
		{option: MeteorRate, want: "FASTER METEOR"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.option.Title(tt.beamKillBonus); got != tt.want {
				t.Fatalf("title = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultBeamHitWidthIsMoreForgivingThanVisualBeam(t *testing.T) {
	if got, want := DefaultTuning().BeamHitWidth, 24.0; got != want {
		t.Fatalf("BeamHitWidth = %v, want %v", got, want)
	}
}

func TestLearnedSkillUpgradeOptionsMatchOriginalProgressionOrder(t *testing.T) {
	tests := []struct {
		skill LearnedSkill
		want  []LevelUpOption
	}{
		{skill: SkillFireball, want: []LevelUpOption{FireRate, ExtraFireball}},
		{skill: SkillLightning, want: []LevelUpOption{LightningBounce, LightningRate}},
		{skill: SkillOrbitalOrb, want: []LevelUpOption{ExtraOrb, OrbitalSpeed}},
		{skill: SkillBeam, want: []LevelUpOption{BeamRate, BeamKillCount}},
		{skill: SkillMeteor, want: []LevelUpOption{ExtraMeteor, MeteorRate}},
	}

	for _, tt := range tests {
		got := tt.skill.UpgradeOptions()
		if len(got) != len(tt.want) {
			t.Fatalf("skill %v option count = %d, want %d", tt.skill, len(got), len(tt.want))
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Fatalf("skill %v options = %v, want %v", tt.skill, got, tt.want)
			}
		}
	}
}

func TestChestTierTitlesMatchOriginalTierModel(t *testing.T) {
	tests := []struct {
		tier ChestTier
		want string
	}{
		{tier: ChestBronze, want: "BRONZE"},
		{tier: ChestSilver, want: "SILVER"},
		{tier: ChestGold, want: "GOLD"},
	}

	for _, tt := range tests {
		if got := tt.tier.Title(); got != tt.want {
			t.Fatalf("tier title = %q, want %q", got, tt.want)
		}
		if got := chestRewardTitle(tt.tier); got != tt.want+" CHEST" {
			t.Fatalf("chest reward title = %q, want %q", got, tt.want+" CHEST")
		}
	}
}

func TestProgressionGainExperienceQueuesMultipleLevels(t *testing.T) {
	p := NewProgression(DefaultTuning())
	if got := p.GainExperience(3); got != 2 {
		t.Fatalf("GainExperience(3) = %d, want 2", got)
	}
	if p.Level != 3 {
		t.Fatalf("level = %d, want 3", p.Level)
	}
	if p.Experience != 0 {
		t.Fatalf("experience remainder = %d, want 0", p.Experience)
	}
}

func TestProgressionGainExperienceMatchesOriginalEightExperienceRollover(t *testing.T) {
	p := NewProgression(DefaultTuning())

	if got := p.GainExperience(8); got != 3 {
		t.Fatalf("GainExperience(8) = %d, want 3", got)
	}
	if p.Level != 4 {
		t.Fatalf("level = %d, want 4", p.Level)
	}
	if p.Experience != 1 {
		t.Fatalf("experience remainder = %d, want 1", p.Experience)
	}
	if p.NextExperience != 8 {
		t.Fatalf("next experience = %d, want 8", p.NextExperience)
	}
}

func TestBeamKillUpgradeUsesGrowingBonus(t *testing.T) {
	p := NewProgression(DefaultTuning())
	p.ApplyLevelUpOption(LearnBeam)
	p.ApplyLevelUpOption(BeamKillCount)
	p.ApplyLevelUpOption(BeamKillCount)

	if got, want := p.BeamKillCount(), 6; got != want {
		t.Fatalf("BeamKillCount = %d, want %d", got, want)
	}
	if got, want := p.BeamKillUpgradeBonus(), 4; got != want {
		t.Fatalf("BeamKillUpgradeBonus = %d, want %d", got, want)
	}
}

func TestUnlockAndUpgradeBeamMatchesOriginalKillCountAndInterval(t *testing.T) {
	tuning := DefaultTuning()
	p := NewProgression(tuning)

	p.ApplyLevelUpOption(LearnBeam)
	p.ApplyLevelUpOption(BeamKillCount)
	p.ApplyLevelUpOption(BeamRate)

	if !p.BeamUnlocked {
		t.Fatal("beam unlocked = false, want true")
	}
	if got, want := p.BeamKillCount(), 3; got != want {
		t.Fatalf("BeamKillCount = %d, want %d", got, want)
	}
	if got := p.BeamCastInterval(); got >= tuning.InitialBeamCast {
		t.Fatalf("BeamCastInterval = %v, want less than original %v after rate upgrade", got, tuning.InitialBeamCast)
	}
}

func TestMageRawDPSStartsWithLevelOneFireballOnly(t *testing.T) {
	p := NewProgression(DefaultTuning())

	if got, want := p.MageRawDPS(), 2.0/actualDPSWindow; math.Abs(got-want) > 0.000001 {
		t.Fatalf("MageRawDPS = %v, want %v", got, want)
	}
}

func TestMageRawDPSScalesWithFireballCountAndRate(t *testing.T) {
	p := NewProgression(DefaultTuning())
	p.ApplyLevelUpOption(ExtraFireball)

	if got, want := p.MageRawDPS(), 4.0/actualDPSWindow; math.Abs(got-want) > 0.000001 {
		t.Fatalf("extra-fireball MageRawDPS = %v, want %v", got, want)
	}

	p = NewProgression(DefaultTuning())
	p.ApplyLevelUpOption(FireRate)
	if got, want := p.MageRawDPS(), 2.0/actualDPSWindow; math.Abs(got-want) > 0.000001 {
		t.Fatalf("fire-rate MageRawDPS = %v, want %v", got, want)
	}
}

func TestMageRawDPSIncludesUnlockedWeaponRates(t *testing.T) {
	tuning := DefaultTuning()
	p := NewProgression(tuning)
	p.ApplyLevelUpOption(LearnLightning)
	p.ApplyLevelUpOption(LightningBounce)
	p.ApplyLevelUpOption(LearnOrb)
	p.ApplyLevelUpOption(ExtraOrb)
	p.ApplyLevelUpOption(OrbitalSpeed)
	p.ApplyLevelUpOption(LearnBeam)
	p.ApplyLevelUpOption(BeamKillCount)
	p.ApplyLevelUpOption(LearnMeteor)
	p.ApplyLevelUpOption(ExtraMeteor)

	want := windowedDamageRate(1, p.FireballCastInterval()) +
		windowedDamageRate(2, p.LightningCastInterval()) +
		windowedDamageRate(2, orbitalOrbHitInterval(p.OrbitalAngularSpeed())) +
		windowedDamageRate(3, p.BeamCastInterval()) +
		windowedDamageRate(2, p.MeteorCastInterval())
	if got := p.MageRawDPS(); math.Abs(got-want) > 0.000001 {
		t.Fatalf("unlocked MageRawDPS = %v, want %v", got, want)
	}
}

func TestAttackSpeedOptionsStopBeforeOneSixtiethSecond(t *testing.T) {
	tuning := DefaultTuning()
	tuning.InitialFireballCast = 0.018
	tuning.FireballIntervalMultiplier = 0.9
	tuning.InitialLightningCast = 0.018
	tuning.LightningIntervalMultiplier = 0.9
	tuning.InitialBeamCast = 0.018
	tuning.BeamIntervalMultiplier = 0.9
	tuning.InitialMeteorCast = 0.018
	tuning.MeteorIntervalMultiplier = 0.9
	p := NewProgression(tuning)
	p.ApplyLevelUpOption(LearnLightning)
	p.ApplyLevelUpOption(LearnBeam)
	p.ApplyLevelUpOption(LearnMeteor)

	options := p.AvailableLevelUpOptions()
	for _, option := range []LevelUpOption{FireRate, LightningRate, BeamRate, MeteorRate} {
		if slices.Contains(options, option) {
			t.Fatalf("options = %v, want capped speed option %v removed", options, option)
		}
	}
	for _, option := range []LevelUpOption{ExtraFireball, LightningBounce, ExtraMeteor, ExtraLife, HalveSkeletons, BeamKillCount, LearnOrb} {
		if !slices.Contains(options, option) {
			t.Fatalf("options = %v, want remaining upgrade %v available", options, option)
		}
	}
}

func TestAttackSpeedUpgradeCanLandExactlyOnOneSixtiethSecond(t *testing.T) {
	tuning := DefaultTuning()
	tuning.InitialFireballCast = minAttackSpawnInterval / 0.5
	tuning.FireballIntervalMultiplier = 0.5
	p := NewProgression(tuning)

	if !p.LevelUpOptionAvailable(FireRate) {
		t.Fatal("fire rate upgrade that lands exactly on 1/60s was unavailable")
	}
	p.ApplyLevelUpOption(FireRate)
	if math.Abs(p.FireballCastInterval()-minAttackSpawnInterval) > 0.000001 {
		t.Fatalf("fireball interval = %v, want %v", p.FireballCastInterval(), minAttackSpawnInterval)
	}
	if p.LevelUpOptionAvailable(FireRate) {
		t.Fatal("fire rate upgrade below 1/60s remained available")
	}
}

func TestCappedAttackSpeedOptionsDoNotApplyDirectly(t *testing.T) {
	tuning := DefaultTuning()
	tuning.InitialFireballCast = 0.018
	tuning.FireballIntervalMultiplier = 0.9
	tuning.InitialLightningCast = 0.018
	tuning.LightningIntervalMultiplier = 0.9
	tuning.InitialBeamCast = 0.018
	tuning.BeamIntervalMultiplier = 0.9
	tuning.InitialMeteorCast = 0.018
	tuning.MeteorIntervalMultiplier = 0.9
	p := NewProgression(tuning)
	p.ApplyLevelUpOption(LearnLightning)
	p.ApplyLevelUpOption(LearnBeam)
	p.ApplyLevelUpOption(LearnMeteor)

	p.ApplyLevelUpOption(FireRate)
	p.ApplyLevelUpOption(ExtraFireball)
	p.ApplyLevelUpOption(LightningRate)
	p.ApplyLevelUpOption(LightningBounce)
	p.ApplyLevelUpOption(BeamRate)
	p.ApplyLevelUpOption(MeteorRate)
	p.ApplyLevelUpOption(ExtraMeteor)

	if got, want := p.FireballCastInterval(), tuning.InitialFireballCast; got != want {
		t.Fatalf("fireball interval = %v, want unchanged %v", got, want)
	}
	if got, want := p.SimultaneousFireball, 2; got != want {
		t.Fatalf("fireball count = %d, want upgraded %d", got, want)
	}
	if got, want := p.LightningCastInterval(), tuning.InitialLightningCast; got != want {
		t.Fatalf("lightning interval = %v, want unchanged %v", got, want)
	}
	if got, want := p.LightningStrikeCount(), 2; got != want {
		t.Fatalf("lightning strike count = %d, want upgraded %d", got, want)
	}
	if got, want := p.BeamCastInterval(), tuning.InitialBeamCast; got != want {
		t.Fatalf("beam interval = %v, want unchanged %v", got, want)
	}
	if got, want := p.MeteorCastInterval(), tuning.InitialMeteorCast; got != want {
		t.Fatalf("meteor interval = %v, want unchanged %v", got, want)
	}
	if got, want := p.MeteorCount(), 2; got != want {
		t.Fatalf("meteor count = %d, want upgraded %d", got, want)
	}
}

func TestDynamicSkeletonSpawnOrderUsesHighestHitPointsFirst(t *testing.T) {
	order := dynamicSkeletonSpawnOrder(DefaultTuning())

	for i := 1; i < len(order); i++ {
		previousHP := order[i-1].HitPoints(DefaultTuning())
		currentHP := order[i].HitPoints(DefaultTuning())
		if previousHP < currentHP {
			t.Fatalf("spawn order = %v, want descending HP", order)
		}
	}
	if order[0] != SkeletonBlue {
		t.Fatalf("first spawn kind = %v, want blue", order[0])
	}
}

func TestInitialDynamicSkeletonSpawnRateIsCappedByRawDPS(t *testing.T) {
	g := New()

	want := min(g.tuning.InitialSkeletonHPPerSecond, g.session.Progression.MageRawDPS())
	if got := g.SkeletonHPPerSecond(); math.Abs(got-want) > 0.000001 {
		t.Fatalf("initial skeleton hp/sec = %v, want %v", got, want)
	}
}

func TestDynamicSpawnPressureUsesPeakActualDPSOnLevelUp(t *testing.T) {
	g := New()
	g.tuning.DynamicSpawnPressureFactor = 1.1
	g.skeletonHPPerSecond = 0.1
	g.maxActualDPS = 0.4
	g.session.Progression.GainExperience(1)

	g.queueLevelUpChoices(1)

	if got, want := g.SkeletonHPPerSecond(), 0.1; math.Abs(got-want) > 0.000001 {
		t.Fatalf("queued skeleton hp/sec = %v, want unchanged %v", got, want)
	}

	g.applyLevelUpOption(ExtraFireball)

	want := 0.1 + 1.1*(g.session.Progression.MageRawDPS()-0.4)
	if got := g.SkeletonHPPerSecond(); math.Abs(got-want) > 0.000001 {
		t.Fatalf("dynamic skeleton hp/sec = %v, want %v", got, want)
	}
	if g.maxActualDPS != 0 {
		t.Fatalf("max actual dps after level transition = %v, want reset", g.maxActualDPS)
	}
	if g.pendingSpawnPressureLevels != 0 || g.pendingSpawnPressureActual != 0 {
		t.Fatalf("pending pressure = levels %d actual %v, want cleared", g.pendingSpawnPressureLevels, g.pendingSpawnPressureActual)
	}
}

func TestDynamicSpawnPressureIsCappedByTheoreticalDPS(t *testing.T) {
	g := New()
	g.tuning.DynamicSpawnPressureFactor = 2
	g.skeletonHPPerSecond = 0.1
	g.session.Progression.GainExperience(1)

	g.queueLevelUpChoices(1)
	g.applyLevelUpOption(ExtraLife)

	if got, want := g.SkeletonHPPerSecond(), g.session.Progression.MageRawDPS(); math.Abs(got-want) > 0.000001 {
		t.Fatalf("capped skeleton hp/sec = %v, want raw dps %v", got, want)
	}
}

func TestDynamicSpawnPressureDoesNotDecreaseWhenActualExceedsTheoretical(t *testing.T) {
	g := New()
	g.tuning.DynamicSpawnPressureFactor = 1.5
	g.skeletonHPPerSecond = 0.5
	g.maxActualDPS = 10
	g.session.Progression.GainExperience(1)

	g.queueLevelUpChoices(1)
	g.applyLevelUpOption(ExtraLife)

	if got, want := g.SkeletonHPPerSecond(), 0.5; math.Abs(got-want) > 0.000001 {
		t.Fatalf("skeleton hp/sec = %v, want unchanged %v", got, want)
	}
}

func TestDynamicSpawnPressureNeverReducesExistingRate(t *testing.T) {
	if got, want := increaseSkeletonHPPerSecond(2, 10, 1), 2.0; got != want {
		t.Fatalf("increased hp/sec = %v, want unchanged %v", got, want)
	}
}

func TestDynamicSpawnPlanGreedilyFillsHPBudgetWithLargestSkeletons(t *testing.T) {
	g := New()

	if got := countDynamicSkeletonSpawnPlan(g.tuning, 198, SkeletonBlack); got != 5 {
		t.Fatalf("black skeleton count = %d, want 5", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 198, SkeletonPurple); got != 6 {
		t.Fatalf("purple skeleton count = %d, want 6", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 198, SkeletonRed); got != 2 {
		t.Fatalf("red skeleton count = %d, want 2", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 198, SkeletonRegular); got != 5 {
		t.Fatalf("regular skeleton count = %d, want 5", got)
	}
}

func TestDynamicSpawningFlushesLowLevelBudgetAsRegularSkeletons(t *testing.T) {
	g := New()
	g.skeleton = g.skeleton[:0]
	g.skeletonHPPerSecond = 1

	g.updateSkeletonSpawning(1)

	if got := countSkeletonKind(g.skeleton, SkeletonRed); got != 0 {
		t.Fatalf("red skeleton count = %d, want 0", got)
	}
	if got, want := countSkeletonKind(g.skeleton, SkeletonRegular), 1; got != want {
		t.Fatalf("regular skeleton count = %d, want %d", got, want)
	}
	if got := g.session.Casts.SkeletonSpawn; got != 0 {
		t.Fatalf("remaining spawn budget = %v, want 0", got)
	}
}

func TestDynamicSpawningPacesQueuedSkeletonsByHitPoints(t *testing.T) {
	g := New()
	g.skeleton = g.skeleton[:0]
	g.tuning.MaxSkeletonSpawnsPerTick = 1
	redHP := SkeletonRed.HitPoints(g.tuning)
	g.skeletonHPPerSecond = float64(redHP * 2)

	g.updateSkeletonSpawning(float64(redHP-1) / g.skeletonHPPerSecond)

	if got := len(g.skeleton); got != 0 {
		t.Fatalf("skeleton count before red HP budget = %d, want 0", got)
	}
	if len(g.dynamicSpawnQueue) == 0 || g.dynamicSpawnQueue[0].Kind != SkeletonRed {
		t.Fatalf("dynamic spawn queue = %v, want red first", g.dynamicSpawnQueue)
	}

	g.updateSkeletonSpawning(1 / g.skeletonHPPerSecond)

	if got, want := len(g.skeleton), 1; got != want {
		t.Fatalf("skeleton count after red HP budget = %d, want %d", got, want)
	}
	if got := g.skeleton[0].Kind; got != SkeletonRed {
		t.Fatalf("spawned skeleton kind = %v, want red", got)
	}
}

func TestDynamicSpawnPlanSpendsSingleRedBudgetAsRegularSkeletons(t *testing.T) {
	g := New()
	redHP := SkeletonRed.HitPoints(g.tuning)
	budget := float64(redHP)

	if got := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRed); got != 0 {
		t.Fatalf("red skeleton count at red budget = %d, want 0", got)
	}
	if got, want := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRegular), redHP; got != want {
		t.Fatalf("regular skeleton count at red budget = %d, want %d", got, want)
	}
}

func TestDynamicSpawnPlanSpendsTwoRedBudgetAsOneRed(t *testing.T) {
	g := New()
	redHP := SkeletonRed.HitPoints(g.tuning)
	budget := float64(redHP * 2)

	if got, want := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRed), 1; got != want {
		t.Fatalf("red skeleton count at two-red budget = %d, want %d", got, want)
	}
	if got, want := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRegular), redHP; got != want {
		t.Fatalf("regular skeleton count at two-red budget = %d, want %d", got, want)
	}
}

func TestDynamicSpawnPlanSpendsPartialRedBudgetAsRegularSkeletons(t *testing.T) {
	g := New()
	budget := float64(SkeletonRed.HitPoints(g.tuning) - 1)

	if got := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRed); got != 0 {
		t.Fatalf("red skeleton count = %d, want 0", got)
	}
	if got, want := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRegular), SkeletonRed.HitPoints(g.tuning)-1; got != want {
		t.Fatalf("regular skeleton count = %d, want %d", got, want)
	}
}

func TestDynamicSpawnPlanRequiresTwoPurpleBudgetBeforePurpleSpawn(t *testing.T) {
	g := New()
	budget := 12.5

	if got := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonPurple); got != 0 {
		t.Fatalf("purple skeleton count at 12.5 HP budget = %d, want 0", got)
	}
	if got, want := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRed), 3; got != want {
		t.Fatalf("red skeleton count at 12.5 HP budget = %d, want %d", got, want)
	}
	if got, want := countDynamicSkeletonSpawnPlan(g.tuning, budget, SkeletonRegular), 3; got != want {
		t.Fatalf("regular skeleton count at 12.5 HP budget = %d, want %d", got, want)
	}
}

func TestDynamicSpawningLimitsSpawnCountPerTick(t *testing.T) {
	g := New()
	g.skeleton = g.skeleton[:0]
	g.tuning.MaxSkeletonSpawnsPerTick = 3
	g.tuning.MaxActiveSkeletons = 0
	g.session.Casts.SkeletonSpawn = 20
	g.dynamicSpawnQueue = dynamicSkeletonSpawnPlanEntries(g.tuning, 20)

	g.spawnQueuedDynamicSkeletons()

	if got, want := len(g.skeleton), 3; got != want {
		t.Fatalf("spawned skeleton count = %d, want %d", got, want)
	}
	if got := g.session.Casts.SkeletonSpawn; got <= 0 {
		t.Fatalf("remaining spawn budget = %v, want queued budget", got)
	}
}

func TestDynamicSpawningStopsAtActiveSkeletonLimit(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{
		{ID: 101, HP: 1},
		{ID: 202, HP: 1},
	}
	g.tuning.MaxActiveSkeletons = 4
	g.tuning.MaxSkeletonSpawnsPerTick = 0
	g.session.Casts.SkeletonSpawn = 20
	g.dynamicSpawnQueue = dynamicSkeletonSpawnPlanEntries(g.tuning, 20)

	g.spawnQueuedDynamicSkeletons()

	if got, want := len(g.skeleton), 4; got != want {
		t.Fatalf("active skeleton count = %d, want capped count %d", got, want)
	}
	if got := g.session.Casts.SkeletonSpawn; got <= 0 {
		t.Fatalf("remaining spawn budget = %v, want queued budget", got)
	}
}

func TestDynamicSpawnPlanRepeatsLargestAffordableSkeletonsBeforeRegulars(t *testing.T) {
	g := New()

	if got := countDynamicSkeletonSpawnPlan(g.tuning, 2000, SkeletonBlue); got != 1 {
		t.Fatalf("blue skeleton count = %d, want 1", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 2000, SkeletonBlack); got != 33 {
		t.Fatalf("black skeleton count = %d, want 33", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 2000, SkeletonPurple); got != 5 {
		t.Fatalf("purple skeleton count = %d, want 5", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 2000, SkeletonRed); got != 1 {
		t.Fatalf("red skeleton count = %d, want 1", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 2000, SkeletonRegular); got != 5 {
		t.Fatalf("regular skeleton count = %d, want 5", got)
	}
}

func TestDynamicSpawnPlanSpendsRegularOverflowOnStrongerSkeletons(t *testing.T) {
	g := New()

	if got := countDynamicSkeletonSpawnPlan(g.tuning, 41, SkeletonBlack); got != 0 {
		t.Fatalf("black skeleton count = %d, want 0", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 41, SkeletonPurple); got != 4 {
		t.Fatalf("purple skeleton count = %d, want 4", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 41, SkeletonRed); got != 3 {
		t.Fatalf("red skeleton count = %d, want 3", got)
	}
	if got := countDynamicSkeletonSpawnPlan(g.tuning, 41, SkeletonRegular); got != 4 {
		t.Fatalf("regular skeleton count = %d, want 4", got)
	}
}

func countSkeletonKind(skeletons []Skeleton, kind SkeletonKind) int {
	count := 0
	for _, skeleton := range skeletons {
		if skeleton.Kind == kind {
			count++
		}
	}
	return count
}

func TestBeamDamageBudgetCanPartiallyDamagePurpleSkeletonLikeOriginal(t *testing.T) {
	g := beamTestGame()
	g.session.Progression.ApplyLevelUpOption(BeamKillCount)
	g.skeleton = []Skeleton{{ID: 1, Kind: SkeletonPurple, Pos: Vec2{X: 60}, HP: SkeletonPurple.HitPoints(g.tuning), Reward: SkeletonPurple.ExperienceReward()}}

	g.castBeam()

	if len(g.skeleton) != 1 {
		t.Fatalf("skeleton count = %d, want 1 partially damaged purple skeleton", len(g.skeleton))
	}
	if got, want := g.skeleton[0].HP, 2; got != want {
		t.Fatalf("purple skeleton HP = %d, want %d", got, want)
	}
	if g.session.Kills.Beam != 0 {
		t.Fatalf("beam kills = %d, want 0", g.session.Kills.Beam)
	}
	if g.session.Progression.Experience != 0 {
		t.Fatalf("experience = %d, want 0", g.session.Progression.Experience)
	}
}

func TestBeamDamageBudgetKillsPurpleSkeletonWhenItCoversHitPointsLikeOriginal(t *testing.T) {
	g := beamTestGame()
	g.session.Progression.ApplyLevelUpOption(BeamKillCount)
	g.session.Progression.ApplyLevelUpOption(BeamKillCount)
	g.skeleton = []Skeleton{{ID: 1, Kind: SkeletonPurple, Pos: Vec2{X: 60}, HP: SkeletonPurple.HitPoints(g.tuning), Reward: SkeletonPurple.ExperienceReward()}}

	g.castBeam()

	if len(g.skeleton) != 0 {
		t.Fatalf("skeleton count = %d, want killed purple skeleton", len(g.skeleton))
	}
	if g.session.Kills.Beam != 1 {
		t.Fatalf("beam kills = %d, want 1", g.session.Kills.Beam)
	}
	if got, want := g.session.Progression.Experience, 3; got != want {
		t.Fatalf("experience = %d, want %d", got, want)
	}
}

func TestLevelUpChoicesDefaultToTwoAndChanceAddsThirdLikeOriginal(t *testing.T) {
	tests := []struct {
		name      string
		numerator int
		want      int
	}{
		{name: "default two", numerator: 0, want: 2},
		{name: "chance adds third", numerator: 1, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			g.tuning.ExtraOptionChanceNumerator = tt.numerator
			g.tuning.ExtraOptionChanceDenominator = 1
			g.rng = rand.New(rand.NewSource(3))
			g.skeleton = nil

			if got := len(g.randomLevelUpOptionsCandidate()); got != tt.want {
				t.Fatalf("option count = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestChestTierPolicy(t *testing.T) {
	tuning := DefaultTuning()
	tests := []struct {
		name      string
		milestone int
		level     int
		wantTier  ChestTier
		wantOK    bool
	}{
		{name: "bronze before cap", milestone: 250, level: 20, wantTier: ChestBronze, wantOK: true},
		{name: "bronze after cap suppressed", milestone: 250, level: 34, wantTier: ChestBronze, wantOK: false},
		{name: "silver before cap", milestone: 1000, level: 55, wantTier: ChestSilver, wantOK: true},
		{name: "silver after cap suppressed", milestone: 1000, level: 56, wantTier: ChestSilver, wantOK: false},
		{name: "gold always wins", milestone: 5000, level: 99, wantTier: ChestGold, wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTier, gotOK := chestTier(tuning, tt.milestone, tt.level)
			if gotTier != tt.wantTier || gotOK != tt.wantOK {
				t.Fatalf("chestTier() = (%v, %v), want (%v, %v)", gotTier, gotOK, tt.wantTier, tt.wantOK)
			}
		})
	}
}

func TestCoinSpawnsWellOutsideVisibleViewportOncePerLevelLikeOriginal(t *testing.T) {
	g := New()
	g.screenW = 640
	g.screenH = 480
	g.rng = rand.New(rand.NewSource(9))
	g.coins = nil
	g.session.SpawnedCoinLevels = map[int]bool{}

	g.spawnCoinForLevel(1)
	g.spawnCoinForLevel(1)

	if len(g.coins) != 1 {
		t.Fatalf("coin count = %d, want 1", len(g.coins))
	}
	pos := g.coins[0].Pos
	outsideHorizontalEdge := math.Abs(pos.X-g.player.Pos.X) >= float64(g.screenW)/2+g.tuning.CoinSpawnMargin
	outsideVerticalEdge := math.Abs(pos.Y-g.player.Pos.Y) >= float64(g.screenH)/2+g.tuning.CoinSpawnMargin
	if !outsideHorizontalEdge && !outsideVerticalEdge {
		t.Fatalf("coin position = %+v, want outside horizontal or vertical viewport edge", pos)
	}
}

func TestChestSpawnsOutsideVisibleViewport(t *testing.T) {
	g := New()
	g.screenW = 640
	g.screenH = 480
	g.rng = rand.New(rand.NewSource(11))
	g.chests = nil

	for i := 0; i < 20; i++ {
		g.spawnChest(ChestBronze)
	}

	for _, chest := range g.chests {
		outsideHorizontalEdge := math.Abs(chest.Pos.X-g.player.Pos.X) >= float64(g.screenW)/2+g.tuning.ChestSpawnMargin
		outsideVerticalEdge := math.Abs(chest.Pos.Y-g.player.Pos.Y) >= float64(g.screenH)/2+g.tuning.ChestSpawnMargin
		if !outsideHorizontalEdge && !outsideVerticalEdge {
			t.Fatalf("chest position = %+v, want outside horizontal or vertical viewport edge", chest.Pos)
		}
	}
}

func TestCoinPickupGrantsRewardAndRemovesCoinLikeOriginal(t *testing.T) {
	g := New()
	g.rng = rand.New(rand.NewSource(2))
	g.coins = nil
	g.session.SpawnedCoinLevels = map[int]bool{}
	g.spawnCoinForLevel(1)
	g.coins[0].Pos = g.player.Pos

	reward := g.coins[0].Amount
	g.checkCoinPickups()

	if reward < g.tuning.CoinMinimumReward || reward > g.tuning.CoinMaximumReward {
		t.Fatalf("reward = %d, want within configured range", reward)
	}
	if g.session.CollectedCoins != reward {
		t.Fatalf("collected coins = %d, want %d", g.session.CollectedCoins, reward)
	}
	if len(g.coins) != 0 {
		t.Fatalf("coin count after pickup = %d, want 0", len(g.coins))
	}
}

func TestLevelUpRedrawCostsCurrentPlayerLevelLikeOriginal(t *testing.T) {
	g := New()
	g.session.Progression.GainExperience(15)
	g.session.LevelUpChoiceActive = true
	g.session.PendingLevelUpLevels = []int{g.session.Progression.Level}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraFireball}
	g.rng = rand.New(rand.NewSource(4))

	g.session.CollectedCoins = 4
	g.redrawLevelUpOptions()
	if got, want := g.session.CollectedCoins, 4; got != want {
		t.Fatalf("coins after unaffordable redraw = %d, want %d", got, want)
	}

	g.session.CollectedCoins = 5
	g.redrawLevelUpOptions()
	if got, want := g.session.CollectedCoins, 0; got != want {
		t.Fatalf("coins after affordable redraw = %d, want %d", got, want)
	}
}

func beamTestGame() *Game {
	g := New()
	g.skeleton = g.skeleton[:0]
	g.session.Progression.GainExperience(3)
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.player.Pos = Vec2{}
	g.player.Facing = 1
	g.player.Moving = false
	g.player.MoveDir = Vec2{}
	return g
}

func TestPurpleSkeletonKillGrantsOriginalExperienceReward(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{{ID: 1, Kind: SkeletonPurple, HP: 1, Reward: SkeletonPurple.ExperienceReward()}}
	g.session.Progression.GainExperience(3)

	g.destroySkeleton(0, AttackNone)

	if got, want := g.session.Progression.Level, 3; got != want {
		t.Fatalf("level = %d, want %d", got, want)
	}
	if got, want := g.session.Progression.Experience, 3; got != want {
		t.Fatalf("experience = %d, want %d", got, want)
	}
}

func TestPurpleSkeletonKillsDoNotSpawnBlackSkeletons(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{{ID: 1, Kind: SkeletonPurple, HP: 1, Reward: SkeletonPurple.ExperienceReward()}}

	g.destroySkeleton(0, AttackNone)

	if len(g.skeleton) != 0 {
		t.Fatalf("skeleton count = %d, want no milestone spawn", len(g.skeleton))
	}
}

func TestBlackSkeletonKillGrantsOriginalExperienceReward(t *testing.T) {
	g := New()
	g.skeleton = []Skeleton{{ID: 1, Kind: SkeletonBlack, HP: 1, Reward: SkeletonBlack.ExperienceReward()}}
	g.session.Progression.GainExperience(15)

	g.destroySkeleton(0, AttackNone)

	if got, want := g.session.Progression.Level, 5; got != want {
		t.Fatalf("level = %d, want %d", got, want)
	}
	if got, want := g.session.Progression.Experience, 10; got != want {
		t.Fatalf("experience = %d, want %d", got, want)
	}
}
