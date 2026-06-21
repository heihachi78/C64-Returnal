package game

import (
	"math"
	"testing"
)

func TestActualDPSRecordsAppliedWeaponDamage(t *testing.T) {
	g := newCombatStatsTestGame()
	g.totalTime = 10
	g.skeleton = []Skeleton{{ID: 1, HP: 3, Reward: 1}}

	g.damageSkeleton(0, 10, AttackBeam, false)

	if got, want := g.ActualDPS(), 3.0/actualDPSWindow; math.Abs(got-want) > 0.000001 {
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

func TestActualDPSUsesRollingWindow(t *testing.T) {
	g := newCombatStatsTestGame()
	g.totalTime = 0
	g.recordActualDamage(5)
	g.totalTime = 4
	g.recordActualDamage(5)

	if got, want := g.ActualDPS(), 10.0/actualDPSWindow; math.Abs(got-want) > 0.000001 {
		t.Fatalf("ActualDPS before cutoff = %v, want %v", got, want)
	}

	g.totalTime = 6
	if got, want := g.ActualDPS(), 5.0/actualDPSWindow; math.Abs(got-want) > 0.000001 {
		t.Fatalf("ActualDPS after cutoff = %v, want %v", got, want)
	}
}

func TestActualDPSMaintainsRollingWindowTotal(t *testing.T) {
	g := newCombatStatsTestGame()
	g.totalTime = 0
	g.recordActualDamage(2)
	g.totalTime = 1
	g.recordActualDamage(3)

	if got, want := g.actualDamageWindowTotal, 5; got != want {
		t.Fatalf("actual damage window total = %d, want %d", got, want)
	}

	g.totalTime = actualDPSWindow + 0.5
	g.pruneActualDamageSamples()

	if got, want := g.actualDamageWindowTotal, 3; got != want {
		t.Fatalf("actual damage window total after prune = %d, want %d", got, want)
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
