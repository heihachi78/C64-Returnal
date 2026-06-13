//
//  GameConfiguration.swift
//  C64-Returnal
//

import SpriteKit

struct GameTuning {
    struct Presentation {
        let backgroundColor: SKColor
        let tileSize: CGFloat
    }

    struct Player {
        let speed: CGFloat
        let initialLives: Int
        let hitInvulnerabilityDuration: TimeInterval
        let animationFrameDuration: TimeInterval
    }

    struct Skeleton {
        let speed: CGFloat
        let initialSpawnInterval: TimeInterval
        let intervalMultiplierPerLevel: Double
        let hitDistance: CGFloat
        let spawnMargin: CGFloat
        let animationFrameDuration: TimeInterval
        let redOnlyLevel: Int
        let redOnlySpawnIntervalMultiplier: TimeInterval
        let purpleOnlyLevel: Int
        let purpleOnlySpawnIntervalMultiplier: TimeInterval
        let redHitPoints: Int
        let redKillInterval: Int
        let purpleHitPoints: Int
        let purpleKillInterval: Int
        let blackHitPoints: Int
        let blackPurpleKillInterval: Int
        let spatialIndexCellSize: CGFloat
    }

    struct Fireball {
        let speed: CGFloat
        let initialCastInterval: TimeInterval
        let intervalMultiplierPerUpgrade: Double
        let hitDistance: CGFloat
        let untargetedLifetime: TimeInterval
        let animationFrameDuration: TimeInterval
    }

    struct Lightning {
        let initialCastInterval: TimeInterval
        let intervalMultiplierPerUpgrade: Double
        let effectDuration: TimeInterval
        let branchWidth: CGFloat
    }

    struct OrbitalOrb {
        let radius: CGFloat
        let initialAngularSpeed: CGFloat
        let speedMultiplierPerUpgrade: CGFloat
        let hitDistance: CGFloat
        let animationFrameDuration: TimeInterval
    }

    struct Beam {
        let initialCastInterval: TimeInterval
        let intervalMultiplierPerUpgrade: Double
        let hitWidth: CGFloat
        let effectDuration: TimeInterval
    }

    struct Meteor {
        let initialCastInterval: TimeInterval
        let intervalMultiplierPerUpgrade: Double
        let targetRadiusMultiplier: CGFloat
        let impactRadius: CGFloat
        let fallDuration: TimeInterval
        let fallHeight: CGFloat
        let fallDrift: CGFloat
        let animationFrameDuration: TimeInterval
    }

    struct Chest {
        let spawnMargin: CGFloat
        let pickupDistance: CGFloat
        let bronzeKillInterval: Int
        let silverKillInterval: Int
        let goldKillInterval: Int
        let bronzeMaximumLevel: Int
        let silverMaximumLevel: Int
    }

    struct Coin {
        let spawnMargin: CGFloat
        let pickupDistance: CGFloat
        let minimumReward: Int
        let maximumReward: Int
        let animationFrameDuration: TimeInterval
    }

    struct Progression {
        let halveHordeChanceNumerator: Int
        let halveHordeChanceDenominator: Int
        let extraLevelUpOptionChanceNumerator: Int
        let extraLevelUpOptionChanceDenominator: Int
    }

    let presentation: Presentation
    let player: Player
    let skeleton: Skeleton
    let fireball: Fireball
    let lightning: Lightning
    let orbitalOrb: OrbitalOrb
    let beam: Beam
    let meteor: Meteor
    let chest: Chest
    let coin: Coin
    let progression: Progression
}

enum GameConfiguration {
    static let defaultTuning = GameTuning(
        presentation: GameTuning.Presentation(
            backgroundColor: SKColor(calibratedRed: 0.18, green: 0.37, blue: 0.16, alpha: 1),
            tileSize: 64
        ),
        player: GameTuning.Player(
            speed: 190,
            initialLives: 3,
            hitInvulnerabilityDuration: 1.0,
            animationFrameDuration: 0.18
        ),
        skeleton: GameTuning.Skeleton(
            speed: 82,
            initialSpawnInterval: 0.91,
            intervalMultiplierPerLevel: 0.915,
            hitDistance: 24,
            spawnMargin: 72,
            animationFrameDuration: 0.20,
            redOnlyLevel: 66,
            redOnlySpawnIntervalMultiplier: 3.0,
            purpleOnlyLevel: 75,
            purpleOnlySpawnIntervalMultiplier: 6.0,
            redHitPoints: 2,
            redKillInterval: 100,
            purpleHitPoints: 5,
            purpleKillInterval: 500,
            blackHitPoints: 25,
            blackPurpleKillInterval: 100,
            spatialIndexCellSize: 96
        ),
        fireball: GameTuning.Fireball(
            speed: 280,
            initialCastInterval: 3.0,
            intervalMultiplierPerUpgrade: 0.9,
            hitDistance: 20,
            untargetedLifetime: 3.0,
            animationFrameDuration: 0.08
        ),
        lightning: GameTuning.Lightning(
            initialCastInterval: 3.0,
            intervalMultiplierPerUpgrade: 0.9,
            effectDuration: 0.18,
            branchWidth: 3
        ),
        orbitalOrb: GameTuning.OrbitalOrb(
            radius: 58,
            initialAngularSpeed: 2.4,
            speedMultiplierPerUpgrade: 1.2,
            hitDistance: 22,
            animationFrameDuration: 0.12
        ),
        beam: GameTuning.Beam(
            initialCastInterval: 3.0,
            intervalMultiplierPerUpgrade: 0.9,
            hitWidth: 18,
            effectDuration: 0.16
        ),
        meteor: GameTuning.Meteor(
            initialCastInterval: 3.0,
            intervalMultiplierPerUpgrade: 0.95,
            targetRadiusMultiplier: 8,
            impactRadius: 48,
            fallDuration: 0.55,
            fallHeight: 240,
            fallDrift: 90,
            animationFrameDuration: 0.14
        ),
        chest: GameTuning.Chest(
            spawnMargin: 88,
            pickupDistance: 28,
            bronzeKillInterval: 250,
            silverKillInterval: 1000,
            goldKillInterval: 5_000,
            bronzeMaximumLevel: 33,
            silverMaximumLevel: 55
        ),
        coin: GameTuning.Coin(
            spawnMargin: 240,
            pickupDistance: 30,
            minimumReward: 1,
            maximumReward: 100,
            animationFrameDuration: 0.14
        ),
        progression: GameTuning.Progression(
            halveHordeChanceNumerator: 5,
            halveHordeChanceDenominator: 100,
            extraLevelUpOptionChanceNumerator: 25,
            extraLevelUpOptionChanceDenominator: 100
        )
    )

    static let defaultInputBindings = InputBindings()

    static let backgroundColor = defaultTuning.presentation.backgroundColor

    static let tileSize = defaultTuning.presentation.tileSize
    static let playerSpeed = defaultTuning.player.speed
    static let initialPlayerLives = defaultTuning.player.initialLives
    static let playerHitInvulnerabilityDuration = defaultTuning.player.hitInvulnerabilityDuration
    static let playerAnimationFrameDuration = defaultTuning.player.animationFrameDuration

    static let skeletonSpeed = defaultTuning.skeleton.speed
    static let initialSkeletonSpawnInterval = defaultTuning.skeleton.initialSpawnInterval
    static let skeletonIntervalMultiplierPerLevel = defaultTuning.skeleton.intervalMultiplierPerLevel
    static let skeletonHitDistance = defaultTuning.skeleton.hitDistance
    static let skeletonSpawnMargin = defaultTuning.skeleton.spawnMargin
    static let skeletonAnimationFrameDuration = defaultTuning.skeleton.animationFrameDuration
    static let redOnlySkeletonLevel = defaultTuning.skeleton.redOnlyLevel
    static let redOnlySkeletonSpawnIntervalMultiplier = defaultTuning.skeleton.redOnlySpawnIntervalMultiplier
    static let purpleOnlySkeletonLevel = defaultTuning.skeleton.purpleOnlyLevel
    static let purpleOnlySkeletonSpawnIntervalMultiplier = defaultTuning.skeleton.purpleOnlySpawnIntervalMultiplier
    static let redSkeletonHitPoints = defaultTuning.skeleton.redHitPoints
    static let redSkeletonKillInterval = defaultTuning.skeleton.redKillInterval
    static let purpleSkeletonHitPoints = defaultTuning.skeleton.purpleHitPoints
    static let purpleSkeletonKillInterval = defaultTuning.skeleton.purpleKillInterval
    static let blackSkeletonHitPoints = defaultTuning.skeleton.blackHitPoints
    static let blackSkeletonPurpleKillInterval = defaultTuning.skeleton.blackPurpleKillInterval

    static let fireballSpeed = defaultTuning.fireball.speed
    static let initialFireballCastInterval = defaultTuning.fireball.initialCastInterval
    static let fireballIntervalMultiplierPerUpgrade = defaultTuning.fireball.intervalMultiplierPerUpgrade
    static let fireballHitDistance = defaultTuning.fireball.hitDistance
    static let fireballUntargetedLifetime = defaultTuning.fireball.untargetedLifetime
    static let fireballAnimationFrameDuration = defaultTuning.fireball.animationFrameDuration

    static let initialLightningCastInterval = defaultTuning.lightning.initialCastInterval
    static let lightningIntervalMultiplierPerUpgrade = defaultTuning.lightning.intervalMultiplierPerUpgrade
    static let lightningEffectDuration = defaultTuning.lightning.effectDuration
    static let lightningBranchWidth = defaultTuning.lightning.branchWidth

    static let orbitalOrbRadius = defaultTuning.orbitalOrb.radius
    static let initialOrbitalOrbAngularSpeed = defaultTuning.orbitalOrb.initialAngularSpeed
    static let orbitalOrbSpeedMultiplierPerUpgrade = defaultTuning.orbitalOrb.speedMultiplierPerUpgrade
    static let orbitalOrbHitDistance = defaultTuning.orbitalOrb.hitDistance
    static let orbitalOrbAnimationFrameDuration = defaultTuning.orbitalOrb.animationFrameDuration

    static let initialBeamCastInterval = defaultTuning.beam.initialCastInterval
    static let beamIntervalMultiplierPerUpgrade = defaultTuning.beam.intervalMultiplierPerUpgrade
    static let beamHitWidth = defaultTuning.beam.hitWidth
    static let beamEffectDuration = defaultTuning.beam.effectDuration

    static let initialMeteorCastInterval = defaultTuning.meteor.initialCastInterval
    static let meteorIntervalMultiplierPerUpgrade = defaultTuning.meteor.intervalMultiplierPerUpgrade
    static let meteorTargetRadiusMultiplier = defaultTuning.meteor.targetRadiusMultiplier
    static let meteorImpactRadius = defaultTuning.meteor.impactRadius
    static let meteorFallDuration = defaultTuning.meteor.fallDuration
    static let meteorFallHeight = defaultTuning.meteor.fallHeight
    static let meteorFallDrift = defaultTuning.meteor.fallDrift
    static let meteorAnimationFrameDuration = defaultTuning.meteor.animationFrameDuration

    static let chestSpawnMargin = defaultTuning.chest.spawnMargin
    static let chestPickupDistance = defaultTuning.chest.pickupDistance
    static let bronzeChestKillInterval = defaultTuning.chest.bronzeKillInterval
    static let silverChestKillInterval = defaultTuning.chest.silverKillInterval
    static let goldChestKillInterval = defaultTuning.chest.goldKillInterval
    static let bronzeChestMaximumLevel = defaultTuning.chest.bronzeMaximumLevel
    static let silverChestMaximumLevel = defaultTuning.chest.silverMaximumLevel

    static let coinSpawnMargin = defaultTuning.coin.spawnMargin
    static let coinPickupDistance = defaultTuning.coin.pickupDistance
    static let minimumCoinReward = defaultTuning.coin.minimumReward
    static let maximumCoinReward = defaultTuning.coin.maximumReward
    static let coinAnimationFrameDuration = defaultTuning.coin.animationFrameDuration

    static let halveHordeLevelUpOptionChanceNumerator = defaultTuning.progression.halveHordeChanceNumerator
    static let halveHordeLevelUpOptionChanceDenominator = defaultTuning.progression.halveHordeChanceDenominator
    static let extraLevelUpOptionChanceNumerator = defaultTuning.progression.extraLevelUpOptionChanceNumerator
    static let extraLevelUpOptionChanceDenominator = defaultTuning.progression.extraLevelUpOptionChanceDenominator
}
