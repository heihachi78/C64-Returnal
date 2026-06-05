//
//  GameHUD.swift
//  C64-Returnal
//

import SpriteKit

enum GameOverOption {
    case restart
    case exit
}

enum LevelUpOption: CaseIterable {
    case fireRate
    case extraFireball
    case lightningBounce
    case lightningRate
    case extraOrb
    case orbitalSpeed
    case beamRate
    case beamKillCount

    var title: String {
        title(level: nil)
    }

    func title(level: Int?) -> String {
        switch self {
        case .fireRate:
            return "FASTER FIRE"
        case .extraFireball:
            return "+1 FIREBALL"
        case .lightningBounce:
            return "+1 CHAIN"
        case .lightningRate:
            return "FASTER BOLT"
        case .extraOrb:
            return "+1 ORB"
        case .orbitalSpeed:
            return "FASTER ORB"
        case .beamRate:
            return "FASTER BEAM"
        case .beamKillCount:
            return "+\(level ?? 1) BEAM KILL"
        }
    }
}

final class GameHUD {
    private let gameOverLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let restartLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let exitLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let levelUpLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireRateLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let extraFireballLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningBounceLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningRateLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let extraOrbLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let orbitalSpeedLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamRateLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamKillCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireRateIcon = SKSpriteNode()
    private let extraFireballIcon = SKSpriteNode()
    private let lightningBounceIcon = SKSpriteNode()
    private let lightningRateIcon = SKSpriteNode()
    private let extraOrbIcon = SKSpriteNode()
    private let orbitalSpeedIcon = SKSpriteNode()
    private let beamRateIcon = SKSpriteNode()
    private let beamKillCountIcon = SKSpriteNode()
    private let levelLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let experienceLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let livesLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireballIcon = SKSpriteNode()
    private let fireballCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireballIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningIcon = SKSpriteNode()
    private let lightningCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let orbIcon = SKSpriteNode()
    private let orbCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let orbSpeedLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamIcon = SKSpriteNode()
    private let beamKillLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let skeletonIcon = SKSpriteNode()
    private let skeletonAliveLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let skeletonIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private var activeLevelUpOptions = [LevelUpOption]()

    func add(to parent: SKNode, fireballTexture: SKTexture, lightningTexture: SKTexture, orbTexture: SKTexture, beamTexture: SKTexture, skeletonTexture: SKTexture) {
        setupGameOverLabel()
        setupLevelUpLabels(fireballTexture: fireballTexture, lightningTexture: lightningTexture, orbTexture: orbTexture, beamTexture: beamTexture)
        setupProgressLabels()
        setupFireballStatus(fireballTexture: fireballTexture)
        setupLightningStatus(lightningTexture: lightningTexture)
        setupOrbStatus(orbTexture: orbTexture)
        setupBeamStatus(beamTexture: beamTexture)
        setupSkeletonStatus(skeletonTexture: skeletonTexture)

        parent.addChild(gameOverLabel)
        parent.addChild(restartLabel)
        parent.addChild(exitLabel)
        parent.addChild(levelUpLabel)
        parent.addChild(fireRateLabel)
        parent.addChild(extraFireballLabel)
        parent.addChild(lightningBounceLabel)
        parent.addChild(lightningRateLabel)
        parent.addChild(extraOrbLabel)
        parent.addChild(orbitalSpeedLabel)
        parent.addChild(beamRateLabel)
        parent.addChild(beamKillCountLabel)
        parent.addChild(fireRateIcon)
        parent.addChild(extraFireballIcon)
        parent.addChild(lightningBounceIcon)
        parent.addChild(lightningRateIcon)
        parent.addChild(extraOrbIcon)
        parent.addChild(orbitalSpeedIcon)
        parent.addChild(beamRateIcon)
        parent.addChild(beamKillCountIcon)
        parent.addChild(levelLabel)
        parent.addChild(experienceLabel)
        parent.addChild(livesLabel)
        parent.addChild(fireballIcon)
        parent.addChild(fireballCountLabel)
        parent.addChild(fireballIntervalLabel)
        parent.addChild(lightningIcon)
        parent.addChild(lightningCountLabel)
        parent.addChild(lightningIntervalLabel)
        parent.addChild(orbIcon)
        parent.addChild(orbCountLabel)
        parent.addChild(orbSpeedLabel)
        parent.addChild(beamIcon)
        parent.addChild(beamKillLabel)
        parent.addChild(beamIntervalLabel)
        parent.addChild(skeletonIcon)
        parent.addChild(skeletonAliveLabel)
        parent.addChild(skeletonIntervalLabel)
    }

    func layout(for sceneSize: CGSize) {
        let left = -sceneSize.width / 2 + 18
        let top = sceneSize.height / 2 - 24

        levelLabel.position = CGPoint(x: left, y: top)
        experienceLabel.position = CGPoint(x: left, y: top - 24)
        livesLabel.position = CGPoint(x: left, y: top - 48)

        let bottom = -sceneSize.height / 2 + 18
        fireballIcon.position = CGPoint(x: left + 9, y: bottom + 14)
        fireballCountLabel.position = CGPoint(x: left + 28, y: bottom + 20)
        fireballIntervalLabel.position = CGPoint(x: left + 28, y: bottom)
        lightningIcon.position = CGPoint(x: left + 9, y: bottom + 58)
        lightningCountLabel.position = CGPoint(x: left + 28, y: bottom + 64)
        lightningIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 44)
        orbIcon.position = CGPoint(x: left + 9, y: bottom + 102)
        orbCountLabel.position = CGPoint(x: left + 28, y: bottom + 108)
        orbSpeedLabel.position = CGPoint(x: left + 28, y: bottom + 88)
        beamIcon.position = CGPoint(x: left + 9, y: bottom + 146)
        beamKillLabel.position = CGPoint(x: left + 28, y: bottom + 152)
        beamIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 132)
        skeletonIcon.position = CGPoint(x: left + 9, y: bottom + 190)
        skeletonAliveLabel.position = CGPoint(x: left + 28, y: bottom + 196)
        skeletonIntervalLabel.position = CGPoint(x: left + 28, y: bottom + 176)
    }

    func updateProgress(level: Int, experience: Int, nextExperience: Int) {
        levelLabel.text = "LV \(level)"
        experienceLabel.text = "XP \(experience)/\(nextExperience)"
    }

    func updateLives(_ lives: Int) {
        livesLabel.text = "LIVES \(lives)"
    }

    func updateFireballStatus(count: Int, interval: TimeInterval) {
        fireballCountLabel.text = "x\(count)"
        fireballIntervalLabel.text = "\(Self.formattedSeconds(interval))s"
    }

    func updateLightningStatus(strikeCount: Int, interval: TimeInterval) {
        lightningCountLabel.text = "x\(strikeCount)"
        lightningIntervalLabel.text = "\(Self.formattedSeconds(interval))s"
    }

    func updateOrbStatus(count: Int, angularSpeed: CGFloat) {
        orbCountLabel.text = "x\(count)"
        orbSpeedLabel.text = "\(String(format: "%.1f", angularSpeed))r/s"
    }

    func updateBeamStatus(killCount: Int, interval: TimeInterval) {
        beamKillLabel.text = "x\(killCount)"
        beamIntervalLabel.text = "\(Self.formattedSeconds(interval))s"
    }

    func updateSkeletonStatus(aliveCount: Int, spawnInterval: TimeInterval) {
        skeletonAliveLabel.text = "x\(aliveCount)"
        skeletonIntervalLabel.text = "\(Self.formattedSeconds(spawnInterval))s"
    }

    func showGameOver() {
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

    func showLevelUp(level: Int, options: [LevelUpOption]) {
        hideLevelUpOptions()
        activeLevelUpOptions = Array(options.prefix(2))
        levelUpLabel.text = "LEVEL \(level)"
        levelUpLabel.setScale(0.75)
        levelUpLabel.run(
            SKAction.group([
                SKAction.fadeIn(withDuration: 0.2),
                SKAction.scale(to: 1, duration: 0.2)
            ])
        )

        for (index, option) in activeLevelUpOptions.enumerated() {
            showLevelUpOption(option, yPosition: index == 0 ? -4 : -56, level: level)
        }
    }

    func hideLevelUp() {
        levelUpLabel.removeAllActions()
        levelUpLabel.alpha = 0
        levelUpLabel.setScale(1)
        hideLevelUpOptions()
        activeLevelUpOptions.removeAll()
    }

    private func hideLevelUpOptions() {
        for label in levelUpOptionLabels {
            label.removeAllActions()
            label.alpha = 0
            label.setScale(1)
        }

        for icon in levelUpOptionIcons {
            icon.removeAllActions()
            icon.alpha = 0
            icon.setScale(1)
        }
    }

    func hideGameOver() {
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

    func levelUpOption(at point: CGPoint) -> LevelUpOption? {
        for option in activeLevelUpOptions {
            if hitArea(for: label(for: option)).contains(point) || icon(for: option).frame.insetBy(dx: -14, dy: -14).contains(point) {
                return option
            }
        }

        return nil
    }

    private func setupGameOverLabel() {
        gameOverLabel.text = "GAME OVER"
        gameOverLabel.fontSize = 44
        gameOverLabel.fontColor = Self.primaryTextColor
        gameOverLabel.horizontalAlignmentMode = .center
        gameOverLabel.verticalAlignmentMode = .center
        gameOverLabel.position = CGPoint(x: 0, y: 42)
        gameOverLabel.zPosition = 100
        gameOverLabel.alpha = 0

        setupGameOverOption(restartLabel, text: "RESTART", yPosition: -22)
        setupGameOverOption(exitLabel, text: "EXIT", yPosition: -70)
    }

    private func setupLevelUpLabels(fireballTexture: SKTexture, lightningTexture: SKTexture, orbTexture: SKTexture, beamTexture: SKTexture) {
        levelUpLabel.fontSize = 40
        levelUpLabel.fontColor = Self.primaryTextColor
        levelUpLabel.horizontalAlignmentMode = .center
        levelUpLabel.verticalAlignmentMode = .center
        levelUpLabel.position = CGPoint(x: 0, y: 56)
        levelUpLabel.zPosition = 100
        levelUpLabel.alpha = 0

        setupLevelUpOption(fireRateLabel, text: LevelUpOption.fireRate.title)
        setupLevelUpOption(extraFireballLabel, text: LevelUpOption.extraFireball.title)
        setupLevelUpOption(lightningBounceLabel, text: LevelUpOption.lightningBounce.title)
        setupLevelUpOption(lightningRateLabel, text: LevelUpOption.lightningRate.title)
        setupLevelUpOption(extraOrbLabel, text: LevelUpOption.extraOrb.title)
        setupLevelUpOption(orbitalSpeedLabel, text: LevelUpOption.orbitalSpeed.title)
        setupLevelUpOption(beamRateLabel, text: LevelUpOption.beamRate.title)
        setupLevelUpOption(beamKillCountLabel, text: LevelUpOption.beamKillCount.title)
        setupLevelUpIcon(fireRateIcon, texture: fireballTexture)
        setupLevelUpIcon(extraFireballIcon, texture: fireballTexture)
        setupLevelUpIcon(lightningBounceIcon, texture: lightningTexture)
        setupLevelUpIcon(lightningRateIcon, texture: lightningTexture)
        setupLevelUpIcon(extraOrbIcon, texture: orbTexture)
        setupLevelUpIcon(orbitalSpeedIcon, texture: orbTexture)
        setupLevelUpIcon(beamRateIcon, texture: beamTexture)
        setupLevelUpIcon(beamKillCountIcon, texture: beamTexture)
    }

    private func setupGameOverOption(_ label: SKLabelNode, text: String, yPosition: CGFloat) {
        label.text = text
        label.fontSize = 22
        label.fontColor = Self.primaryTextColor
        label.horizontalAlignmentMode = .center
        label.verticalAlignmentMode = .center
        label.position = CGPoint(x: 0, y: yPosition)
        label.zPosition = 100
        label.alpha = 0
    }

    private func setupLevelUpOption(_ label: SKLabelNode, text: String) {
        label.text = text
        label.fontSize = 22
        label.fontColor = Self.primaryTextColor
        label.horizontalAlignmentMode = .left
        label.verticalAlignmentMode = .center
        label.position = CGPoint(x: -76, y: 0)
        label.zPosition = 100
        label.alpha = 0
    }

    private func setupLevelUpIcon(_ icon: SKSpriteNode, texture: SKTexture) {
        icon.texture = texture
        icon.size = CGSize(width: 24, height: 24)
        icon.position = CGPoint(x: -110, y: 0)
        icon.zPosition = 100
        icon.alpha = 0
    }

    private func setupProgressLabels() {
        for label in [levelLabel, experienceLabel, livesLabel] {
            label.fontSize = 16
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    private func setupFireballStatus(fireballTexture: SKTexture) {
        fireballIcon.texture = fireballTexture
        fireballIcon.size = CGSize(width: 18, height: 18)
        fireballIcon.zPosition = 90

        for label in [fireballCountLabel, fireballIntervalLabel] {
            label.fontSize = 14
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    private func setupSkeletonStatus(skeletonTexture: SKTexture) {
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

    private func setupLightningStatus(lightningTexture: SKTexture) {
        lightningIcon.texture = lightningTexture
        lightningIcon.size = CGSize(width: 18, height: 18)
        lightningIcon.zPosition = 90

        for label in [lightningCountLabel, lightningIntervalLabel] {
            label.fontSize = 14
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    private func setupOrbStatus(orbTexture: SKTexture) {
        orbIcon.texture = orbTexture
        orbIcon.size = CGSize(width: 18, height: 18)
        orbIcon.zPosition = 90

        for label in [orbCountLabel, orbSpeedLabel] {
            label.fontSize = 14
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    private func setupBeamStatus(beamTexture: SKTexture) {
        beamIcon.texture = beamTexture
        beamIcon.size = CGSize(width: 18, height: 18)
        beamIcon.zPosition = 90

        for label in [beamKillLabel, beamIntervalLabel] {
            label.fontSize = 14
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    private func hitArea(for label: SKLabelNode) -> CGRect {
        label.frame.insetBy(dx: -26, dy: -12)
    }

    private func showLevelUpOption(_ option: LevelUpOption, yPosition: CGFloat, level: Int) {
        let label = label(for: option)
        let icon = icon(for: option)

        label.text = option.title(level: level)
        label.position = CGPoint(x: -76, y: yPosition)
        icon.position = CGPoint(x: -110, y: yPosition)

        label.run(SKAction.fadeIn(withDuration: 0.2))
        icon.run(SKAction.fadeIn(withDuration: 0.2))
    }

    private func label(for option: LevelUpOption) -> SKLabelNode {
        switch option {
        case .fireRate:
            return fireRateLabel
        case .extraFireball:
            return extraFireballLabel
        case .lightningBounce:
            return lightningBounceLabel
        case .lightningRate:
            return lightningRateLabel
        case .extraOrb:
            return extraOrbLabel
        case .orbitalSpeed:
            return orbitalSpeedLabel
        case .beamRate:
            return beamRateLabel
        case .beamKillCount:
            return beamKillCountLabel
        }
    }

    private func icon(for option: LevelUpOption) -> SKSpriteNode {
        switch option {
        case .fireRate:
            return fireRateIcon
        case .extraFireball:
            return extraFireballIcon
        case .lightningBounce:
            return lightningBounceIcon
        case .lightningRate:
            return lightningRateIcon
        case .extraOrb:
            return extraOrbIcon
        case .orbitalSpeed:
            return orbitalSpeedIcon
        case .beamRate:
            return beamRateIcon
        case .beamKillCount:
            return beamKillCountIcon
        }
    }

    private var levelUpOptionLabels: [SKLabelNode] {
        [fireRateLabel, extraFireballLabel, lightningBounceLabel, lightningRateLabel, extraOrbLabel, orbitalSpeedLabel, beamRateLabel, beamKillCountLabel]
    }

    private var levelUpOptionIcons: [SKSpriteNode] {
        [fireRateIcon, extraFireballIcon, lightningBounceIcon, lightningRateIcon, extraOrbIcon, orbitalSpeedIcon, beamRateIcon, beamKillCountIcon]
    }

    private static func formattedSeconds(_ value: TimeInterval) -> String {
        if value >= 1 {
            return String(format: "%.1f", value)
        }

        return String(format: "%.2f", value)
    }

    private static let primaryTextColor = SKColor(calibratedRed: 0.96, green: 0.93, blue: 0.83, alpha: 1)
}
