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
	halfW := math.Max(1, float64(g.screenW)/2)
	halfH := math.Max(1, float64(g.screenH)/2)
	margin := math.Max(1, g.tuning.ChestSpawnMargin)
	switch g.rng.Intn(4) {
	case 0:
		return Vec2{X: g.player.Pos.X - halfW - margin - g.randRange(0, halfW), Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
	case 1:
		return Vec2{X: g.player.Pos.X + halfW + margin + g.randRange(0, halfW), Y: g.player.Pos.Y + g.randRange(-halfH, halfH)}
	case 2:
		return Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y - halfH - margin - g.randRange(0, halfH)}
	default:
		return Vec2{X: g.player.Pos.X + g.randRange(-halfW, halfW), Y: g.player.Pos.Y + halfH + margin + g.randRange(0, halfH)}
	}
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
	skills := g.availableChestRewardSkills()
	items := []ChestRewardDisplayItem{}
	if len(skills) == 0 {
		return
	}
	switch tier {
	case ChestBronze:
		skill := skills[g.rng.Intn(len(skills))]
		options := g.session.Progression.AvailableUpgradeOptionsForSkill(skill)
		option := options[g.rng.Intn(len(options))]
		items = append(items, g.chestRewardItemForSkill(option, skill))
		g.applyUpgradeEffect(option)
	case ChestSilver:
		skill := skills[g.rng.Intn(len(skills))]
		for _, option := range g.session.Progression.UpgradeAllProperties(skill) {
			items = append(items, g.chestRewardItemForSkill(option, skill))
		}
	case ChestGold:
		g.rng.Shuffle(len(skills), func(i, j int) { skills[i], skills[j] = skills[j], skills[i] })
		for _, skill := range skills[:min(2, len(skills))] {
			for _, option := range g.session.Progression.UpgradeAllProperties(skill) {
				items = append(items, g.chestRewardItemForSkill(option, skill))
			}
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
func (g *Game) availableChestRewardSkills() []LearnedSkill {
	learned := g.session.Progression.LearnedSkills()
	skills := make([]LearnedSkill, 0, len(learned))
	for _, skill := range learned {
		if len(g.session.Progression.AvailableUpgradeOptionsForSkill(skill)) > 0 {
			skills = append(skills, skill)
		}
	}
	return skills
}
func (g *Game) chestRewardItemsForSkill(skill LearnedSkill) []ChestRewardDisplayItem {
	options := g.session.Progression.AvailableUpgradeOptionsForSkill(skill)
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
