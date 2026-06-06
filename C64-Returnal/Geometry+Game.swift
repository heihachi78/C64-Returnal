//
//  Geometry+Game.swift
//  C64-Returnal
//

import CoreGraphics

extension CGPoint {
    func distanceSquared(to other: CGPoint) -> CGFloat {
        let dx = x - other.x
        let dy = y - other.y

        return dx * dx + dy * dy
    }

    func distance(to other: CGPoint) -> CGFloat {
        sqrt(distanceSquared(to: other))
    }
}

extension CGVector {
    init(from start: CGPoint, to end: CGPoint) {
        self.init(dx: end.x - start.x, dy: end.y - start.y)
    }

    var normalized: CGVector {
        let length = sqrt(dx * dx + dy * dy)

        guard length > 0 else {
            return CGVector(dx: 0, dy: 0)
        }

        return CGVector(dx: dx / length, dy: dy / length)
    }
}
