package game

import (
	"errors"
	"image"
	"image/color"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func TestCoverageRenderPassDrawsWorldHUDAndAllOverlays(t *testing.T) {
	g := New()
	g.screenW = ScreenWidth
	g.screenH = ScreenHeight
	g.player = Player{
		Pos:           Vec2{X: 12, Y: -8},
		Facing:        -1,
		Moving:        true,
		MoveDir:       Vec2{X: -1, Y: 1},
		AnimFrame:     1,
		HitFlash:      playerHitFlashDuration / 2,
		DeathRotation: -math.Pi / 4,
	}
	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.session.Progression.ApplyLevelUpOption(LightningBounce)
	g.session.Progression.ApplyLevelUpOption(LearnOrb)
	g.session.Progression.ApplyLevelUpOption(ExtraOrb)
	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.session.Progression.ApplyLevelUpOption(BeamKillCount)
	g.session.Progression.ApplyLevelUpOption(LearnMeteor)
	g.session.Progression.ApplyLevelUpOption(ExtraMeteor)
	g.session.PlayerLives = 13
	g.session.CollectedCoins = 7
	g.session.Kills = KillCounts{Fireball: 1, Lightning: 2, OrbitalOrb: 3, Beam: 4, Meteor: 5}
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraLife, LearnLightning, LearnBeam}
	g.session.CurrentLevelUpPresentation = 5
	g.session.LevelUpOverlayTimer = modalFadeDuration
	g.session.LevelUpTitleScaleTimer = modalFadeDuration / 2
	g.session.LevelUpOptionFadeTimer = modalFadeDuration
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.LevelUpRedrawCoinFadeTimer = redrawStatusFadeDuration / 2
	g.session.LevelUpRedrawStatusTimer = redrawFailurePulseDuration / 2
	g.session.ChestRewardActive = true
	g.session.ActiveChestTier = ChestGold
	g.session.ActiveChestRewardItems = []ChestRewardDisplayItem{
		{Option: ExtraMeteor, Title: ExtraMeteor.Title(0)},
		{Option: BeamKillCount, Title: BeamKillCount.Title(g.session.Progression.BeamKillUpgradeBonus())},
	}
	g.session.ChestRewardOverlayTimer = modalFadeDuration
	g.session.GameOver = true
	g.session.GameOverOverlayTimer = modalFadeDuration
	g.skeleton = []Skeleton{
		{ID: 101, Pos: Vec2{X: -40, Y: 10}, Kind: SkeletonRegular, HP: 1, Reward: 1, Facing: 1, AnimFrame: 0},
		{ID: 102, Pos: Vec2{X: 30, Y: 30}, Kind: SkeletonRed, HP: 2, Reward: 1, Facing: -1, HitFlash: skeletonDamageFlashDuration / 2, AnimFrame: 1},
		{ID: 103, Pos: Vec2{X: 70, Y: -20}, Kind: SkeletonPurple, HP: 3, Reward: 3, Facing: 1, AnimFrame: 0},
		{ID: 104, Pos: Vec2{X: -80, Y: -60}, Kind: SkeletonBlack, HP: 5, Reward: 10, Facing: -1, AnimFrame: 1},
	}
	g.spatial.Rebuild(g.skeleton)
	g.fireball = []Fireball{{Pos: Vec2{X: 10, Y: 10}, Velocity: Vec2{X: 1, Y: -1}, AnimFrame: 1}}
	g.orbs = []OrbitalOrb{{Pos: Vec2{X: 25, Y: 0}, Active: true, AnimFrame: 1}, {Pos: Vec2{X: -25, Y: 0}, Active: false}}
	g.meteors = []MeteorProjectile{{Pos: Vec2{X: 10, Y: 120}, Start: Vec2{X: 0, Y: 140}, Impact: Vec2{X: 40, Y: 30}, AnimFrame: 1}}
	g.chests = []Chest{{Pos: Vec2{X: -30, Y: -30}, Tier: ChestBronze}, {Pos: Vec2{X: 20, Y: -50}, Tier: ChestSilver}, {Pos: Vec2{X: 60, Y: 20}, Tier: ChestGold}}
	g.coins = []Coin{{Pos: Vec2{X: 50, Y: 0}, Amount: 2, Phase: 0.21}}
	g.effects = []Effect{
		{Kind: EffectLightningHit, Pos: Vec2{X: -40, Y: 10}, Frame: 1, Facing: -1, TTL: g.tuning.LightningEffectDuration / 2, MaxTTL: g.tuning.LightningEffectDuration},
		{Kind: EffectLightning, Start: Vec2{}, End: Vec2{X: 60, Y: 40}, Points: []Vec2{{}, {X: 30, Y: 20}, {X: 60, Y: 40}}, TTL: g.tuning.LightningEffectDuration, MaxTTL: g.tuning.LightningEffectDuration},
		{Kind: EffectBeam, Start: Vec2{X: -20}, End: Vec2{X: 120}, TTL: g.tuning.BeamEffectDuration, MaxTTL: g.tuning.BeamEffectDuration},
		{Kind: EffectMeteorImpact, Pos: Vec2{X: 15, Y: 15}, Radius: g.tuning.MeteorImpactRadius, TTL: meteorImpactEffectDuration / 2, MaxTTL: meteorImpactEffectDuration},
	}

	screen := ebiten.NewImage(ScreenWidth, ScreenHeight)
	g.Draw(screen)
}

func TestCoverageDrawHelpersAndOptionIcons(t *testing.T) {
	g := New()
	screen := ebiten.NewImage(96, 96)
	img := ebiten.NewImage(8, 8)
	img.Fill(color.White)
	white := color.RGBA{255, 255, 255, 255}

	g.drawSprite(screen, img, -1000, -1000, 10, 10, false, white)
	g.drawSpriteScreen(screen, img, 12, 12, 10, 10, true, color.RGBA{255, 0, 0, 128})
	g.drawSpriteRotated(screen, img, 24, 24, 10, 10, math.Pi/4, false, white)
	g.drawSpriteRotatedBlend(screen, img, 36, 36, 10, 10, 0, false, color.RGBA{0, 255, 0, 128}, 0.5)
	g.panel(screen, 2, 2, 40, 20)
	g.panelWithAlpha(screen, 3, 30, 40, 20, 90)
	drawFilledRoundedRect(screen, 0, 0, 0, 10, 3, color.White)
	drawFilledRoundedRect(screen, 0, 0, 10, 10, 0, color.White)

	for _, option := range []LevelUpOption{
		FireRate, ExtraFireball, ExtraLife, HalveSkeletons,
		LearnLightning, LightningBounce, LightningRate,
		LearnOrb, ExtraOrb, OrbitalSpeed,
		LearnBeam, BeamRate, BeamKillCount,
		LearnMeteor, ExtraMeteor, MeteorRate,
		LevelUpOption(999),
	} {
		g.drawOptionIcon(screen, option, 48, 48)
		g.drawOptionIconTinted(screen, option, 64, 64, color.RGBA{255, 255, 255, 128})
	}
	g.drawBolt(screen, []Vec2{{}}, 1, white)
	g.drawBolt(screen, []Vec2{{}, {X: 10}}, 1, white)
	g.drawScaledTextImage(screen, "skip", 10, 10, 12, 0, false, color.White)
	g.drawTextSize(screen, "text", 4, 8, 12, color.White)
	g.drawTextSizeScaled(screen, "scaled", 4, 20, 12, 1, color.White)
	g.drawCenteredTextSize(screen, "center", 48, 32, 12, color.White)
	g.drawCenteredTextSizeScaled(screen, "center scaled", 48, 48, 12, 1, color.White)
}

func TestCoveragePureHelpersCoverDefensiveBranches(t *testing.T) {
	if got := (Vec2{}).Normalized(); got != (Vec2{}) {
		t.Fatalf("zero normalized = %+v", got)
	}
	if got := Clamp(-1, 0, 2); got != 0 {
		t.Fatalf("Clamp below = %v, want 0", got)
	}
	if got := Clamp(3, 0, 2); got != 2 {
		t.Fatalf("Clamp above = %v, want 2", got)
	}
	if got := Clamp(1, 0, 2); got != 1 {
		t.Fatalf("Clamp middle = %v, want 1", got)
	}
	if clamp01(-0.5) != 0 || clamp01(1.5) != 1 || clamp01(0.25) != 0.25 {
		t.Fatal("clamp01 did not clamp all branches")
	}
	if rgb(-1, 0.5, 2) != (color.RGBA{0, 128, 255, 255}) {
		t.Fatal("rgb did not clamp and round channels")
	}

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	px(img, 1, 1, 0, 2, color.RGBA{255, 255, 255, 255})
	px(img, -1, -1, 3, 3, color.RGBA{1, 2, 3, 4})
	fill(img, color.RGBA{5, 6, 7, 8})

	g := New()
	g.skeleton = []Skeleton{{ID: 42}, {ID: 99}}
	if g.skeletonIndexByID(0) != -1 || g.skeletonIndexByID(42) != 0 || g.skeletonIndexByID(7) != -1 {
		t.Fatal("skeletonIndexByID missed a branch")
	}
	if sameOptionSet([]LevelUpOption{FireRate}, []LevelUpOption{FireRate, ExtraLife}) {
		t.Fatal("sameOptionSet accepted mismatched lengths")
	}
	if sameOptionSet([]LevelUpOption{FireRate, FireRate}, []LevelUpOption{FireRate, ExtraLife}) {
		t.Fatal("sameOptionSet accepted mismatched multiset")
	}
	if !sameOptionSet([]LevelUpOption{FireRate, ExtraLife}, []LevelUpOption{ExtraLife, FireRate}) {
		t.Fatal("sameOptionSet rejected equal multiset")
	}
}

func TestCoverageProgressionBranches(t *testing.T) {
	p := NewProgression(DefaultTuning())
	if got := p.LightningStrikeCount(); got != 0 {
		t.Fatalf("locked LightningStrikeCount = %d, want 0", got)
	}
	if got := p.OrbitalOrbCount(); got != 0 {
		t.Fatalf("locked OrbitalOrbCount = %d, want 0", got)
	}
	if got := p.BeamKillCount(); got != 0 {
		t.Fatalf("locked BeamKillCount = %d, want 0", got)
	}
	if got := p.MeteorCount(); got != 0 {
		t.Fatalf("locked MeteorCount = %d, want 0", got)
	}
	if got := p.LearnedSkills(); !reflect.DeepEqual(got, []LearnedSkill{SkillFireball}) {
		t.Fatalf("locked learned skills = %v", got)
	}
	lockedOptions := p.AvailableLevelUpOptions()
	for _, option := range []LevelUpOption{LearnLightning, LearnOrb, LearnBeam, LearnMeteor} {
		if !slices.Contains(lockedOptions, option) {
			t.Fatalf("locked options missing %v: %v", option, lockedOptions)
		}
	}

	for _, option := range []LevelUpOption{
		FireRate, ExtraFireball, LearnLightning, LightningBounce, LightningRate,
		LearnOrb, ExtraOrb, OrbitalSpeed, LearnBeam, BeamRate, BeamKillCount,
		LearnMeteor, ExtraMeteor, MeteorRate,
	} {
		p.ApplyLevelUpOption(option)
	}
	if got := p.LearnedSkills(); !reflect.DeepEqual(got, []LearnedSkill{SkillFireball, SkillLightning, SkillOrbitalOrb, SkillBeam, SkillMeteor}) {
		t.Fatalf("unlocked learned skills = %v", got)
	}
	unlockedOptions := p.AvailableLevelUpOptions()
	for _, option := range []LevelUpOption{LightningBounce, LightningRate, ExtraOrb, OrbitalSpeed, BeamRate, BeamKillCount, ExtraMeteor, MeteorRate} {
		if !slices.Contains(unlockedOptions, option) {
			t.Fatalf("unlocked options missing %v: %v", option, unlockedOptions)
		}
	}

	before := p
	p.UpgradeAllProperties(SkillFireball)
	p.UpgradeAllProperties(SkillLightning)
	p.UpgradeAllProperties(SkillOrbitalOrb)
	p.UpgradeAllProperties(SkillBeam)
	p.UpgradeAllProperties(SkillMeteor)
	if p.SimultaneousFireball != before.SimultaneousFireball+1 ||
		p.LightningBounceCount != before.LightningBounceCount+1 ||
		p.upgradedOrbitalOrbCount != before.upgradedOrbitalOrbCount+1 ||
		p.upgradedBeamKillCount <= before.upgradedBeamKillCount ||
		p.upgradedMeteorCount != before.upgradedMeteorCount+1 {
		t.Fatalf("UpgradeAllProperties did not upgrade unlocked skills: before %+v after %+v", before, p)
	}

	locked := NewProgression(DefaultTuning())
	locked.UpgradeAllProperties(SkillLightning)
	locked.UpgradeAllProperties(SkillOrbitalOrb)
	locked.UpgradeAllProperties(SkillBeam)
	locked.UpgradeAllProperties(SkillMeteor)
	if locked.LightningBounceCount != 0 || locked.upgradedOrbitalOrbCount != 1 || locked.upgradedBeamKillCount != 1 || locked.upgradedMeteorCount != 1 {
		t.Fatalf("locked skill upgrades changed progression: %+v", locked)
	}
}

func TestCoverageChestBranches(t *testing.T) {
	g := New()
	g.rng = rand.New(rand.NewSource(1))
	g.spawnChest(ChestBronze)
	if len(g.chests) == 0 {
		t.Fatal("spawnChest did not append a chest")
	}
	g.tuning.ChestPickupDistance = 1000
	g.chests = []Chest{{Pos: g.player.Pos, Tier: ChestBronze}, {Pos: Vec2{X: 10000}, Tier: ChestSilver}}
	g.checkChestPickups()
	if len(g.chests) != 1 || !g.session.ChestRewardActive {
		t.Fatalf("checkChestPickups did not collect one chest: chests=%v active=%v", g.chests, g.session.ChestRewardActive)
	}

	gold := New()
	gold.rng = rand.New(rand.NewSource(2))
	for _, option := range []LevelUpOption{LearnLightning, LearnOrb, LearnBeam, LearnMeteor} {
		gold.session.Progression.ApplyLevelUpOption(option)
	}
	gold.applyChestReward(ChestGold)
	if !gold.session.ChestRewardActive || len(gold.session.ActiveChestRewardItems) != 4 {
		t.Fatalf("gold chest items = %v active=%v, want two skills worth of rewards", gold.session.ActiveChestRewardItems, gold.session.ChestRewardActive)
	}

	silver := New()
	silver.rng = rand.New(rand.NewSource(3))
	silver.session.Progression.ApplyLevelUpOption(LearnBeam)
	silver.applyChestReward(ChestSilver)
	if !silver.session.ChestRewardActive || len(silver.session.ActiveChestRewardItems) != 2 {
		t.Fatalf("silver chest items = %v active=%v, want one full skill reward", silver.session.ActiveChestRewardItems, silver.session.ChestRewardActive)
	}
}

func TestCoveragePlayerAndSessionBranches(t *testing.T) {
	g := New()
	pressed := map[ebiten.Key]bool{
		ebiten.KeyArrowLeft: true,
		ebiten.KeyArrowUp:   true,
	}
	g.suppressedMovement[ebiten.KeyArrowLeft] = true
	move := g.playerMovementVector(func(key ebiten.Key) bool { return pressed[key] })
	if move != (Vec2{Y: 1}) {
		t.Fatalf("suppressed movement = %+v, want up only", move)
	}
	pressed[ebiten.KeyArrowLeft] = false
	_ = g.playerMovementVector(func(key ebiten.Key) bool { return pressed[key] })
	if g.suppressedMovement[ebiten.KeyArrowLeft] {
		t.Fatal("released suppressed key was not cleared")
	}

	g.updateInvulnerability(0)
	g.session.PlayerHitInvulnerability = 0.1
	g.updateInvulnerability(1)
	if g.session.PlayerHitInvulnerability != 0 {
		t.Fatalf("invulnerability = %v, want 0", g.session.PlayerHitInvulnerability)
	}

	g.session.PlayerLives = 2
	g.damagePlayer()
	if g.session.PlayerLives != 1 || g.session.GameOver {
		t.Fatalf("nonfatal damage state = lives %d gameOver %v", g.session.PlayerLives, g.session.GameOver)
	}
	g.damagePlayer()
	if !g.session.GameOver {
		t.Fatal("fatal damage did not trigger game over")
	}
	g.triggerGameOver()

	var s Session
	s.Reset(DefaultTuning())
	for _, attack := range []AttackKind{AttackFireball, AttackLightning, AttackOrbitalOrb, AttackBeam, AttackMeteor, AttackNone} {
		s.RegisterAttackKill(attack)
	}
	if s.Kills.Fireball != 1 || s.Kills.Lightning != 1 || s.Kills.OrbitalOrb != 1 || s.Kills.Beam != 1 || s.Kills.Meteor != 1 {
		t.Fatalf("attack kill counts = %+v", s.Kills)
	}
}

func TestCoverageAnimationAndSpawnBranches(t *testing.T) {
	g := New()
	g.effects = []Effect{{TTL: 0.1}, {TTL: 0.3}}
	g.updateEffects(0.2)
	if len(g.effects) != 1 {
		t.Fatalf("effects after expiry = %d, want 1", len(g.effects))
	}
	g.updatePausedAnimations(0)
	g.session.LevelUpChoiceActive = true
	g.session.ChestRewardActive = true
	g.session.GameOver = true
	g.session.LevelUpRedrawStatusTimer = 0.1
	g.updatePausedAnimations(0.2)
	if g.session.LevelUpRedrawStatusTimer != 0 {
		t.Fatalf("redraw status timer = %v, want 0", g.session.LevelUpRedrawStatusTimer)
	}
	g.session.LevelUpRedrawStatusTimer = 0.1
	g.updateNewlyPresentedOverlayActions(0.2)
	if g.session.LevelUpRedrawStatusTimer != 0 {
		t.Fatalf("new overlay redraw status timer = %v, want 0", g.session.LevelUpRedrawStatusTimer)
	}

	g.skeleton = nil
	g.skeletonAnimTimer = 1
	g.updateSkeletonAnimation(0)
	if g.skeletonAnimTimer != 0 {
		t.Fatalf("empty skeleton animation timer = %v, want 0", g.skeletonAnimTimer)
	}

	g.tuning.RedKillInterval = 1
	g.tuning.PurpleKillInterval = 1
	g.session.Kills.TotalSkeletons = 1
	g.spawnMilestoneSkeletonsIfNeeded()
	if len(g.skeleton) != 2 {
		t.Fatalf("milestone skeletons = %v, want red and purple", g.skeleton)
	}
	g.session.Progression.Level = g.tuning.RedOnlyLevel
	g.spawnMilestoneSkeletonsIfNeeded()

	g.tuning.BronzeKillInterval = 1
	g.tuning.SilverKillInterval = 2
	g.tuning.GoldKillInterval = 3
	g.session.NextChestMilestone = 1
	g.session.Kills.TotalSkeletons = 3
	g.session.Progression.Level = 1
	g.spawnChestsForMilestones()
	if len(g.chests) < 3 {
		t.Fatalf("milestone chests = %v, want bronze/silver/gold", g.chests)
	}
}

func TestCoverageWeaponBranches(t *testing.T) {
	g := New()
	g.skeleton = nil
	g.updateFireballCasting(1)
	if g.session.Casts.Fireball != 0 {
		t.Fatal("fireball timer was not reset without skeletons")
	}
	g.session.Progression.SimultaneousFireball = 0
	if got := g.closestSkeletons(Vec2{}, nil, 0); got != nil {
		t.Fatalf("closestSkeletons with zero limit = %v, want nil", got)
	}
	g.skeleton = []Skeleton{{ID: 21, Pos: Vec2{X: 100}}, {ID: 22, Pos: Vec2{X: 80}}, {ID: 23, Pos: Vec2{X: 10}}}
	if got := g.closestSkeletons(Vec2{}, nil, 2); !reflect.DeepEqual(got, []int{2, 1}) {
		t.Fatalf("closestSkeletons sorted targets = %v, want [2 1]", got)
	}

	g.session.Progression.SimultaneousFireball = 1
	g.skeleton = []Skeleton{{ID: 1, Pos: Vec2{X: 100}, HP: 1, Reward: 1}, {ID: 2, Pos: Vec2{X: 20}, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)
	g.spawnFireballs()
	if len(g.fireball) != 1 || g.fireball[0].TargetID != 2 {
		t.Fatalf("spawned fireball = %v, want nearest unreserved target 2", g.fireball)
	}
	g.updateHomingFireball(0, 1, 0)
	if len(g.fireball) != 0 {
		t.Fatalf("homing fireball count = %d, want removed on hit", len(g.fireball))
	}
	g.fireball = []Fireball{{Pos: Vec2{}, TargetID: 1}}
	g.skeleton = []Skeleton{{ID: 1, Pos: Vec2{X: 100}, HP: 2, Reward: 1}}
	g.updateHomingFireball(0, 0, 0.01)
	if len(g.fireball) != 1 || g.fireball[0].Pos == (Vec2{}) {
		t.Fatalf("homing fireball did not move toward distant target: %v", g.fireball)
	}

	g.fireball = []Fireball{{Pos: Vec2{}, Velocity: Vec2{X: 1}}}
	g.skeleton = []Skeleton{{ID: 3, Pos: Vec2{X: 5}, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)
	g.updateUntargetedFireball(0, 0.1)
	if len(g.fireball) != 0 {
		t.Fatal("untargeted fireball did not hit segment target")
	}
	g.fireball = []Fireball{{Pos: Vec2{}, Velocity: Vec2{X: 1}}}
	g.skeleton = nil
	g.spatial.Rebuild(g.skeleton)
	g.updateUntargetedFireball(0, g.tuning.FireballUntargetedLifetime)
	if len(g.fireball) != 0 {
		t.Fatal("untargeted fireball did not expire")
	}
	g.fireball = []Fireball{{Pos: Vec2{}, TargetID: 999, Velocity: Vec2{X: 1}}}
	g.updateFireballs(0.01)
	if len(g.fireball) != 1 || g.fireball[0].TargetID != 0 {
		t.Fatalf("updateFireballs untargeted branch = %v, want retained untargeted fireball", g.fireball)
	}
	g.session.LevelUpChoiceActive = true
	g.fireball = []Fireball{{Pos: Vec2{}, TargetID: 999, Velocity: Vec2{X: 1}}}
	g.updateFireballs(0.01)

	g.session.Progression.ApplyLevelUpOption(LearnLightning)
	g.skeleton = nil
	g.updateLightningCasting(1)
	if g.session.Casts.Lightning != 0 {
		t.Fatal("lightning timer was not reset without skeletons")
	}
	g.session.Progression.LightningUnlocked = false
	g.skeleton = []Skeleton{{ID: 4, Pos: Vec2{X: 1}, HP: 1, Reward: 1}}
	g.updateLightningCasting(1)
	if g.session.Casts.Lightning != 0 {
		t.Fatal("lightning timer was not reset while locked")
	}
	g.session.Progression.LightningUnlocked = true
	g.session.Progression.LightningBounceCount = 1
	g.fireball = []Fireball{{TargetID: 4}}
	if targets := g.chainLightningTargets(); len(targets) != 0 {
		t.Fatalf("reserved lightning targets = %v, want none", targets)
	}
	g.session.Progression.LightningUnlocked = false
	if targets := g.chainLightningTargets(); targets != nil {
		t.Fatalf("locked lightning targets = %v, want nil", targets)
	}
	g.session.Progression.LightningUnlocked = true
	g.fireball = nil
	g.skeleton = []Skeleton{{ID: 6, Pos: Vec2{X: 100}}, {ID: 7, Pos: Vec2{X: 10}}}
	if targets := g.chainLightningTargets(); len(targets) == 0 || targets[0].targetID != 7 {
		t.Fatalf("chain lightning nearest target = %v, want skeleton 7 first", targets)
	}
	if got := g.applyLightningStrikes([]lightningStrikeTarget{{targetID: 999, end: Vec2{X: 1}}}); got != 0 {
		t.Fatalf("missing lightning target levelups = %d, want 0", got)
	}

	g.session.Progression.ApplyLevelUpOption(LearnBeam)
	g.skeleton = nil
	g.updateBeamCasting(1)
	if g.session.Casts.Beam != 0 {
		t.Fatal("beam timer was not reset without skeletons")
	}
	g.session.Progression.BeamUnlocked = false
	g.skeleton = []Skeleton{{ID: 5, Pos: Vec2{X: 1}, HP: 1, Reward: 1}}
	g.updateBeamCasting(1)
	if g.session.Casts.Beam != 0 {
		t.Fatal("beam timer was not reset while locked")
	}
	g.session.Progression.BeamUnlocked = true
	g.player.MoveDir = Vec2{X: 0, Y: 2}
	if got := g.playerBeamDirection(); got != (Vec2{Y: 1}) {
		t.Fatalf("beam direction = %+v, want normalized move dir", got)
	}
	g.player.MoveDir = Vec2{}
	g.player.Facing = 1
	g.skeleton = []Skeleton{
		{ID: 8, Pos: Vec2{X: -1}},
		{ID: 9, Pos: Vec2{X: 10, Y: g.tuning.BeamHitWidth + 1}},
		{ID: 10, Pos: Vec2{X: 50}},
		{ID: 11, Pos: Vec2{X: 30}},
		{ID: 12, Pos: Vec2{X: 20}},
	}
	if got := g.beamTargets(Vec2{X: 1}, 40, g.tuning.BeamHitWidth, 2); !reflect.DeepEqual(got, []int{12, 11}) {
		t.Fatalf("beam targets = %v, want closest two in beam", got)
	}
	g.skeleton = []Skeleton{
		{ID: 13, Pos: Vec2{X: 80}, HP: 1},
		{ID: 14, Pos: Vec2{X: 70}, HP: 1},
		{ID: 15, Pos: Vec2{X: 20}, HP: 1},
	}
	if got := g.beamTargets(Vec2{X: 1}, 100, g.tuning.BeamHitWidth, 2); !reflect.DeepEqual(got, []int{15, 14}) {
		t.Fatalf("replacement beam targets = %v, want [15 14]", got)
	}
	if got := g.applyBeamDamage([]int{999}, 1); got != 0 {
		t.Fatalf("missing beam target levelups = %d, want 0", got)
	}
	if got := g.applyBeamDamage([]int{13}, 0); got != 0 {
		t.Fatalf("zero-budget beam levelups = %d, want 0", got)
	}

	g.session.Progression.ApplyLevelUpOption(LearnOrb)
	g.orbs = []OrbitalOrb{{Active: false, MissingOrbitProgress: math.Pi*2 - 0.01}}
	g.updateOrbitalOrbs(1)
	if !g.orbs[0].Active {
		t.Fatal("inactive orbital orb did not reactivate after a full orbit")
	}
	locked := New()
	locked.updateOrbitalOrbs(1)
	noTarget := New()
	noTarget.session.Progression.OrbitalOrbUnlocked = true
	noTarget.session.Progression.upgradedOrbitalOrbCount = 0
	noTarget.updateOrbitalOrbs(1)
	shrink := New()
	shrink.session.Progression.OrbitalOrbUnlocked = true
	shrink.session.Progression.upgradedOrbitalOrbCount = 1
	shrink.orbs = []OrbitalOrb{{Active: true}, {Active: true}}
	shrink.syncOrbitalOrbCount()
	if len(shrink.orbs) != 1 {
		t.Fatalf("shrunk orbital orb count = %d, want 1", len(shrink.orbs))
	}
	g.orbs = []OrbitalOrb{{Active: false}}
	g.checkOrbitalOrbCollisions()
	g.orbs = []OrbitalOrb{{Active: false}}
	g.updateOrbAnimation(1)
	if g.orbAnimTimer != 0 {
		t.Fatalf("orb animation timer = %v, want 0 without active orbs", g.orbAnimTimer)
	}

	g.session.Progression.ApplyLevelUpOption(LearnMeteor)
	g.skeleton = nil
	g.updateMeteorCasting(1)
	if g.session.Casts.Meteor != 0 {
		t.Fatal("meteor timer was not reset without skeletons")
	}
	g.session.GameOver = true
	g.impactMeteor(Vec2{})
	if len(g.effects) != 0 {
		t.Fatal("game-over meteor impact added an effect")
	}
}

func TestCoverageSpatialAndHitHelpers(t *testing.T) {
	index := NewSpatialIndex(10)
	skeletons := []Skeleton{{Pos: Vec2{X: 0}}, {Pos: Vec2{X: 30}}}
	index.Rebuild(skeletons)
	visited := []int{}
	index.ForEachRect(Vec2{X: -1, Y: -1}, Vec2{X: 40, Y: 1}, func(i int) bool {
		visited = append(visited, i)
		return false
	})
	if !reflect.DeepEqual(visited, []int{0}) {
		t.Fatalf("ForEachRect early stop visited %v", visited)
	}

	g := New()
	g.skeleton = skeletons
	g.spatial.Rebuild(g.skeleton)
	if got := g.firstSkeletonHitBySegment(Vec2{}, Vec2{}, 1); got != 0 {
		t.Fatalf("zero-length segment hit = %d, want 0", got)
	}
	if got := g.pointInRect(5, 5, 0, 0, 10, 10); !got {
		t.Fatal("pointInRect missed inside point")
	}
	if got := g.pointInRect(15, 5, 0, 0, 10, 10); got {
		t.Fatal("pointInRect accepted outside point")
	}
}

func TestCoveragePresentationEdgeValues(t *testing.T) {
	if got := effectFadeAlpha(1, 0); got != 0 {
		t.Fatalf("effectFadeAlpha zero max = %d, want 0", got)
	}
	if got := lightningHitEffectAlpha(1, 0); got != 0 {
		t.Fatalf("lightningHitEffectAlpha zero max = %d, want 0", got)
	}
	if got := scaleAlpha(100, 1.5); got != 100 {
		t.Fatalf("scaleAlpha above one = %d, want 100", got)
	}
	if got := linearPingPong(0.75, 0.5); math.Abs(got-0.5) > 0.0001 {
		t.Fatalf("linearPingPong falling = %v, want 0.5", got)
	}
	if got := linearPingPong(-0.25, 0.5); math.Abs(got-0.5) > 0.0001 {
		t.Fatalf("linearPingPong negative = %v, want 0.5", got)
	}
	if got := linearPingPong(1, 0); got != 0 {
		t.Fatalf("linearPingPong zero period = %v, want 0", got)
	}
	if got := meteorImpactPresentation(Effect{TTL: 1, MaxTTL: 0}); got != (meteorImpactStyle{Scale: 1.25, Alpha: 0}) {
		t.Fatalf("meteorImpactPresentation zero max = %+v", got)
	}
	if got := meteorImpactPresentation(Effect{TTL: 0.2, MaxTTL: 0.3}); got.Alpha != 1 || got.Scale != 1 {
		t.Fatalf("meteorImpactPresentation hold = %+v", got)
	}
	if got := meteorImpactPresentation(Effect{TTL: 0, MaxTTL: 0.4}); got.Alpha != 0 || got.Scale != 1.25 {
		t.Fatalf("meteorImpactPresentation fade = %+v", got)
	}
	if got := redrawPulseScale(0); got != 1 {
		t.Fatalf("redrawPulseScale idle = %v, want 1", got)
	}
	if got := redrawPulseScale(redrawFailurePulseDuration); math.Abs(got-1) > 0.0001 {
		t.Fatalf("redrawPulseScale start = %v, want 1", got)
	}
	if got := redrawCoinAlpha(true, 0); got != 255 {
		t.Fatalf("redrawCoinAlpha can redraw = %d, want 255", got)
	}
	if got := withAlpha(c64Text, 12); got.A != 12 {
		t.Fatalf("withAlpha alpha = %d, want 12", got.A)
	}
}

func TestCoverageFontAndIconBranches(t *testing.T) {
	if fontFaceForSize(0) != basicFontFace() {
		t.Fatal("fontFaceForSize(0) did not return the fallback basic face")
	}
	oldFont := hudFont
	hudFont = nil
	if fontFaceForSize(12) != basicFontFace() {
		t.Fatal("fontFaceForSize with nil font did not return fallback basic face")
	}
	hudFont = oldFont
	if _, name := loadSystemFontByFullName([]string{"/definitely/not/a/font"}, []string{"missing"}); name != "" {
		t.Fatalf("missing system font name = %q, want empty", name)
	}
	icons := WindowIcons()
	if len(icons) != len(appIconPaths) {
		t.Fatalf("WindowIcons count = %d, want %d", len(icons), len(appIconPaths))
	}
}

func basicFontFace() font.Face {
	return fontFaceForSize(-1)
}

func TestCoverageSelectAndAdvanceBranches(t *testing.T) {
	g := New()
	if consumed, err := g.selectGameOverOption(""); consumed || err != nil {
		t.Fatalf("default game-over option = consumed %v err %v", consumed, err)
	}
	if consumed, err := g.selectGameOverOption("exit"); !consumed || !errors.Is(err, ebiten.Termination) {
		t.Fatalf("exit option = consumed %v err %v", consumed, err)
	}

	if g.advanceChestReward() {
		t.Fatal("advanceChestReward succeeded while inactive")
	}
	g.session.ChestRewardActive = true
	g.session.ActiveChestRewardItems = []ChestRewardDisplayItem{{Option: FireRate, Title: "x"}}
	if !g.advanceChestReward() || g.session.ChestRewardActive || g.session.ActiveChestRewardItems != nil {
		t.Fatalf("advanceChestReward state = active %v items %v", g.session.ChestRewardActive, g.session.ActiveChestRewardItems)
	}

	if g.selectLevelUpOptionAt(0) {
		t.Fatal("selectLevelUpOptionAt succeeded while inactive")
	}
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{ExtraLife}
	if g.selectLevelUpOptionAt(-1) || g.selectLevelUpOptionAt(1) {
		t.Fatal("selectLevelUpOptionAt accepted invalid index")
	}
	if !g.selectLevelUpOptionAt(0) {
		t.Fatal("selectLevelUpOptionAt rejected valid option")
	}
	if !isKillAllAndGrantExperienceKey(ebiten.KeyDigit1) || !isKillAllAndGrantExperienceKey(ebiten.KeyNumpad1) || isKillAllAndGrantExperienceKey(ebiten.KeyQ) {
		t.Fatal("kill-all key predicate mismatch")
	}
}

func TestCoverageInputOverridesDriveOverlayInputBranches(t *testing.T) {
	oldKeyPressed := ebitenIsKeyPressed
	oldCursor := ebitenCursorPosition
	oldKeyJustPressed := inpututilIsKeyJustPressed
	oldMouseJustPressed := inpututilIsMouseButtonJustPressed
	defer func() {
		ebitenIsKeyPressed = oldKeyPressed
		ebitenCursorPosition = oldCursor
		inpututilIsKeyJustPressed = oldKeyJustPressed
		inpututilIsMouseButtonJustPressed = oldMouseJustPressed
	}()

	pressed := map[ebiten.Key]bool{}
	justPressed := map[ebiten.Key]bool{}
	mousePressed := false
	cursorX, cursorY := 0, 0
	ebitenIsKeyPressed = func(key ebiten.Key) bool { return pressed[key] }
	inpututilIsKeyJustPressed = func(key ebiten.Key) bool { return justPressed[key] }
	inpututilIsMouseButtonJustPressed = func(button ebiten.MouseButton) bool {
		return button == ebiten.MouseButtonLeft && mousePressed
	}
	ebitenCursorPosition = func() (int, int) { return cursorX, cursorY }

	g := New()
	g.session.GameOver = true
	mousePressed = true
	cursorX, cursorY = 400, 370
	consumed, err := g.updateOverlayInput()
	if !consumed || !errors.Is(err, ebiten.Termination) {
		t.Fatalf("game-over mouse exit = consumed %v err %v", consumed, err)
	}

	g = New()
	g.session.GameOver = true
	mousePressed = false
	if consumed, err = g.updateOverlayInput(); consumed || err != nil {
		t.Fatalf("game-over no input = consumed %v err %v", consumed, err)
	}

	g = New()
	g.session.ChestRewardActive = true
	pressed[ebiten.KeyArrowLeft] = true
	justPressed = map[ebiten.Key]bool{chestRewardAdvanceKey(): true}
	consumed, err = g.updateOverlayInput()
	if !consumed || err != nil || g.session.ChestRewardActive {
		t.Fatalf("chest advance = consumed %v err %v active %v", consumed, err, g.session.ChestRewardActive)
	}

	g = New()
	g.skeleton = []Skeleton{{ID: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)
	justPressed = map[ebiten.Key]bool{ebiten.KeyDigit1: true}
	consumed, err = g.updateOverlayInput()
	if consumed || err != nil || len(g.skeleton) != 0 {
		t.Fatalf("kill-all key = consumed %v err %v skeletons %d", consumed, err, len(g.skeleton))
	}

	g = New()
	g.session.LevelUpChoiceActive = true
	g.session.PendingLevelUpLevels = []int{1}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraLife}
	g.session.CollectedCoins = 10
	justPressed = map[ebiten.Key]bool{levelUpRedrawKey(): true}
	consumed, err = g.updateOverlayInput()
	if !consumed || err != nil {
		t.Fatalf("redraw key = consumed %v err %v", consumed, err)
	}

	g = New()
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{ExtraLife}
	justPressed = map[ebiten.Key]bool{ebiten.KeyQ: true}
	consumed, err = g.updateOverlayInput()
	if !consumed || err != nil || g.session.PlayerLives != g.tuning.InitialPlayerLives+1 {
		t.Fatalf("level option key = consumed %v err %v lives %d", consumed, err, g.session.PlayerLives)
	}

	g = New()
	g.session.LevelUpChoiceActive = true
	g.session.PendingLevelUpLevels = []int{1}
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate, ExtraLife}
	g.session.CollectedCoins = 10
	mousePressed = true
	justPressed = nil
	cursorX, cursorY = 400-150, int(float64(g.screenH)/2-90+94+float64(max(2, len(g.session.ActiveLevelUpOptions)))*52+14)
	consumed, err = g.updateOverlayInput()
	if !consumed || err != nil {
		t.Fatalf("redraw mouse = consumed %v err %v", consumed, err)
	}

	g = New()
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{ExtraLife}
	cursorX, cursorY = 400-78, int(float64(g.screenH)/2-90+94)
	consumed, err = g.updateOverlayInput()
	if !consumed || err != nil || g.session.PlayerLives != g.tuning.InitialPlayerLives+1 {
		t.Fatalf("level option mouse = consumed %v err %v lives %d", consumed, err, g.session.PlayerLives)
	}
}

func TestCoverageUpdateEarlyExitBranches(t *testing.T) {
	oldMouseJustPressed := inpututilIsMouseButtonJustPressed
	oldCursor := ebitenCursorPosition
	defer func() {
		inpututilIsMouseButtonJustPressed = oldMouseJustPressed
		ebitenCursorPosition = oldCursor
	}()

	g := New()
	g.session.GameOver = true
	inpututilIsMouseButtonJustPressed = func(ebiten.MouseButton) bool { return true }
	ebitenCursorPosition = func() (int, int) { return 400, 370 }
	if err := g.Update(); !errors.Is(err, ebiten.Termination) {
		t.Fatalf("Update game-over exit err = %v, want termination", err)
	}

	g = New()
	g.session.GameOver = true
	ebitenCursorPosition = func() (int, int) { return 400, 322 }
	if err := g.Update(); err != nil {
		t.Fatalf("Update game-over restart err = %v", err)
	}

	inpututilIsMouseButtonJustPressed = func(ebiten.MouseButton) bool { return false }
	for _, configure := range []func(*Game){
		func(g *Game) { g.session.LevelUpChoiceActive = true },
		func(g *Game) { g.session.ChestRewardActive = true },
		func(g *Game) { g.session.GameOver = true },
	} {
		g = New()
		g.hasUpdated = true
		configure(g)
		if err := g.Update(); err != nil {
			t.Fatalf("paused Update err = %v", err)
		}
	}

	g = New()
	if w, h := g.Layout(0, -4); w != 1 || h != 1 {
		t.Fatalf("Layout clamped size = %dx%d, want 1x1", w, h)
	}
}

func TestCoverageUpdateWeaponLevelUpExitBranches(t *testing.T) {
	run := func(name string, configure func(*Game)) {
		t.Run(name, func(t *testing.T) {
			g := New()
			g.hasUpdated = true
			g.tuning.RedKillInterval = 0
			g.tuning.PurpleKillInterval = 0
			g.session.NextChestMilestone = 1_000_000
			configure(g)
			if err := g.Update(); err != nil {
				t.Fatalf("Update err = %v", err)
			}
			if !g.session.LevelUpChoiceActive {
				t.Fatal("Update did not stop on newly active level-up choice")
			}
		})
	}

	run("orbital", func(g *Game) {
		g.session.Progression.ApplyLevelUpOption(LearnOrb)
		dt := 1.0 / float64(TargetTPS)
		angle := g.session.Progression.OrbitalAngularSpeed() * dt
		g.skeleton = []Skeleton{{ID: 100, Pos: Vec2{X: math.Cos(angle) * g.tuning.OrbitalOrbRadius, Y: math.Sin(angle) * g.tuning.OrbitalOrbRadius}, HP: 1, Reward: 1}}
		g.spatial.Rebuild(g.skeleton)
	})
	run("lightning", func(g *Game) {
		g.session.Progression.ApplyLevelUpOption(LearnLightning)
		g.session.Casts.Lightning = g.session.Progression.LightningCastInterval()
		g.skeleton = []Skeleton{{ID: 101, Pos: Vec2{X: 80}, HP: 1, Reward: 1}}
		g.spatial.Rebuild(g.skeleton)
	})
	run("fireball", func(g *Game) {
		g.session.Casts.Fireball = g.session.Progression.FireballCastInterval()
		g.skeleton = []Skeleton{{ID: 102, Pos: Vec2{X: 1}, HP: 1, Reward: 1}}
		g.spatial.Rebuild(g.skeleton)
	})
	run("meteor", func(g *Game) {
		g.session.Progression.ApplyLevelUpOption(LearnMeteor)
		g.meteors = []MeteorProjectile{{Pos: Vec2{X: 100}, Start: Vec2{X: 100}, Impact: Vec2{X: 100}, Age: g.tuning.MeteorFallDuration}}
		g.skeleton = []Skeleton{{ID: 103, Pos: Vec2{X: 100}, HP: 1, Reward: 1}}
		g.spatial.Rebuild(g.skeleton)
	})
}

func TestCoverageMoreGameLogicBranches(t *testing.T) {
	oldKeyPressed := ebitenIsKeyPressed
	defer func() { ebitenIsKeyPressed = oldKeyPressed }()

	pressed := map[ebiten.Key]bool{ebiten.KeyArrowLeft: true, ebiten.KeyArrowUp: true}
	ebitenIsKeyPressed = func(key ebiten.Key) bool { return pressed[key] }
	g := New()
	g.updatePlayer(1)
	if !g.player.Moving || g.player.Facing != -1 || g.player.Pos == (Vec2{}) {
		t.Fatalf("left/up movement state = %+v", g.player)
	}
	pressed = map[ebiten.Key]bool{ebiten.KeyArrowRight: true, ebiten.KeyArrowDown: true}
	g.updatePlayer(1)
	if g.player.Facing != 1 {
		t.Fatalf("right/down facing = %v, want 1", g.player.Facing)
	}
	pressed = nil
	g.updatePlayer(1)
	if g.player.Moving || g.player.AnimFrame != 0 {
		t.Fatalf("idle player state = %+v", g.player)
	}

	g.suppressedMovement = nil
	pressed = map[ebiten.Key]bool{ebiten.KeyArrowRight: true}
	g.suppressHeldMovementKeys(func(key ebiten.Key) bool { return pressed[key] })
	if !g.suppressedMovement[ebiten.KeyArrowRight] {
		t.Fatal("suppressHeldMovementKeys did not initialize and store right key")
	}
	pressed = map[ebiten.Key]bool{ebiten.KeyArrowRight: true, ebiten.KeyArrowDown: true}
	move := g.playerMovementVector(func(key ebiten.Key) bool { return pressed[key] })
	if move != (Vec2{Y: -1}) {
		t.Fatalf("suppressed right/down movement = %+v, want down only", move)
	}

	p := NewProgression(DefaultTuning())
	p.Level = p.tuning.RedOnlyLevel
	redInterval := p.SkeletonSpawnInterval()
	p.Level = p.tuning.PurpleOnlyLevel
	purpleInterval := p.SkeletonSpawnInterval()
	if redInterval <= 0 || purpleInterval <= 0 {
		t.Fatalf("spawn intervals red=%v purple=%v, want positive values", redInterval, purpleInterval)
	}

	g = New()
	g.tuning.ChestPickupDistance = 10_000
	pos := g.randomChestPosition()
	if pos != (Vec2{X: g.player.Pos.X + math.Max(48, float64(g.screenW)/2-g.tuning.ChestSpawnMargin), Y: g.player.Pos.Y}) {
		t.Fatalf("fallback chest position = %+v", pos)
	}

	g.applyUpgradeEffect(ExtraLife)
	if g.session.PlayerLives != g.tuning.InitialPlayerLives+1 {
		t.Fatalf("extra life lives = %d", g.session.PlayerLives)
	}
	g.skeleton = []Skeleton{{ID: 1, Reward: 1}}
	g.applyUpgradeEffect(HalveSkeletons)
	if len(g.skeleton) != 1 {
		t.Fatalf("halve one skeleton should leave it alive, got %d", len(g.skeleton))
	}

	g.skeleton = []Skeleton{{ID: 10, Pos: g.player.Pos, HP: 1, Reward: 1}}
	g.spatial.Rebuild(g.skeleton)
	g.session.PlayerHitInvulnerability = 1
	beforeLives := g.session.PlayerLives
	g.checkSkeletonCollisions()
	if g.session.PlayerLives != beforeLives {
		t.Fatalf("invulnerable collision changed lives to %d", g.session.PlayerLives)
	}
	if got := g.damageSkeleton(-1, 1, AttackFireball, true); got != 0 {
		t.Fatalf("invalid damage levelups = %d", got)
	}
	if got := g.destroySkeleton(-1, AttackFireball); got != 0 {
		t.Fatalf("invalid destroy levelups = %d", got)
	}

	if got := LevelUpOption(999).Title(0); got != "UNKNOWN" {
		t.Fatalf("unknown level-up title = %q", got)
	}
}

func TestCoverageMoreRenderAndPlatformBranches(t *testing.T) {
	if spriteBoundsVisible(0, 600, 10, 10, 10, 10, 0) {
		t.Fatal("spriteBoundsVisible accepted zero screen width")
	}
	screen := ebiten.NewImage(32, 32)
	oldPanelEdge := c64PanelEdge
	c64PanelEdge.A = 255
	New().panelWithAlpha(screen, 1, 1, 10, 10, 255)
	c64PanelEdge = oldPanelEdge

	g := New()
	g.session.LevelUpChoiceActive = true
	g.session.ActiveLevelUpOptions = []LevelUpOption{FireRate}
	g.session.LevelUpOverlayTimer = modalFadeDuration
	g.session.LevelUpOptionFadeTimer = modalFadeDuration
	g.session.LevelUpRedrawFadeTimer = redrawStatusFadeDuration
	g.session.CollectedCoins = 0
	g.drawLevelUpOverlay(screen)

	oldPaths := hudFontPaths
	oldFont := hudFont
	oldFaces := hudFontFaces
	oldNewOpenTypeFace := newOpenTypeFace
	defer func() {
		hudFontPaths = oldPaths
		hudFont = oldFont
		hudFontFaces = oldFaces
		newOpenTypeFace = oldNewOpenTypeFace
	}()
	hudFontPaths = []string{"/definitely/not/a/font"}
	loadedFont, name := loadHUDFont()
	if loadedFont == nil || name != "Go Mono Bold" {
		t.Fatalf("fallback HUD font = %v %q", loadedFont, name)
	}
	invalidPath := filepath.Join(t.TempDir(), "bad-font.ttf")
	if err := os.WriteFile(invalidPath, []byte("not a font"), 0o600); err != nil {
		t.Fatal(err)
	}
	if loadedFont, name := loadSystemFontByFullName([]string{invalidPath}, []string{"bad"}); loadedFont != nil || name != "" {
		t.Fatalf("invalid system font = %v %q", loadedFont, name)
	}
	newOpenTypeFace = func(*opentype.Font, *opentype.FaceOptions) (font.Face, error) {
		return nil, errors.New("face")
	}
	hudFontFaces = map[int]font.Face{}
	if fontFaceForSize(12) != basicFontFace() {
		t.Fatal("font face creation error did not return fallback face")
	}

	oldReadIconFile := readIconFile
	oldDecodeIcon := decodeIcon
	readIconFile = func(string) ([]byte, error) { return nil, errors.New("boom") }
	if _, err := loadWindowIcons(); err == nil {
		t.Fatal("loadWindowIcons read error = nil")
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("WindowIcons did not panic on icon read failure")
			}
		}()
		WindowIcons()
	}()
	readIconFile = func(string) ([]byte, error) { return []byte("bad"), nil }
	decodeIcon = func(io.Reader) (image.Image, string, error) { return nil, "", errors.New("decode") }
	if _, err := loadWindowIcons(); err == nil {
		t.Fatal("loadWindowIcons decode error = nil")
	}
	readIconFile = oldReadIconFile
	decodeIcon = oldDecodeIcon
}
