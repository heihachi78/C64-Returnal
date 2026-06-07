import XCTest

final class ProgressionTests: XCTestCase {
    func testExperienceCanAdvanceMultipleLevels() {
        var progression = Progression()

        let levelUpCount = progression.gainExperience(8)

        XCTAssertEqual(levelUpCount, 3)
        XCTAssertEqual(progression.level, 4)
        XCTAssertEqual(progression.experience, 1)
        XCTAssertEqual(progression.nextExperience, 8)
    }

    func testUnlockAndUpgradeBeamChangesKillCountAndInterval() {
        var progression = Progression()

        progression.applyLevelUpOption(.learnBeam)
        progression.applyLevelUpOption(.beamKillCount)
        progression.applyLevelUpOption(.beamRate)

        XCTAssertTrue(progression.isBeamUnlocked)
        XCTAssertEqual(progression.beamKillCount, 3)
        XCTAssertLessThan(progression.beamCastInterval, GameConfiguration.defaultTuning.beam.initialCastInterval)
    }
}
