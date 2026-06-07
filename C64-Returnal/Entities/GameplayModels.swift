import SpriteKit

enum AttackKind {
    case fireball
    case lightning
    case orbitalOrb
    case beam
    case meteor
}

enum SkeletonKind {
    case regular
    case red
    case purple

    var hitPoints: Int {
        switch self {
        case .regular:
            return 1
        case .red:
            return max(1, GameConfiguration.redSkeletonHitPoints)
        case .purple:
            return max(1, GameConfiguration.purpleSkeletonHitPoints)
        }
    }

    var tint: SKColor {
        switch self {
        case .regular:
            return .white
        case .red:
            return SKColor(calibratedRed: 0.95, green: 0.05, blue: 0.04, alpha: 1)
        case .purple:
            return SKColor(calibratedRed: 0.58, green: 0.12, blue: 0.95, alpha: 1)
        }
    }

    var tintBlendFactor: CGFloat {
        switch self {
        case .regular:
            return 0
        case .red:
            return 0.68
        case .purple:
            return 0.72
        }
    }
}

enum SkeletonUserDataKey {
    static let hitPoints = "skeletonHitPoints"
}
