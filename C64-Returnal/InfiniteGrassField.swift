//
//  InfiniteGrassField.swift
//  C64-Returnal
//

import SpriteKit

final class InfiniteGrassField {
    let node = SKNode()

    private let tileSize: CGFloat
    private let textures: [SKTexture]
    private var tiles = [SKSpriteNode]()
    private var columns = 0
    private var rows = 0

    init(tileSize: CGFloat, textures: [SKTexture]) {
        self.tileSize = tileSize
        self.textures = textures
    }

    func rebuild(for sceneSize: CGSize) {
        let neededColumns = max(6, Int(ceil(sceneSize.width / tileSize)) + 4)
        let neededRows = max(6, Int(ceil(sceneSize.height / tileSize)) + 4)
        let neededTileCount = neededColumns * neededRows

        guard neededColumns != columns || neededRows != rows || tiles.count != neededTileCount else {
            return
        }

        tiles.forEach { $0.removeFromParent() }
        tiles.removeAll()
        columns = neededColumns
        rows = neededRows

        for _ in 0..<neededTileCount {
            let tile = SKSpriteNode(texture: textures[0])
            tile.size = CGSize(width: tileSize, height: tileSize)
            tile.zPosition = -20
            tile.colorBlendFactor = 0.22
            node.addChild(tile)
            tiles.append(tile)
        }
    }

    func update(around center: CGPoint) {
        guard columns > 0, rows > 0 else {
            return
        }

        let centerColumn = Int(floor(center.x / tileSize))
        let centerRow = Int(floor(center.y / tileSize))
        let startColumn = centerColumn - columns / 2
        let startRow = centerRow - rows / 2

        for row in 0..<rows {
            for column in 0..<columns {
                updateTile(
                    tiles[row * columns + column],
                    column: startColumn + column,
                    row: startRow + row
                )
            }
        }
    }

    private func updateTile(_ tile: SKSpriteNode, column: Int, row: Int) {
        tile.position = CGPoint(
            x: CGFloat(column) * tileSize + tileSize / 2,
            y: CGFloat(row) * tileSize + tileSize / 2
        )
        tile.texture = textures[grassHash(column: column, row: row, salt: 17) % textures.count]
        tile.color = grassTint(column: column, row: row)
        tile.xScale = grassFlipValue(column: column, row: row, salt: 31)
        tile.yScale = grassFlipValue(column: column, row: row, salt: 67)
        tile.zRotation = CGFloat(grassHash(column: column, row: row, salt: 101) % 4) * .pi / 2
    }

    private func grassTint(column: Int, row: Int) -> SKColor {
        let palette = [
            SKColor(calibratedRed: 0.27, green: 0.55, blue: 0.21, alpha: 1),
            SKColor(calibratedRed: 0.20, green: 0.47, blue: 0.18, alpha: 1),
            SKColor(calibratedRed: 0.34, green: 0.62, blue: 0.24, alpha: 1),
            SKColor(calibratedRed: 0.16, green: 0.40, blue: 0.17, alpha: 1),
            SKColor(calibratedRed: 0.42, green: 0.67, blue: 0.25, alpha: 1)
        ]
        let index = grassHash(column: column, row: row, salt: 0) % palette.count

        return palette[index]
    }

    private func grassFlipValue(column: Int, row: Int, salt: Int) -> CGFloat {
        grassHash(column: column, row: row, salt: salt).isMultiple(of: 2) ? 1 : -1
    }

    private func grassHash(column: Int, row: Int, salt: Int) -> Int {
        let hash = (column &* 73_856_093) ^ (row &* 19_349_663) ^ (salt &* 83_492_791)
        return hash == Int.min ? 0 : abs(hash)
    }
}
