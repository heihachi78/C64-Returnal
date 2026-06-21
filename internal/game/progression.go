package game

import "math"

type Progression struct {
	tuning Tuning

	Level                   int
	Experience              int
	NextExperience          int
	SimultaneousFireball    int
	LightningBounceCount    int
	LightningUnlocked       bool
	OrbitalOrbUnlocked      bool
	BeamUnlocked            bool
	MeteorUnlocked          bool
	fireRateUpgrades        int
	lightningRateUpgrades   int
	orbitalSpeedUpgrades    int
	beamRateUpgrades        int
	meteorRateUpgrades      int
	upgradedOrbitalOrbCount int
	beamKillLevel           int
	upgradedBeamKillCount   int
	upgradedMeteorCount     int
	skeletonSpawnRate       float64
}

func NewProgression(t Tuning) Progression {
	p := Progression{tuning: t}
	p.Reset()
	return p
}

func (p *Progression) Reset() {
	p.Level = 1
	p.Experience = 0
	p.NextExperience = ExperienceRequirement(p.Level)
	p.SimultaneousFireball = 1
	p.LightningBounceCount = 0
	p.LightningUnlocked = false
	p.OrbitalOrbUnlocked = false
	p.BeamUnlocked = false
	p.MeteorUnlocked = false
	p.fireRateUpgrades = 0
	p.lightningRateUpgrades = 0
	p.orbitalSpeedUpgrades = 0
	p.beamRateUpgrades = 0
	p.meteorRateUpgrades = 0
	p.upgradedOrbitalOrbCount = 1
	p.beamKillLevel = 1
	p.upgradedBeamKillCount = 1
	p.upgradedMeteorCount = 1
	p.skeletonSpawnRate = 1
}

func ExperienceRequirement(level int) int {
	return max(1, level*level/2)
}

func (p Progression) SkeletonSpawnInterval() float64 {
	interval := p.tuning.InitialSkeletonSpawn * math.Pow(p.tuning.SkeletonIntervalMultiplier, float64(p.Level-1))
	if p.Level >= p.tuning.RedOnlyLevel {
		interval *= p.tuning.RedOnlySpawnMultiplier
	}
	if p.Level >= p.tuning.PurpleOnlyLevel {
		interval *= p.tuning.PurpleOnlySpawnMultiplier
	}
	if p.Level >= p.tuning.BlackOnlyLevel {
		interval *= p.tuning.BlackOnlySpawnMultiplier
	}
	if p.skeletonSpawnRate <= 0 {
		return math.Inf(1)
	}
	return interval / p.skeletonSpawnRate
}

func (p *Progression) ScaleSkeletonSpawnRate(factor float64) {
	p.skeletonSpawnRate = math.Max(0, p.skeletonSpawnRate) * math.Max(0, factor)
}

func (p Progression) FireballCastInterval() float64 {
	return p.tuning.InitialFireballCast * math.Pow(p.tuning.FireballIntervalMultiplier, float64(p.fireRateUpgrades))
}

func (p Progression) MageRawDPS() float64 {
	return damageRate(p.SimultaneousFireball, p.FireballCastInterval()) +
		damageRate(p.LightningStrikeCount(), p.LightningCastInterval()) +
		float64(p.OrbitalOrbCount())*p.OrbitalAngularSpeed()/(math.Pi*2) +
		damageRate(p.BeamKillCount(), p.BeamCastInterval()) +
		damageRate(p.MeteorCount(), p.MeteorCastInterval())
}

func damageRate(count int, interval float64) float64 {
	if count <= 0 || interval <= 0 {
		return 0
	}
	return float64(count) / interval
}

func (p Progression) LightningCastInterval() float64 {
	return p.tuning.InitialLightningCast * math.Pow(p.tuning.LightningIntervalMultiplier, float64(p.lightningRateUpgrades))
}

func (p Progression) LightningStrikeCount() int {
	if !p.LightningUnlocked {
		return 0
	}
	return p.LightningBounceCount + 1
}

func (p Progression) OrbitalAngularSpeed() float64 {
	return p.tuning.InitialOrbitalAngularSpeed * math.Pow(p.tuning.OrbitalSpeedUpgradeMultipler, float64(p.orbitalSpeedUpgrades))
}

func (p Progression) OrbitalOrbCount() int {
	if !p.OrbitalOrbUnlocked {
		return 0
	}
	return p.upgradedOrbitalOrbCount
}

func (p Progression) BeamCastInterval() float64 {
	return p.tuning.InitialBeamCast * math.Pow(p.tuning.BeamIntervalMultiplier, float64(p.beamRateUpgrades))
}

func (p Progression) BeamKillCount() int {
	if !p.BeamUnlocked {
		return 0
	}
	return p.upgradedBeamKillCount
}

func (p Progression) BeamKillUpgradeBonus() int {
	return p.beamKillLevel + 1
}

func (p Progression) MeteorCastInterval() float64 {
	return p.tuning.InitialMeteorCast * math.Pow(p.tuning.MeteorIntervalMultiplier, float64(p.meteorRateUpgrades))
}

func (p Progression) MeteorSpawnInterval() float64 {
	count := p.MeteorCount()
	if count <= 0 {
		return math.Inf(1)
	}
	return p.MeteorCastInterval() / float64(count)
}

func (p Progression) MeteorCount() int {
	if !p.MeteorUnlocked {
		return 0
	}
	return p.upgradedMeteorCount
}

func (p Progression) LearnedSkills() []LearnedSkill {
	skills := []LearnedSkill{SkillFireball}
	if p.LightningUnlocked {
		skills = append(skills, SkillLightning)
	}
	if p.OrbitalOrbUnlocked {
		skills = append(skills, SkillOrbitalOrb)
	}
	if p.BeamUnlocked {
		skills = append(skills, SkillBeam)
	}
	if p.MeteorUnlocked {
		skills = append(skills, SkillMeteor)
	}
	return skills
}

func (p Progression) AvailableLevelUpOptions() []LevelUpOption {
	options := []LevelUpOption{}
	for _, option := range []LevelUpOption{FireRate, ExtraFireball, ExtraLife, HalveSkeletons} {
		if p.LevelUpOptionAvailable(option) {
			options = append(options, option)
		}
	}
	if p.LightningUnlocked {
		options = appendAvailableLevelUpOptions(p, options, LightningBounce, LightningRate)
	} else {
		options = append(options, LearnLightning)
	}
	if p.OrbitalOrbUnlocked {
		options = appendAvailableLevelUpOptions(p, options, ExtraOrb, OrbitalSpeed)
	} else {
		options = append(options, LearnOrb)
	}
	if p.BeamUnlocked {
		options = appendAvailableLevelUpOptions(p, options, BeamRate, BeamKillCount)
	} else {
		options = append(options, LearnBeam)
	}
	if p.MeteorUnlocked {
		options = appendAvailableLevelUpOptions(p, options, ExtraMeteor, MeteorRate)
	} else {
		options = append(options, LearnMeteor)
	}
	return options
}

func appendAvailableLevelUpOptions(p Progression, options []LevelUpOption, candidates ...LevelUpOption) []LevelUpOption {
	for _, option := range candidates {
		if p.LevelUpOptionAvailable(option) {
			options = append(options, option)
		}
	}
	return options
}

func (p Progression) AvailableUpgradeOptionsForSkill(skill LearnedSkill) []LevelUpOption {
	options := skill.UpgradeOptions()
	filtered := make([]LevelUpOption, 0, len(options))
	for _, option := range options {
		if p.LevelUpOptionAvailable(option) {
			filtered = append(filtered, option)
		}
	}
	return filtered
}

func (p Progression) LevelUpOptionAvailable(option LevelUpOption) bool {
	interval, ok := p.attackSpawnIntervalAfterOption(option)
	return !ok || interval >= minAttackSpawnInterval
}

func (p Progression) attackSpawnIntervalAfterOption(option LevelUpOption) (float64, bool) {
	switch option {
	case FireRate:
		return p.tuning.InitialFireballCast * math.Pow(p.tuning.FireballIntervalMultiplier, float64(p.fireRateUpgrades+1)), true
	case LightningRate:
		return p.tuning.InitialLightningCast * math.Pow(p.tuning.LightningIntervalMultiplier, float64(p.lightningRateUpgrades+1)), true
	case BeamRate:
		return p.tuning.InitialBeamCast * math.Pow(p.tuning.BeamIntervalMultiplier, float64(p.beamRateUpgrades+1)), true
	case MeteorRate:
		return p.tuning.InitialMeteorCast * math.Pow(p.tuning.MeteorIntervalMultiplier, float64(p.meteorRateUpgrades+1)), true
	default:
		return 0, false
	}
}

func (p *Progression) GainExperience(amount int) int {
	p.Experience += max(0, amount)
	levelUps := 0
	for p.Experience >= p.NextExperience {
		p.Experience -= p.NextExperience
		p.Level++
		p.NextExperience = ExperienceRequirement(p.Level)
		levelUps++
	}
	return levelUps
}

func (p *Progression) GainExperienceToLevel(targetLevel int) int {
	if targetLevel <= p.Level {
		return 0
	}
	amount := p.NextExperience - p.Experience
	for level := p.Level + 1; level < targetLevel; level++ {
		amount += ExperienceRequirement(level)
	}
	return p.GainExperience(amount)
}

func (p *Progression) ApplyLevelUpOption(option LevelUpOption) {
	switch option {
	case FireRate:
		if p.LevelUpOptionAvailable(option) {
			p.fireRateUpgrades++
		}
	case ExtraFireball:
		p.SimultaneousFireball++
	case LearnLightning:
		p.LightningUnlocked = true
	case LightningBounce:
		p.LightningBounceCount++
	case LightningRate:
		if p.LevelUpOptionAvailable(option) {
			p.lightningRateUpgrades++
		}
	case LearnOrb:
		p.OrbitalOrbUnlocked = true
	case ExtraOrb:
		p.upgradedOrbitalOrbCount++
	case OrbitalSpeed:
		p.orbitalSpeedUpgrades++
	case LearnBeam:
		p.BeamUnlocked = true
	case BeamRate:
		if p.LevelUpOptionAvailable(option) {
			p.beamRateUpgrades++
		}
	case BeamKillCount:
		p.beamKillLevel++
		p.upgradedBeamKillCount += p.beamKillLevel
	case LearnMeteor:
		p.MeteorUnlocked = true
	case ExtraMeteor:
		p.upgradedMeteorCount++
	case MeteorRate:
		if p.LevelUpOptionAvailable(option) {
			p.meteorRateUpgrades++
		}
	}
}

func (p *Progression) UpgradeAllProperties(skill LearnedSkill) []LevelUpOption {
	if !p.skillUnlocked(skill) {
		return nil
	}
	applied := []LevelUpOption{}
	for _, option := range skill.UpgradeOptions() {
		if !p.LevelUpOptionAvailable(option) {
			continue
		}
		p.ApplyLevelUpOption(option)
		applied = append(applied, option)
	}
	return applied
}

func (p Progression) skillUnlocked(skill LearnedSkill) bool {
	switch skill {
	case SkillFireball:
		return true
	case SkillLightning:
		return p.LightningUnlocked
	case SkillOrbitalOrb:
		return p.OrbitalOrbUnlocked
	case SkillBeam:
		return p.BeamUnlocked
	case SkillMeteor:
		return p.MeteorUnlocked
	default:
		return false
	}
}
