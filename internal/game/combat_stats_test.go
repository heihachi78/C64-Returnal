package game

import (
	"math"
	"testing"
)

func TestActualDPSRecordsAppliedWeaponDamage(t *testing.T) {
	g := newCombatStatsTestGame()
	g.actualDamageLevelStartTime = 4
	g.totalTime = 10
	g.skeleton = []Skeleton{{ID: 1, HP: 3, Reward: 1}}

	g.damageSkeleton(0, 10, AttackBeam, false)

	if got, want := g.ActualDPS(), 3.0/6.0; math.Abs(got-want) > 0.000001 {
		t.Fatalf("ActualDPS = %v, want %v", got, want)
	}
}

func TestActualDPSSkipsNonWeaponDamage(t *testing.T) {
	g := newCombatStatsTestGame()
	g.totalTime = 1
	g.skeleton = []Skeleton{{ID: 1, HP: 2, Reward: 1}}

	g.damageSkeleton(0, 1, AttackNone, false)

	if got := g.ActualDPS(); got != 0 {
		t.Fatalf("ActualDPS after AttackNone = %v, want 0", got)
	}
}

func TestActualDPSUsesElapsedLevelTime(t *testing.T) {
	g := newCombatStatsTestGame()
	g.actualDamageLevelStartTime = 2
	g.totalTime = 4
	g.recordActualDamage(5)

	if got, want := g.ActualDPS(), 5.0/2.0; math.Abs(got-want) > 0.000001 {
		t.Fatalf("ActualDPS before later level time = %v, want %v", got, want)
	}

	g.totalTime = 12
	g.recordActualDamage(5)
	if got, want := g.ActualDPS(), 10.0/10.0; math.Abs(got-want) > 0.000001 {
		t.Fatalf("ActualDPS after more elapsed level time = %v, want %v", got, want)
	}
}

func TestActualDPSExcludesPausedLevelTime(t *testing.T) {
	g := newCombatStatsTestGame()
	g.actualDamageLevelStartTime = 2
	g.totalTime = 9
	g.recordActualDamage(12)
	g.pauseActualDamageLevelStats(3)

	if got, want := g.ActualDPS(), 12.0/4.0; math.Abs(got-want) > 0.000001 {
		t.Fatalf("ActualDPS with paused level time = %v, want %v", got, want)
	}
}

func TestActualDPSUsesOneSecondMinimumElapsedTime(t *testing.T) {
	g := newCombatStatsTestGame()
	g.actualDamageLevelStartTime = 10
	g.totalTime = 10 + 1.0/float64(TargetTPS)
	g.recordActualDamage(1)

	if got, want := g.ActualDPS(), 1.0; math.Abs(got-want) > 0.000001 {
		t.Fatalf("ActualDPS with first-tick damage = %v, want %v", got, want)
	}
}

func TestActualDPSMaintainsLevelTotalUntilReset(t *testing.T) {
	g := newCombatStatsTestGame()
	g.totalTime = 0
	g.recordActualDamage(2)
	g.totalTime = 1
	g.recordActualDamage(3)

	if got, want := g.actualDamageLevelTotal, 5; got != want {
		t.Fatalf("actual damage level total = %d, want %d", got, want)
	}

	g.totalTime = 10
	g.resetActualDamageLevelStats()

	if got, want := g.actualDamageLevelTotal, 0; got != want {
		t.Fatalf("actual damage level total after reset = %d, want %d", got, want)
	}
	if got, want := g.actualDamageLevelStartTime, 10.0; got != want {
		t.Fatalf("actual damage level start after reset = %v, want %v", got, want)
	}
	if got, want := g.actualDamageLevelPausedTime, 0.0; got != want {
		t.Fatalf("actual damage paused time after reset = %v, want %v", got, want)
	}
}

func newCombatStatsTestGame() *Game {
	tuning := DefaultTuning()
	return &Game{
		tuning:  tuning,
		session: NewSession(tuning),
		spatial: NewSpatialIndex(tuning.SpatialIndexCellSize),
	}
}
