//
//  GameConfiguration.swift
//  C64-Returnal
//

import SpriteKit

enum GameConfiguration {
    static let backgroundColor = SKColor(calibratedRed: 0.18, green: 0.37, blue: 0.16, alpha: 1)

    static let tileSize: CGFloat = 64
    static let playerSpeed: CGFloat = 190
    static let initialPlayerLives = 3
    static let playerHitInvulnerabilityDuration: TimeInterval = 1.0
    static let playerAnimationFrameDuration: TimeInterval = 0.18

    static let skeletonSpeed: CGFloat = 82
    static let initialSkeletonSpawnInterval: TimeInterval = 0.85
    static let skeletonIntervalMultiplierPerLevel = 0.9
    static let skeletonHitDistance: CGFloat = 24
    static let skeletonSpawnMargin: CGFloat = 72
    static let skeletonAnimationFrameDuration: TimeInterval = 0.20

    static let fireballSpeed: CGFloat = 280
    static let initialFireballCastInterval: TimeInterval = 3.0
    static let fireballIntervalMultiplierPerUpgrade = 0.9
    static let fireballHitDistance: CGFloat = 20
    static let fireballUntargetedLifetime: TimeInterval = 3.0
    static let fireballAnimationFrameDuration: TimeInterval = 0.08

    static let initialLightningCastInterval: TimeInterval = 3.0
    static let lightningIntervalMultiplierPerUpgrade = 0.9
    static let lightningEffectDuration: TimeInterval = 0.18
    static let lightningBranchWidth: CGFloat = 3

    static let orbitalOrbRadius: CGFloat = 58
    static let initialOrbitalOrbAngularSpeed: CGFloat = 2.4
    static let orbitalOrbSpeedMultiplierPerUpgrade: CGFloat = 1.2
    static let orbitalOrbHitDistance: CGFloat = 22
    static let orbitalOrbAnimationFrameDuration: TimeInterval = 0.12

    static let initialBeamCastInterval: TimeInterval = 3.0
    static let beamIntervalMultiplierPerUpgrade = 0.9
    static let beamHitWidth: CGFloat = 18
    static let beamEffectDuration: TimeInterval = 0.16

    static let initialMeteorCastInterval: TimeInterval = 3.0
    static let meteorIntervalMultiplierPerUpgrade = 0.95
    static let meteorTargetRadiusMultiplier: CGFloat = 8
    static let meteorImpactRadius: CGFloat = 24
    static let meteorFallDuration: TimeInterval = 0.55
    static let meteorFallHeight: CGFloat = 240
    static let meteorFallDrift: CGFloat = 90
    static let meteorAnimationFrameDuration: TimeInterval = 0.14

    static let chestSpawnMargin: CGFloat = 88
    static let chestPickupDistance: CGFloat = 28
    static let bronzeChestKillInterval = 250
    static let silverChestKillInterval = 1000
    static let goldChestKillInterval = 5_000

    static let thirdLevelUpOptionChanceNumerator = 5
    static let thirdLevelUpOptionChanceDenominator = 100
}

enum ArrowKey {
    static let left: UInt16 = 123
    static let right: UInt16 = 124
    static let down: UInt16 = 125
    static let up: UInt16 = 126

    static func contains(_ keyCode: UInt16) -> Bool {
        keyCode == left || keyCode == right || keyCode == down || keyCode == up
    }
}

enum DebugKey {
    static let one: UInt16 = 18
    static let keypadOne: UInt16 = 83

    static func isLevelSetupShortcut(_ keyCode: UInt16) -> Bool {
        keyCode == one || keyCode == keypadOne
    }
}

enum LevelUpSelectionKey {
    static let firstOption: UInt16 = 12
    static let secondOption: UInt16 = 0
    static let thirdOption: UInt16 = 16

    static func optionIndex(for keyCode: UInt16) -> Int? {
        switch keyCode {
        case firstOption:
            return 0
        case secondOption:
            return 1
        case thirdOption:
            return 2
        default:
            return nil
        }
    }
}

enum ChestRewardKey {
    static let advance: UInt16 = 12

    static func isAdvance(_ keyCode: UInt16) -> Bool {
        keyCode == advance
    }
}
