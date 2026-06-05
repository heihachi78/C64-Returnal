//
//  Beam.swift
//  C64-Returnal
//

import SpriteKit

struct Beam {
    static let effectName = "beamEffect"

    let start: CGPoint
    let end: CGPoint
    let targets: [SKSpriteNode]

    init(origin: CGPoint, direction: CGVector, length: CGFloat, hitWidth: CGFloat, killLimit: Int, targets: [SKSpriteNode]) {
        let normalizedDirection = direction.normalized
        start = origin
        end = CGPoint(
            x: origin.x + normalizedDirection.dx * length,
            y: origin.y + normalizedDirection.dy * length
        )

        self.targets = targets
            .compactMap { target -> (target: SKSpriteNode, progress: CGFloat)? in
                guard let progress = Self.progressAlongBeam(
                    targetPosition: target.position,
                    origin: origin,
                    direction: normalizedDirection,
                    length: length,
                    hitWidth: hitWidth
                ) else {
                    return nil
                }

                return (target, progress)
            }
            .sorted { $0.progress < $1.progress }
            .prefix(killLimit)
            .map(\.target)
    }

    static func makeEffectNode(from start: CGPoint, to end: CGPoint) -> SKNode {
        let effect = SKNode()
        effect.name = effectName
        effect.zPosition = 13

        let outerBeam = makeBeamLine(from: start, to: end, lineWidth: 9)
        outerBeam.strokeColor = SKColor(calibratedRed: 1.0, green: 0.72, blue: 0.08, alpha: 1)
        outerBeam.glowWidth = 7

        let coreBeam = makeBeamLine(from: start, to: end, lineWidth: 4)
        coreBeam.strokeColor = SKColor(calibratedRed: 1.0, green: 0.94, blue: 0.22, alpha: 1)
        coreBeam.glowWidth = 3

        let whiteCore = makeBeamLine(from: start, to: end, lineWidth: 1)
        whiteCore.strokeColor = SKColor.white

        effect.addChild(outerBeam)
        effect.addChild(coreBeam)
        effect.addChild(whiteCore)
        effect.run(
            SKAction.sequence([
                SKAction.fadeOut(withDuration: GameConfiguration.beamEffectDuration),
                SKAction.removeFromParent()
            ])
        )

        return effect
    }

    private static func progressAlongBeam(
        targetPosition: CGPoint,
        origin: CGPoint,
        direction: CGVector,
        length: CGFloat,
        hitWidth: CGFloat
    ) -> CGFloat? {
        let targetVector = CGVector(from: origin, to: targetPosition)
        let progress = targetVector.dx * direction.dx + targetVector.dy * direction.dy

        guard progress >= 0, progress <= length else {
            return nil
        }

        let closestPoint = CGPoint(
            x: origin.x + direction.dx * progress,
            y: origin.y + direction.dy * progress
        )

        guard closestPoint.distance(to: targetPosition) <= hitWidth else {
            return nil
        }

        return progress
    }

    private static func makeBeamLine(from start: CGPoint, to end: CGPoint, lineWidth: CGFloat) -> SKShapeNode {
        let path = CGMutablePath()
        path.move(to: start)
        path.addLine(to: end)

        let beam = SKShapeNode(path: path)
        beam.lineCap = .square
        beam.lineJoin = .miter
        beam.lineWidth = lineWidth
        return beam
    }
}
