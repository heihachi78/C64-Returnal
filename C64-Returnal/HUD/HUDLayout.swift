import SpriteKit

extension GameHUD {
    func setupBackgroundPanels() {
        for panel in [topStatusBackground, combatStatusBackground] {
            panel.fillColor = Self.panelColor
            panel.strokeColor = .clear
            panel.alpha = Self.panelAlpha
            panel.zPosition = 80
        }

        for panel in [levelUpBackground, chestRewardBackground, gameOverBackground] {
            panel.fillColor = Self.panelColor
            panel.strokeColor = .clear
            panel.alpha = 0
            panel.zPosition = 99
        }
    }

    func layout(for sceneSize: CGSize) {
        currentSceneSize = sceneSize
        let left = -sceneSize.width / 2 + 18
        let top = sceneSize.height / 2 - 24

        levelLabel.position = CGPoint(x: left, y: top)
        experienceLabel.position = CGPoint(x: left, y: top - 24)
        lifeIconOrigin = CGPoint(x: left + 7, y: top - 50)
        layoutLifeIcons()
        layoutTopStatusBackground(left: left, top: top)

        let bottom = -sceneSize.height / 2 + 18
        fireballIcon.position = CGPoint(x: left + 9, y: bottom + 20)
        fireballCountLabel.position = CGPoint(x: left + 28, y: bottom + 30)
        fireballIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 14)
        fireballKillsLabel.position = CGPoint(x: left + 28, y: bottom)
        lightningIcon.position = CGPoint(x: left + 9, y: bottom + 74)
        lightningCountLabel.position = CGPoint(x: left + 28, y: bottom + 84)
        lightningIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 68)
        lightningKillsLabel.position = CGPoint(x: left + 28, y: bottom + 54)
        orbIcon.position = CGPoint(x: left + 9, y: bottom + 128)
        orbCountLabel.position = CGPoint(x: left + 28, y: bottom + 138)
        orbSpeedLabel.position = CGPoint(x: left + 28, y: bottom + 122)
        orbKillsLabel.position = CGPoint(x: left + 28, y: bottom + 108)
        beamIcon.position = CGPoint(x: left + 9, y: bottom + 182)
        beamKillLabel.position = CGPoint(x: left + 28, y: bottom + 192)
        beamIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 176)
        beamKillsLabel.position = CGPoint(x: left + 28, y: bottom + 162)
        meteorIcon.position = CGPoint(x: left + 9, y: bottom + 236)
        meteorCountLabel.position = CGPoint(x: left + 28, y: bottom + 246)
        meteorIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 230)
        meteorKillsLabel.position = CGPoint(x: left + 28, y: bottom + 216)
        skeletonIcon.position = CGPoint(x: left + 9, y: bottom + 290)
        skeletonAliveLabel.position = CGPoint(x: left + 28, y: bottom + 300)
        skeletonIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 284)

        layoutChestRewardItems()

        setPanel(
            combatStatusBackground,
            rect: CGRect(x: left - 10, y: bottom - 15, width: 176, height: 330),
            cornerRadius: Self.panelCornerRadius
        )

        layoutLevelUpBackground(optionCount: max(2, activeLevelUpOptions.count))
        layoutChestRewardBackground(itemCount: activeChestRewardItems.count)
        setPanel(
            gameOverBackground,
            rect: CGRect(x: -centeredPanelWidth / 2, y: -100, width: centeredPanelWidth, height: 190),
            cornerRadius: Self.panelCornerRadius
        )
    }


    func layoutLifeIcons() {
        for (index, icon) in lifeIcons.enumerated() {
            let column = index % Self.lifeIconsPerRow
            let row = index / Self.lifeIconsPerRow
            icon.position = CGPoint(
                x: lifeIconOrigin.x + CGFloat(column) * Self.lifeIconSpacing,
                y: lifeIconOrigin.y - CGFloat(row) * Self.lifeIconSpacing
            )
        }
    }

    func layoutTopStatusBackground(left: CGFloat? = nil, top: CGFloat? = nil) {
        let rows = max(1, Int(ceil(Double(max(1, currentLives)) / Double(Self.lifeIconsPerRow))))
        let panelLeft = (left ?? levelLabel.position.x) - 10
        let panelTop = (top ?? levelLabel.position.y) + 15
        let lowestHeartBottom = lifeIconOrigin.y - CGFloat(rows - 1) * Self.lifeIconSpacing - Self.lifeIconSize / 2
        let height = panelTop - lowestHeartBottom + 8

        setPanel(
            topStatusBackground,
            rect: CGRect(x: panelLeft, y: panelTop - height, width: 210, height: height),
            cornerRadius: Self.panelCornerRadius
        )
    }


    func layoutChestRewardItems() {
        let itemCount = chestRewardNodes.count / 2

        guard itemCount > 0 else {
            return
        }

        for index in 0..<itemCount {
            let yPosition = Self.chestRewardItemYPosition(for: index, itemCount: itemCount)
            chestRewardNodes[index * 2].position = CGPoint(x: -104, y: yPosition)
            chestRewardNodes[index * 2 + 1].position = CGPoint(x: -66, y: yPosition)
        }
    }

    var centeredPanelWidth: CGFloat {
        min(max(360, currentSceneSize.width - 48), 620)
    }

    func layoutLevelUpBackground(optionCount: Int) {
        let clampedOptionCount = max(2, optionCount)
        let height = CGFloat(178 + max(0, clampedOptionCount - 2) * 52)

        setPanel(
            levelUpBackground,
            rect: CGRect(x: -centeredPanelWidth / 2, y: 90 - height, width: centeredPanelWidth, height: height),
            cornerRadius: Self.panelCornerRadius
        )
    }

    func layoutChestRewardBackground(itemCount: Int) {
        let visibleItemCount = max(1, itemCount)
        let height = min(currentSceneSize.height - 64, CGFloat(174 + visibleItemCount * 34))

        chestRewardLabel.position = CGPoint(x: 0, y: height / 2 - 54)
        chestRewardContinueLabel.position = CGPoint(x: 0, y: -height / 2 + 36)
        setPanel(
            chestRewardBackground,
            rect: CGRect(x: -centeredPanelWidth / 2, y: -height / 2, width: centeredPanelWidth, height: height),
            cornerRadius: Self.panelCornerRadius
        )
    }

    func hitArea(for label: SKLabelNode) -> CGRect {
        label.frame.insetBy(dx: -26, dy: -12)
    }


    func setPanel(_ panel: SKShapeNode, rect: CGRect, cornerRadius: CGFloat) {
        panel.path = CGPath(
            roundedRect: rect,
            cornerWidth: cornerRadius,
            cornerHeight: cornerRadius,
            transform: nil
        )
    }
}
