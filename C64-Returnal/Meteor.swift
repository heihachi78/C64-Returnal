//
//  Meteor.swift
//  C64-Returnal
//

import SpriteKit

struct Meteor {
    static let projectileName = "meteorProjectile"
    static let effectName = "meteorImpactEffect"

    let spawnPosition: CGPoint
    let impactPosition: CGPoint

    init(origin: CGPoint, targetRadius: CGFloat) {
        let angle = CGFloat.random(in: 0..<(CGFloat.pi * 2))
        let distance = sqrt(CGFloat.random(in: 0...1)) * targetRadius

        impactPosition = CGPoint(
            x: origin.x + cos(angle) * distance,
            y: origin.y + sin(angle) * distance
        )

        spawnPosition = CGPoint(
            x: impactPosition.x + CGFloat.random(in: -GameConfiguration.meteorFallDrift...GameConfiguration.meteorFallDrift),
            y: impactPosition.y + GameConfiguration.meteorFallHeight
        )
    }

    static func makeImpactEffectNode(at position: CGPoint, radius: CGFloat) -> SKNode {
        let effect = SKNode()
        effect.name = effectName
        effect.position = position
        effect.zPosition = 13

        let crater = SKShapeNode(circleOfRadius: radius)
        crater.fillColor = SKColor(calibratedRed: 0.29, green: 0.17, blue: 0.08, alpha: 0.45)
        crater.strokeColor = SKColor(calibratedRed: 0.76, green: 0.50, blue: 0.24, alpha: 0.85)
        crater.lineWidth = 2
        crater.glowWidth = 3

        let core = SKShapeNode(circleOfRadius: radius * 0.35)
        core.fillColor = SKColor(calibratedRed: 0.61, green: 0.38, blue: 0.17, alpha: 0.65)
        core.strokeColor = .clear

        effect.addChild(crater)
        effect.addChild(core)
        effect.setScale(0.25)
        effect.run(
            SKAction.sequence([
                SKAction.group([
                    SKAction.scale(to: 1, duration: 0.08),
                    SKAction.fadeIn(withDuration: 0.04)
                ]),
                SKAction.wait(forDuration: 0.08),
                SKAction.group([
                    SKAction.scale(to: 1.25, duration: 0.16),
                    SKAction.fadeOut(withDuration: 0.16)
                ]),
                SKAction.removeFromParent()
            ])
        )

        return effect
    }
}
