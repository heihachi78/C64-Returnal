//
//  Chest.swift
//  C64-Returnal
//

import SpriteKit

enum ChestTier: Hashable {
    case bronze
    case silver
    case gold

    var title: String {
        switch self {
        case .bronze:
            return "BRONZE"
        case .silver:
            return "SILVER"
        case .gold:
            return "GOLD"
        }
    }
}

struct Chest {
    let node: SKSpriteNode
    let tier: ChestTier
}
