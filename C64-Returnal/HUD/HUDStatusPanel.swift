import SpriteKit

extension GameHUD {
    func updateProgress(level: Int, experience: Int, nextExperience: Int) {
        levelLabel.text = "LV \(level)"
        experienceLabel.text = "XP \(experience)/\(nextExperience)"
    }

    func updateLives(_ lives: Int) {
        currentLives = max(0, lives)
        syncLifeIcons()
        layoutLifeIcons()
        layoutTopStatusBackground()
    }

    func updateFireballStatus(count: Int, interval: TimeInterval) {
        fireballCountLabel.text = "x\(count)"
        fireballIntervalLabel.text = "\(Self.formattedSeconds(interval))s"
    }

    func updateLightningStatus(isUnlocked: Bool, strikeCount: Int, interval: TimeInterval) {
        guard isUnlocked else {
            lightningCountLabel.text = "LOCKED"
            lightningIntervalLabel.text = "--"
            return
        }

        lightningCountLabel.text = "x\(strikeCount)"
        lightningIntervalLabel.text = "\(Self.formattedSeconds(interval))s"
    }

    func updateOrbStatus(isUnlocked: Bool, count: Int, angularSpeed: CGFloat) {
        guard isUnlocked else {
            orbCountLabel.text = "LOCKED"
            orbSpeedLabel.text = "--"
            return
        }

        orbCountLabel.text = "x\(count)"
        orbSpeedLabel.text = "\(String(format: "%.1f", angularSpeed))r/s"
    }

    func updateBeamStatus(isUnlocked: Bool, killCount: Int, interval: TimeInterval) {
        guard isUnlocked else {
            beamKillLabel.text = "LOCKED"
            beamIntervalLabel.text = "--"
            return
        }

        beamKillLabel.text = "x\(killCount)"
        beamIntervalLabel.text = "\(Self.formattedSeconds(interval))s"
    }

    func updateMeteorStatus(isUnlocked: Bool, count: Int, interval: TimeInterval) {
        guard isUnlocked else {
            meteorCountLabel.text = "LOCKED"
            meteorIntervalLabel.text = "--"
            return
        }

        meteorCountLabel.text = "x\(count)"
        meteorIntervalLabel.text = "\(Self.formattedSeconds(interval))s"
    }

    func updateSkeletonStatus(aliveCount: Int, spawnInterval: TimeInterval) {
        skeletonAliveLabel.text = "x\(aliveCount)"
        skeletonIntervalLabel.text = "\(Self.formattedSeconds(spawnInterval))s"
    }

    func updateAttackKillCounts(fireball: Int, lightning: Int, orb: Int, beam: Int, meteor: Int) {
        fireballKillsLabel.text = "KILLS \(fireball)"
        lightningKillsLabel.text = "KILLS \(lightning)"
        orbKillsLabel.text = "KILLS \(orb)"
        beamKillsLabel.text = "KILLS \(beam)"
        meteorKillsLabel.text = "KILLS \(meteor)"
    }


    func setupProgressLabels() {
        for label in [levelLabel, experienceLabel] {
            label.fontSize = 16
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    func syncLifeIcons() {
        guard let parentNode = parentNode, let lifeTexture = lifeTexture else {
            return
        }

        while lifeIcons.count < currentLives {
            let icon = SKSpriteNode(texture: lifeTexture)
            icon.size = CGSize(width: Self.lifeIconSize, height: Self.lifeIconSize)
            icon.zPosition = 90
            parentNode.addChild(icon)
            lifeIcons.append(icon)
        }

        while lifeIcons.count > currentLives {
            lifeIcons.removeLast().removeFromParent()
        }
    }


    func setupFireballStatus(fireballTexture: SKTexture) {
        fireballIcon.texture = fireballTexture
        fireballIcon.size = CGSize(width: 18, height: 18)
        fireballIcon.zPosition = 90

        setupCombatStatusLabels([fireballCountLabel, fireballIntervalLabel, fireballKillsLabel])
    }

    func setupSkeletonStatus(skeletonTexture: SKTexture) {
        skeletonIcon.texture = skeletonTexture
        skeletonIcon.size = CGSize(width: 16, height: 22)
        skeletonIcon.zPosition = 90

        for label in [skeletonAliveLabel, skeletonIntervalLabel] {
            label.fontSize = 14
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    func setupLightningStatus(lightningTexture: SKTexture) {
        lightningIcon.texture = lightningTexture
        lightningIcon.size = CGSize(width: 18, height: 18)
        lightningIcon.zPosition = 90

        setupCombatStatusLabels([lightningCountLabel, lightningIntervalLabel, lightningKillsLabel])
    }

    func setupOrbStatus(orbTexture: SKTexture) {
        orbIcon.texture = orbTexture
        orbIcon.size = CGSize(width: 18, height: 18)
        orbIcon.zPosition = 90

        setupCombatStatusLabels([orbCountLabel, orbSpeedLabel, orbKillsLabel])
    }

    func setupBeamStatus(beamTexture: SKTexture) {
        beamIcon.texture = beamTexture
        beamIcon.size = CGSize(width: 18, height: 18)
        beamIcon.zPosition = 90

        setupCombatStatusLabels([beamKillLabel, beamIntervalLabel, beamKillsLabel])
    }

    func setupMeteorStatus(meteorTexture: SKTexture) {
        meteorIcon.texture = meteorTexture
        meteorIcon.size = CGSize(width: 18, height: 18)
        meteorIcon.zPosition = 90

        setupCombatStatusLabels([meteorCountLabel, meteorIntervalLabel, meteorKillsLabel])
    }

    func setupCombatStatusLabels(_ labels: [SKLabelNode]) {
        for label in labels {
            label.fontSize = 14
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    static func formattedSeconds(_ value: TimeInterval) -> String {
        if value >= 1 {
            return String(format: "%.1f", value)
        }

        return String(format: "%.2f", value)
    }

}
