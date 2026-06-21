package game

import "math"

func (p Progression) FireballCastInterval() float64 {
	return p.tuning.InitialFireballCast * math.Pow(p.tuning.FireballIntervalMultiplier, float64(p.fireRateUpgrades))
}

func orbitalOrbHitInterval(angularSpeed float64) float64 {
	if angularSpeed <= 0 {
		return 0
	}
	return math.Pi * 2 / angularSpeed
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
