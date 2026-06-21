package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

type Game struct {
	tuning Tuning
	rng    *rand.Rand
	assets *Assets

	screenW int
	screenH int
	nextID  int

	player                      Player
	session                     Session
	spatial                     SpatialIndex
	skeleton                    []Skeleton
	fireball                    []Fireball
	lightningTargetReservations map[int]bool
	fireballTargetReservations  map[int]bool
	orbs                        []OrbitalOrb
	meteors                     []MeteorProjectile
	deathWaves                  []DeathWave
	chests                      []Chest
	coins                       []Coin
	effects                     []Effect
	actualDamage                []actualDamageSample
	actualDamageWindowTotal     int
	maxActualDPS                float64
	pendingSpawnPressureActual  float64
	pendingSpawnPressureLevels  int
	skeletonHPPerSecond         float64
	skeletonSpatialDirty        bool
	dynamicSpawnQueue           []dynamicSpawnPlanEntry
	closestSkeletonScratch      []closestSkeletonPick
	closestSkeletonResult       []int
	beamTargetScratch           []beamTargetHit
	beamTargetResult            []int
	lightningRemainingTargets   []int
	lightningTargetScratch      []lightningStrikeTarget

	skeletonAnimTimer  float64
	skeletonAnimFrame  int
	fireAnimTimer      float64
	fireAnimFrame      int
	orbAnimTimer       float64
	orbAnimFrame       int
	meteorAnimTimer    float64
	meteorAnimFrame    int
	totalTime          float64
	hasUpdated         bool
	lastParallelJobs   int
	suppressedMovement map[ebiten.Key]bool
	scaledTextCache    map[scaledTextCacheKey]scaledTextCacheEntry
}

func New() *Game {
	tuning := DefaultTuning()
	g := &Game{
		tuning:                      tuning,
		rng:                         rand.New(rand.NewSource(rand.Int63())),
		assets:                      NewAssets(int(tuning.TileSize)),
		screenW:                     ScreenWidth,
		screenH:                     ScreenHeight,
		spatial:                     NewSpatialIndex(tuning.SpatialIndexCellSize),
		lightningTargetReservations: map[int]bool{},
		fireballTargetReservations:  map[int]bool{},
		suppressedMovement:          map[ebiten.Key]bool{},
		scaledTextCache:             map[scaledTextCacheKey]scaledTextCacheEntry{},
	}
	g.reset()
	return g
}
func (g *Game) reset() {
	g.nextID = 1
	g.player = Player{Facing: 1}
	g.session = NewSession(g.tuning)
	g.skeleton = g.skeleton[:0]
	g.fireball = g.fireball[:0]
	clear(g.lightningTargetReservations)
	clear(g.fireballTargetReservations)
	g.orbs = g.orbs[:0]
	g.meteors = g.meteors[:0]
	g.deathWaves = g.deathWaves[:0]
	g.chests = g.chests[:0]
	g.coins = g.coins[:0]
	g.effects = g.effects[:0]
	g.actualDamage = g.actualDamage[:0]
	g.actualDamageWindowTotal = 0
	g.maxActualDPS = 0
	g.pendingSpawnPressureActual = 0
	g.pendingSpawnPressureLevels = 0
	g.skeletonHPPerSecond = initialSkeletonHPPerSecond(g.tuning, g.session.Progression)
	g.dynamicSpawnQueue = g.dynamicSpawnQueue[:0]
	g.spatial.Rebuild(g.skeleton)
	g.skeletonSpatialDirty = false
	g.skeletonAnimTimer = 0
	g.skeletonAnimFrame = 0
	g.fireAnimTimer = 0
	g.fireAnimFrame = 0
	g.orbAnimTimer = 0
	g.orbAnimFrame = 0
	g.meteorAnimTimer = 0
	g.meteorAnimFrame = 0
	g.totalTime = 0
	g.hasUpdated = false
	g.lastParallelJobs = 0
	clear(g.suppressedMovement)
	g.spawnSkeleton(SkeletonRegular)
	g.spawnCoinForLevel(g.session.Progression.Level)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenW = max(1, outsideWidth)
	g.screenH = max(1, outsideHeight)
	return g.screenW, g.screenH
}
