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
	orbs                        []OrbitalOrb
	meteors                     []MeteorProjectile
	chests                      []Chest
	coins                       []Coin
	effects                     []Effect
	actualDamage                []actualDamageSample

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
	g.orbs = g.orbs[:0]
	g.meteors = g.meteors[:0]
	g.chests = g.chests[:0]
	g.coins = g.coins[:0]
	g.effects = g.effects[:0]
	g.actualDamage = g.actualDamage[:0]
	g.spatial.Rebuild(g.skeleton)
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
func (g *Game) Update() error {
	dt := 0.0
	if g.hasUpdated {
		dt = 1.0 / float64(TargetTPS)
	}
	g.hasUpdated = true
	g.totalTime += dt
	clear(g.lightningTargetReservations)
	consumedFrame, err := g.updateOverlayInput()
	if err != nil {
		return err
	}
	if consumedFrame {
		return nil
	}
	g.updatePausedAnimations(dt)

	if g.session.GameOver || g.session.LevelUpChoiceActive || g.session.ChestRewardActive {
		return nil
	}

	g.updatePlayer(dt)
	g.checkCoinPickups()
	g.checkChestPickups()
	if g.session.ChestRewardActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateSkeletons(dt)
	g.updateOrbitalOrbs(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateLightningCasting(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateFireballCasting(dt)
	g.updateFireballs(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateBeamCasting(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.updateMeteorCasting(dt)
	g.updateMeteors(dt)
	g.updateInvulnerability(dt)
	if g.session.LevelUpChoiceActive {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}

	g.checkSkeletonCollisions()
	g.updateSkeletonSpawning(dt)
	if g.session.GameOver {
		g.updateNewlyPresentedOverlayActions(dt)
		return nil
	}
	g.updatePlayerWalkAnimation(dt)
	g.updatePlayerHitFlash(dt)
	g.updateSkeletonHitFlashes(dt)
	g.updateCoins(dt)
	g.updateEffects(dt)
	return nil
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenW = max(1, outsideWidth)
	g.screenH = max(1, outsideHeight)
	return g.screenW, g.screenH
}
