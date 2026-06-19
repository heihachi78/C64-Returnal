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
	return interval
}

func (p Progression) FireballCastInterval() float64 {
	return p.tuning.InitialFireballCast * math.Pow(p.tuning.FireballIntervalMultiplier, float64(p.fireRateUpgrades))
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
	options := []LevelUpOption{FireRate, ExtraFireball, ExtraLife, HalveSkeletons}
	if p.LightningUnlocked {
		options = append(options, LightningBounce, LightningRate)
	} else {
		options = append(options, LearnLightning)
	}
	if p.OrbitalOrbUnlocked {
		options = append(options, ExtraOrb, OrbitalSpeed)
	} else {
		options = append(options, LearnOrb)
	}
	if p.BeamUnlocked {
		options = append(options, BeamRate, BeamKillCount)
	} else {
		options = append(options, LearnBeam)
	}
	if p.MeteorUnlocked {
		options = append(options, ExtraMeteor, MeteorRate)
	} else {
		options = append(options, LearnMeteor)
	}
	return options
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

func (p *Progression) ApplyLevelUpOption(option LevelUpOption) {
	switch option {
	case FireRate:
		p.fireRateUpgrades++
	case ExtraFireball:
		p.SimultaneousFireball++
	case LearnLightning:
		p.LightningUnlocked = true
	case LightningBounce:
		p.LightningBounceCount++
	case LightningRate:
		p.lightningRateUpgrades++
	case LearnOrb:
		p.OrbitalOrbUnlocked = true
	case ExtraOrb:
		p.upgradedOrbitalOrbCount++
	case OrbitalSpeed:
		p.orbitalSpeedUpgrades++
	case LearnBeam:
		p.BeamUnlocked = true
	case BeamRate:
		p.beamRateUpgrades++
	case BeamKillCount:
		p.beamKillLevel++
		p.upgradedBeamKillCount += p.beamKillLevel
	case LearnMeteor:
		p.MeteorUnlocked = true
	case ExtraMeteor:
		p.upgradedMeteorCount++
	case MeteorRate:
		p.meteorRateUpgrades++
	}
}

func (p *Progression) UpgradeAllProperties(skill LearnedSkill) {
	switch skill {
	case SkillFireball:
		p.fireRateUpgrades++
		p.SimultaneousFireball++
	case SkillLightning:
		if p.LightningUnlocked {
			p.LightningBounceCount++
			p.lightningRateUpgrades++
		}
	case SkillOrbitalOrb:
		if p.OrbitalOrbUnlocked {
			p.upgradedOrbitalOrbCount++
			p.orbitalSpeedUpgrades++
		}
	case SkillBeam:
		if p.BeamUnlocked {
			p.beamRateUpgrades++
			p.beamKillLevel++
			p.upgradedBeamKillCount += p.beamKillLevel
		}
	case SkillMeteor:
		if p.MeteorUnlocked {
			p.upgradedMeteorCount++
			p.meteorRateUpgrades++
		}
	}
}
