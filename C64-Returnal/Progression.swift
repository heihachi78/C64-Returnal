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
    private(set) var isLightningUnlocked = false
    private(set) var isOrbitalOrbUnlocked = false
    private(set) var isBeamUnlocked = false
    private(set) var isMeteorUnlocked = false

    private var fireRateUpgradeCount = 0
    private var lightningRateUpgradeCount = 0
    private var orbitalOrbSpeedUpgradeCount = 0
    private var beamRateUpgradeCount = 0
    private var meteorRateUpgradeCount = 0
    private var upgradedOrbitalOrbCount = 1
    private var beamKillLevel = 1
    private var upgradedBeamKillCount = 1
    private var upgradedMeteorCount = 1

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
        isLightningUnlocked ? lightningBounceCount + 1 : 0
    }

    var orbitalOrbAngularSpeed: CGFloat {
        GameConfiguration.initialOrbitalOrbAngularSpeed
            * pow(GameConfiguration.orbitalOrbSpeedMultiplierPerUpgrade, CGFloat(orbitalOrbSpeedUpgradeCount))
    }

    var orbitalOrbCount: Int {
        isOrbitalOrbUnlocked ? upgradedOrbitalOrbCount : 0
    }

    var beamCastInterval: TimeInterval {
        GameConfiguration.initialBeamCastInterval
            * pow(GameConfiguration.beamIntervalMultiplierPerUpgrade, Double(beamRateUpgradeCount))
    }

    var beamKillCount: Int {
        isBeamUnlocked ? upgradedBeamKillCount : 0
    }

    var beamKillUpgradeBonus: Int {
        beamKillLevel + 1
    }

    var meteorCastInterval: TimeInterval {
        GameConfiguration.initialMeteorCastInterval
            * pow(GameConfiguration.meteorIntervalMultiplierPerUpgrade, Double(meteorRateUpgradeCount))
    }

    var meteorCount: Int {
        isMeteorUnlocked ? upgradedMeteorCount : 0
    }

    var maximumSkeletons: Int {
        (GameConfiguration.baseMaximumSkeletons + level - 1) * 10
    }

    var availableLevelUpOptions: [LevelUpOption] {
        var options: [LevelUpOption] = [.fireRate, .extraFireball, .extraLife]

        if isLightningUnlocked {
            options.append(contentsOf: [.lightningBounce, .lightningRate])
        } else {
            options.append(.learnLightning)
        }

        if isOrbitalOrbUnlocked {
            options.append(contentsOf: [.extraOrb, .orbitalSpeed])
        } else {
            options.append(.learnOrb)
        }

        if isBeamUnlocked {
            options.append(contentsOf: [.beamRate, .beamKillCount])
        } else {
            options.append(.learnBeam)
        }

        if isMeteorUnlocked {
            options.append(contentsOf: [.extraMeteor, .meteorRate])
        } else {
            options.append(.learnMeteor)
        }

        return options
    }

    mutating func reset() {
        level = 1
        experience = 0
        nextExperience = Self.experienceRequirement(for: level)
        fireRateUpgradeCount = 0
        lightningRateUpgradeCount = 0
        orbitalOrbSpeedUpgradeCount = 0
        beamRateUpgradeCount = 0
        meteorRateUpgradeCount = 0
        simultaneousFireballCount = 1
        lightningBounceCount = 0
        upgradedOrbitalOrbCount = 1
        beamKillLevel = 1
        upgradedBeamKillCount = 1
        upgradedMeteorCount = 1
        isLightningUnlocked = false
        isOrbitalOrbUnlocked = false
        isBeamUnlocked = false
        isMeteorUnlocked = false
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

    mutating func advanceToOneKillBeforeNextLevel() {
        experience = max(experience, nextExperience - 1)
    }

    mutating func applyLevelUpOption(_ option: LevelUpOption) {
        switch option {
        case .fireRate:
            fireRateUpgradeCount += 1
        case .extraFireball:
            simultaneousFireballCount += 1
        case .extraLife:
            break
        case .learnLightning:
            isLightningUnlocked = true
        case .lightningBounce:
            lightningBounceCount += 1
        case .lightningRate:
            lightningRateUpgradeCount += 1
        case .learnOrb:
            isOrbitalOrbUnlocked = true
        case .extraOrb:
            upgradedOrbitalOrbCount += 1
        case .orbitalSpeed:
            orbitalOrbSpeedUpgradeCount += 1
        case .learnBeam:
            isBeamUnlocked = true
        case .beamRate:
            beamRateUpgradeCount += 1
        case .beamKillCount:
            beamKillLevel += 1
            upgradedBeamKillCount += beamKillLevel
        case .learnMeteor:
            isMeteorUnlocked = true
        case .extraMeteor:
            upgradedMeteorCount += 1
        case .meteorRate:
            meteorRateUpgradeCount += 1
        }
    }

    private mutating func levelUp() {
        level += 1
        nextExperience = Self.experienceRequirement(for: level)
    }

    private static func experienceRequirement(for level: Int) -> Int {
        level * level
    }
}
