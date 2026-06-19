package game

type AttackKind int

const (
	AttackFireball AttackKind = iota
	AttackLightning
	AttackOrbitalOrb
	AttackBeam
	AttackMeteor
	AttackNone
)

type SkeletonKind int

const (
	SkeletonRegular SkeletonKind = iota
	SkeletonRed
	SkeletonPurple
	SkeletonBlack
)

func (k SkeletonKind) HitPoints(t Tuning) int {
	switch k {
	case SkeletonRed:
		return max(1, t.RedHitPoints)
	case SkeletonPurple:
		return max(1, t.PurpleHitPoints)
	case SkeletonBlack:
		return max(1, t.BlackHitPoints)
	default:
		return 1
	}
}

func (k SkeletonKind) ExperienceReward() int {
	switch k {
	case SkeletonPurple:
		return 3
	case SkeletonBlack:
		return 10
	default:
		return 1
	}
}

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

type Player struct {
	Pos           Vec2
	Facing        float64
	Moving        bool
	MoveDir       Vec2
	AnimTimer     float64
	AnimFrame     int
	HitFlash      float64
	DeathTimer    float64
	DeathRotation float64
}

type Skeleton struct {
	ID        int
	Pos       Vec2
	Kind      SkeletonKind
	HP        int
	Reward    int
	Facing    float64
	HitFlash  float64
	AnimFrame int
}

type Fireball struct {
	Pos               Vec2
	TargetID          int
	Velocity          Vec2
	TimeWithoutTarget float64
	AnimFrame         int
}

type OrbitalOrb struct {
	Pos                  Vec2
	Active               bool
	MissingOrbitProgress float64
	AnimFrame            int
}

type MeteorProjectile struct {
	Pos       Vec2
	Start     Vec2
	Impact    Vec2
	Age       float64
	AnimFrame int
}

type Coin struct {
	Pos    Vec2
	Amount int
	Level  int
	Phase  float64
}

type Chest struct {
	Pos  Vec2
	Tier ChestTier
}

type ChestRewardDisplayItem struct {
	Option LevelUpOption
	Title  string
}

type EffectKind int

const (
	EffectLightning EffectKind = iota
	EffectLightningHit
	EffectBeam
	EffectMeteorImpact
)

type Effect struct {
	Kind        EffectKind
	Start       Vec2
	End         Vec2
	Pos         Vec2
	Points      []Vec2
	InnerPoints []Vec2
	Frame       int
	Facing      float64
	Radius      float64
	TTL         float64
	MaxTTL      float64
}
