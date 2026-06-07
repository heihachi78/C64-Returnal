import SpriteKit
import XCTest

final class CombatTargetingTests: XCTestCase {
    func testBeamSelectsClosestTargetsAlongBeamUpToDamageLimit() {
        let near = skeleton(at: CGPoint(x: 30, y: 0))
        let far = skeleton(at: CGPoint(x: 80, y: 0))
        let outsideWidth = skeleton(at: CGPoint(x: 20, y: 40))

        let beam = Beam(
            origin: .zero,
            direction: CGVector(dx: 1, dy: 0),
            length: 100,
            hitWidth: 10,
            damageLimit: 1,
            targets: [far, outsideWidth, near]
        )

        XCTAssertEqual(beam.targets, [near])
    }

    func testChainLightningChainsToNearestRemainingTargets() {
        let first = skeleton(at: CGPoint(x: 10, y: 0))
        let second = skeleton(at: CGPoint(x: 20, y: 0))
        let distant = skeleton(at: CGPoint(x: 200, y: 0))

        let lightning = ChainLightning(
            origin: .zero,
            strikeCount: 2,
            targets: [distant, second, first]
        )

        XCTAssertEqual(lightning.strikes.map(\.target), [first, second])
    }

    func testSpatialIndexFindsCandidatesInRadiusAndRect() {
        let nearby = skeleton(at: CGPoint(x: 10, y: 0))
        let far = skeleton(at: CGPoint(x: 150, y: 0))
        let index = SkeletonSpatialIndex(cellSize: 32)

        index.rebuild(with: [nearby, far])

        XCTAssertEqual(
            index.firstCandidate(near: .zero, radius: 16, isValid: { _ in true }, matches: { _ in true }),
            nearby
        )

        var rectMatches = [SKSpriteNode]()
        index.forEachCandidate(in: CGRect(x: 120, y: -10, width: 60, height: 20), isValid: { _ in true }) {
            rectMatches.append($0)
        }

        XCTAssertEqual(rectMatches, [far])
    }

    private func skeleton(at position: CGPoint) -> SKSpriteNode {
        let node = SKSpriteNode()
        node.position = position
        return node
    }
}
