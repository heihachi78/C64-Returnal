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

    static let skeletonSpeed: CGFloat = 82
    static let initialSkeletonSpawnInterval: TimeInterval = 0.4
    static let skeletonIntervalMultiplierPerLevel = 0.95
    static let skeletonHitDistance: CGFloat = 24
    static let skeletonSpawnMargin: CGFloat = 72
    static let baseMaximumSkeletons = 16

    static let fireballSpeed: CGFloat = 280
    static let initialFireballCastInterval: TimeInterval = 3.0
    static let fireballIntervalMultiplierPerUpgrade = 0.85
    static let fireballHitDistance: CGFloat = 20
    static let fireballUntargetedLifetime: TimeInterval = 3.0
    static let fireballAnimationFrameDuration: TimeInterval = 0.08

    static let initialLightningCastInterval: TimeInterval = 3.0
    static let lightningIntervalMultiplierPerUpgrade = 0.85
    static let lightningEffectDuration: TimeInterval = 0.18
    static let lightningBranchWidth: CGFloat = 3

    static let orbitalOrbRadius: CGFloat = 58
    static let initialOrbitalOrbAngularSpeed: CGFloat = 2.4
    static let orbitalOrbSpeedMultiplierPerUpgrade: CGFloat = 1.2
    static let orbitalOrbHitDistance: CGFloat = 22
    static let orbitalOrbAnimationFrameDuration: TimeInterval = 0.12

    static let initialBeamCastInterval: TimeInterval = 3.0
    static let beamIntervalMultiplierPerUpgrade = 0.8
    static let beamHitWidth: CGFloat = 18
    static let beamEffectDuration: TimeInterval = 0.16
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
