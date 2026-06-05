//
//  OrbitalOrb.swift
//  C64-Returnal
//

import SpriteKit

struct OrbitalOrb {
    var node: SKSpriteNode?
    private var missingOrbitProgress: CGFloat = 0

    var isActive: Bool {
        node != nil
    }

    init(node: SKSpriteNode? = nil) {
        self.node = node
    }

    mutating func attach(_ node: SKSpriteNode) {
        self.node = node
        missingOrbitProgress = 0
    }

    mutating func deactivate() {
        node?.removeAllActions()
        node?.removeFromParent()
        node = nil
        missingOrbitProgress = 0
    }

    mutating func updateMissingOrbitProgress(by angleDelta: CGFloat) -> Bool {
        guard node == nil else {
            return false
        }

        missingOrbitProgress += abs(angleDelta)
        return missingOrbitProgress >= CGFloat.pi * 2
    }

    func updatePosition(around center: CGPoint, angle: CGFloat, radius: CGFloat) {
        node?.position = CGPoint(
            x: center.x + cos(angle) * radius,
            y: center.y + sin(angle) * radius
        )
    }
}
