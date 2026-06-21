package game

type KillCounts struct {
	TotalSkeletons int
	Fireball       int
	Lightning      int
	OrbitalOrb     int
	Beam           int
	Meteor         int
}

type CastTimers struct {
	SkeletonSpawn float64
	Fireball      float64
	Lightning     float64
	Beam          float64
	Meteor        float64
	DeathWave     float64
}

type Session struct {
	Progression                Progression
	Casts                      CastTimers
	PlayerHitInvulnerability   float64
	PlayerLives                int
	OrbitalOrbAngle            float64
	PendingLevelUpLevels       []int
	Kills                      KillCounts
	NextChestMilestone         int
	CollectedCoins             int
	SpawnedCoinLevels          map[int]bool
	GameOver                   bool
	LevelUpChoiceActive        bool
	ChestRewardActive          bool
	ActiveLevelUpOptions       []LevelUpOption
	ActiveChestTier            ChestTier
	ActiveChestRewardItems     []ChestRewardDisplayItem
	LevelUpRedrawStatusTimer   float64
	LevelUpRedrawFadeTimer     float64
	LevelUpRedrawCoinFadeTimer float64
	LevelUpOverlayTimer        float64
	LevelUpTitleScaleTimer     float64
	LevelUpOptionFadeTimer     float64
	ChestRewardOverlayTimer    float64
	GameOverOverlayTimer       float64
	CurrentLevelUpPresentation int
}

func NewSession(t Tuning) Session {
	return Session{
		Progression:        NewProgression(t),
		PlayerLives:        t.InitialPlayerLives,
		NextChestMilestone: t.BronzeKillInterval,
		SpawnedCoinLevels:  map[int]bool{},
	}
}

func (s *Session) Reset(t Tuning) {
	*s = NewSession(t)
}

func (s *Session) RegisterAttackKill(kind AttackKind) {
	switch kind {
	case AttackFireball:
		s.Kills.Fireball++
	case AttackLightning:
		s.Kills.Lightning++
	case AttackOrbitalOrb:
		s.Kills.OrbitalOrb++
	case AttackBeam:
		s.Kills.Beam++
	case AttackMeteor:
		s.Kills.Meteor++
	}
}
