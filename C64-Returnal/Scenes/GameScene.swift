//
//  GameScene.swift
//  C64-Returnal
//
//  Created by Tóth István on 2026. 06. 02..
//

import SpriteKit

final class GameScene: SKScene {
    let tuning: GameTuning
    let inputController: InputController
    let worldNode = SKNode()
    let cameraNode = SKCameraNode()
    let mageTextures = PixelArtFactory.makeMageTextures()
    lazy var player = SKSpriteNode(texture: mageTextures[0])
    let hud = GameHUD()
    let skeletonSpatialIndex: SkeletonSpatialIndex
    let grassField: InfiniteGrassField

    let skeletonTextures = PixelArtFactory.makeSkeletonTextures()
    let fireballTextures = PixelArtFactory.makeFireballTextures()
    let lightningTexture = PixelArtFactory.makeLightningTexture()
    let orbitalOrbTextures = PixelArtFactory.makeOrbitalOrbTextures()
    let beamTexture = PixelArtFactory.makeBeamTexture()
    let meteorTextures = PixelArtFactory.makeMeteorTextures()
    let lifeTexture = PixelArtFactory.makeLifeTexture()
    let coinTextures = PixelArtFactory.makeCoinTextures()
    let chestTextures: [ChestTier: SKTexture] = [
        .bronze: PixelArtFactory.makeChestTexture(tier: .bronze),
        .silver: PixelArtFactory.makeChestTexture(tier: .silver),
        .gold: PixelArtFactory.makeChestTexture(tier: .gold)
    ]

    var session: GameSessionState
    var skeletons = [SKSpriteNode]()
    var skeletonIdentifiers = Set<ObjectIdentifier>()
    var skeletonIndices = [ObjectIdentifier: Int]()
    var fireballs = [Fireball]()
    var meteors = [MeteorProjectile]()
    var chests = [Chest]()
    var coins = [Coin]()
    var orbitalOrbs = [OrbitalOrb]()

    init(
        size: CGSize,
        tuning: GameTuning = GameConfiguration.defaultTuning,
        inputBindings: InputBindings = GameConfiguration.defaultInputBindings
    ) {
        self.tuning = tuning
        inputController = InputController(bindings: inputBindings)
        skeletonSpatialIndex = SkeletonSpatialIndex(cellSize: tuning.skeleton.spatialIndexCellSize)
        grassField = InfiniteGrassField(
            tileSize: tuning.presentation.tileSize,
            textures: PixelArtFactory.makeGrassTextures(tileSize: tuning.presentation.tileSize)
        )
        session = GameSessionState(tuning: tuning)
        super.init(size: size)
    }

    required init?(coder aDecoder: NSCoder) {
        tuning = GameConfiguration.defaultTuning
        inputController = InputController(bindings: GameConfiguration.defaultInputBindings)
        skeletonSpatialIndex = SkeletonSpatialIndex(cellSize: GameConfiguration.defaultTuning.skeleton.spatialIndexCellSize)
        grassField = InfiniteGrassField(
            tileSize: GameConfiguration.defaultTuning.presentation.tileSize,
            textures: PixelArtFactory.makeGrassTextures(tileSize: GameConfiguration.defaultTuning.presentation.tileSize)
        )
        session = GameSessionState(tuning: GameConfiguration.defaultTuning)
        super.init(coder: aDecoder)
    }

    override func didMove(to view: SKView) {
        if !session.isSceneConfigured {
            configureScene()
        }

        view.ignoresSiblingOrder = true
        view.shouldCullNonVisibleNodes = true
        view.window?.makeFirstResponder(view)
        layoutViewportContent()
    }

    override func didChangeSize(_ oldSize: CGSize) {
        layoutViewportContent()
    }

    override func keyDown(with event: NSEvent) {
        guard !session.isGameOver else {
            return
        }

        if session.isChestRewardActive {
            advanceChestReward(with: event.keyCode)
            return
        }

        if session.isLevelUpChoiceActive {
            selectLevelUpOption(with: event.keyCode)
            return
        }

        if inputController.isKillAllAndGrantExperience(event.keyCode) {
            killAllEnemiesAndGrantExperience()
            return
        }

        guard inputController.isMovementKey(event.keyCode) else {
            super.keyDown(with: event)
            return
        }

        session.pressedKeys.insert(event.keyCode)
    }

    override func keyUp(with event: NSEvent) {
        guard !session.isGameOver else {
            return
        }

        guard !session.isChestRewardActive else {
            return
        }

        guard !session.isLevelUpChoiceActive else {
            return
        }

        guard inputController.isMovementKey(event.keyCode) else {
            super.keyUp(with: event)
            return
        }

        session.pressedKeys.remove(event.keyCode)
    }

    override func mouseDown(with event: NSEvent) {
        guard session.isGameOver || session.isLevelUpChoiceActive else {
            return
        }

        let cameraPoint = cameraNode.convert(event.location(in: self), from: self)

        if session.isGameOver {
            switch hud.option(at: cameraPoint) {
            case .restart:
                restartGame()
            case .exit:
                NSApp.terminate(nil)
            case .none:
                break
            }
        } else if hud.isLevelUpRedraw(at: cameraPoint) {
            redrawLevelUpOptions()
        } else if let option = hud.levelUpOption(at: cameraPoint) {
            applyLevelUpOption(option)
        }
    }

    override func update(_ currentTime: TimeInterval) {
        if session.lastUpdateTime == 0 {
            session.lastUpdateTime = currentTime
        }

        let deltaTime = currentTime - session.lastUpdateTime

        if !session.isGameOver && !session.isLevelUpChoiceActive && !session.isChestRewardActive {
            updatePlayer(deltaTime: deltaTime)
            checkCoinPickups()
            checkChestPickups()

            guard !session.isChestRewardActive else {
                session.lastUpdateTime = currentTime
                return
            }

            updateSkeletons(deltaTime: deltaTime)
            updateOrbitalOrbs(deltaTime: deltaTime)

            guard !session.isLevelUpChoiceActive else {
                session.lastUpdateTime = currentTime
                return
            }

            updateLightningCasting(deltaTime: deltaTime)

            guard !session.isLevelUpChoiceActive else {
                session.lastUpdateTime = currentTime
                return
            }

            updateFireballCasting(deltaTime: deltaTime)
            updateFireballs(deltaTime: deltaTime)

            guard !session.isLevelUpChoiceActive else {
                session.lastUpdateTime = currentTime
                return
            }

            updateBeamCasting(deltaTime: deltaTime)

            guard !session.isLevelUpChoiceActive else {
                session.lastUpdateTime = currentTime
                return
            }

            updateMeteorCasting(deltaTime: deltaTime)
            updateMeteors(deltaTime: deltaTime)
            updatePlayerHitInvulnerability(deltaTime: deltaTime)

            guard !session.isLevelUpChoiceActive else {
                session.lastUpdateTime = currentTime
                return
            }

            checkSkeletonCollisions()
            updateSkeletonSpawning(deltaTime: deltaTime)
        }

        session.lastUpdateTime = currentTime
    }

    func configureScene() {
        session.isSceneConfigured = true
        backgroundColor = tuning.presentation.backgroundColor
        anchorPoint = CGPoint(x: 0.5, y: 0.5)

        addChild(worldNode)
        worldNode.addChild(grassField.node)
        configurePlayer()

        camera = cameraNode
        addChild(cameraNode)
        hud.add(
            to: cameraNode,
            fireballTexture: fireballTextures[0],
            lightningTexture: lightningTexture,
            orbTexture: orbitalOrbTextures[0],
            beamTexture: beamTexture,
            meteorTexture: meteorTextures[0],
            lifeTexture: lifeTexture,
            coinTexture: coinTextures[0],
            skeletonTexture: skeletonTextures[0]
        )
        syncOrbitalOrbCount()
        updateHUDProgress()

        spawnSkeleton()
        spawnCoin(for: session.progression.level)
    }

    static let skeletonDamageFlashActionKey = "skeletonDamageFlash"
    static let playerHitFlashActionKey = "playerHitFlash"
    static let playerAnimationActionKey = "playerAnimation"
}
