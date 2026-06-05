//
//  Progression.swift
//  C64-Returnal
//

import CoreGraphics
import Foundation

struct Progression {
    private(set) var level = 1
    private(set) var experience = 0
    private(set) var nextExperience = 1
    private(set) var simultaneousFireballCount = 1
    private(set) var lightningBounceCount = 0
    private(set) var orbitalOrbCount = 1
    private(set) var beamKillCount = 1

    private var previousFibonacciExperience = 1
    private var currentFibonacciExperience = 1
    private var fireRateUpgradeCount = 0
    private var lightningRateUpgradeCount = 0
    private var orbitalOrbSpeedUpgradeCount = 0
    private var beamRateUpgradeCount = 0

    var skeletonSpawnInterval: TimeInterval {
        GameConfiguration.initialSkeletonSpawnInterval
            * pow(GameConfiguration.skeletonIntervalMultiplierPerLevel, Double(level - 1))
    }

    var fireballCastInterval: TimeInterval {
        GameConfiguration.initialFireballCastInterval
            * pow(GameConfiguration.fireballIntervalMultiplierPerUpgrade, Double(fireRateUpgradeCount))
    }

    var lightningCastInterval: TimeInterval {
        GameConfiguration.initialLightningCastInterval
            * pow(GameConfiguration.lightningIntervalMultiplierPerUpgrade, Double(lightningRateUpgradeCount))
    }

    var lightningStrikeCount: Int {
        lightningBounceCount + 1
    }

    var orbitalOrbAngularSpeed: CGFloat {
        GameConfiguration.initialOrbitalOrbAngularSpeed
            * pow(GameConfiguration.orbitalOrbSpeedMultiplierPerUpgrade, CGFloat(orbitalOrbSpeedUpgradeCount))
    }

    var beamCastInterval: TimeInterval {
        GameConfiguration.initialBeamCastInterval
            * pow(GameConfiguration.beamIntervalMultiplierPerUpgrade, Double(beamRateUpgradeCount))
    }

    var maximumSkeletons: Int {
        (GameConfiguration.baseMaximumSkeletons + level - 1) * 10
    }

    mutating func reset() {
        level = 1
        experience = 0
        previousFibonacciExperience = 1
        currentFibonacciExperience = 1
        nextExperience = Self.scaledFibonacciRequirement(currentFibonacciExperience)
        fireRateUpgradeCount = 0
        lightningRateUpgradeCount = 0
        orbitalOrbSpeedUpgradeCount = 0
        beamRateUpgradeCount = 0
        simultaneousFireballCount = 1
        lightningBounceCount = 0
        orbitalOrbCount = 1
        beamKillCount = 1
    }

    @discardableResult
    mutating func gainExperience() -> Bool {
        experience += 1

        guard experience >= nextExperience else {
            return false
        }

        experience -= nextExperience
        levelUp()
        return true
    }

    mutating func applyLevelUpOption(_ option: LevelUpOption) {
        switch option {
        case .fireRate:
            fireRateUpgradeCount += 1
        case .extraFireball:
            simultaneousFireballCount += 1
        case .lightningBounce:
            lightningBounceCount += 1
        case .lightningRate:
            lightningRateUpgradeCount += 1
        case .extraOrb:
            orbitalOrbCount += 1
        case .orbitalSpeed:
            orbitalOrbSpeedUpgradeCount += 1
        case .beamRate:
            beamRateUpgradeCount += 1
        case .beamKillCount:
            beamKillCount += level
        }
    }

    private mutating func levelUp() {
        level += 1

        let newThreshold = previousFibonacciExperience + currentFibonacciExperience
        previousFibonacciExperience = currentFibonacciExperience
        currentFibonacciExperience = newThreshold
        nextExperience = Self.scaledFibonacciRequirement(currentFibonacciExperience)
    }

    private static func scaledFibonacciRequirement(_ fibonacciValue: Int) -> Int {
        Int(ceil(Double(fibonacciValue) / 2))
    }
}
