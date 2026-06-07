import SpriteKit

extension GameScene {
    func updateHUDProgress() {
        hud.updateProgress(
            level: session.progression.level,
            experience: session.progression.experience,
            nextExperience: session.progression.nextExperience
        )
        hud.updateLives(session.playerLives)
        hud.updateFireballStatus(
            count: session.progression.simultaneousFireballCount,
            interval: session.progression.fireballCastInterval
        )
        hud.updateLightningStatus(
            isUnlocked: session.progression.isLightningUnlocked,
            strikeCount: session.progression.lightningStrikeCount,
            interval: session.progression.lightningCastInterval
        )
        hud.updateOrbStatus(
            isUnlocked: session.progression.isOrbitalOrbUnlocked,
            count: session.progression.orbitalOrbCount,
            angularSpeed: session.progression.orbitalOrbAngularSpeed
        )
        hud.updateBeamStatus(
            isUnlocked: session.progression.isBeamUnlocked,
            killCount: session.progression.beamKillCount,
            interval: session.progression.beamCastInterval
        )
        hud.updateMeteorStatus(
            isUnlocked: session.progression.isMeteorUnlocked,
            count: session.progression.meteorCount,
            interval: session.progression.meteorCastInterval
        )
        hud.updateAttackKillCounts(
            fireball: session.kills.fireball,
            lightning: session.kills.lightning,
            orb: session.kills.orbitalOrb,
            beam: session.kills.beam,
            meteor: session.kills.meteor
        )
        updateHUDCombatStatus()
    }

    func updateHUDCombatStatus() {
        hud.updateSkeletonStatus(
            aliveCount: skeletons.count,
            spawnInterval: session.progression.skeletonSpawnInterval
        )
    }


}
