import SpriteKit

extension GameScene {
    func triggerGameOver() {
        guard !session.isGameOver else {
            return
        }

        session.isGameOver = true
        session.isLevelUpChoiceActive = false
        session.isChestRewardActive = false
        session.playerHitInvulnerabilityTimer = 0
        session.pendingLevelUpLevels.removeAll(keepingCapacity: true)
        worldNode.isPaused = false
        session.pressedKeys.removeAll()
        hud.hideLevelUp()
        hud.hideChestReward()

        player.removeAction(forKey: Self.playerHitFlashActionKey)
        stopPlayerAnimation()
        player.color = SKColor(calibratedRed: 0.85, green: 0.05, blue: 0.08, alpha: 1)
        player.colorBlendFactor = 0.65
        player.alpha = 0.45
        player.run(SKAction.rotate(byAngle: -.pi / 2, duration: 0.16))

        hud.showGameOver(level: session.progression.level)

        for fireball in fireballs {
            fireball.node.removeAllActions()
        }
        removeLightningEffects()
        removeBeamEffects()
        removeMeteorEffects()
    }

    func restartGame() {
        session.reset()
        worldNode.isPaused = false

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


    func removeAllEnemiesAndProjectiles() {
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

}
