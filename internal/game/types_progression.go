package game

type LevelUpOption int

const (
	FireRate LevelUpOption = iota
	ExtraFireball
	ExtraLife
	HalveSkeletons
	LearnLightning
	LightningBounce
	LightningRate
	LearnOrb
	ExtraOrb
	OrbitalSpeed
	LearnBeam
	BeamRate
	BeamKillCount
	LearnMeteor
	ExtraMeteor
	MeteorRate
	BuyDeathWaveScroll
)

func (o LevelUpOption) Title(beamKillBonus int) string {
	switch o {
	case FireRate:
		return "FASTER FIRE"
	case ExtraFireball:
		return "+1 FIREBALL"
	case ExtraLife:
		return "+1 LIFE"
	case HalveSkeletons:
		return "HALVE HORDE"
	case LearnLightning:
		return "LEARN BOLT"
	case LightningBounce:
		return "+1 CHAIN"
	case LightningRate:
		return "FASTER BOLT"
	case LearnOrb:
		return "LEARN ORB"
	case ExtraOrb:
		return "+1 ORB"
	case OrbitalSpeed:
		return "FASTER ORB"
	case LearnBeam:
		return "LEARN BEAM"
	case BeamRate:
		return "FASTER BEAM"
	case BeamKillCount:
		return "+" + itoa(max(1, beamKillBonus)) + " BEAM KILL"
	case LearnMeteor:
		return "LEARN METEOR"
	case ExtraMeteor:
		return "+1 METEOR"
	case MeteorRate:
		return "FASTER METEOR"
	case BuyDeathWaveScroll:
		return "DEATH SCROLL"
	default:
		return "UNKNOWN"
	}
}

type LearnedSkill int

const (
	SkillFireball LearnedSkill = iota
	SkillLightning
	SkillOrbitalOrb
	SkillBeam
	SkillMeteor
)

func (s LearnedSkill) UpgradeOptions() []LevelUpOption {
	switch s {
	case SkillLightning:
		return []LevelUpOption{LightningBounce, LightningRate}
	case SkillOrbitalOrb:
		return []LevelUpOption{ExtraOrb, OrbitalSpeed}
	case SkillBeam:
		return []LevelUpOption{BeamRate, BeamKillCount}
	case SkillMeteor:
		return []LevelUpOption{ExtraMeteor, MeteorRate}
	default:
		return []LevelUpOption{FireRate, ExtraFireball}
	}
}

type ChestTier int

const (
	ChestBronze ChestTier = iota
	ChestSilver
	ChestGold
)

func (t ChestTier) Title() string {
	switch t {
	case ChestSilver:
		return "SILVER"
	case ChestGold:
		return "GOLD"
	default:
		return "BRONZE"
	}
}

type ChestRewardDisplayItem struct {
	Option LevelUpOption
	Title  string
}
