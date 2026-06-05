//
//  ChainLightning.swift
//  C64-Returnal
//

import SpriteKit

struct ChainLightning {
    struct Strike {
        let start: CGPoint
        let end: CGPoint
        let target: SKSpriteNode
    }

    static let effectName = "lightningEffect"

    let strikes: [Strike]

    init(origin: CGPoint, strikeCount: Int, targets: [SKSpriteNode]) {
        var remainingTargets = targets
        var strikeOrigin = origin
        var builtStrikes = [Strike]()

        for _ in 0..<strikeCount {
            guard let targetIndex = Self.closestTargetIndex(to: strikeOrigin, in: remainingTargets) else {
                break
            }

            let target = remainingTargets.remove(at: targetIndex)
            let targetPosition = target.position
            builtStrikes.append(Strike(start: strikeOrigin, end: targetPosition, target: target))
            strikeOrigin = targetPosition
        }

        strikes = builtStrikes
    }

    static func makeEffectNode(from start: CGPoint, to end: CGPoint, texture: SKTexture) -> SKNode {
        let effect = SKNode()
        effect.name = effectName
        effect.zPosition = 13

        let outerBolt = makeBolt(from: start, to: end, lineWidth: GameConfiguration.lightningBranchWidth)
        outerBolt.strokeColor = SKColor(calibratedRed: 0.13, green: 0.64, blue: 1.0, alpha: 1)
        outerBolt.glowWidth = 5

        let innerBolt = makeBolt(from: start, to: end, lineWidth: 1)
        innerBolt.strokeColor = SKColor.white
        innerBolt.glowWidth = 1

        let impact = SKSpriteNode(texture: texture)
        impact.position = end
        impact.size = CGSize(width: 24, height: 24)
        impact.zPosition = 1

        effect.addChild(outerBolt)
        effect.addChild(innerBolt)
        effect.addChild(impact)
        effect.run(
            SKAction.sequence([
                SKAction.fadeOut(withDuration: GameConfiguration.lightningEffectDuration),
                SKAction.removeFromParent()
            ])
        )

        return effect
    }

    private static func closestTargetIndex(to position: CGPoint, in targets: [SKSpriteNode]) -> Int? {
        targets.indices.min {
            targets[$0].position.distance(to: position) < targets[$1].position.distance(to: position)
        }
    }

    private static func makeBolt(from start: CGPoint, to end: CGPoint, lineWidth: CGFloat) -> SKShapeNode {
        let path = CGMutablePath()
        let delta = CGVector(dx: end.x - start.x, dy: end.y - start.y)
        let distance = max(1, start.distance(to: end))
        let normal = CGVector(dx: -delta.dy / distance, dy: delta.dx / distance)
        let segmentCount = max(3, min(9, Int(distance / 30)))

        path.move(to: start)

        for segment in 1..<segmentCount {
            let progress = CGFloat(segment) / CGFloat(segmentCount)
            let basePoint = CGPoint(
                x: start.x + delta.dx * progress,
                y: start.y + delta.dy * progress
            )
            let offset = CGFloat.random(in: -8...8)
            path.addLine(
                to: CGPoint(
                    x: basePoint.x + normal.dx * offset,
                    y: basePoint.y + normal.dy * offset
                )
            )
        }

        path.addLine(to: end)

        let bolt = SKShapeNode(path: path)
        bolt.lineCap = .square
        bolt.lineJoin = .miter
        bolt.lineWidth = lineWidth
        return bolt
    }
}
