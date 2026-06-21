package game

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
	DeathWaveScrolls        int
	DeathWaveUnlocked       bool
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
	p.DeathWaveScrolls = 0
	p.DeathWaveUnlocked = false
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
