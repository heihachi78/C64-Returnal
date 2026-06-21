package game

const (
	ScreenWidth  = 1200
	ScreenHeight = 900
	TargetTPS    = 120

	grassTintBlendFactor         = 0.22
	modalFadeDuration            = 0.20
	panelCornerRadius            = 6.0
	redrawStatusFadeDuration     = 0.14
	redrawFailurePulseDuration   = 0.16
	playerHitFlashDuration       = 0.96
	playerDeathRotationDuration  = 0.16
	skeletonDamageFlashDuration  = 0.24
	fireballImpactEffectDuration = 0.16
	meteorImpactEffectDuration   = 0.32
	debugLevelJumpTarget         = 100
	debugLevelJumpCoins          = 5000
	debugLevelJumpGoldChests     = 40
	debugLevelJumpChestRadius    = 56
	deathWaveScrollCost          = 1000
	deathWaveRequiredScrolls     = 5
	dynamicSpawnBatchInterval    = 1.0
	minAttackSpawnInterval       = 1.0 / 60.0
)

type Tuning struct {
	TileSize float64

	PlayerSpeed                     float64
	InitialPlayerLives              int
	PlayerHitInvulnerability        float64
	PlayerAnimationFrameTime        float64
	SkeletonSpeed                   float64
	InitialSkeletonHPPerSecond      float64
	SkeletonHPPerSecondLevelUpBonus float64
	SkeletonHitDistance             float64
	SkeletonSpawnMargin             float64
	SkeletonAnimationFrameTime      float64
	RedHitPoints                    int
	PurpleHitPoints                 int
	BlackHitPoints                  int
	MaxActiveSkeletons              int
	MaxSkeletonSpawnsPerTick        int
	SpatialIndexCellSize            float64

	FireballSpeed              float64
	InitialFireballCast        float64
	FireballIntervalMultiplier float64
	FireballHitDistance        float64
	FireballUntargetedLifetime float64
	FireballAnimationFrameTime float64

	InitialLightningCast        float64
	LightningIntervalMultiplier float64
	LightningEffectDuration     float64
	LightningBranchWidth        float32

	OrbitalOrbRadius             float64
	InitialOrbitalAngularSpeed   float64
	OrbitalSpeedUpgradeMultipler float64
	OrbitalHitDistance           float64
	OrbitalAnimationFrameTime    float64

	InitialBeamCast        float64
	BeamIntervalMultiplier float64
	BeamHitWidth           float64
	BeamEffectDuration     float64

	InitialMeteorCast        float64
	MeteorIntervalMultiplier float64
	MeteorTargetMultiplier   float64
	MeteorImpactRadius       float64
	MeteorFallDuration       float64
	MeteorFallHeight         float64
	MeteorFallDrift          float64
	MeteorAnimationFrameTime float64

	DeathWaveInterval  float64
	DeathWaveSpeed     float64
	DeathWaveWidth     float64
	DeathWaveMaxRadius float64

	ChestSpawnMargin       float64
	ChestPickupDistance    float64
	BronzeKillInterval     int
	SilverKillInterval     int
	GoldKillInterval       int
	BronzeMaximumLevel     int
	SilverMaximumLevel     int
	CoinSpawnMargin        float64
	CoinPickupDistance     float64
	CoinMinimumReward      int
	CoinMaximumReward      int
	CoinAnimationFrameTime float64

	HalveHordeChanceNumerator       int
	HalveHordeChanceDenominator     int
	ExtraOptionChanceNumerator      int
	ExtraOptionChanceDenominator    int
	ParallelSkeletonUpdateThreshold int
}

func DefaultTuning() Tuning {
	return Tuning{
		TileSize: 64,

		PlayerSpeed:                     190,
		InitialPlayerLives:              3,
		PlayerHitInvulnerability:        1,
		PlayerAnimationFrameTime:        0.18,
		SkeletonSpeed:                   82,
		InitialSkeletonHPPerSecond:      1.0,
		SkeletonHPPerSecondLevelUpBonus: 0.375,
		SkeletonHitDistance:             24,
		SkeletonSpawnMargin:             72,
		SkeletonAnimationFrameTime:      0.20,
		RedHitPoints:                    3,
		PurpleHitPoints:                 7,
		BlackHitPoints:                  29,
		MaxActiveSkeletons:              999,
		MaxSkeletonSpawnsPerTick:        1,
		SpatialIndexCellSize:            96,

		FireballSpeed:              280,
		InitialFireballCast:        3,
		FireballIntervalMultiplier: 0.9,
		FireballHitDistance:        20,
		FireballUntargetedLifetime: 3,
		FireballAnimationFrameTime: 0.08,

		InitialLightningCast:        3,
		LightningIntervalMultiplier: 0.9,
		LightningEffectDuration:     0.18,
		LightningBranchWidth:        3,

		OrbitalOrbRadius:             58,
		InitialOrbitalAngularSpeed:   2.4,
		OrbitalSpeedUpgradeMultipler: 1.2,
		OrbitalHitDistance:           22,
		OrbitalAnimationFrameTime:    0.12,

		InitialBeamCast:        3,
		BeamIntervalMultiplier: 0.9,
		BeamHitWidth:           24,
		BeamEffectDuration:     0.16,

		InitialMeteorCast:        3,
		MeteorIntervalMultiplier: 0.95,
		MeteorTargetMultiplier:   8,
		MeteorImpactRadius:       48,
		MeteorFallDuration:       0.55,
		MeteorFallHeight:         240,
		MeteorFallDrift:          90,
		MeteorAnimationFrameTime: 0.14,

		DeathWaveInterval:  30,
		DeathWaveSpeed:     320,
		DeathWaveWidth:     28,
		DeathWaveMaxRadius: 760,

		ChestSpawnMargin:       88,
		ChestPickupDistance:    28,
		BronzeKillInterval:     250,
		SilverKillInterval:     1000,
		GoldKillInterval:       2500,
		BronzeMaximumLevel:     33,
		SilverMaximumLevel:     55,
		CoinSpawnMargin:        240,
		CoinPickupDistance:     30,
		CoinMinimumReward:      1,
		CoinMaximumReward:      100,
		CoinAnimationFrameTime: 0.14,

		HalveHordeChanceNumerator:       5,
		HalveHordeChanceDenominator:     100,
		ExtraOptionChanceNumerator:      25,
		ExtraOptionChanceDenominator:    100,
		ParallelSkeletonUpdateThreshold: 1024,
	}
}
