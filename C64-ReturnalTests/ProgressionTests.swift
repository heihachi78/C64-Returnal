import SpriteKit
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

    func testSkeletonSpawnIntervalAddsLevelThresholdMultipliers() {
        let tuning = GameConfiguration.defaultTuning
        var progression = Progression(tuning: tuning)

        progression.advance(to: tuning.skeleton.redOnlyLevel - 1)
        XCTAssertEqual(
            progression.skeletonSpawnInterval,
            expectedSkeletonSpawnInterval(at: tuning.skeleton.redOnlyLevel - 1, tuning: tuning),
            accuracy: 0.0001
        )

        progression.advance(to: tuning.skeleton.redOnlyLevel)
        XCTAssertEqual(
            progression.skeletonSpawnInterval,
            expectedSkeletonSpawnInterval(
                at: tuning.skeleton.redOnlyLevel,
                tuning: tuning,
                multiplier: tuning.skeleton.redOnlySpawnIntervalMultiplier
            ),
            accuracy: 0.0001
        )

        progression.advance(to: tuning.skeleton.purpleOnlyLevel)
        XCTAssertEqual(
            progression.skeletonSpawnInterval,
            expectedSkeletonSpawnInterval(
                at: tuning.skeleton.purpleOnlyLevel,
                tuning: tuning,
                multiplier: tuning.skeleton.redOnlySpawnIntervalMultiplier
                    * tuning.skeleton.purpleOnlySpawnIntervalMultiplier
            ),
            accuracy: 0.0001
        )
    }

    func testTimedSkeletonSpawnsTurnPurpleAtLevel75() {
        let tuning = GameConfiguration.defaultTuning
        let scene = GameScene(size: CGSize(width: 640, height: 480), tuning: tuning)

        scene.session.progression.advance(to: tuning.skeleton.redOnlyLevel)
        XCTAssertEqual(scene.timedSkeletonSpawnKind, .red)

        scene.session.progression.advance(to: tuning.skeleton.purpleOnlyLevel)
        XCTAssertEqual(scene.timedSkeletonSpawnKind, .purple)
    }

    func testPurpleSkeletonKillGrantsThreeExperience() {
        let scene = GameScene(size: CGSize(width: 640, height: 480))
        scene.session.progression.advance(to: 3)

        scene.spawnSkeleton(kind: .purple, shouldUpdateHUD: false)
        scene.destroySkeleton(scene.skeletons[0], shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)

        XCTAssertEqual(scene.session.progression.level, 3)
        XCTAssertEqual(scene.session.progression.experience, 3)
    }

    func testBlackSkeletonSpawnsAfterOneHundredPurpleSkeletonKills() {
        let scene = GameScene(size: CGSize(width: 640, height: 480))
        scene.session.kills.purpleSkeletons = GameConfiguration.blackSkeletonPurpleKillInterval - 1

        scene.spawnSkeleton(kind: .purple, shouldUpdateHUD: false)
        scene.destroySkeleton(scene.skeletons[0], shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)

        XCTAssertEqual(scene.skeletons.count, 1)
        XCTAssertEqual(scene.skeletonKind(for: scene.skeletons[0]), .black)
        XCTAssertEqual(scene.skeletonHitPoints(for: scene.skeletons[0]), 25)
        XCTAssertEqual(scene.session.kills.purpleSkeletons, GameConfiguration.blackSkeletonPurpleKillInterval)
    }

    func testBlackSkeletonKillGrantsTenExperience() {
        let scene = GameScene(size: CGSize(width: 640, height: 480))
        scene.session.progression.advance(to: 5)

        scene.spawnSkeleton(kind: .black, shouldUpdateHUD: false)
        scene.destroySkeleton(scene.skeletons[0], shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)

        XCTAssertEqual(scene.session.progression.level, 5)
        XCTAssertEqual(scene.session.progression.experience, 10)
    }

    func testBeamDamageBudgetCanPartiallyDamagePurpleSkeleton() {
        let scene = beamTestScene()

        scene.session.progression.applyLevelUpOption(.beamKillCount)
        scene.spawnSkeleton(kind: .purple, shouldUpdateHUD: false)
        scene.skeletons[0].position = CGPoint(x: 60, y: 0)

        scene.castBeam()

        XCTAssertEqual(scene.skeletons.count, 1)
        XCTAssertEqual(scene.skeletonHitPoints(for: scene.skeletons[0]), 2)
        XCTAssertEqual(scene.session.kills.beam, 0)
        XCTAssertEqual(scene.session.progression.experience, 0)
    }

    func testBeamDamageBudgetKillsPurpleSkeletonWhenItCoversHitPoints() {
        let scene = beamTestScene()

        scene.session.progression.applyLevelUpOption(.beamKillCount)
        scene.session.progression.applyLevelUpOption(.beamKillCount)
        scene.spawnSkeleton(kind: .purple, shouldUpdateHUD: false)
        scene.skeletons[0].position = CGPoint(x: 60, y: 0)

        scene.castBeam()

        XCTAssertTrue(scene.skeletons.isEmpty)
        XCTAssertEqual(scene.session.kills.beam, 1)
        XCTAssertEqual(scene.session.progression.experience, 3)
    }

    private func expectedSkeletonSpawnInterval(
        at level: Int,
        tuning: GameTuning,
        multiplier: TimeInterval = 1
    ) -> TimeInterval {
        tuning.skeleton.initialSpawnInterval
            * pow(tuning.skeleton.intervalMultiplierPerLevel, Double(level - 1))
            * multiplier
    }

    private func beamTestScene() -> GameScene {
        let scene = GameScene(size: CGSize(width: 640, height: 480))
        scene.session.progression.advance(to: 3)
        scene.session.progression.applyLevelUpOption(.learnBeam)
        scene.player.position = .zero
        scene.player.xScale = 1
        return scene
    }
}

private extension Progression {
    mutating func advance(to targetLevel: Int) {
        while level < targetLevel {
            gainExperience(nextExperience)
        }
    }
}
