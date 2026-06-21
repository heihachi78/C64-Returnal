package game

import "math"

var baseLevelUpOptions = [...]LevelUpOption{
	FireRate,
	ExtraFireball,
	ExtraLife,
	HalveSkeletons,
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
	options := p.availableOptions(baseLevelUpOptions[:])
	options = p.appendWeaponOptions(options, p.LightningUnlocked, LearnLightning, LightningBounce, LightningRate)
	options = p.appendWeaponOptions(options, p.OrbitalOrbUnlocked, LearnOrb, ExtraOrb, OrbitalSpeed)
	options = p.appendWeaponOptions(options, p.BeamUnlocked, LearnBeam, BeamRate, BeamKillCount)
	options = p.appendWeaponOptions(options, p.MeteorUnlocked, LearnMeteor, ExtraMeteor, MeteorRate)
	return options
}

func (p Progression) appendWeaponOptions(options []LevelUpOption, unlocked bool, learn LevelUpOption, upgrades ...LevelUpOption) []LevelUpOption {
	if !unlocked {
		return append(options, learn)
	}
	return append(options, p.availableOptions(upgrades)...)
}

func (p Progression) AvailableUpgradeOptionsForSkill(skill LearnedSkill) []LevelUpOption {
	return p.availableOptions(skill.UpgradeOptions())
}

func (p Progression) availableOptions(candidates []LevelUpOption) []LevelUpOption {
	options := make([]LevelUpOption, 0, len(candidates))
	for _, option := range candidates {
		if p.LevelUpOptionAvailable(option) {
			options = append(options, option)
		}
	}
	return options
}

func (p Progression) LevelUpOptionAvailable(option LevelUpOption) bool {
	if option == BuyDeathWaveScroll {
		return p.DeathWaveScrolls < deathWaveRequiredScrolls
	}
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

func (p *Progression) ApplyLevelUpOption(option LevelUpOption) {
	switch option {
	case BuyDeathWaveScroll:
		p.applyDeathWaveScroll()
	case FireRate:
		p.applyRateUpgrade(option, &p.fireRateUpgrades)
	case ExtraFireball:
		p.SimultaneousFireball++
	case LearnLightning:
		p.LightningUnlocked = true
	case LightningBounce:
		p.LightningBounceCount++
	case LightningRate:
		p.applyRateUpgrade(option, &p.lightningRateUpgrades)
	case LearnOrb:
		p.OrbitalOrbUnlocked = true
	case ExtraOrb:
		p.upgradedOrbitalOrbCount++
	case OrbitalSpeed:
		p.orbitalSpeedUpgrades++
	case LearnBeam:
		p.BeamUnlocked = true
	case BeamRate:
		p.applyRateUpgrade(option, &p.beamRateUpgrades)
	case BeamKillCount:
		p.beamKillLevel++
		p.upgradedBeamKillCount += p.beamKillLevel
	case LearnMeteor:
		p.MeteorUnlocked = true
	case ExtraMeteor:
		p.upgradedMeteorCount++
	case MeteorRate:
		p.applyRateUpgrade(option, &p.meteorRateUpgrades)
	}
}

func (p *Progression) applyDeathWaveScroll() {
	if p.DeathWaveScrolls >= deathWaveRequiredScrolls {
		return
	}
	p.DeathWaveScrolls++
	p.DeathWaveUnlocked = p.DeathWaveScrolls >= deathWaveRequiredScrolls
}

func (p *Progression) applyRateUpgrade(option LevelUpOption, upgrades *int) {
	if p.LevelUpOptionAvailable(option) {
		(*upgrades)++
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
