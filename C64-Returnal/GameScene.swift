//
//  GameScene.swift
//  C64-Returnal
//
//  Created by Tóth István on 2026. 06. 02..
//

import SpriteKit

final class GameScene: SKScene {
    private let worldNode = SKNode()
    private let cameraNode = SKCameraNode()
    private let mageTextures = PixelArtFactory.makeMageTextures()
    private lazy var player = SKSpriteNode(texture: mageTextures[0])
    private let hud = GameHUD()
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

    private var progression = Progression()
    private var skeletons = [SKSpriteNode]()
    private var fireballs = [Fireball]()
    private var orbitalOrbs = [OrbitalOrb]()
    private var pressedKeys = Set<UInt16>()
    private var lastUpdateTime: TimeInterval = 0
    private var skeletonSpawnTimer: TimeInterval = 0
    private var fireballCastTimer: TimeInterval = 0
    private var lightningCastTimer: TimeInterval = 0
    private var beamCastTimer: TimeInterval = 0
    private var meteorCastTimer: TimeInterval = 0
    private var playerHitInvulnerabilityTimer: TimeInterval = 0
    private var playerLives = GameConfiguration.initialPlayerLives
    private var currentPlayerMovementDirection: CGVector?
    private var orbitalOrbAngle: CGFloat = 0
    private var isGameOver = false
    private var isLevelUpChoiceActive = false
    private var isSceneConfigured = false

    override func didMove(to view: SKView) {
        if !isSceneConfigured {
            configureScene()
        }

        view.window?.makeFirstResponder(view)
        layoutViewportContent()
    }

    override func didChangeSize(_ oldSize: CGSize) {
        layoutViewportContent()
    }

    override func keyDown(with event: NSEvent) {
        guard !isGameOver, !isLevelUpChoiceActive else {
            return
        }

        if DebugKey.isLevelSetupShortcut(event.keyCode) {
            advanceDebugExperience()
            return
        }

        guard ArrowKey.contains(event.keyCode) else {
            super.keyDown(with: event)
            return
        }

        pressedKeys.insert(event.keyCode)
    }

    override func keyUp(with event: NSEvent) {
        guard !isGameOver, !isLevelUpChoiceActive else {
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

        if !isGameOver && !isLevelUpChoiceActive {
            updatePlayer(deltaTime: deltaTime)
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
            let movement = CGVector(from: skeleton.position, to: player.position).normalized

            guard movement.dx != 0 || movement.dy != 0 else {
                continue
            }

            skeleton.position.x += movement.dx * GameConfiguration.skeletonSpeed * CGFloat(deltaTime)
            skeleton.position.y += movement.dy * GameConfiguration.skeletonSpeed * CGFloat(deltaTime)
            updateFacing(for: skeleton, movement: movement)
        }
    }

    private func updateSkeletonSpawning(deltaTime: TimeInterval) {
        let skeletonLimit = progression.maximumSkeletons

        guard skeletons.count < skeletonLimit else {
            return
        }

        skeletonSpawnTimer += deltaTime

        while skeletonSpawnTimer >= progression.skeletonSpawnInterval && skeletons.count < skeletonLimit {
            skeletonSpawnTimer -= progression.skeletonSpawnInterval
            spawnSkeleton()
        }
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

        while fireballCastTimer >= progression.fireballCastInterval {
            fireballCastTimer -= progression.fireballCastInterval
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

        while lightningCastTimer >= progression.lightningCastInterval {
            lightningCastTimer -= progression.lightningCastInterval
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

        while beamCastTimer >= progression.beamCastInterval {
            beamCastTimer -= progression.beamCastInterval
            castBeam()

            if isLevelUpChoiceActive {
                return
            }
        }
    }

    private func updateMeteorCasting(deltaTime: TimeInterval) {
        guard progression.isMeteorUnlocked else {
            meteorCastTimer = 0
            return
        }

        guard !skeletons.isEmpty else {
            meteorCastTimer = 0
            return
        }

        meteorCastTimer += deltaTime

        while meteorCastTimer >= progression.meteorCastInterval {
            meteorCastTimer -= progression.meteorCastInterval
            castMeteors()
        }
    }

    private func updateFireballs(deltaTime: TimeInterval) {
        for index in fireballs.indices.reversed() {
            if let target = fireballs[index].target, !skeletons.contains(where: { $0 === target }) {
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
        let distanceToTarget = fireballs[index].node.position.distance(to: target.position)
        let travelDistance = GameConfiguration.fireballSpeed * CGFloat(deltaTime)
        let movement = CGVector(from: fireballs[index].node.position, to: target.position).normalized

        guard movement.dx != 0 || movement.dy != 0 else {
            destroySkeleton(target)
            removeFireball(at: index)
            return
        }

        fireballs[index].velocity = movement

        if distanceToTarget <= GameConfiguration.fireballHitDistance + travelDistance {
            fireballs[index].node.position = target.position
            destroySkeleton(target)
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
            destroySkeleton(target)
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

    private func spawnSkeleton() {
        let skeleton = SKSpriteNode(texture: skeletonTextures[0])
        skeleton.size = CGSize(width: 30, height: 42)
        skeleton.position = skeletonSpawnPosition()
        skeleton.zPosition = 9
        startSkeletonAnimation(skeleton)

        skeletons.append(skeleton)
        worldNode.addChild(skeleton)
        updateHUDCombatStatus()
    }

    private func spawnFireballs() {
        let targets = availableSkeletonTargets()
            .prefix(progression.simultaneousFireballCount)

        for target in targets {
            spawnFireball(targeting: target)
        }
    }

    private func spawnFireball(targeting target: SKSpriteNode) {
        let fireballNode = SKSpriteNode(texture: fireballTextures[0])
        fireballNode.size = CGSize(width: 18, height: 18)
        fireballNode.position = player.position
        fireballNode.zPosition = 12
        fireballNode.run(
            SKAction.repeatForever(
                SKAction.animate(
                    with: fireballTextures,
                    timePerFrame: GameConfiguration.fireballAnimationFrameDuration
                )
            )
        )

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
        orbNode.run(
            SKAction.repeatForever(
                SKAction.animate(
                    with: orbitalOrbTextures,
                    timePerFrame: GameConfiguration.orbitalOrbAnimationFrameDuration
                )
            )
        )

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
        var didLevelUp = false

        for strike in chainLightning.strikes {
            showLightningStrike(from: strike.start, to: strike.end)
            didLevelUp = destroySkeleton(strike.target, shouldTriggerLevelUpChoice: false) || didLevelUp
        }

        if didLevelUp {
            triggerLevelUpChoice()
        }
    }

    private func showLightningStrike(from start: CGPoint, to end: CGPoint) {
        let effect = ChainLightning.makeEffectNode(from: start, to: end, texture: lightningTexture)
        worldNode.addChild(effect)
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
        var didLevelUp = false

        showBeam(from: beam.start, to: beam.end)

        for target in beam.targets {
            didLevelUp = destroySkeleton(target, shouldTriggerLevelUpChoice: false) || didLevelUp
        }

        if didLevelUp {
            triggerLevelUpChoice()
        }
    }

    private func showBeam(from start: CGPoint, to end: CGPoint) {
        worldNode.addChild(Beam.makeEffectNode(from: start, to: end))
    }

    private func castMeteors() {
        for _ in 0..<progression.meteorCount {
            castMeteor()
        }
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
        meteorNode.run(
            SKAction.sequence([
                SKAction.move(to: meteor.impactPosition, duration: GameConfiguration.meteorFallDuration),
                SKAction.removeFromParent()
            ])
        ) { [weak self] in
            self?.impactMeteor(at: meteor.impactPosition)
        }
    }

    private func makeMeteorNode() -> SKSpriteNode {
        let meteorNode = SKSpriteNode(texture: meteorTextures[0])
        meteorNode.name = Meteor.projectileName
        meteorNode.size = CGSize(width: 24, height: 24)
        meteorNode.zPosition = 14
        meteorNode.run(
            SKAction.repeatForever(
                SKAction.animate(
                    with: meteorTextures,
                    timePerFrame: GameConfiguration.meteorAnimationFrameDuration
                )
            )
        )

        return meteorNode
    }

    private func impactMeteor(at position: CGPoint) {
        guard !isGameOver else {
            return
        }

        showMeteorImpact(at: position)

        let targets = skeletons.filter { skeleton in
            skeleton.position.distance(to: position) <= GameConfiguration.meteorImpactRadius
        }
        var didLevelUp = false

        for target in targets {
            didLevelUp = destroySkeleton(target, shouldTriggerLevelUpChoice: false) || didLevelUp
        }

        if didLevelUp {
            triggerLevelUpChoice()
        }
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

        for skeleton in skeletons where skeleton.position.distance(to: player.position) <= GameConfiguration.skeletonHitDistance {
            damagePlayer()
            return
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
        var didLevelUp = false

        for index in orbitalOrbs.indices where orbitalOrbs[index].isActive {
            guard let skeleton = firstSkeletonTouchingOrb(at: index) else {
                continue
            }

            orbitalOrbs[index].deactivate()
            didLevelUp = destroySkeleton(skeleton, shouldTriggerLevelUpChoice: false) || didLevelUp
        }

        if didLevelUp {
            triggerLevelUpChoice()
        }
    }

    private func firstSkeletonTouchingOrb(at index: Int) -> SKSpriteNode? {
        guard let orbNode = orbitalOrbs[index].node else {
            return nil
        }

        return skeletons.first { skeleton in
            orbNode.position.distance(to: skeleton.position) <= GameConfiguration.orbitalOrbHitDistance
        }
    }

    private func triggerGameOver() {
        guard !isGameOver else {
            return
        }

        isGameOver = true
        isLevelUpChoiceActive = false
        playerHitInvulnerabilityTimer = 0
        worldNode.isPaused = false
        pressedKeys.removeAll()
        hud.hideLevelUp()

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
        worldNode.isPaused = false
        pressedKeys.removeAll()
        skeletonSpawnTimer = 0
        fireballCastTimer = 0
        lightningCastTimer = 0
        beamCastTimer = 0
        meteorCastTimer = 0
        playerHitInvulnerabilityTimer = 0
        playerLives = GameConfiguration.initialPlayerLives
        currentPlayerMovementDirection = nil
        orbitalOrbAngle = 0
        progression.reset()

        resetPlayer()
        removeAllEnemiesAndProjectiles()
        syncOrbitalOrbCount()

        cameraNode.position = player.position
        grassField.update(around: player.position)
        updateHUDProgress()
        hud.hideLevelUp()
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

        for fireball in fireballs {
            fireball.node.removeAllActions()
            fireball.node.removeFromParent()
        }
        fireballs.removeAll()

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
    private func destroySkeleton(_ skeleton: SKSpriteNode, shouldTriggerLevelUpChoice: Bool = true) -> Bool {
        guard let index = skeletons.firstIndex(where: { $0 === skeleton }) else {
            return false
        }

        skeleton.removeAllActions()
        skeleton.removeFromParent()
        skeletons.remove(at: index)

        let didLevelUp = progression.gainExperience()
        updateHUDProgress()

        if didLevelUp && shouldTriggerLevelUpChoice {
            triggerLevelUpChoice()
        }

        return didLevelUp
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

        for skeleton in skeletons {
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

            guard closestPoint.distance(to: skeleton.position) <= GameConfiguration.fireballHitDistance else {
                continue
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

    private func triggerLevelUpChoice() {
        guard !isGameOver else {
            return
        }

        isLevelUpChoiceActive = true
        worldNode.isPaused = true
        pressedKeys.removeAll()
        stopPlayerAnimation()
        hud.showLevelUp(
            level: progression.level,
            options: randomLevelUpOptions(),
            beamKillUpgradeBonus: progression.beamKillUpgradeBonus
        )
    }

    private func applyLevelUpOption(_ option: LevelUpOption) {
        if option == .extraLife {
            playerLives += 1
        } else {
            progression.applyLevelUpOption(option)
        }

        syncOrbitalOrbCount()
        updateHUDProgress()

        isLevelUpChoiceActive = false
        worldNode.isPaused = false
        hud.hideLevelUp()
    }

    private func advanceDebugExperience() {
        progression.advanceToOneKillBeforeNextLevel()
        updateHUDProgress()
    }

    private func availableSkeletonTargets() -> [SKSpriteNode] {
        let reservedTargets = fireballTargetIdentifiers()
        return skeletonsByDistance(to: player.position).filter { skeleton in
            !reservedTargets.contains(ObjectIdentifier(skeleton))
        }
    }

    private func availableLightningTargets() -> [SKSpriteNode] {
        let reservedTargets = fireballTargetIdentifiers()
        return skeletons.filter { skeleton in
            !reservedTargets.contains(ObjectIdentifier(skeleton))
        }
    }

    private func fireballTargetIdentifiers() -> Set<ObjectIdentifier> {
        Set(
            fireballs.compactMap { fireball in
                fireball.target.map(ObjectIdentifier.init)
            }
        )
    }

    private func skeletonsByDistance(to position: CGPoint) -> [SKSpriteNode] {
        skeletons.sorted {
            $0.position.distance(to: position) < $1.position.distance(to: position)
        }
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

    private func startSkeletonAnimation(_ skeleton: SKSpriteNode) {
        skeleton.run(
            SKAction.repeatForever(
                SKAction.animate(
                    with: skeletonTextures,
                    timePerFrame: GameConfiguration.skeletonAnimationFrameDuration
                )
            ),
            withKey: Self.skeletonAnimationActionKey
        )
    }

    private func randomLevelUpOptions() -> [LevelUpOption] {
        Array(progression.availableLevelUpOptions.shuffled().prefix(2))
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

    private static let playerHitFlashActionKey = "playerHitFlash"
    private static let playerAnimationActionKey = "playerAnimation"
    private static let skeletonAnimationActionKey = "skeletonAnimation"
}
