package game

import (
	"math"
	"math/rand"
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

func TestTimedSkeletonSpawnKindTurnsRedThenPurpleThenBlackAtThresholds(t *testing.T) {
	g := New()
	g.session.Progression.Level = g.tuning.RedOnlyLevel
	if got := g.timedSkeletonSpawnKind(); got != SkeletonRed {
		t.Fatalf("spawn kind at red threshold = %v, want red", got)
	}

	g.session.Progression.Level = g.tuning.PurpleOnlyLevel
	if got := g.timedSkeletonSpawnKind(); got != SkeletonPurple {
		t.Fatalf("spawn kind at purple threshold = %v, want purple", got)
	}

	g.session.Progression.Level = g.tuning.BlackOnlyLevel - 1
	if got := g.timedSkeletonSpawnKind(); got != SkeletonPurple {
		t.Fatalf("spawn kind before black threshold = %v, want purple", got)
	}

	g.session.Progression.Level = g.tuning.BlackOnlyLevel
	if got := g.timedSkeletonSpawnKind(); got != SkeletonBlack {
		t.Fatalf("spawn kind at black threshold = %v, want black", got)
	}
}

func TestBlackOnlyLevelIncreasesTimedSkeletonSpawnInterval(t *testing.T) {
	p := NewProgression(DefaultTuning())

	p.Level = p.tuning.BlackOnlyLevel - 1
	beforeBlack := p.SkeletonSpawnInterval()
	p.Level = p.tuning.BlackOnlyLevel
	black := p.SkeletonSpawnInterval()

	if black <= beforeBlack {
		t.Fatalf("black spawn interval = %v, want greater than level-before-black interval %v", black, beforeBlack)
	}
}

func TestBlueMonsterOnlySpawnsAfterLevel100WithLargeHorde(t *testing.T) {
	g := New()
	g.session.Progression.Level = g.tuning.BlueMonsterMinimumLevel
	g.skeleton = makeSkeletonHorde(g.tuning.BlueMonsterMinimumEnemies)
	g.nextID = 10_000

	g.spawnTimedSkeleton()

	if got, want := len(g.skeleton), g.tuning.BlueMonsterMinimumEnemies+1; got != want {
		t.Fatalf("skeleton count below blue threshold = %d, want %d", got, want)
	}
	if got := countSkeletonKind(g.skeleton, SkeletonBlue); got != 0 {
		t.Fatalf("blue monsters below threshold = %d, want 0", got)
	}
	if got := countSkeletonKind(g.skeleton, SkeletonBlack); got != 1 {
		t.Fatalf("black monsters below threshold = %d, want 1", got)
	}

	g = New()
	g.session.Progression.Level = g.tuning.BlueMonsterMinimumLevel - 1
	g.skeleton = makeSkeletonHorde(g.tuning.BlueMonsterMinimumEnemies + 1)
	g.nextID = 10_000

	g.spawnTimedSkeleton()

	if got := countSkeletonKind(g.skeleton, SkeletonBlue); got != 0 {
		t.Fatalf("blue monsters at level 100 = %d, want 0", got)
	}
}

func TestBlueMonsterSpawnSlowsSpawnRateAndCullsConfiguredShareOfHorde(t *testing.T) {
	g := New()
	g.session.Progression.Level = g.tuning.BlueMonsterMinimumLevel
	g.skeleton = makeSkeletonHorde(g.tuning.BlueMonsterMinimumEnemies + 1)
	g.nextID = 10_000
	beforeInterval := g.session.Progression.SkeletonSpawnInterval()

	g.spawnTimedSkeleton()

	if got, want := countSkeletonKind(g.skeleton, SkeletonBlue), 1; got != want {
		t.Fatalf("blue monsters = %d, want %d", got, want)
	}
	wantCount := g.tuning.BlueMonsterMinimumEnemies + 2 - (g.tuning.BlueMonsterMinimumEnemies+1)/g.tuning.BlueMonsterCullDivisor
	if got, want := len(g.skeleton), wantCount; got != want {
		t.Fatalf("skeleton count after blue cull = %d, want %d", got, want)
	}
	afterInterval := g.session.Progression.SkeletonSpawnInterval()
	wantInterval := beforeInterval / g.tuning.BlueMonsterSpawnRateFactor
	if math.Abs(afterInterval-wantInterval) > 0.0001 {
		t.Fatalf("spawn interval after blue = %v, want %v", afterInterval, wantInterval)
	}
	for _, skeleton := range g.skeleton {
		if skeleton.Kind == SkeletonBlue {
			if got, want := skeleton.HP, g.tuning.BlueMonsterHitPoints; got != want {
				t.Fatalf("blue monster HP = %d, want %d", got, want)
			}
			if got, want := skeleton.Reward, 75; got != want {
				t.Fatalf("blue monster reward = %d, want %d", got, want)
			}
			return
		}
	}
	t.Fatal("blue monster was not preserved after cull")
}

func TestBlueMonsterCullDivisorIsConfigurable(t *testing.T) {
	g := New()
	g.tuning.BlueMonsterCullDivisor = 4
	g.skeleton = makeSkeletonHorde(12)
	g.nextID = 100

	g.spawnBlueMonster()

	if got, want := len(g.skeleton), 10; got != want {
		t.Fatalf("skeleton count after blue cull divisor 4 = %d, want %d", got, want)
	}
	if got, want := countSkeletonKind(g.skeleton, SkeletonBlue), 1; got != want {
		t.Fatalf("blue monsters after configured cull = %d, want %d", got, want)
	}
}

func TestBlueMonsterCullDivisorCanDisableCull(t *testing.T) {
	g := New()
	g.tuning.BlueMonsterCullDivisor = 0
	g.skeleton = makeSkeletonHorde(12)

	g.spawnBlueMonster()

	if got, want := len(g.skeleton), 13; got != want {
		t.Fatalf("skeleton count after disabled blue cull = %d, want %d", got, want)
	}
}

func TestBlueMonsterHitPointsGrowByTwentyFivePercentPerSpawn(t *testing.T) {
	tests := []struct {
		spawns int
		want   int
	}{
		{spawns: 1, want: 1000},
		{spawns: 2, want: 1250},
		{spawns: 3, want: 1563},
		{spawns: 4, want: 1954},
	}
	for _, tt := range tests {
		if got := blueMonsterHitPoints(1000, tt.spawns); got != tt.want {
			t.Fatalf("blue monster HP after %d spawns = %d, want %d", tt.spawns, got, tt.want)
		}
	}
}

func TestRepeatedBlueMonsterSpawnsUseIncreasingHitPoints(t *testing.T) {
	g := New()
	g.skeleton = nil
	g.nextID = 1
	wantHP := []int{1000, 1250, 1563}

	for _, want := range wantHP {
		g.spawnBlueMonster()
		newID := g.nextID - 1
		idx := findSkeletonByID(g.skeleton, newID)
		if idx < 0 {
			t.Fatalf("new blue monster %d was not preserved", newID)
		}
		if got := g.skeleton[idx].HP; got != want {
			t.Fatalf("blue monster %d HP = %d, want %d", newID, got, want)
		}
	}
}

func TestBlueMonsterSpawnRateFactorIsConfigurable(t *testing.T) {
	g := New()
	g.tuning.BlueMonsterSpawnRateFactor = 0.25
	g.session.Progression = NewProgression(g.tuning)
	g.session.Progression.Level = g.tuning.BlueMonsterMinimumLevel
	g.skeleton = makeSkeletonHorde(g.tuning.BlueMonsterMinimumEnemies + 1)
	beforeInterval := g.session.Progression.SkeletonSpawnInterval()

	g.spawnTimedSkeleton()

	afterInterval := g.session.Progression.SkeletonSpawnInterval()
	if math.Abs(afterInterval-beforeInterval*4) > 0.0001 {
		t.Fatalf("spawn interval after blue = %v, want quadruple %v", afterInterval, beforeInterval)
	}
}

func makeSkeletonHorde(count int) []Skeleton {
	skeletons := make([]Skeleton, count)
	for i := range skeletons {
		skeletons[i] = Skeleton{ID: i + 1, Kind: SkeletonRegular, HP: 1, Reward: 1}
	}
	return skeletons
}

func findSkeletonByID(skeletons []Skeleton, id int) int {
	for i, skeleton := range skeletons {
		if skeleton.ID == id {
			return i
		}
	}
	return -1
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

func TestBlackSkeletonSpawnsAfterOriginalPurpleKillMilestone(t *testing.T) {
	g := New()
	g.session.Kills.PurpleSkeletons = g.tuning.BlackPurpleKillInterval - 1
	g.skeleton = []Skeleton{{ID: 1, Kind: SkeletonPurple, HP: 1, Reward: SkeletonPurple.ExperienceReward()}}

	g.destroySkeleton(0, AttackNone)

	if len(g.skeleton) != 1 {
		t.Fatalf("skeleton count = %d, want 1 spawned black skeleton", len(g.skeleton))
	}
	if got, want := g.skeleton[0].Kind, SkeletonBlack; got != want {
		t.Fatalf("spawned skeleton kind = %v, want %v", got, want)
	}
	if got, want := g.skeleton[0].HP, g.tuning.BlackHitPoints; got != want {
		t.Fatalf("black skeleton HP = %d, want %d", got, want)
	}
	if got, want := g.session.Kills.PurpleSkeletons, g.tuning.BlackPurpleKillInterval; got != want {
		t.Fatalf("purple kills = %d, want %d", got, want)
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
