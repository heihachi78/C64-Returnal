import SpriteKit

extension GameScene {
    func spawnCoin(for level: Int) {
        guard level > 0, !session.spawnedCoinLevels.contains(level), !coinTextures.isEmpty else {
            return
        }

        session.spawnedCoinLevels.insert(level)

        let rewardRange = coinRewardRange
        let node = SKSpriteNode(texture: coinTextures[0])
        node.size = CGSize(width: 28, height: 28)
        node.position = randomCoinPosition()
        node.zPosition = 8.75

        animateCoin(node)
        worldNode.addChild(node)
        coins.append(
            Coin(
                node: node,
                amount: Int.random(in: rewardRange),
                level: level
            )
        )
    }

    var coinRewardRange: ClosedRange<Int> {
        let minimum = max(1, tuning.coin.minimumReward)
        let maximum = max(minimum, tuning.coin.maximumReward)
        return minimum...maximum
    }

    func randomCoinPosition() -> CGPoint {
        let halfWidth = max(1, size.width / 2)
        let halfHeight = max(1, size.height / 2)
        let margin = max(1, tuning.coin.spawnMargin)

        switch Int.random(in: 0..<4) {
        case 0:
            return CGPoint(
                x: player.position.x - halfWidth - margin - CGFloat.random(in: 0...halfWidth),
                y: player.position.y + CGFloat.random(in: -halfHeight...halfHeight)
            )
        case 1:
            return CGPoint(
                x: player.position.x + halfWidth + margin + CGFloat.random(in: 0...halfWidth),
                y: player.position.y + CGFloat.random(in: -halfHeight...halfHeight)
            )
        case 2:
            return CGPoint(
                x: player.position.x + CGFloat.random(in: -halfWidth...halfWidth),
                y: player.position.y - halfHeight - margin - CGFloat.random(in: 0...halfHeight)
            )
        default:
            return CGPoint(
                x: player.position.x + CGFloat.random(in: -halfWidth...halfWidth),
                y: player.position.y + halfHeight + margin + CGFloat.random(in: 0...halfHeight)
            )
        }
    }

    func animateCoin(_ coin: SKSpriteNode) {
        coin.removeAllActions()

        let spin = SKAction.repeatForever(
            SKAction.animate(
                with: coinTextures,
                timePerFrame: tuning.coin.animationFrameDuration
            )
        )
        let float = SKAction.repeatForever(
            SKAction.sequence([
                SKAction.moveBy(x: 0, y: 5, duration: 0.42),
                SKAction.moveBy(x: 0, y: -5, duration: 0.42)
            ])
        )
        let shimmer = SKAction.repeatForever(
            SKAction.sequence([
                SKAction.fadeAlpha(to: 0.72, duration: 0.28),
                SKAction.fadeAlpha(to: 1, duration: 0.28)
            ])
        )

        coin.run(SKAction.group([spin, float, shimmer]))
    }

    func checkCoinPickups() {
        let pickupDistanceSquared = tuning.coin.pickupDistance * tuning.coin.pickupDistance

        for index in coins.indices.reversed() {
            guard coins[index].node.position.distanceSquared(to: player.position) <= pickupDistanceSquared else {
                continue
            }

            collectCoin(at: index)
            return
        }
    }

    func collectCoin(at index: Int) {
        let coin = coins[index]
        coin.node.removeAllActions()
        coin.node.removeFromParent()
        coins.remove(at: index)

        session.collectedCoins += coin.amount
        updateHUDProgress()
    }

    var levelUpRedrawCost: Int {
        max(1, session.progression.level)
    }

    func canRedrawLevelUpOptions() -> Bool {
        session.collectedCoins >= levelUpRedrawCost
    }

    func spendCoinsForLevelUpRedraw() -> Bool {
        let cost = levelUpRedrawCost

        guard session.collectedCoins >= cost else {
            return false
        }

        session.collectedCoins -= cost
        return true
    }
}
