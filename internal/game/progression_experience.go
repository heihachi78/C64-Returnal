package game

func ExperienceRequirement(level int) int {
	return max(1, level*level/2)
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
