import SpriteKit

extension GameHUD {
    func showGameOver(level: Int) {
        gameOverLabel.text = "YOU DIED AT LEVEL \(level)"
        gameOverBackground.run(SKAction.fadeAlpha(to: Self.panelAlpha, duration: 0.2))
        gameOverLabel.setScale(0.75)
        gameOverLabel.run(
            SKAction.group([
                SKAction.fadeIn(withDuration: 0.2),
                SKAction.scale(to: 1, duration: 0.2)
            ])
        )

        restartLabel.run(SKAction.fadeIn(withDuration: 0.2))
        exitLabel.run(SKAction.fadeIn(withDuration: 0.2))
    }


    func hideGameOver() {
        gameOverBackground.removeAllActions()
        gameOverBackground.alpha = 0

        for label in [gameOverLabel, restartLabel, exitLabel] {
            label.removeAllActions()
            label.alpha = 0
            label.setScale(1)
        }
    }

    func option(at point: CGPoint) -> GameOverOption? {
        if hitArea(for: restartLabel).contains(point) {
            return .restart
        }

        if hitArea(for: exitLabel).contains(point) {
            return .exit
        }

        return nil
    }


    func setupGameOverLabel() {
        gameOverLabel.text = "YOU DIED AT LEVEL 1"
        gameOverLabel.fontName = "Menlo-Bold"
        gameOverLabel.fontSize = 40
        gameOverLabel.fontColor = Self.deathTextColor
        gameOverLabel.horizontalAlignmentMode = .center
        gameOverLabel.verticalAlignmentMode = .center
        gameOverLabel.position = CGPoint(x: 0, y: 42)
        gameOverLabel.zPosition = 100
        gameOverLabel.alpha = 0

        setupGameOverOption(restartLabel, text: "RESTART", yPosition: -22)
        setupGameOverOption(exitLabel, text: "EXIT", yPosition: -70)
    }


    func setupGameOverOption(_ label: SKLabelNode, text: String, yPosition: CGFloat) {
        label.text = text
        label.fontSize = 22
        label.fontColor = Self.deathTextColor
        label.horizontalAlignmentMode = .center
        label.verticalAlignmentMode = .center
        label.position = CGPoint(x: 0, y: yPosition)
        label.zPosition = 100
        label.alpha = 0
    }


}
