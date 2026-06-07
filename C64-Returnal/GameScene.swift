//
//  GameScene.swift
//  C64-Returnal
//
//  Created by Tóth István on 2026. 06. 02..
//

import SpriteKit

private enum AttackKind {
    case fireball
    case lightning
    case orbitalOrb
    case beam
    case meteor
}

private enum SkeletonKind {
    case regular
    case red
    case purple

    var hitPoints: Int {
        switch self {
        case .regular:
            return 1
        case .red:
            return max(1, GameConfiguration.redSkeletonHitPoints)
        case .purple:
            return max(1, GameConfiguration.purpleSkeletonHitPoints)
        }
    }

    var tint: SKColor {
        switch self {
        case .regular:
            return .white
        case .red:
            return SKColor(calibratedRed: 0.95, green: 0.05, blue: 0.04, alpha: 1)
        case .purple:
            return SKColor(calibratedRed: 0.58, green: 0.12, blue: 0.95, alpha: 1)
        }
    }

    var tintBlendFactor: CGFloat {
        switch self {
        case .regular:
            return 0
        case .red:
            return 0.68
        case .purple:
            return 0.72
        }
    }
}

private enum SkeletonUserDataKey {
    static let hitPoints = "skeletonHitPoints"
}

final class GameScene: SKScene {
    private let worldNode = SKNode()
    private let cameraNode = SKCameraNode()
    private let mageTextures = PixelArtFactory.makeMageTextures()
    private lazy var player = SKSpriteNode(texture: mageTextures[0])
    private let hud = GameHUD()
    private let skeletonSpatialIndex = SkeletonSpatialIndex(cellSize: 96)
    private let grassField = InfiniteGrassField(
        tileSize: GameConfiguration.tileSize,
        textures: PixelArtFactory.makeGrassTextures(tileSize: GameConfiguration.tileSize)
    )

    private let skeletonTextures = PixelArtFactory.makeSkeletonTextures()
    private let fireballTextures = PixelArtFactory.makeFireballTextures()
    private let lightningTexture = PixelArtFactory.makeLightningTexture()
    private let orbitalOrbTextures = PixelArtFactory.makeOrbitalOrbTextures()
    private let beamTexture = PixelArtFactory.makeBeamTexture()
    private let meteorTextures = PixelArtFactory.makeMeteorTextures()
    private let lifeTexture = PixelArtFactory.makeLifeTexture()
    private let chestTextures: [ChestTier: SKTexture] = [
        .bronze: PixelArtFactory.makeChestTexture(tier: .bronze),
        .silver: PixelArtFactory.makeChestTexture(tier: .silver),
        .gold: PixelArtFactory.makeChestTexture(tier: .gold)
    ]

    private var progression = Progression()
    private var skeletons = [SKSpriteNode]()
    private var skeletonIdentifiers = Set<ObjectIdentifier>()
    private var skeletonIndices = [ObjectIdentifier: Int]()
    private var fireballs = [Fireball]()
    private var meteors = [MeteorProjectile]()
    private var chests = [Chest]()
    private var orbitalOrbs = [OrbitalOrb]()
    private var pressedKeys = Set<UInt16>()
    private var lastUpdateTime: TimeInterval = 0
    private var skeletonAnimationTimer: TimeInterval = 0
    private var skeletonAnimationFrameIndex = 0
    private var fireballAnimationTimer: TimeInterval = 0
    private var fireballAnimationFrameIndex = 0
    private var orbitalOrbAnimationTimer: TimeInterval = 0
    private var orbitalOrbAnimationFrameIndex = 0
    private var meteorAnimationTimer: TimeInterval = 0
    private var meteorAnimationFrameIndex = 0
    private var skeletonSpawnTimer: TimeInterval = 0
    private var fireballCastTimer: TimeInterval = 0
    private var lightningCastTimer: TimeInterval = 0
    private var beamCastTimer: TimeInterval = 0
    private var meteorCastTimer: TimeInterval = 0
    private var playerHitInvulnerabilityTimer: TimeInterval = 0
    private var playerLives = GameConfiguration.initialPlayerLives
    private var currentPlayerMovementDirection: CGVector?
    private var orbitalOrbAngle: CGFloat = 0
    private var pendingLevelUpLevels = [Int]()
    private var totalSkeletonKills = 0
    private var fireballKillCount = 0
    private var lightningKillCount = 0
    private var orbitalOrbKillCount = 0
    private var beamKillCount = 0
    private var meteorKillCount = 0
    private var nextChestMilestone = GameConfiguration.bronzeChestKillInterval
    private var isGameOver = false
    private var isLevelUpChoiceActive = false
    private var isChestRewardActive = false
    private var isSceneConfigured = false

    override func didMove(to view: SKView) {
        if !isSceneConfigured {
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
        guard !isGameOver else {
            return
        }

        if isChestRewardActive {
            advanceChestReward(with: event.keyCode)
            return
        }

        if isLevelUpChoiceActive {
            selectLevelUpOption(with: event.keyCode)
            return
        }

        if ExperienceDebugKey.isKillAllAndGrantExperience(event.keyCode) {
            killAllEnemiesAndGrantExperience()
            return
        }

        guard ArrowKey.contains(event.keyCode) else {
            super.keyDown(with: event)
            return
        }

        pressedKeys.insert(event.keyCode)
    }

    override func keyUp(with event: NSEvent) {
        guard !isGameOver else {
            return
        }

        guard !isChestRewardActive else {
            return
        }

        guard !isLevelUpChoiceActive else {
            return
        }

        guard ArrowKey.contains(event.keyCode) else {
            super.keyUp(with: event)
            return
        }

        pressedKeys.remove(event.keyCode)
    }

    override func mouseDown(with event: NSEvent) {
        guard isGameOver || isLevelUpChoiceActive else {
            return
        }

        let cameraPoint = cameraNode.convert(event.location(in: self), from: self)

        if isGameOver {
            switch hud.option(at: cameraPoint) {
            case .restart:
                restartGame()
            case .exit:
                NSApp.terminate(nil)
            case .none:
                break
            }
        } else if let option = hud.levelUpOption(at: cameraPoint) {
            applyLevelUpOption(option)
        }
    }

    override func update(_ currentTime: TimeInterval) {
        if lastUpdateTime == 0 {
            lastUpdateTime = currentTime
        }

        let deltaTime = currentTime - lastUpdateTime

        if !isGameOver && !isLevelUpChoiceActive && !isChestRewardActive {
            updatePlayer(deltaTime: deltaTime)
            checkChestPickups()

            guard !isChestRewardActive else {
                lastUpdateTime = currentTime
                return
            }

            updateSkeletons(deltaTime: deltaTime)
            updateOrbitalOrbs(deltaTime: deltaTime)

            guard !isLevelUpChoiceActive else {
                lastUpdateTime = currentTime
                return
            }

            updateLightningCasting(deltaTime: deltaTime)

            guard !isLevelUpChoiceActive else {
                lastUpdateTime = currentTime
                return
            }

            updateFireballCasting(deltaTime: deltaTime)
            updateFireballs(deltaTime: deltaTime)

            guard !isLevelUpChoiceActive else {
                lastUpdateTime = currentTime
                return
            }

            updateBeamCasting(deltaTime: deltaTime)

            guard !isLevelUpChoiceActive else {
                lastUpdateTime = currentTime
                return
            }

            updateMeteorCasting(deltaTime: deltaTime)
            updateMeteors(deltaTime: deltaTime)
            updatePlayerHitInvulnerability(deltaTime: deltaTime)

            guard !isLevelUpChoiceActive else {
                lastUpdateTime = currentTime
                return
            }

            checkSkeletonCollisions()
            updateSkeletonSpawning(deltaTime: deltaTime)
        }

        lastUpdateTime = currentTime
    }

    private func configureScene() {
        isSceneConfigured = true
        backgroundColor = GameConfiguration.backgroundColor
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
            skeletonTexture: skeletonTextures[0]
        )
        syncOrbitalOrbCount()
        updateHUDProgress()

        spawnSkeleton()
    }

    private func configurePlayer() {
        player.size = CGSize(width: 32, height: 44)
        player.zPosition = 10
        stopPlayerAnimation()
        worldNode.addChild(player)
    }

    private func layoutViewportContent() {
        hud.layout(for: size)
        grassField.rebuild(for: size)
        grassField.update(around: player.position)
    }

    private func updatePlayer(deltaTime: TimeInterval) {
        let horizontal = directionValue(negative: ArrowKey.left, positive: ArrowKey.right)
        let vertical = directionValue(negative: ArrowKey.down, positive: ArrowKey.up)
        var movement = CGVector(dx: horizontal, dy: vertical)

        if movement.dx != 0 || movement.dy != 0 {
            movement = movement.normalized
            currentPlayerMovementDirection = movement
            player.position.x += movement.dx * GameConfiguration.playerSpeed * CGFloat(deltaTime)
            player.position.y += movement.dy * GameConfiguration.playerSpeed * CGFloat(deltaTime)
            updateFacing(for: player, movement: movement)
            startPlayerAnimation()
        } else {
            currentPlayerMovementDirection = nil
            stopPlayerAnimation()
        }

        cameraNode.position = player.position
        grassField.update(around: player.position)
    }

    private func updateSkeletons(deltaTime: TimeInterval) {
        for skeleton in skeletons {
            let dx = player.position.x - skeleton.position.x
            let dy = player.position.y - skeleton.position.y
            let distanceSquared = dx * dx + dy * dy

            guard distanceSquared > 0 else {
                continue
            }

            let distance = sqrt(distanceSquared)
            let movement = CGVector(dx: dx / distance, dy: dy / distance)
            skeleton.position.x += movement.dx * GameConfiguration.skeletonSpeed * CGFloat(deltaTime)
            skeleton.position.y += movement.dy * GameConfiguration.skeletonSpeed * CGFloat(deltaTime)
            updateFacing(for: skeleton, movement: movement)
        }

        updateSkeletonAnimation(deltaTime: deltaTime)
        skeletonSpatialIndex.rebuild(with: skeletons)
    }

    private func updateSkeletonSpawning(deltaTime: TimeInterval) {
        skeletonSpawnTimer += deltaTime
        let spawnInterval = progression.skeletonSpawnInterval
        var didSpawn = false

        while skeletonSpawnTimer >= spawnInterval {
            skeletonSpawnTimer -= spawnInterval
            spawnSkeleton(kind: timedSkeletonSpawnKind, shouldUpdateHUD: false)
            didSpawn = true
        }

        if didSpawn {
            updateHUDCombatStatus()
        }
    }

    private var timedSkeletonSpawnKind: SkeletonKind {
        usesRedOnlySkeletonSpawns ? .red : .regular
    }

    private var usesRedOnlySkeletonSpawns: Bool {
        progression.level >= GameConfiguration.redOnlySkeletonLevel
    }

    private func updateOrbitalOrbs(deltaTime: TimeInterval) {
        guard progression.isOrbitalOrbUnlocked else {
            return
        }

        guard !orbitalOrbs.isEmpty else {
            return
        }

        let angleDelta = progression.orbitalOrbAngularSpeed * CGFloat(deltaTime)
        orbitalOrbAngle += angleDelta
        updateOrbitalOrbAnimation(deltaTime: deltaTime)
        respawnOrbitalOrbs(angleDelta: angleDelta)
        alignOrbitalOrbs()
        checkOrbitalOrbCollisions()
    }

    private func updateFireballCasting(deltaTime: TimeInterval) {
        guard !skeletons.isEmpty else {
            fireballCastTimer = 0
            return
        }

        fireballCastTimer += deltaTime
        let castInterval = progression.fireballCastInterval

        while fireballCastTimer >= castInterval {
            fireballCastTimer -= castInterval
            spawnFireballs()
        }
    }

    private func updateLightningCasting(deltaTime: TimeInterval) {
        guard progression.isLightningUnlocked else {
            lightningCastTimer = 0
            return
        }

        guard !skeletons.isEmpty else {
            lightningCastTimer = 0
            return
        }

        lightningCastTimer += deltaTime
        let castInterval = progression.lightningCastInterval

        while lightningCastTimer >= castInterval {
            lightningCastTimer -= castInterval
            castLightning()

            if isLevelUpChoiceActive {
                return
            }
        }
    }

    private func updateBeamCasting(deltaTime: TimeInterval) {
        guard progression.isBeamUnlocked else {
            beamCastTimer = 0
            return
        }

        guard !skeletons.isEmpty else {
            beamCastTimer = 0
            return
        }

        beamCastTimer += deltaTime
        let castInterval = progression.beamCastInterval

        while beamCastTimer >= castInterval {
            beamCastTimer -= castInterval
            castBeam()

            if isLevelUpChoiceActive {
                return
            }
        }
    }

    private func updateMeteorCasting(deltaTime: TimeInterval) {
        guard progression.isMeteorUnlocked, progression.meteorCount > 0 else {
            meteorCastTimer = 0
            return
        }

        guard !skeletons.isEmpty else {
            meteorCastTimer = 0
            return
        }

        meteorCastTimer += deltaTime
        let spawnInterval = progression.meteorCastInterval / TimeInterval(progression.meteorCount)

        while meteorCastTimer >= spawnInterval {
            meteorCastTimer -= spawnInterval
            castMeteor()
        }
    }

    private func updateFireballs(deltaTime: TimeInterval) {
        updateFireballAnimation(deltaTime: deltaTime)

        for index in fireballs.indices.reversed() {
            if let target = fireballs[index].target, !isSkeletonAlive(target) {
                fireballs[index].target = nil
            }

            if let target = fireballs[index].target {
                updateHomingFireball(at: index, target: target, deltaTime: deltaTime)
            } else {
                updateUntargetedFireball(at: index, deltaTime: deltaTime)
            }

            if isLevelUpChoiceActive {
                return
            }
        }
    }

    private func updateHomingFireball(at index: Int, target: SKSpriteNode, deltaTime: TimeInterval) {
        let fireballPosition = fireballs[index].node.position
        let dx = target.position.x - fireballPosition.x
        let dy = target.position.y - fireballPosition.y
        let distanceSquared = dx * dx + dy * dy
        let travelDistance = GameConfiguration.fireballSpeed * CGFloat(deltaTime)

        guard distanceSquared > 0 else {
            damageSkeleton(target, killedBy: .fireball)
            removeFireball(at: index)
            return
        }

        let distance = sqrt(distanceSquared)
        let movement = CGVector(dx: dx / distance, dy: dy / distance)
        fireballs[index].velocity = movement

        let hitDistance = GameConfiguration.fireballHitDistance + travelDistance

        if distanceSquared <= hitDistance * hitDistance {
            fireballs[index].node.position = target.position
            damageSkeleton(target, killedBy: .fireball)
            removeFireball(at: index)
            return
        }

        moveFireball(at: index, deltaTime: deltaTime)
    }

    private func updateUntargetedFireball(at index: Int, deltaTime: TimeInterval) {
        let startPosition = fireballs[index].node.position
        fireballs[index].timeWithoutTarget += deltaTime
        moveFireball(at: index, deltaTime: deltaTime)

        if let target = firstSkeletonHitByFireball(from: startPosition, to: fireballs[index].node.position) {
            fireballs[index].node.position = target.position
            damageSkeleton(target, killedBy: .fireball)
            removeFireball(at: index)
            return
        }

        if fireballs[index].timeWithoutTarget >= GameConfiguration.fireballUntargetedLifetime {
            removeFireball(at: index)
        }
    }

    private func moveFireball(at index: Int, deltaTime: TimeInterval) {
        let movement = fireballs[index].velocity

        fireballs[index].node.position.x += movement.dx * GameConfiguration.fireballSpeed * CGFloat(deltaTime)
        fireballs[index].node.position.y += movement.dy * GameConfiguration.fireballSpeed * CGFloat(deltaTime)
        fireballs[index].node.zRotation = atan2(movement.dy, movement.dx)
    }

    private func spawnSkeleton(kind: SkeletonKind = .regular, shouldUpdateHUD: Bool = true) {
        let skeleton = SKSpriteNode(texture: skeletonTextures[0])
        skeleton.size = CGSize(width: 30, height: 42)
        skeleton.position = skeletonSpawnPosition()
        skeleton.zPosition = 9
        skeleton.color = kind.tint
        skeleton.colorBlendFactor = kind.tintBlendFactor
        skeleton.userData = NSMutableDictionary(
            object: NSNumber(value: kind.hitPoints),
            forKey: SkeletonUserDataKey.hitPoints as NSString
        )

        let identifier = ObjectIdentifier(skeleton)
        skeletonIndices[identifier] = skeletons.count
        skeletonIdentifiers.insert(identifier)
        skeletons.append(skeleton)
        worldNode.addChild(skeleton)

        if shouldUpdateHUD {
            updateHUDCombatStatus()
        }
    }

    private func spawnFireballs() {
        let targets = availableSkeletonTargets(limit: progression.simultaneousFireballCount)

        for target in targets {
            spawnFireball(targeting: target)
        }
    }

    private func spawnFireball(targeting target: SKSpriteNode) {
        let fireballNode = SKSpriteNode(texture: fireballTextures[0])
        fireballNode.size = CGSize(width: 18, height: 18)
        fireballNode.position = player.position
        fireballNode.zPosition = 12

        worldNode.addChild(fireballNode)
        fireballs.append(
            Fireball(
                node: fireballNode,
                target: target,
                velocity: CGVector(from: player.position, to: target.position).normalized,
                timeWithoutTarget: 0
            )
        )
    }

    private func syncOrbitalOrbCount() {
        while orbitalOrbs.count < progression.orbitalOrbCount {
            let node = makeOrbitalOrbNode()
            worldNode.addChild(node)
            orbitalOrbs.append(OrbitalOrb(node: node))
        }

        while orbitalOrbs.count > progression.orbitalOrbCount {
            removeOrbitalOrb(at: orbitalOrbs.count - 1)
        }

        alignOrbitalOrbs()
    }

    private func makeOrbitalOrbNode() -> SKSpriteNode {
        let orbNode = SKSpriteNode(texture: orbitalOrbTextures[0])
        orbNode.size = CGSize(width: 20, height: 20)
        orbNode.zPosition = 11

        return orbNode
    }

    private func alignOrbitalOrbs() {
        guard !orbitalOrbs.isEmpty else {
            return
        }

        let spacing = CGFloat.pi * 2 / CGFloat(orbitalOrbs.count)
        for (index, orb) in orbitalOrbs.enumerated() {
            orb.updatePosition(
                around: player.position,
                angle: orbitalOrbAngle + spacing * CGFloat(index),
                radius: GameConfiguration.orbitalOrbRadius
            )
        }
    }

    private func removeOrbitalOrb(at index: Int) {
        orbitalOrbs[index].deactivate()
        orbitalOrbs.remove(at: index)
    }

    private func respawnOrbitalOrbs(angleDelta: CGFloat) {
        for index in orbitalOrbs.indices where orbitalOrbs[index].updateMissingOrbitProgress(by: angleDelta) {
            let orbNode = makeOrbitalOrbNode()
            worldNode.addChild(orbNode)
            orbitalOrbs[index].attach(orbNode)
        }
    }

    private func castLightning() {
        let chainLightning = ChainLightning(
            origin: player.position,
            strikeCount: progression.lightningStrikeCount,
            targets: availableLightningTargets()
        )
        var levelUpCount = 0
        var didHit = false

        for strike in chainLightning.strikes {
            guard isSkeletonAlive(strike.target) else {
                continue
            }

            showLightningStrike(from: strike.start, to: strike.end)
            showLightningTargetHit(strike.target)
            didHit = true
            levelUpCount += damageSkeleton(strike.target, killedBy: .lightning, shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)
        }

        if didHit {
            updateHUDProgress()
        }

        queueLevelUpChoices(levelUpCount)
    }

    private func showLightningStrike(from start: CGPoint, to end: CGPoint) {
        let effect = ChainLightning.makeEffectNode(from: start, to: end, texture: lightningTexture)
        worldNode.addChild(effect)
    }

    private func showLightningTargetHit(_ target: SKSpriteNode) {
        let hitSprite = SKSpriteNode(texture: target.texture)
        hitSprite.name = ChainLightning.effectName
        hitSprite.position = target.position
        hitSprite.size = target.size
        hitSprite.xScale = target.xScale
        hitSprite.yScale = target.yScale
        hitSprite.zPosition = target.zPosition + 0.5
        hitSprite.color = SKColor(calibratedRed: 0.35, green: 0.86, blue: 1.0, alpha: 1)
        hitSprite.colorBlendFactor = 0.8
        hitSprite.alpha = 0.85
        hitSprite.run(
            SKAction.sequence([
                SKAction.fadeOut(withDuration: GameConfiguration.lightningEffectDuration),
                SKAction.removeFromParent()
            ])
        )
        worldNode.addChild(hitSprite)
    }

    private func castBeam() {
        let beam = Beam(
            origin: player.position,
            direction: playerBeamDirection(),
            length: beamLength,
            hitWidth: GameConfiguration.beamHitWidth,
            killLimit: progression.beamKillCount,
            targets: skeletons
        )
        var levelUpCount = 0
        var didHit = false

        showBeam(from: beam.start, to: beam.end)

        for target in beam.targets {
            didHit = true
            levelUpCount += damageSkeleton(target, killedBy: .beam, shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)
        }

        if didHit {
            updateHUDProgress()
        }

        queueLevelUpChoices(levelUpCount)
    }

    private func showBeam(from start: CGPoint, to end: CGPoint) {
        worldNode.addChild(Beam.makeEffectNode(from: start, to: end))
    }

    private func castMeteor() {
        let meteor = Meteor(
            origin: player.position,
            targetRadius: GameConfiguration.orbitalOrbRadius * GameConfiguration.meteorTargetRadiusMultiplier
        )
        let meteorNode = makeMeteorNode()
        meteorNode.position = meteor.spawnPosition
        meteorNode.zRotation = atan2(
            meteor.impactPosition.y - meteor.spawnPosition.y,
            meteor.impactPosition.x - meteor.spawnPosition.x
        )

        worldNode.addChild(meteorNode)
        meteors.append(
            MeteorProjectile(
                node: meteorNode,
                startPosition: meteor.spawnPosition,
                impactPosition: meteor.impactPosition
            )
        )
    }

    private func makeMeteorNode() -> SKSpriteNode {
        let meteorNode = SKSpriteNode(texture: meteorTextures[0])
        meteorNode.name = Meteor.projectileName
        meteorNode.size = CGSize(width: 24, height: 24)
        meteorNode.zPosition = 14

        return meteorNode
    }

    private func updateMeteors(deltaTime: TimeInterval) {
        updateMeteorAnimation(deltaTime: deltaTime)

        for index in meteors.indices.reversed() {
            if meteors[index].update(deltaTime: deltaTime) {
                let impactPosition = meteors[index].impactPosition
                meteors[index].node.removeFromParent()
                meteors.remove(at: index)
                impactMeteor(at: impactPosition)

                if isLevelUpChoiceActive {
                    return
                }
            }
        }
    }

    private func impactMeteor(at position: CGPoint) {
        guard !isGameOver else {
            return
        }

        showMeteorImpact(at: position)

        var targets = [SKSpriteNode]()
        skeletonSpatialIndex.forEachCandidate(
            near: position,
            radius: GameConfiguration.meteorImpactRadius,
            isValid: isSkeletonAlive
        ) { skeleton in
            targets.append(skeleton)
        }
        var levelUpCount = 0

        for target in targets {
            levelUpCount += damageSkeleton(target, killedBy: .meteor, shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)
        }

        if !targets.isEmpty {
            updateHUDProgress()
        }

        queueLevelUpChoices(levelUpCount)
    }

    private func showMeteorImpact(at position: CGPoint) {
        worldNode.addChild(
            Meteor.makeImpactEffectNode(
                at: position,
                radius: GameConfiguration.meteorImpactRadius
            )
        )
    }

    private func skeletonSpawnPosition() -> CGPoint {
        let halfWidth = size.width / 2
        let halfHeight = size.height / 2
        let spawnDistance = hypot(halfWidth, halfHeight) + GameConfiguration.skeletonSpawnMargin
        let directionTarget: CGPoint

        switch Int.random(in: 0..<4) {
        case 0:
            directionTarget = CGPoint(
                x: player.position.x - halfWidth,
                y: player.position.y + CGFloat.random(in: -halfHeight...halfHeight)
            )
        case 1:
            directionTarget = CGPoint(
                x: player.position.x + halfWidth,
                y: player.position.y + CGFloat.random(in: -halfHeight...halfHeight)
            )
        case 2:
            directionTarget = CGPoint(
                x: player.position.x + CGFloat.random(in: -halfWidth...halfWidth),
                y: player.position.y - halfHeight
            )
        default:
            directionTarget = CGPoint(
                x: player.position.x + CGFloat.random(in: -halfWidth...halfWidth),
                y: player.position.y + halfHeight
            )
        }

        let direction = CGVector(from: player.position, to: directionTarget).normalized
        return CGPoint(
            x: player.position.x + direction.dx * spawnDistance,
            y: player.position.y + direction.dy * spawnDistance
        )
    }

    private func checkSkeletonCollisions() {
        guard playerHitInvulnerabilityTimer <= 0 else {
            return
        }

        if skeletonSpatialIndex.firstCandidate(
            near: player.position,
            radius: GameConfiguration.skeletonHitDistance,
            isValid: isSkeletonAlive,
            matches: { _ in true }
        ) != nil {
            damagePlayer()
        }
    }

    private func damagePlayer() {
        playerLives = max(0, playerLives - 1)
        hud.updateLives(playerLives)

        guard playerLives > 0 else {
            triggerGameOver()
            return
        }

        playerHitInvulnerabilityTimer = GameConfiguration.playerHitInvulnerabilityDuration
        showPlayerHitFeedback()
    }

    private func showPlayerHitFeedback() {
        player.removeAction(forKey: Self.playerHitFlashActionKey)
        player.alpha = 1

        let flash = SKAction.sequence([
            SKAction.fadeAlpha(to: 0.35, duration: 0.08),
            SKAction.fadeAlpha(to: 1, duration: 0.08)
        ])
        player.run(SKAction.repeat(flash, count: 6), withKey: Self.playerHitFlashActionKey)
    }

    private func checkOrbitalOrbCollisions() {
        var levelUpCount = 0
        var didHit = false

        for index in orbitalOrbs.indices where orbitalOrbs[index].isActive {
            guard let skeleton = firstSkeletonTouchingOrb(at: index) else {
                continue
            }

            orbitalOrbs[index].deactivate()
            didHit = true
            levelUpCount += damageSkeleton(skeleton, killedBy: .orbitalOrb, shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)
        }

        if didHit {
            updateHUDProgress()
        }

        queueLevelUpChoices(levelUpCount)
    }

    private func firstSkeletonTouchingOrb(at index: Int) -> SKSpriteNode? {
        guard let orbNode = orbitalOrbs[index].node else {
            return nil
        }

        return skeletonSpatialIndex.firstCandidate(
            near: orbNode.position,
            radius: GameConfiguration.orbitalOrbHitDistance,
            isValid: isSkeletonAlive,
            matches: { _ in true }
        )
    }

    private func triggerGameOver() {
        guard !isGameOver else {
            return
        }

        isGameOver = true
        isLevelUpChoiceActive = false
        isChestRewardActive = false
        playerHitInvulnerabilityTimer = 0
        pendingLevelUpLevels.removeAll(keepingCapacity: true)
        worldNode.isPaused = false
        pressedKeys.removeAll()
        hud.hideLevelUp()
        hud.hideChestReward()

        player.removeAction(forKey: Self.playerHitFlashActionKey)
        stopPlayerAnimation()
        player.color = SKColor(calibratedRed: 0.85, green: 0.05, blue: 0.08, alpha: 1)
        player.colorBlendFactor = 0.65
        player.alpha = 0.45
        player.run(SKAction.rotate(byAngle: -.pi / 2, duration: 0.16))

        hud.showGameOver(level: progression.level)

        for fireball in fireballs {
            fireball.node.removeAllActions()
        }
        removeLightningEffects()
        removeBeamEffects()
        removeMeteorEffects()
    }

    private func restartGame() {
        isGameOver = false
        isLevelUpChoiceActive = false
        isChestRewardActive = false
        worldNode.isPaused = false
        pressedKeys.removeAll()
        skeletonAnimationTimer = 0
        skeletonAnimationFrameIndex = 0
        fireballAnimationTimer = 0
        fireballAnimationFrameIndex = 0
        orbitalOrbAnimationTimer = 0
        orbitalOrbAnimationFrameIndex = 0
        meteorAnimationTimer = 0
        meteorAnimationFrameIndex = 0
        skeletonSpawnTimer = 0
        fireballCastTimer = 0
        lightningCastTimer = 0
        beamCastTimer = 0
        meteorCastTimer = 0
        playerHitInvulnerabilityTimer = 0
        playerLives = GameConfiguration.initialPlayerLives
        currentPlayerMovementDirection = nil
        orbitalOrbAngle = 0
        pendingLevelUpLevels.removeAll(keepingCapacity: true)
        totalSkeletonKills = 0
        fireballKillCount = 0
        lightningKillCount = 0
        orbitalOrbKillCount = 0
        beamKillCount = 0
        meteorKillCount = 0
        nextChestMilestone = GameConfiguration.bronzeChestKillInterval
        progression.reset()

        resetPlayer()
        removeAllEnemiesAndProjectiles()
        syncOrbitalOrbCount()

        cameraNode.position = player.position
        grassField.update(around: player.position)
        updateHUDProgress()
        hud.hideLevelUp()
        hud.hideChestReward()
        hud.hideGameOver()
        spawnSkeleton()
    }

    private func resetPlayer() {
        player.removeAllActions()
        player.position = .zero
        player.zRotation = 0
        player.xScale = 1
        player.yScale = 1
        player.alpha = 1
        player.colorBlendFactor = 0
        player.texture = mageTextures[0]
        stopPlayerAnimation()
    }

    private func removeAllEnemiesAndProjectiles() {
        skeletons.forEach { $0.removeFromParent() }
        skeletons.removeAll()
        skeletonIdentifiers.removeAll(keepingCapacity: true)
        skeletonIndices.removeAll(keepingCapacity: true)
        skeletonSpatialIndex.removeAll()

        for fireball in fireballs {
            fireball.node.removeAllActions()
            fireball.node.removeFromParent()
        }
        fireballs.removeAll()

        for meteor in meteors {
            meteor.node.removeAllActions()
            meteor.node.removeFromParent()
        }
        meteors.removeAll()

        for chest in chests {
            chest.node.removeAllActions()
            chest.node.removeFromParent()
        }
        chests.removeAll()

        for index in orbitalOrbs.indices {
            orbitalOrbs[index].deactivate()
        }
        orbitalOrbs.removeAll()

        removeLightningEffects()
        removeBeamEffects()
        removeMeteorEffects()
        updateHUDCombatStatus()
    }

    @discardableResult
    private func damageSkeleton(_ skeleton: SKSpriteNode, killedBy attackKind: AttackKind? = nil, shouldTriggerLevelUpChoice: Bool = true, shouldUpdateHUD: Bool = true) -> Int {
        guard isSkeletonAlive(skeleton) else {
            return 0
        }

        let remainingHitPoints = skeletonHitPoints(for: skeleton)

        guard remainingHitPoints > 1 else {
            return destroySkeleton(
                skeleton,
                killedBy: attackKind,
                shouldTriggerLevelUpChoice: shouldTriggerLevelUpChoice,
                shouldUpdateHUD: shouldUpdateHUD
            )
        }

        setSkeletonHitPoints(remainingHitPoints - 1, for: skeleton)
        showSkeletonDamageFeedback(skeleton)
        return 0
    }

    private func skeletonHitPoints(for skeleton: SKSpriteNode) -> Int {
        guard let hitPoints = skeleton.userData?[SkeletonUserDataKey.hitPoints] as? NSNumber else {
            return 1
        }

        return max(1, hitPoints.intValue)
    }

    private func setSkeletonHitPoints(_ hitPoints: Int, for skeleton: SKSpriteNode) {
        if skeleton.userData == nil {
            skeleton.userData = NSMutableDictionary()
        }

        skeleton.userData?[SkeletonUserDataKey.hitPoints] = NSNumber(value: max(1, hitPoints))
    }

    private func showSkeletonDamageFeedback(_ skeleton: SKSpriteNode) {
        skeleton.removeAction(forKey: Self.skeletonDamageFlashActionKey)
        skeleton.alpha = 1

        let flash = SKAction.sequence([
            SKAction.fadeAlpha(to: 0.35, duration: 0.06),
            SKAction.fadeAlpha(to: 1, duration: 0.06)
        ])
        skeleton.run(SKAction.repeat(flash, count: 2), withKey: Self.skeletonDamageFlashActionKey)
    }

    @discardableResult
    private func destroySkeleton(_ skeleton: SKSpriteNode, killedBy attackKind: AttackKind? = nil, shouldTriggerLevelUpChoice: Bool = true, shouldUpdateHUD: Bool = true) -> Int {
        let identifier = ObjectIdentifier(skeleton)

        guard let index = skeletonIndices[identifier] else {
            return 0
        }

        skeleton.removeAllActions()
        skeleton.removeFromParent()
        removeSkeletonFromTracking(identifier: identifier, at: index)
        registerSkeletonKill()
        registerAttackKill(attackKind)

        let levelUpCount = progression.gainExperience()

        if shouldUpdateHUD {
            updateHUDProgress()
        }

        if shouldTriggerLevelUpChoice {
            queueLevelUpChoices(levelUpCount)
        }

        return levelUpCount
    }

    private func registerAttackKill(_ attackKind: AttackKind?) {
        switch attackKind {
        case .fireball:
            fireballKillCount += 1
        case .lightning:
            lightningKillCount += 1
        case .orbitalOrb:
            orbitalOrbKillCount += 1
        case .beam:
            beamKillCount += 1
        case .meteor:
            meteorKillCount += 1
        case .none:
            break
        }
    }

    private func removeSkeletonFromTracking(identifier: ObjectIdentifier, at index: Int) {
        let lastIndex = skeletons.count - 1

        if index != lastIndex {
            let lastSkeleton = skeletons[lastIndex]
            skeletons[index] = lastSkeleton
            skeletonIndices[ObjectIdentifier(lastSkeleton)] = index
        }

        skeletons.removeLast()
        skeletonIndices[identifier] = nil
        skeletonIdentifiers.remove(identifier)
    }

    private func registerSkeletonKill() {
        totalSkeletonKills += 1
        spawnMilestoneSkeletonsIfNeeded()

        while totalSkeletonKills >= nextChestMilestone {
            if let tier = chestTier(for: nextChestMilestone) {
                spawnChest(tier: tier)
            }
            nextChestMilestone += GameConfiguration.bronzeChestKillInterval
        }
    }

    private func spawnMilestoneSkeletonsIfNeeded() {
        guard !usesRedOnlySkeletonSpawns else {
            return
        }

        spawnSkeleton(kind: .red, afterEveryKills: GameConfiguration.redSkeletonKillInterval)
        spawnSkeleton(kind: .purple, afterEveryKills: GameConfiguration.purpleSkeletonKillInterval)
    }

    private func spawnSkeleton(kind: SkeletonKind, afterEveryKills killInterval: Int) {
        guard killInterval > 0, totalSkeletonKills.isMultiple(of: killInterval) else {
            return
        }

        spawnSkeleton(kind: kind, shouldUpdateHUD: false)
    }

    private func chestTier(for milestone: Int) -> ChestTier? {
        if milestone.isMultiple(of: GameConfiguration.goldChestKillInterval) {
            return .gold
        }

        if milestone.isMultiple(of: GameConfiguration.silverChestKillInterval) {
            guard progression.level <= GameConfiguration.silverChestMaximumLevel else {
                return nil
            }

            return .silver
        }

        guard progression.level <= GameConfiguration.bronzeChestMaximumLevel else {
            return nil
        }

        return .bronze
    }

    private func spawnChest(tier: ChestTier) {
        guard let texture = chestTextures[tier] else {
            return
        }

        let node = SKSpriteNode(texture: texture)
        node.size = CGSize(width: 32, height: 28)
        node.position = randomChestPosition()
        node.zPosition = 8.5

        worldNode.addChild(node)
        chests.append(Chest(node: node, tier: tier))
    }

    private func randomChestPosition() -> CGPoint {
        let halfWidth = max(48, size.width / 2 - GameConfiguration.chestSpawnMargin)
        let halfHeight = max(48, size.height / 2 - GameConfiguration.chestSpawnMargin)
        let minimumDistanceSquared = GameConfiguration.chestPickupDistance * GameConfiguration.chestPickupDistance * 4

        for _ in 0..<12 {
            let position = CGPoint(
                x: player.position.x + CGFloat.random(in: -halfWidth...halfWidth),
                y: player.position.y + CGFloat.random(in: -halfHeight...halfHeight)
            )

            if position.distanceSquared(to: player.position) >= minimumDistanceSquared {
                return position
            }
        }

        return CGPoint(x: player.position.x + halfWidth, y: player.position.y)
    }

    private func checkChestPickups() {
        let pickupDistanceSquared = GameConfiguration.chestPickupDistance * GameConfiguration.chestPickupDistance

        for index in chests.indices.reversed() {
            guard chests[index].node.position.distanceSquared(to: player.position) <= pickupDistanceSquared else {
                continue
            }

            collectChest(at: index)
            return
        }
    }

    private func collectChest(at index: Int) {
        let chest = chests[index]
        chest.node.removeAllActions()
        chest.node.removeFromParent()
        chests.remove(at: index)

        let items = applyChestReward(chest.tier)
        showChestReward(tier: chest.tier, items: items)
    }

    private func applyChestReward(_ tier: ChestTier) -> [ChestRewardDisplayItem] {
        let learnedSkills = progression.learnedSkills
        let items: [ChestRewardDisplayItem]

        switch tier {
        case .bronze:
            guard let skill = learnedSkills.randomElement(),
                  let option = skill.upgradeOptions.randomElement() else {
                return []
            }

            items = [chestRewardItem(for: option, skill: skill)]
            applyUpgradeEffect(option)
        case .silver:
            guard let skill = learnedSkills.randomElement() else {
                return []
            }

            items = chestRewardItems(for: [skill])
            progression.upgradeAllProperties(for: skill)
        case .gold:
            let rewardedSkills = Array(learnedSkills.shuffled().prefix(2))
            items = chestRewardItems(for: rewardedSkills)
            progression.upgradeAllProperties(for: rewardedSkills)
        }

        syncOrbitalOrbCount()
        updateHUDProgress()

        return items
    }

    private func chestRewardItems(for skills: [LearnedSkill]) -> [ChestRewardDisplayItem] {
        var items = [ChestRewardDisplayItem]()

        for skill in skills {
            for option in skill.upgradeOptions {
                items.append(chestRewardItem(for: option, skill: skill))
            }
        }

        return items
    }

    private func chestRewardItem(for option: LevelUpOption, skill: LearnedSkill) -> ChestRewardDisplayItem {
        ChestRewardDisplayItem(
            option: option,
            title: option.title(beamKillBonus: skill == .beam ? progression.beamKillUpgradeBonus : nil)
        )
    }

    private func showChestReward(tier: ChestTier, items: [ChestRewardDisplayItem]) {
        guard !items.isEmpty else {
            return
        }

        isChestRewardActive = true
        worldNode.isPaused = true
        pressedKeys.removeAll()
        stopPlayerAnimation()
        hud.showChestReward(tier: tier, items: items)
    }

    private func advanceChestReward(with keyCode: UInt16) {
        guard ChestRewardKey.isAdvance(keyCode) else {
            return
        }

        isChestRewardActive = false
        worldNode.isPaused = false
        hud.hideChestReward()
        presentNextLevelUpChoiceIfNeeded()
    }

    private func removeFireball(at index: Int) {
        fireballs[index].node.removeAllActions()
        fireballs[index].node.removeFromParent()
        fireballs.remove(at: index)
    }

    private func firstSkeletonHitByFireball(from start: CGPoint, to end: CGPoint) -> SKSpriteNode? {
        let dx = end.x - start.x
        let dy = end.y - start.y
        let lengthSquared = dx * dx + dy * dy
        var closestHit: (skeleton: SKSpriteNode, progress: CGFloat)?

        let hitRadius = GameConfiguration.fireballHitDistance
        let searchRect = CGRect(
            x: min(start.x, end.x) - hitRadius,
            y: min(start.y, end.y) - hitRadius,
            width: abs(dx) + hitRadius * 2,
            height: abs(dy) + hitRadius * 2
        )
        let hitRadiusSquared = hitRadius * hitRadius

        skeletonSpatialIndex.forEachCandidate(in: searchRect, isValid: isSkeletonAlive) { skeleton in
            let progress: CGFloat

            if lengthSquared > 0 {
                let skeletonDx = skeleton.position.x - start.x
                let skeletonDy = skeleton.position.y - start.y
                let rawProgress = (skeletonDx * dx + skeletonDy * dy) / lengthSquared
                progress = min(1, max(0, rawProgress))
            } else {
                progress = 0
            }

            let closestPoint = CGPoint(
                x: start.x + dx * progress,
                y: start.y + dy * progress
            )

            guard closestPoint.distanceSquared(to: skeleton.position) <= hitRadiusSquared else {
                return
            }

            if closestHit == nil || progress < closestHit!.progress {
                closestHit = (skeleton, progress)
            }
        }

        return closestHit?.skeleton
    }

    private func updateHUDProgress() {
        hud.updateProgress(
            level: progression.level,
            experience: progression.experience,
            nextExperience: progression.nextExperience
        )
        hud.updateLives(playerLives)
        hud.updateFireballStatus(
            count: progression.simultaneousFireballCount,
            interval: progression.fireballCastInterval
        )
        hud.updateLightningStatus(
            isUnlocked: progression.isLightningUnlocked,
            strikeCount: progression.lightningStrikeCount,
            interval: progression.lightningCastInterval
        )
        hud.updateOrbStatus(
            isUnlocked: progression.isOrbitalOrbUnlocked,
            count: progression.orbitalOrbCount,
            angularSpeed: progression.orbitalOrbAngularSpeed
        )
        hud.updateBeamStatus(
            isUnlocked: progression.isBeamUnlocked,
            killCount: progression.beamKillCount,
            interval: progression.beamCastInterval
        )
        hud.updateMeteorStatus(
            isUnlocked: progression.isMeteorUnlocked,
            count: progression.meteorCount,
            interval: progression.meteorCastInterval
        )
        hud.updateAttackKillCounts(
            fireball: fireballKillCount,
            lightning: lightningKillCount,
            orb: orbitalOrbKillCount,
            beam: beamKillCount,
            meteor: meteorKillCount
        )
        updateHUDCombatStatus()
    }

    private func updateHUDCombatStatus() {
        hud.updateSkeletonStatus(
            aliveCount: skeletons.count,
            spawnInterval: progression.skeletonSpawnInterval
        )
    }

    private func updatePlayerHitInvulnerability(deltaTime: TimeInterval) {
        guard playerHitInvulnerabilityTimer > 0 else {
            return
        }

        playerHitInvulnerabilityTimer = max(0, playerHitInvulnerabilityTimer - deltaTime)
    }

    private func killAllEnemiesAndGrantExperience() {
        let defeatedCount = skeletons.count

        guard defeatedCount > 0 else {
            return
        }

        skeletons.forEach {
            $0.removeAllActions()
            $0.removeFromParent()
        }
        skeletons.removeAll()
        skeletonIdentifiers.removeAll(keepingCapacity: true)
        skeletonIndices.removeAll(keepingCapacity: true)
        skeletonSpatialIndex.removeAll()

        let levelUpCount = progression.gainExperience(defeatedCount)
        updateHUDProgress()
        queueLevelUpChoices(levelUpCount)
    }

    private func queueLevelUpChoices(_ count: Int) {
        guard count > 0, !isGameOver else {
            return
        }

        let firstQueuedLevel = progression.level - count + 1
        pendingLevelUpLevels.append(contentsOf: firstQueuedLevel...progression.level)
        presentNextLevelUpChoiceIfNeeded()
    }

    private func presentNextLevelUpChoiceIfNeeded() {
        guard !isGameOver, !isLevelUpChoiceActive, let level = pendingLevelUpLevels.first else {
            return
        }

        isLevelUpChoiceActive = true
        worldNode.isPaused = true
        pressedKeys.removeAll()
        stopPlayerAnimation()
        hud.showLevelUp(
            level: level,
            options: randomLevelUpOptions(),
            beamKillUpgradeBonus: progression.beamKillUpgradeBonus
        )
    }

    private func applyLevelUpOption(_ option: LevelUpOption) {
        applyUpgradeEffect(option)

        syncOrbitalOrbCount()
        updateHUDProgress()

        if !pendingLevelUpLevels.isEmpty {
            pendingLevelUpLevels.removeFirst()
        }

        isLevelUpChoiceActive = false
        worldNode.isPaused = false
        hud.hideLevelUp()
        presentNextLevelUpChoiceIfNeeded()
    }

    private func applyUpgradeEffect(_ option: LevelUpOption) {
        switch option {
        case .extraLife:
            playerLives += 1
        case .halveSkeletons:
            halveSkeletons()
        default:
            progression.applyLevelUpOption(option)
        }
    }

    private func halveSkeletons() {
        let killCount = skeletons.count / 2

        guard killCount > 0 else {
            return
        }

        let targets = Array(skeletons.shuffled().prefix(killCount))
        var levelUpCount = 0

        for target in targets {
            levelUpCount += destroySkeleton(target, shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)
        }

        updateHUDProgress()
        queueLevelUpChoices(levelUpCount)
    }

    private func selectLevelUpOption(with keyCode: UInt16) {
        guard let index = LevelUpSelectionKey.optionIndex(for: keyCode),
              let option = hud.levelUpOption(atIndex: index) else {
            return
        }

        applyLevelUpOption(option)
    }

    private func availableSkeletonTargets(limit: Int) -> [SKSpriteNode] {
        let reservedTargets = fireballTargetIdentifiers()
        return closestSkeletons(to: player.position, excluding: reservedTargets, limit: limit)
    }

    private func availableLightningTargets() -> [SKSpriteNode] {
        let reservedTargets = fireballTargetIdentifiers()
        return skeletons.filter { skeleton in
            isSkeletonAlive(skeleton) && !reservedTargets.contains(ObjectIdentifier(skeleton))
        }
    }

    private func fireballTargetIdentifiers() -> Set<ObjectIdentifier> {
        Set(
            fireballs.compactMap { fireball in
                fireball.target.map(ObjectIdentifier.init)
            }
        )
    }

    private func closestSkeletons(to position: CGPoint, excluding excludedIdentifiers: Set<ObjectIdentifier>, limit: Int) -> [SKSpriteNode] {
        guard limit > 0 else {
            return []
        }

        var selected = [(skeleton: SKSpriteNode, distanceSquared: CGFloat)]()
        selected.reserveCapacity(limit)

        for skeleton in skeletons {
            guard isSkeletonAlive(skeleton), !excludedIdentifiers.contains(ObjectIdentifier(skeleton)) else {
                continue
            }

            let distanceSquared = skeleton.position.distanceSquared(to: position)

            if selected.count < limit {
                selected.append((skeleton, distanceSquared))
                var index = selected.count - 1

                while index > 0 && selected[index].distanceSquared < selected[index - 1].distanceSquared {
                    selected.swapAt(index, index - 1)
                    index -= 1
                }
            } else if distanceSquared < selected[selected.count - 1].distanceSquared {
                selected[selected.count - 1] = (skeleton, distanceSquared)
                var index = selected.count - 1

                while index > 0 && selected[index].distanceSquared < selected[index - 1].distanceSquared {
                    selected.swapAt(index, index - 1)
                    index -= 1
                }
            }
        }

        return selected.map(\.skeleton)
    }

    private func removeLightningEffects() {
        worldNode.enumerateChildNodes(withName: ChainLightning.effectName) { node, _ in
            node.removeAllActions()
            node.removeFromParent()
        }
    }

    private func removeBeamEffects() {
        worldNode.enumerateChildNodes(withName: Beam.effectName) { node, _ in
            node.removeAllActions()
            node.removeFromParent()
        }
    }

    private func removeMeteorEffects() {
        meteors.removeAll()

        worldNode.enumerateChildNodes(withName: Meteor.projectileName) { node, _ in
            node.removeAllActions()
            node.removeFromParent()
        }
        worldNode.enumerateChildNodes(withName: Meteor.effectName) { node, _ in
            node.removeAllActions()
            node.removeFromParent()
        }
    }

    private func updateFacing(for node: SKSpriteNode, movement: CGVector) {
        if movement.dx < 0 {
            node.xScale = -abs(node.xScale)
        } else if movement.dx > 0 {
            node.xScale = abs(node.xScale)
        }
    }

    private func startPlayerAnimation() {
        guard player.action(forKey: Self.playerAnimationActionKey) == nil else {
            return
        }

        player.removeAction(forKey: Self.playerAnimationActionKey)
        player.run(
            SKAction.repeatForever(
                SKAction.animate(
                    with: mageTextures,
                    timePerFrame: GameConfiguration.playerAnimationFrameDuration
                )
            ),
            withKey: Self.playerAnimationActionKey
        )
    }

    private func stopPlayerAnimation() {
        player.removeAction(forKey: Self.playerAnimationActionKey)
        player.texture = mageTextures[0]
    }

    private func updateSkeletonAnimation(deltaTime: TimeInterval) {
        guard !skeletons.isEmpty else {
            skeletonAnimationTimer = 0
            return
        }

        skeletonAnimationTimer += deltaTime

        guard skeletonAnimationTimer >= GameConfiguration.skeletonAnimationFrameDuration else {
            return
        }

        skeletonAnimationTimer.formTruncatingRemainder(dividingBy: GameConfiguration.skeletonAnimationFrameDuration)
        skeletonAnimationFrameIndex = (skeletonAnimationFrameIndex + 1) % skeletonTextures.count
        let texture = skeletonTextures[skeletonAnimationFrameIndex]

        for skeleton in skeletons {
            skeleton.texture = texture
        }
    }

    private func updateFireballAnimation(deltaTime: TimeInterval) {
        guard !fireballs.isEmpty else {
            fireballAnimationTimer = 0
            return
        }

        fireballAnimationTimer += deltaTime

        guard fireballAnimationTimer >= GameConfiguration.fireballAnimationFrameDuration else {
            return
        }

        fireballAnimationTimer.formTruncatingRemainder(dividingBy: GameConfiguration.fireballAnimationFrameDuration)
        fireballAnimationFrameIndex = (fireballAnimationFrameIndex + 1) % fireballTextures.count
        let texture = fireballTextures[fireballAnimationFrameIndex]

        for fireball in fireballs {
            fireball.node.texture = texture
        }
    }

    private func updateOrbitalOrbAnimation(deltaTime: TimeInterval) {
        guard orbitalOrbs.contains(where: { $0.isActive }) else {
            orbitalOrbAnimationTimer = 0
            return
        }

        orbitalOrbAnimationTimer += deltaTime

        guard orbitalOrbAnimationTimer >= GameConfiguration.orbitalOrbAnimationFrameDuration else {
            return
        }

        orbitalOrbAnimationTimer.formTruncatingRemainder(dividingBy: GameConfiguration.orbitalOrbAnimationFrameDuration)
        orbitalOrbAnimationFrameIndex = (orbitalOrbAnimationFrameIndex + 1) % orbitalOrbTextures.count
        let texture = orbitalOrbTextures[orbitalOrbAnimationFrameIndex]

        for orb in orbitalOrbs {
            orb.node?.texture = texture
        }
    }

    private func updateMeteorAnimation(deltaTime: TimeInterval) {
        guard !meteors.isEmpty else {
            meteorAnimationTimer = 0
            return
        }

        meteorAnimationTimer += deltaTime

        guard meteorAnimationTimer >= GameConfiguration.meteorAnimationFrameDuration else {
            return
        }

        meteorAnimationTimer.formTruncatingRemainder(dividingBy: GameConfiguration.meteorAnimationFrameDuration)
        meteorAnimationFrameIndex = (meteorAnimationFrameIndex + 1) % meteorTextures.count
        let texture = meteorTextures[meteorAnimationFrameIndex]

        for meteor in meteors {
            meteor.node.texture = texture
        }
    }

    private func isSkeletonAlive(_ skeleton: SKSpriteNode) -> Bool {
        skeletonIdentifiers.contains(ObjectIdentifier(skeleton))
    }

    private func randomLevelUpOptions() -> [LevelUpOption] {
        let optionCount = shouldShowThirdLevelUpOption() ? 3 : 2
        let availableOptions = progression.availableLevelUpOptions.filter { option in
            option != .halveSkeletons
        }
        var selectedOptions = Array(availableOptions.shuffled().prefix(optionCount))

        if shouldShowHalveHordeOption(), !selectedOptions.isEmpty {
            selectedOptions[Int.random(in: selectedOptions.indices)] = .halveSkeletons
        }

        return selectedOptions
    }

    private func shouldShowHalveHordeOption() -> Bool {
        guard !skeletons.isEmpty else {
            return false
        }

        let numerator = GameConfiguration.halveHordeLevelUpOptionChanceNumerator
        let denominator = GameConfiguration.halveHordeLevelUpOptionChanceDenominator

        guard numerator > 0, denominator > 0 else {
            return false
        }

        return Int.random(in: 1...denominator) <= numerator
    }

    private func shouldShowThirdLevelUpOption() -> Bool {
        let numerator = GameConfiguration.thirdLevelUpOptionChanceNumerator
        let denominator = GameConfiguration.thirdLevelUpOptionChanceDenominator

        guard numerator > 0, denominator > 0 else {
            return false
        }

        return Int.random(in: 1...denominator) <= numerator
    }

    private func playerBeamDirection() -> CGVector {
        if let movementDirection = currentPlayerMovementDirection {
            return movementDirection
        }

        return CGVector(dx: player.xScale < 0 ? -1 : 1, dy: 0)
    }

    private var beamLength: CGFloat {
        max(700, hypot(size.width, size.height) / 2 + GameConfiguration.skeletonSpawnMargin)
    }

    private func directionValue(negative: UInt16, positive: UInt16) -> CGFloat {
        var value: CGFloat = 0

        if pressedKeys.contains(negative) {
            value -= 1
        }

        if pressedKeys.contains(positive) {
            value += 1
        }

        return value
    }

    private static let skeletonDamageFlashActionKey = "skeletonDamageFlash"
    private static let playerHitFlashActionKey = "playerHitFlash"
    private static let playerAnimationActionKey = "playerAnimation"
}
