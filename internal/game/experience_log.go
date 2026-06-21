package game

import (
	"fmt"
	"os"
)

func (g *Game) logExperienceAward(attack AttackKind, skeleton Skeleton, deaths, xp, beforeLevel, beforeXP, beforeNextXP, levelUps int) {
	if g.experienceLogPath == "" || deaths <= 0 || xp <= 0 {
		return
	}
	file, err := os.OpenFile(g.experienceLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	after := g.session.Progression
	fmt.Fprintf(
		file,
		"t=%.3f attack=%s deaths=%d xp=%d enemy=%s enemy_id=%d level=%d->%d xp_bar=%d/%d->%d/%d level_ups=%d total_kills=%d\n",
		g.totalTime,
		attack.LogName(),
		deaths,
		xp,
		skeleton.Kind.LogName(),
		skeleton.ID,
		beforeLevel,
		after.Level,
		beforeXP,
		beforeNextXP,
		after.Experience,
		after.NextExperience,
		levelUps,
		g.session.Kills.TotalSkeletons,
	)
}

func (k AttackKind) LogName() string {
	switch k {
	case AttackFireball:
		return "fireball"
	case AttackLightning:
		return "lightning"
	case AttackOrbitalOrb:
		return "orbital_orb"
	case AttackBeam:
		return "beam"
	case AttackMeteor:
		return "meteor"
	case AttackNone:
		return "none"
	default:
		return "unknown"
	}
}

func (k SkeletonKind) LogName() string {
	switch k {
	case SkeletonRed:
		return "red"
	case SkeletonPurple:
		return "purple"
	case SkeletonBlack:
		return "black"
	case SkeletonBlue:
		return "blue"
	default:
		return "white"
	}
}
