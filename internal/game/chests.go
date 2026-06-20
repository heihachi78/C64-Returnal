package game

import (
	"math"
	"slices"
)

func (g *Game) spawnChest(tier ChestTier) {
	g.chests = append(g.chests, Chest{Pos: g.randomChestPosition(), Tier: tier})
}
func (g *Game) spawnGoldChestsAroundPlayer(count int, radius float64) {
	if count <= 0 {
		return
	}
	radius = math.Max(0, radius)
	for i := 0; i < count; i++ {
		angle := math.Pi * 2 * float64(i) / float64(count)
		g.chests = append(g.chests, Chest{
			Pos: Vec2{
				X: g.player.Pos.X + math.Cos(angle)*radius,
				Y: g.player.Pos.Y + math.Sin(angle)*radius,
			},
			Tier: ChestGold,
		})
	}
}
func (g *Game) randomChestPosition() Vec2 {
	halfW := math.Max(48, float64(g.screenW)/2-g.tuning.ChestSpawnMargin)
	halfH := math.Max(48, float64(g.screenH)/2-g.tuning.ChestSpawnMargin)
	minDistSq := g.tuning.ChestPickupDistance * g.tuning.ChestPickupDistance * 4
	for range 12 {
		pos := Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
		if DistanceSq(pos, g.player.Pos) >= minDistSq {
			return pos
		}
	}
	return Vec2{X: g.player.Pos.X + halfW, Y: g.player.Pos.Y}
}
func (g *Game) checkChestPickups() {
	distSq := g.tuning.ChestPickupDistance * g.tuning.ChestPickupDistance
	for i := len(g.chests) - 1; i >= 0; i-- {
		if DistanceSq(g.chests[i].Pos, g.player.Pos) <= distSq {
			chest := g.chests[i]
			g.chests = slices.Delete(g.chests, i, i+1)
			g.applyChestReward(chest.Tier)
			return
		}
	}
}
func (g *Game) applyChestReward(tier ChestTier) {
	skills := g.session.Progression.LearnedSkills()
	items := []ChestRewardDisplayItem{}
	switch tier {
	case ChestBronze:
		skill := skills[g.rng.Intn(len(skills))]
		options := skill.UpgradeOptions()
		option := options[g.rng.Intn(len(options))]
		items = append(items, g.chestRewardItemForSkill(option, skill))
		g.applyUpgradeEffect(option)
	case ChestSilver:
		skill := skills[g.rng.Intn(len(skills))]
		items = append(items, g.chestRewardItemsForSkill(skill)...)
		g.session.Progression.UpgradeAllProperties(skill)
	case ChestGold:
		g.rng.Shuffle(len(skills), func(i, j int) { skills[i], skills[j] = skills[j], skills[i] })
		for _, skill := range skills[:min(2, len(skills))] {
			items = append(items, g.chestRewardItemsForSkill(skill)...)
			g.session.Progression.UpgradeAllProperties(skill)
		}
	}
	g.syncOrbitalOrbCount()
	if len(items) > 0 {
		g.session.ChestRewardActive = true
		g.session.ActiveChestTier = tier
		g.session.ActiveChestRewardItems = items
		g.session.ChestRewardOverlayTimer = 0
		g.suppressHeldMovementKeys(ebitenIsKeyPressed)
		g.stopPlayerAnimation()
	}
}
func (g *Game) chestRewardItemsForSkill(skill LearnedSkill) []ChestRewardDisplayItem {
	options := skill.UpgradeOptions()
	items := make([]ChestRewardDisplayItem, 0, len(options))
	for _, option := range options {
		items = append(items, g.chestRewardItemForSkill(option, skill))
	}
	return items
}
func (g *Game) chestRewardItemForSkill(option LevelUpOption, skill LearnedSkill) ChestRewardDisplayItem {
	beamKillBonus := 0
	if skill == SkillBeam {
		beamKillBonus = g.session.Progression.BeamKillUpgradeBonus()
	}
	return ChestRewardDisplayItem{
		Option: option,
		Title:  option.Title(beamKillBonus),
	}
}
