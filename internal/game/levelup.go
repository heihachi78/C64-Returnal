package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"slices"
)

func (g *Game) queueLevelUpChoices(count int) {
	if count <= 0 || g.session.GameOver {
		return
	}
	first := g.session.Progression.Level - count + 1
	for level := first; level <= g.session.Progression.Level; level++ {
		g.session.PendingLevelUpLevels = append(g.session.PendingLevelUpLevels, level)
		g.spawnCoinForLevel(level)
	}
	g.presentNextLevelUpChoiceIfNeeded()
}
func (g *Game) presentNextLevelUpChoiceIfNeeded() {
	if g.session.GameOver || g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) == 0 {
		return
	}
	g.session.LevelUpChoiceActive = true
	g.session.CurrentLevelUpPresentation = g.session.PendingLevelUpLevels[0]
	g.session.ActiveLevelUpOptions = visibleLevelUpOptions(g.randomLevelUpOptions(nil))
	g.session.LevelUpOverlayTimer = 0
	g.session.LevelUpTitleScaleTimer = 0
	g.session.LevelUpOptionFadeTimer = 0
	g.showLevelUpRedrawPresentation(true)
	g.suppressHeldMovementKeys(ebiten.IsKeyPressed)
	g.stopPlayerAnimation()
}
func (g *Game) randomLevelUpOptions(excluding []LevelUpOption) []LevelUpOption {
	selected := g.randomLevelUpOptionsCandidate()
	if len(excluding) > 0 {
		for tries := 0; tries < 8 && sameOptionSet(selected, excluding); tries++ {
			selected = g.randomLevelUpOptionsCandidate()
		}
	}
	return selected
}
func (g *Game) randomLevelUpOptionsCandidate() []LevelUpOption {
	hasSkeletons := len(g.skeleton) > 0
	count := 2
	if chance(g.rng, g.tuning.ExtraOptionChanceNumerator, g.tuning.ExtraOptionChanceDenominator) {
		count = 3
	}
	available := slices.Clone(g.session.Progression.AvailableLevelUpOptions())
	available = slices.DeleteFunc(available, func(o LevelUpOption) bool { return o == HalveSkeletons })
	g.rng.Shuffle(len(available), func(i, j int) { available[i], available[j] = available[j], available[i] })
	selected := slices.Clone(available[:min(count, len(available))])
	if hasSkeletons && len(selected) > 0 && chance(g.rng, g.tuning.HalveHordeChanceNumerator, g.tuning.HalveHordeChanceDenominator) {
		selected[g.rng.Intn(len(selected))] = HalveSkeletons
	}
	return selected
}
func (g *Game) applyLevelUpOption(option LevelUpOption) {
	g.applyUpgradeEffect(option)
	g.syncOrbitalOrbCount()
	if len(g.session.PendingLevelUpLevels) > 0 {
		g.session.PendingLevelUpLevels = g.session.PendingLevelUpLevels[1:]
	}
	g.session.LevelUpChoiceActive = false
	g.session.ActiveLevelUpOptions = nil
	g.hideLevelUpPresentation()
	g.presentNextLevelUpChoiceIfNeeded()
}
func (g *Game) applyUpgradeEffect(option LevelUpOption) {
	switch option {
	case ExtraLife:
		g.session.PlayerLives++
	case HalveSkeletons:
		g.halveSkeletons()
	default:
		g.session.Progression.ApplyLevelUpOption(option)
	}
}
func (g *Game) halveSkeletons() {
	killCount := len(g.skeleton) / 2
	if killCount <= 0 {
		return
	}
	targetIDs := make([]int, len(g.skeleton))
	for i, skeleton := range g.skeleton {
		targetIDs[i] = skeleton.ID
	}
	g.rng.Shuffle(len(targetIDs), func(i, j int) { targetIDs[i], targetIDs[j] = targetIDs[j], targetIDs[i] })
	levelUps := 0
	for _, id := range targetIDs[:killCount] {
		if idx := g.skeletonIndexByID(id); idx >= 0 {
			levelUps += g.destroySkeleton(idx, AttackNone)
		}
	}
	g.queueLevelUpChoices(levelUps)
}
func (g *Game) redrawLevelUpOptions() {
	if !g.session.LevelUpChoiceActive || len(g.session.PendingLevelUpLevels) == 0 {
		return
	}
	cost := g.levelUpRedrawCost()
	if g.session.CollectedCoins < cost {
		g.session.LevelUpRedrawStatusTimer = redrawFailurePulseDuration
		g.session.LevelUpRedrawCoinFadeTimer = 0
		return
	}
	g.session.CollectedCoins -= cost
	previous := slices.Clone(g.session.ActiveLevelUpOptions)
	g.session.ActiveLevelUpOptions = visibleLevelUpOptions(g.randomLevelUpOptions(previous))
	g.session.LevelUpTitleScaleTimer = 0
	g.session.LevelUpOptionFadeTimer = 0
	g.showLevelUpRedrawPresentation(false)
}
func visibleLevelUpOptions(options []LevelUpOption) []LevelUpOption {
	return slices.Clone(options[:min(4, len(options))])
}
func (g *Game) showLevelUpRedrawPresentation(resetTextFade bool) {
	g.session.LevelUpRedrawStatusTimer = 0
	if g.session.CollectedCoins < g.levelUpRedrawCost() {
		g.session.LevelUpRedrawStatusTimer = redrawFailurePulseDuration
		g.session.LevelUpRedrawCoinFadeTimer = 0
	} else {
		g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration
	}
	if resetTextFade {
		g.session.LevelUpRedrawFadeTimer = 0
	}
}
func (g *Game) levelUpRedrawCost() int {
	return max(1, g.session.Progression.Level)
}
func (g *Game) hideLevelUpPresentation() {
	g.session.LevelUpRedrawStatusTimer = 0
	g.session.LevelUpRedrawFadeTimer = 0
	g.session.LevelUpRedrawCoinFadeTimer = 0
	g.session.LevelUpOverlayTimer = 0
	g.session.LevelUpTitleScaleTimer = 0
	g.session.LevelUpOptionFadeTimer = 0
}
