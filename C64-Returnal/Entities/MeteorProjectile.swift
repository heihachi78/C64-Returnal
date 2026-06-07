//
//  MeteorProjectile.swift
//  C64-Returnal
//

import SpriteKit

struct MeteorProjectile {
    let node: SKSpriteNode
    let startPosition: CGPoint
    let impactPosition: CGPoint
    private var elapsedTime: TimeInterval = 0

    init(node: SKSpriteNode, startPosition: CGPoint, impactPosition: CGPoint) {
        self.node = node
        self.startPosition = startPosition
        self.impactPosition = impactPosition
    }

    mutating func update(deltaTime: TimeInterval) -> Bool {
        elapsedTime += deltaTime
        let progress = min(1, CGFloat(elapsedTime / GameConfiguration.meteorFallDuration))
        node.position = CGPoint(
            x: startPosition.x + (impactPosition.x - startPosition.x) * progress,
            y: startPosition.y + (impactPosition.y - startPosition.y) * progress
        )

        return progress >= 1
    }
}
