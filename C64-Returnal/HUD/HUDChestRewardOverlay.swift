import SpriteKit

extension GameHUD {
    func showChestReward(tier: ChestTier, items: [ChestRewardDisplayItem]) {
        hideChestReward()
        activeChestRewardItems = items
        chestRewardLabel.text = "\(tier.title) CHEST"
        layoutChestRewardBackground(itemCount: items.count)
        chestRewardBackground.run(SKAction.fadeAlpha(to: Self.panelAlpha, duration: 0.2))
        chestRewardLabel.setScale(0.75)
        chestRewardLabel.run(
            SKAction.group([
                SKAction.fadeIn(withDuration: 0.2),
                SKAction.scale(to: 1, duration: 0.2)
            ])
        )
        chestRewardContinueLabel.run(SKAction.fadeIn(withDuration: 0.2))

        showChestRewardItems(items)
    }

    func hideChestReward() {
        chestRewardBackground.removeAllActions()
        chestRewardBackground.alpha = 0
        chestRewardLabel.removeAllActions()
        chestRewardLabel.alpha = 0
        chestRewardLabel.setScale(1)
        chestRewardContinueLabel.removeAllActions()
        chestRewardContinueLabel.alpha = 0
        activeChestRewardItems.removeAll()

        for node in chestRewardNodes {
            node.removeAllActions()
            node.removeFromParent()
        }
        chestRewardNodes.removeAll()
    }


    func setupChestRewardLabels() {
        chestRewardLabel.text = "BRONZE CHEST"
        chestRewardLabel.fontSize = 34
        chestRewardLabel.fontColor = Self.primaryTextColor
        chestRewardLabel.horizontalAlignmentMode = .center
        chestRewardLabel.verticalAlignmentMode = .center
        chestRewardLabel.position = CGPoint(x: 0, y: 88)
        chestRewardLabel.zPosition = 100
        chestRewardLabel.alpha = 0

        chestRewardContinueLabel.text = "[Q] CONTINUE"
        chestRewardContinueLabel.fontSize = 18
        chestRewardContinueLabel.fontColor = Self.keyHintTextColor
        chestRewardContinueLabel.horizontalAlignmentMode = .center
        chestRewardContinueLabel.verticalAlignmentMode = .center
        chestRewardContinueLabel.position = CGPoint(x: 0, y: -116)
        chestRewardContinueLabel.zPosition = 100
        chestRewardContinueLabel.alpha = 0
    }


    func showChestRewardItems(_ items: [ChestRewardDisplayItem]) {
        for (index, item) in items.enumerated() {
            showChestRewardItem(item, index: index, itemCount: items.count)
        }
    }

    func showChestRewardItem(_ item: ChestRewardDisplayItem, index: Int, itemCount: Int) {
        guard let parentNode = parentNode else {
            return
        }

        let yPosition = Self.chestRewardItemYPosition(for: index, itemCount: itemCount)
        let sourceIcon = icon(for: item.option)
        let icon = SKSpriteNode(texture: sourceIcon.texture)
        icon.size = sourceIcon.size
        icon.position = CGPoint(x: -104, y: yPosition)
        icon.zPosition = 100
        icon.alpha = 0

        let label = SKLabelNode(fontNamed: "Menlo-Bold")
        label.text = item.title
        label.fontSize = 20
        label.fontColor = Self.primaryTextColor
        label.horizontalAlignmentMode = .left
        label.verticalAlignmentMode = .center
        label.position = CGPoint(x: -66, y: yPosition)
        label.zPosition = 100
        label.alpha = 0

        parentNode.addChild(icon)
        parentNode.addChild(label)
        chestRewardNodes.append(icon)
        chestRewardNodes.append(label)

        icon.run(SKAction.fadeIn(withDuration: 0.2))
        label.run(SKAction.fadeIn(withDuration: 0.2))
    }


    static func chestRewardItemYPosition(for index: Int, itemCount: Int) -> CGFloat {
        let spacing: CGFloat = 30
        let totalHeight = CGFloat(max(0, itemCount - 1)) * spacing

        return totalHeight / 2 - CGFloat(index) * spacing - 12
    }


}
