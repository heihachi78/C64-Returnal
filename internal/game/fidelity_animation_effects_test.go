package game

import (
	"math"
	"testing"
)

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
