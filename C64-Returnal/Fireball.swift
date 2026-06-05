//
//  Fireball.swift
//  C64-Returnal
//

import SpriteKit

struct Fireball {
    let node: SKSpriteNode
    var target: SKSpriteNode?
    var velocity: CGVector
    var timeWithoutTarget: TimeInterval
}
