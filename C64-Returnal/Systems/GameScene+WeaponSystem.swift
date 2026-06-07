import SpriteKit

extension GameScene {
    func updateOrbitalOrbs(deltaTime: TimeInterval) {
        guard session.progression.isOrbitalOrbUnlocked else {
            return
        }

        guard !orbitalOrbs.isEmpty else {
            return
        }

        let angleDelta = session.progression.orbitalOrbAngularSpeed * CGFloat(deltaTime)
        session.orbitalOrbAngle += angleDelta
        updateOrbitalOrbAnimation(deltaTime: deltaTime)
        respawnOrbitalOrbs(angleDelta: angleDelta)
        alignOrbitalOrbs()
        checkOrbitalOrbCollisions()
    }

    func updateFireballCasting(deltaTime: TimeInterval) {
        guard !skeletons.isEmpty else {
            session.casts.fireball = 0
            return
        }

        session.casts.fireball += deltaTime
        let castInterval = session.progression.fireballCastInterval

        while session.casts.fireball >= castInterval {
            session.casts.fireball -= castInterval
            spawnFireballs()
        }
    }

    func updateLightningCasting(deltaTime: TimeInterval) {
        guard session.progression.isLightningUnlocked else {
            session.casts.lightning = 0
            return
        }

        guard !skeletons.isEmpty else {
            session.casts.lightning = 0
            return
        }

        session.casts.lightning += deltaTime
        let castInterval = session.progression.lightningCastInterval

        while session.casts.lightning >= castInterval {
            session.casts.lightning -= castInterval
            castLightning()

            if session.isLevelUpChoiceActive {
                return
            }
        }
    }

    func updateBeamCasting(deltaTime: TimeInterval) {
        guard session.progression.isBeamUnlocked else {
            session.casts.beam = 0
            return
        }

        guard !skeletons.isEmpty else {
            session.casts.beam = 0
            return
        }

        session.casts.beam += deltaTime
        let castInterval = session.progression.beamCastInterval

        while session.casts.beam >= castInterval {
            session.casts.beam -= castInterval
            castBeam()

            if session.isLevelUpChoiceActive {
                return
            }
        }
    }

    func updateMeteorCasting(deltaTime: TimeInterval) {
        guard session.progression.isMeteorUnlocked, session.progression.meteorCount > 0 else {
            session.casts.meteor = 0
            return
        }

        guard !skeletons.isEmpty else {
            session.casts.meteor = 0
            return
        }

        session.casts.meteor += deltaTime
        let spawnInterval = session.progression.meteorCastInterval / TimeInterval(session.progression.meteorCount)

        while session.casts.meteor >= spawnInterval {
            session.casts.meteor -= spawnInterval
            castMeteor()
        }
    }

    func updateFireballs(deltaTime: TimeInterval) {
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

            if session.isLevelUpChoiceActive {
                return
            }
        }
    }

    func updateHomingFireball(at index: Int, target: SKSpriteNode, deltaTime: TimeInterval) {
        let fireballPosition = fireballs[index].node.position
        let dx = target.position.x - fireballPosition.x
        let dy = target.position.y - fireballPosition.y
        let distanceSquared = dx * dx + dy * dy
        let travelDistance = tuning.fireball.speed * CGFloat(deltaTime)

        guard distanceSquared > 0 else {
            _ = damageSkeleton(target, killedBy: .fireball)
            removeFireball(at: index)
            return
        }

        let distance = sqrt(distanceSquared)
        let movement = CGVector(dx: dx / distance, dy: dy / distance)
        fireballs[index].velocity = movement

        let hitDistance = tuning.fireball.hitDistance + travelDistance

        if distanceSquared <= hitDistance * hitDistance {
            fireballs[index].node.position = target.position
            _ = damageSkeleton(target, killedBy: .fireball)
            removeFireball(at: index)
            return
        }

        moveFireball(at: index, deltaTime: deltaTime)
    }

    func updateUntargetedFireball(at index: Int, deltaTime: TimeInterval) {
        let startPosition = fireballs[index].node.position
        fireballs[index].timeWithoutTarget += deltaTime
        moveFireball(at: index, deltaTime: deltaTime)

        if let target = firstSkeletonHitByFireball(from: startPosition, to: fireballs[index].node.position) {
            fireballs[index].node.position = target.position
            _ = damageSkeleton(target, killedBy: .fireball)
            removeFireball(at: index)
            return
        }

        if fireballs[index].timeWithoutTarget >= tuning.fireball.untargetedLifetime {
            removeFireball(at: index)
        }
    }

    func moveFireball(at index: Int, deltaTime: TimeInterval) {
        let movement = fireballs[index].velocity

        fireballs[index].node.position.x += movement.dx * tuning.fireball.speed * CGFloat(deltaTime)
        fireballs[index].node.position.y += movement.dy * tuning.fireball.speed * CGFloat(deltaTime)
        fireballs[index].node.zRotation = atan2(movement.dy, movement.dx)
    }


    func spawnFireballs() {
        let targets = availableSkeletonTargets(limit: session.progression.simultaneousFireballCount)

        for target in targets {
            spawnFireball(targeting: target)
        }
    }

    func spawnFireball(targeting target: SKSpriteNode) {
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

    func syncOrbitalOrbCount() {
        while orbitalOrbs.count < session.progression.orbitalOrbCount {
            let node = makeOrbitalOrbNode()
            worldNode.addChild(node)
            orbitalOrbs.append(OrbitalOrb(node: node))
        }

        while orbitalOrbs.count > session.progression.orbitalOrbCount {
            removeOrbitalOrb(at: orbitalOrbs.count - 1)
        }

        alignOrbitalOrbs()
    }

    func makeOrbitalOrbNode() -> SKSpriteNode {
        let orbNode = SKSpriteNode(texture: orbitalOrbTextures[0])
        orbNode.size = CGSize(width: 20, height: 20)
        orbNode.zPosition = 11

        return orbNode
    }

    func alignOrbitalOrbs() {
        guard !orbitalOrbs.isEmpty else {
            return
        }

        let spacing = CGFloat.pi * 2 / CGFloat(orbitalOrbs.count)
        for (index, orb) in orbitalOrbs.enumerated() {
            orb.updatePosition(
                around: player.position,
                angle: session.orbitalOrbAngle + spacing * CGFloat(index),
                radius: tuning.orbitalOrb.radius
            )
        }
    }

    func removeOrbitalOrb(at index: Int) {
        orbitalOrbs[index].deactivate()
        orbitalOrbs.remove(at: index)
    }

    func respawnOrbitalOrbs(angleDelta: CGFloat) {
        for index in orbitalOrbs.indices where orbitalOrbs[index].updateMissingOrbitProgress(by: angleDelta) {
            let orbNode = makeOrbitalOrbNode()
            worldNode.addChild(orbNode)
            orbitalOrbs[index].attach(orbNode)
        }
    }

    func castLightning() {
        let chainLightning = ChainLightning(
            origin: player.position,
            strikeCount: session.progression.lightningStrikeCount,
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

    func showLightningStrike(from start: CGPoint, to end: CGPoint) {
        let effect = ChainLightning.makeEffectNode(from: start, to: end, texture: lightningTexture)
        worldNode.addChild(effect)
    }

    func showLightningTargetHit(_ target: SKSpriteNode) {
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
                SKAction.fadeOut(withDuration: tuning.lightning.effectDuration),
                SKAction.removeFromParent()
            ])
        )
        worldNode.addChild(hitSprite)
    }

    func castBeam() {
        let beam = Beam(
            origin: player.position,
            direction: playerBeamDirection(),
            length: beamLength,
            hitWidth: tuning.beam.hitWidth,
            killLimit: session.progression.beamKillCount,
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

    func showBeam(from start: CGPoint, to end: CGPoint) {
        worldNode.addChild(Beam.makeEffectNode(from: start, to: end))
    }

    func castMeteor() {
        let meteor = Meteor(
            origin: player.position,
            targetRadius: tuning.orbitalOrb.radius * tuning.meteor.targetRadiusMultiplier
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

    func makeMeteorNode() -> SKSpriteNode {
        let meteorNode = SKSpriteNode(texture: meteorTextures[0])
        meteorNode.name = Meteor.projectileName
        meteorNode.size = CGSize(width: 24, height: 24)
        meteorNode.zPosition = 14

        return meteorNode
    }

    func updateMeteors(deltaTime: TimeInterval) {
        updateMeteorAnimation(deltaTime: deltaTime)

        for index in meteors.indices.reversed() {
            if meteors[index].update(deltaTime: deltaTime) {
                let impactPosition = meteors[index].impactPosition
                meteors[index].node.removeFromParent()
                meteors.remove(at: index)
                impactMeteor(at: impactPosition)

                if session.isLevelUpChoiceActive {
                    return
                }
            }
        }
    }

    func impactMeteor(at position: CGPoint) {
        guard !session.isGameOver else {
            return
        }

        showMeteorImpact(at: position)

        var targets = [SKSpriteNode]()
        skeletonSpatialIndex.forEachCandidate(
            near: position,
            radius: tuning.meteor.impactRadius,
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

    func showMeteorImpact(at position: CGPoint) {
        worldNode.addChild(
            Meteor.makeImpactEffectNode(
                at: position,
                radius: tuning.meteor.impactRadius
            )
        )
    }


    func checkOrbitalOrbCollisions() {
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

    func firstSkeletonTouchingOrb(at index: Int) -> SKSpriteNode? {
        guard let orbNode = orbitalOrbs[index].node else {
            return nil
        }

        return skeletonSpatialIndex.firstCandidate(
            near: orbNode.position,
            radius: tuning.orbitalOrb.hitDistance,
            isValid: isSkeletonAlive,
            matches: { _ in true }
        )
    }


    func removeFireball(at index: Int) {
        fireballs[index].node.removeAllActions()
        fireballs[index].node.removeFromParent()
        fireballs.remove(at: index)
    }

    func firstSkeletonHitByFireball(from start: CGPoint, to end: CGPoint) -> SKSpriteNode? {
        let dx = end.x - start.x
        let dy = end.y - start.y
        let lengthSquared = dx * dx + dy * dy
        var closestHit: (skeleton: SKSpriteNode, progress: CGFloat)?

        let hitRadius = tuning.fireball.hitDistance
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


    func removeLightningEffects() {
        worldNode.enumerateChildNodes(withName: ChainLightning.effectName) { node, _ in
            node.removeAllActions()
            node.removeFromParent()
        }
    }

    func removeBeamEffects() {
        worldNode.enumerateChildNodes(withName: Beam.effectName) { node, _ in
            node.removeAllActions()
            node.removeFromParent()
        }
    }

    func removeMeteorEffects() {
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


    func updateFireballAnimation(deltaTime: TimeInterval) {
        guard !fireballs.isEmpty else {
            session.animations.fireballTimer = 0
            return
        }

        session.animations.fireballTimer += deltaTime

        guard session.animations.fireballTimer >= tuning.fireball.animationFrameDuration else {
            return
        }

        session.animations.fireballTimer.formTruncatingRemainder(dividingBy: tuning.fireball.animationFrameDuration)
        session.animations.fireballFrameIndex = (session.animations.fireballFrameIndex + 1) % fireballTextures.count
        let texture = fireballTextures[session.animations.fireballFrameIndex]

        for fireball in fireballs {
            fireball.node.texture = texture
        }
    }

    func updateOrbitalOrbAnimation(deltaTime: TimeInterval) {
        guard orbitalOrbs.contains(where: { $0.isActive }) else {
            session.animations.orbitalOrbTimer = 0
            return
        }

        session.animations.orbitalOrbTimer += deltaTime

        guard session.animations.orbitalOrbTimer >= tuning.orbitalOrb.animationFrameDuration else {
            return
        }

        session.animations.orbitalOrbTimer.formTruncatingRemainder(dividingBy: tuning.orbitalOrb.animationFrameDuration)
        session.animations.orbitalOrbFrameIndex = (session.animations.orbitalOrbFrameIndex + 1) % orbitalOrbTextures.count
        let texture = orbitalOrbTextures[session.animations.orbitalOrbFrameIndex]

        for orb in orbitalOrbs {
            orb.node?.texture = texture
        }
    }

    func updateMeteorAnimation(deltaTime: TimeInterval) {
        guard !meteors.isEmpty else {
            session.animations.meteorTimer = 0
            return
        }

        session.animations.meteorTimer += deltaTime

        guard session.animations.meteorTimer >= tuning.meteor.animationFrameDuration else {
            return
        }

        session.animations.meteorTimer.formTruncatingRemainder(dividingBy: tuning.meteor.animationFrameDuration)
        session.animations.meteorFrameIndex = (session.animations.meteorFrameIndex + 1) % meteorTextures.count
        let texture = meteorTextures[session.animations.meteorFrameIndex]

        for meteor in meteors {
            meteor.node.texture = texture
        }
    }


    var beamLength: CGFloat {
        max(700, hypot(size.width, size.height) / 2 + tuning.skeleton.spawnMargin)
    }


}
