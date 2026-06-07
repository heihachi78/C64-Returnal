//
//  SkeletonSpatialIndex.swift
//  C64-Returnal
//

import SpriteKit

final class SkeletonSpatialIndex {
    private struct Cell: Hashable {
        let x: Int
        let y: Int
    }

    private let cellSize: CGFloat
    private var buckets = [Cell: [SKSpriteNode]]()

    init(cellSize: CGFloat) {
        self.cellSize = cellSize
    }

    func removeAll() {
        buckets.removeAll(keepingCapacity: true)
    }

    func rebuild(with skeletons: [SKSpriteNode]) {
        buckets.removeAll(keepingCapacity: true)

        for skeleton in skeletons {
            buckets[cell(for: skeleton.position), default: []].append(skeleton)
        }
    }

    func firstCandidate(
        near position: CGPoint,
        radius: CGFloat,
        isValid: (SKSpriteNode) -> Bool,
        matches: (SKSpriteNode) -> Bool
    ) -> SKSpriteNode? {
        let radiusSquared = radius * radius
        var result: SKSpriteNode?

        forEachCandidate(near: position, radius: radius) { skeleton in
            guard isValid(skeleton), position.distanceSquared(to: skeleton.position) <= radiusSquared, matches(skeleton) else {
                return true
            }

            result = skeleton
            return false
        }

        return result
    }

    func forEachCandidate(
        near position: CGPoint,
        radius: CGFloat,
        isValid: (SKSpriteNode) -> Bool,
        body: (SKSpriteNode) -> Void
    ) {
        let radiusSquared = radius * radius

        forEachCandidate(near: position, radius: radius) { skeleton in
            guard isValid(skeleton), position.distanceSquared(to: skeleton.position) <= radiusSquared else {
                return true
            }

            body(skeleton)
            return true
        }
    }

    func forEachCandidate(
        in rect: CGRect,
        isValid: (SKSpriteNode) -> Bool,
        body: (SKSpriteNode) -> Void
    ) {
        let minCell = cell(for: CGPoint(x: rect.minX, y: rect.minY))
        let maxCell = cell(for: CGPoint(x: rect.maxX, y: rect.maxY))

        for y in minCell.y...maxCell.y {
            for x in minCell.x...maxCell.x {
                guard let bucket = buckets[Cell(x: x, y: y)] else {
                    continue
                }

                for skeleton in bucket where isValid(skeleton) {
                    body(skeleton)
                }
            }
        }
    }

    private func forEachCandidate(
        near position: CGPoint,
        radius: CGFloat,
        body: (SKSpriteNode) -> Bool
    ) {
        let minCell = cell(for: CGPoint(x: position.x - radius, y: position.y - radius))
        let maxCell = cell(for: CGPoint(x: position.x + radius, y: position.y + radius))

        for y in minCell.y...maxCell.y {
            for x in minCell.x...maxCell.x {
                guard let bucket = buckets[Cell(x: x, y: y)] else {
                    continue
                }

                for skeleton in bucket where !body(skeleton) {
                    return
                }
            }
        }
    }

    private func cell(for position: CGPoint) -> Cell {
        Cell(
            x: Int(floor(position.x / cellSize)),
            y: Int(floor(position.y / cellSize))
        )
    }
}
