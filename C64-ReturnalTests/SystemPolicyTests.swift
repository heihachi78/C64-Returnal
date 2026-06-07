import XCTest

final class SystemPolicyTests: XCTestCase {
    func testInputControllerMapsMovementAndModalKeys() {
        let bindings = InputBindings(
            moveLeft: 1,
            moveRight: 2,
            moveDown: 3,
            moveUp: 4,
            firstLevelUpOption: 10,
            secondLevelUpOption: 11,
            thirdLevelUpOption: 12,
            fourthLevelUpOption: 13,
            advanceChestReward: 99,
            killAllAndGrantExperience: [42]
        )
        let controller = InputController(bindings: bindings)

        XCTAssertTrue(controller.isMovementKey(1))
        XCTAssertTrue(controller.isMovementKey(4))
        XCTAssertEqual(controller.levelUpOptionIndex(for: 11), 1)
        XCTAssertEqual(controller.levelUpOptionIndex(for: 13), 3)
        XCTAssertTrue(controller.isChestRewardAdvance(99))
        XCTAssertTrue(controller.isKillAllAndGrantExperience(42))
        XCTAssertFalse(controller.isMovementKey(99))
    }

    func testLevelUpChoicesDefaultToThreeAndChanceAddsFourth() {
        let alwaysThree = progressionSystem(extraOptionNumerator: 0, extraOptionDenominator: 1)
        let alwaysFour = progressionSystem(extraOptionNumerator: 1, extraOptionDenominator: 1)
        let options: [LevelUpOption] = [.fireRate, .extraFireball, .extraLife, .learnLightning]

        XCTAssertEqual(
            alwaysThree.randomLevelUpOptions(from: options, hasSkeletons: false).count,
            3
        )
        XCTAssertEqual(
            alwaysFour.randomLevelUpOptions(from: options, hasSkeletons: false).count,
            4
        )
    }

    func testChestSystemChoosesTierAndRespectsLevelCaps() {
        let system = ChestSystem(tuning: GameConfiguration.defaultTuning)

        XCTAssertEqual(system.tier(for: 250, playerLevel: 1), .bronze)
        XCTAssertEqual(system.tier(for: 1000, playerLevel: 1), .silver)
        XCTAssertEqual(system.tier(for: 5000, playerLevel: 90), .gold)
        XCTAssertNil(system.tier(for: 250, playerLevel: 90))
        XCTAssertNil(system.tier(for: 1000, playerLevel: 90))
    }

    func testChestRewardItemsUseBeamUpgradeBonusForBeamOnly() {
        let items = ChestSystem(tuning: GameConfiguration.defaultTuning).rewardItems(
            for: [.fireball, .beam],
            beamKillUpgradeBonus: 4
        )

        XCTAssertTrue(items.contains { item in
            if case .extraFireball = item.option {
                return item.title == "+1 FIREBALL"
            }

            return false
        })
        XCTAssertTrue(items.contains { item in
            if case .beamKillCount = item.option {
                return item.title == "+4 BEAM KILL"
            }

            return false
        })
    }

    private func progressionSystem(
        extraOptionNumerator: Int,
        extraOptionDenominator: Int
    ) -> ProgressionSystem {
        let tuning = GameConfiguration.defaultTuning

        return ProgressionSystem(
            tuning: GameTuning(
                presentation: tuning.presentation,
                player: tuning.player,
                skeleton: tuning.skeleton,
                fireball: tuning.fireball,
                lightning: tuning.lightning,
                orbitalOrb: tuning.orbitalOrb,
                beam: tuning.beam,
                meteor: tuning.meteor,
                chest: tuning.chest,
                progression: GameTuning.Progression(
                    halveHordeChanceNumerator: tuning.progression.halveHordeChanceNumerator,
                    halveHordeChanceDenominator: tuning.progression.halveHordeChanceDenominator,
                    extraLevelUpOptionChanceNumerator: extraOptionNumerator,
                    extraLevelUpOptionChanceDenominator: extraOptionDenominator
                )
            )
        )
    }
}
