package game

import (
	"math"
	"math/rand"
	"testing"
)

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
