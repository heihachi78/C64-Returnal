package game

import (
	"golang.org/x/image/font"
	"math"
	"math/rand"
	"slices"
	"testing"
)

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
		t.Fatalf("player move direction = %+v, want original preserved movement direction", g.player.MoveDir)
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
		t.Fatalf("player move direction = %+v, want original preserved movement direction", g.player.MoveDir)
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
		t.Fatalf("redrawn option count = %d, want retry to rerun extra-option chance and keep 3", len(redrawn))
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

func TestLevelUpMouseHitAreasMatchOriginalLabelAndIconOnly(t *testing.T) {
	g := New()
	g.screenW = 800
	g.screenH = 600
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate}

	if got := g.levelUpOptionAt(250, 304); got != -1 {
		t.Fatalf("option at key hint center = %d, want no selection like original behavior", got)
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
		t.Fatal("redraw hit at long label tail = false, want true like original behavior")
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
