import SpriteKit

extension GameScene {
    func updateSkeletons(deltaTime: TimeInterval) {
        for skeleton in skeletons {
            let dx = player.position.x - skeleton.position.x
            let dy = player.position.y - skeleton.position.y
            let distanceSquared = dx * dx + dy * dy

            guard distanceSquared > 0 else {
                continue
            }

            let distance = sqrt(distanceSquared)
            let movement = CGVector(dx: dx / distance, dy: dy / distance)
            skeleton.position.x += movement.dx * tuning.skeleton.speed * CGFloat(deltaTime)
            skeleton.position.y += movement.dy * tuning.skeleton.speed * CGFloat(deltaTime)
            updateFacing(for: skeleton, movement: movement)
        }

        updateSkeletonAnimation(deltaTime: deltaTime)
        skeletonSpatialIndex.rebuild(with: skeletons)
    }

    func updateSkeletonSpawning(deltaTime: TimeInterval) {
        session.casts.skeletonSpawn += deltaTime
        let spawnInterval = session.progression.skeletonSpawnInterval
        var didSpawn = false

        while session.casts.skeletonSpawn >= spawnInterval {
            session.casts.skeletonSpawn -= spawnInterval
            spawnSkeleton(kind: timedSkeletonSpawnKind, shouldUpdateHUD: false)
            didSpawn = true
        }

        if didSpawn {
            updateHUDCombatStatus()
        }
    }

    var timedSkeletonSpawnKind: SkeletonKind {
        usesRedOnlySkeletonSpawns ? .red : .regular
    }

    var usesRedOnlySkeletonSpawns: Bool {
        session.progression.level >= tuning.skeleton.redOnlyLevel
    }


    func spawnSkeleton(kind: SkeletonKind = .regular, shouldUpdateHUD: Bool = true) {
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


    func skeletonSpawnPosition() -> CGPoint {
        let halfWidth = size.width / 2
        let halfHeight = size.height / 2
        let spawnDistance = hypot(halfWidth, halfHeight) + tuning.skeleton.spawnMargin
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

    func checkSkeletonCollisions() {
        guard session.playerHitInvulnerabilityTimer <= 0 else {
            return
        }

        if skeletonSpatialIndex.firstCandidate(
            near: player.position,
            radius: tuning.skeleton.hitDistance,
            isValid: isSkeletonAlive,
            matches: { _ in true }
        ) != nil {
            damagePlayer()
        }
    }


    func damageSkeleton(_ skeleton: SKSpriteNode, killedBy attackKind: AttackKind? = nil, shouldTriggerLevelUpChoice: Bool = true, shouldUpdateHUD: Bool = true) -> Int {
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

    func skeletonHitPoints(for skeleton: SKSpriteNode) -> Int {
        guard let hitPoints = skeleton.userData?[SkeletonUserDataKey.hitPoints] as? NSNumber else {
            return 1
        }

        return max(1, hitPoints.intValue)
    }

    func setSkeletonHitPoints(_ hitPoints: Int, for skeleton: SKSpriteNode) {
        if skeleton.userData == nil {
            skeleton.userData = NSMutableDictionary()
        }

        skeleton.userData?[SkeletonUserDataKey.hitPoints] = NSNumber(value: max(1, hitPoints))
    }

    func showSkeletonDamageFeedback(_ skeleton: SKSpriteNode) {
        skeleton.removeAction(forKey: Self.skeletonDamageFlashActionKey)
        skeleton.alpha = 1

        let flash = SKAction.sequence([
            SKAction.fadeAlpha(to: 0.35, duration: 0.06),
            SKAction.fadeAlpha(to: 1, duration: 0.06)
        ])
        skeleton.run(SKAction.repeat(flash, count: 2), withKey: Self.skeletonDamageFlashActionKey)
    }

    @discardableResult
    func destroySkeleton(_ skeleton: SKSpriteNode, killedBy attackKind: AttackKind? = nil, shouldTriggerLevelUpChoice: Bool = true, shouldUpdateHUD: Bool = true) -> Int {
        let identifier = ObjectIdentifier(skeleton)

        guard let index = skeletonIndices[identifier] else {
            return 0
        }

        skeleton.removeAllActions()
        skeleton.removeFromParent()
        removeSkeletonFromTracking(identifier: identifier, at: index)
        registerSkeletonKill()
        registerAttackKill(attackKind)

        let levelUpCount = session.progression.gainExperience()

        if shouldUpdateHUD {
            updateHUDProgress()
        }

        if shouldTriggerLevelUpChoice {
            queueLevelUpChoices(levelUpCount)
        }

        return levelUpCount
    }

    func registerAttackKill(_ attackKind: AttackKind?) {
        switch attackKind {
        case .fireball:
            session.kills.fireball += 1
        case .lightning:
            session.kills.lightning += 1
        case .orbitalOrb:
            session.kills.orbitalOrb += 1
        case .beam:
            session.kills.beam += 1
        case .meteor:
            session.kills.meteor += 1
        case .none:
            break
        }
    }

    func removeSkeletonFromTracking(identifier: ObjectIdentifier, at index: Int) {
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

    func registerSkeletonKill() {
        session.kills.totalSkeletons += 1
        spawnMilestoneSkeletonsIfNeeded()

        while session.kills.totalSkeletons >= session.nextChestMilestone {
            if let tier = chestTier(for: session.nextChestMilestone) {
                spawnChest(tier: tier)
            }
            session.nextChestMilestone += tuning.chest.bronzeKillInterval
        }
    }

    func spawnMilestoneSkeletonsIfNeeded() {
        guard !usesRedOnlySkeletonSpawns else {
            return
        }

        spawnSkeleton(kind: .red, afterEveryKills: tuning.skeleton.redKillInterval)
        spawnSkeleton(kind: .purple, afterEveryKills: tuning.skeleton.purpleKillInterval)
    }

    func spawnSkeleton(kind: SkeletonKind, afterEveryKills killInterval: Int) {
        guard killInterval > 0, session.kills.totalSkeletons.isMultiple(of: killInterval) else {
            return
        }

        spawnSkeleton(kind: kind, shouldUpdateHUD: false)
    }


    func killAllEnemiesAndGrantExperience() {
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

        let levelUpCount = session.progression.gainExperience(defeatedCount)
        updateHUDProgress()
        queueLevelUpChoices(levelUpCount)
    }


    func availableSkeletonTargets(limit: Int) -> [SKSpriteNode] {
        let reservedTargets = fireballTargetIdentifiers()
        return closestSkeletons(to: player.position, excluding: reservedTargets, limit: limit)
    }

    func availableLightningTargets() -> [SKSpriteNode] {
        let reservedTargets = fireballTargetIdentifiers()
        return skeletons.filter { skeleton in
            isSkeletonAlive(skeleton) && !reservedTargets.contains(ObjectIdentifier(skeleton))
        }
    }

    func fireballTargetIdentifiers() -> Set<ObjectIdentifier> {
        Set(
            fireballs.compactMap { fireball in
                fireball.target.map(ObjectIdentifier.init)
            }
        )
    }

    func closestSkeletons(to position: CGPoint, excluding excludedIdentifiers: Set<ObjectIdentifier>, limit: Int) -> [SKSpriteNode] {
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


    func updateSkeletonAnimation(deltaTime: TimeInterval) {
        guard !skeletons.isEmpty else {
            session.animations.skeletonTimer = 0
            return
        }

        session.animations.skeletonTimer += deltaTime

        guard session.animations.skeletonTimer >= tuning.skeleton.animationFrameDuration else {
            return
        }

        session.animations.skeletonTimer.formTruncatingRemainder(dividingBy: tuning.skeleton.animationFrameDuration)
        session.animations.skeletonFrameIndex = (session.animations.skeletonFrameIndex + 1) % skeletonTextures.count
        let texture = skeletonTextures[session.animations.skeletonFrameIndex]

        for skeleton in skeletons {
            skeleton.texture = texture
        }
    }


    func isSkeletonAlive(_ skeleton: SKSpriteNode) -> Bool {
        skeletonIdentifiers.contains(ObjectIdentifier(skeleton))
    }


}
