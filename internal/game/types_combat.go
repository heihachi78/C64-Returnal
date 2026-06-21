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

type actualDamageSample struct {
	Time   float64
	Amount int
}

type SkeletonKind int

const (
	SkeletonRegular SkeletonKind = iota
	SkeletonRed
	SkeletonPurple
	SkeletonBlack
	SkeletonBlue
)

const blueMonsterHitPoints = 1000

func (k SkeletonKind) HitPoints(t Tuning) int {
	switch k {
	case SkeletonRed:
		return max(1, t.RedHitPoints)
	case SkeletonPurple:
		return max(1, t.PurpleHitPoints)
	case SkeletonBlack:
		return max(1, t.BlackHitPoints)
	case SkeletonBlue:
		return blueMonsterHitPoints
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
	case SkeletonBlue:
		return 75
	default:
		return 1
	}
}

type EffectKind int

const (
	EffectLightning EffectKind = iota
	EffectLightningHit
	EffectBeam
	EffectMeteorImpact
)
