import SpriteKit

extension GameScene {
    func configurePlayer() {
        player.size = CGSize(width: 32, height: 44)
        player.zPosition = 10
        stopPlayerAnimation()
        worldNode.addChild(player)
    }

    func layoutViewportContent() {
        hud.layout(for: size)
        grassField.rebuild(for: size)
        grassField.update(around: player.position)
    }

    func updatePlayer(deltaTime: TimeInterval) {
        let horizontal = directionValue(negative: inputController.bindings.moveLeft, positive: inputController.bindings.moveRight)
        let vertical = directionValue(negative: inputController.bindings.moveDown, positive: inputController.bindings.moveUp)
        var movement = CGVector(dx: horizontal, dy: vertical)

        if movement.dx != 0 || movement.dy != 0 {
            movement = movement.normalized
            session.currentPlayerMovementDirection = movement
            player.position.x += movement.dx * tuning.player.speed * CGFloat(deltaTime)
            player.position.y += movement.dy * tuning.player.speed * CGFloat(deltaTime)
            updateFacing(for: player, movement: movement)
            startPlayerAnimation()
        } else {
            session.currentPlayerMovementDirection = nil
            stopPlayerAnimation()
        }

        cameraNode.position = player.position
        grassField.update(around: player.position)
    }


    func damagePlayer() {
        session.playerLives = max(0, session.playerLives - 1)
        hud.updateLives(session.playerLives)

        guard session.playerLives > 0 else {
            triggerGameOver()
            return
        }

        session.playerHitInvulnerabilityTimer = tuning.player.hitInvulnerabilityDuration
        showPlayerHitFeedback()
    }

    func showPlayerHitFeedback() {
        player.removeAction(forKey: Self.playerHitFlashActionKey)
        player.alpha = 1

        let flash = SKAction.sequence([
            SKAction.fadeAlpha(to: 0.35, duration: 0.08),
            SKAction.fadeAlpha(to: 1, duration: 0.08)
        ])
        player.run(SKAction.repeat(flash, count: 6), withKey: Self.playerHitFlashActionKey)
    }


    func resetPlayer() {
        player.removeAllActions()
        player.position = .zero
        player.zRotation = 0
        player.xScale = 1
        player.yScale = 1
        player.alpha = 1
        player.colorBlendFactor = 0
        player.texture = mageTextures[0]
        stopPlayerAnimation()
    }


    func updatePlayerHitInvulnerability(deltaTime: TimeInterval) {
        guard session.playerHitInvulnerabilityTimer > 0 else {
            return
        }

        session.playerHitInvulnerabilityTimer = max(0, session.playerHitInvulnerabilityTimer - deltaTime)
    }


    func updateFacing(for node: SKSpriteNode, movement: CGVector) {
        if movement.dx < 0 {
            node.xScale = -abs(node.xScale)
        } else if movement.dx > 0 {
            node.xScale = abs(node.xScale)
        }
    }

    func startPlayerAnimation() {
        guard player.action(forKey: Self.playerAnimationActionKey) == nil else {
            return
        }

        player.removeAction(forKey: Self.playerAnimationActionKey)
        player.run(
            SKAction.repeatForever(
                SKAction.animate(
                    with: mageTextures,
                    timePerFrame: tuning.player.animationFrameDuration
                )
            ),
            withKey: Self.playerAnimationActionKey
        )
    }

    func stopPlayerAnimation() {
        player.removeAction(forKey: Self.playerAnimationActionKey)
        player.texture = mageTextures[0]
    }


    func playerBeamDirection() -> CGVector {
        if let movementDirection = session.currentPlayerMovementDirection {
            return movementDirection
        }

        return CGVector(dx: player.xScale < 0 ? -1 : 1, dy: 0)
    }

    func directionValue(negative: UInt16, positive: UInt16) -> CGFloat {
        var value: CGFloat = 0

        if session.pressedKeys.contains(negative) {
            value -= 1
        }

        if session.pressedKeys.contains(positive) {
            value += 1
        }

        return value
    }


}
