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
    case extraLife
    case halveSkeletons
    case learnLightning
    case lightningBounce
    case lightningRate
    case learnOrb
    case extraOrb
    case orbitalSpeed
    case learnBeam
    case beamRate
    case beamKillCount
    case learnMeteor
    case extraMeteor
    case meteorRate

    var title: String {
        title(beamKillBonus: nil)
    }

    func title(beamKillBonus: Int?) -> String {
        switch self {
        case .fireRate:
            return "FASTER FIRE"
        case .extraFireball:
            return "+1 FIREBALL"
        case .extraLife:
            return "+1 LIFE"
        case .halveSkeletons:
            return "HALVE HORDE"
        case .learnLightning:
            return "LEARN BOLT"
        case .lightningBounce:
            return "+1 CHAIN"
        case .lightningRate:
            return "FASTER BOLT"
        case .learnOrb:
            return "LEARN ORB"
        case .extraOrb:
            return "+1 ORB"
        case .orbitalSpeed:
            return "FASTER ORB"
        case .learnBeam:
            return "LEARN BEAM"
        case .beamRate:
            return "FASTER BEAM"
        case .beamKillCount:
            return "+\(beamKillBonus ?? 1) BEAM KILL"
        case .learnMeteor:
            return "LEARN METEOR"
        case .extraMeteor:
            return "+1 METEOR"
        case .meteorRate:
            return "FASTER METEOR"
        }
    }
}

struct ChestRewardDisplayItem {
    let option: LevelUpOption
    let title: String
}

final class GameHUD {
    private let topStatusBackground = SKShapeNode()
    private let combatStatusBackground = SKShapeNode()
    private let levelUpBackground = SKShapeNode()
    private let chestRewardBackground = SKShapeNode()
    private let gameOverBackground = SKShapeNode()
    private let gameOverLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let restartLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let exitLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let levelUpLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let chestRewardLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let chestRewardContinueLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireRateLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let extraFireballLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let extraLifeLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let halveSkeletonsLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let learnLightningLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningBounceLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningRateLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let learnOrbLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let extraOrbLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let orbitalSpeedLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let learnBeamLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamRateLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamKillCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let learnMeteorLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let extraMeteorLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let meteorRateLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let firstLevelUpKeyLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let secondLevelUpKeyLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let thirdLevelUpKeyLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireRateIcon = SKSpriteNode()
    private let extraFireballIcon = SKSpriteNode()
    private let extraLifeIcon = SKSpriteNode()
    private let halveSkeletonsIcon = SKSpriteNode()
    private let learnLightningIcon = SKSpriteNode()
    private let lightningBounceIcon = SKSpriteNode()
    private let lightningRateIcon = SKSpriteNode()
    private let learnOrbIcon = SKSpriteNode()
    private let extraOrbIcon = SKSpriteNode()
    private let orbitalSpeedIcon = SKSpriteNode()
    private let learnBeamIcon = SKSpriteNode()
    private let beamRateIcon = SKSpriteNode()
    private let beamKillCountIcon = SKSpriteNode()
    private let learnMeteorIcon = SKSpriteNode()
    private let extraMeteorIcon = SKSpriteNode()
    private let meteorRateIcon = SKSpriteNode()
    private let levelLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let experienceLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireballIcon = SKSpriteNode()
    private let fireballCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireballIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let fireballKillsLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningIcon = SKSpriteNode()
    private let lightningCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let lightningKillsLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let orbIcon = SKSpriteNode()
    private let orbCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let orbSpeedLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let orbKillsLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamIcon = SKSpriteNode()
    private let beamKillLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let beamKillsLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let meteorIcon = SKSpriteNode()
    private let meteorCountLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let meteorIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let meteorKillsLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let skeletonIcon = SKSpriteNode()
    private let skeletonAliveLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private let skeletonIntervalLabel = SKLabelNode(fontNamed: "Menlo-Bold")
    private var activeLevelUpOptions = [LevelUpOption]()
    private var activeChestRewardItems = [ChestRewardDisplayItem]()
    private var chestRewardNodes = [SKNode]()
    private var lifeIcons = [SKSpriteNode]()
    private var lifeTexture: SKTexture?
    private var lifeIconOrigin = CGPoint.zero
    private var currentLives = 0
    private var currentSceneSize = CGSize(width: 800, height: 600)
    private weak var parentNode: SKNode?

    func add(to parent: SKNode, fireballTexture: SKTexture, lightningTexture: SKTexture, orbTexture: SKTexture, beamTexture: SKTexture, meteorTexture: SKTexture, lifeTexture: SKTexture, skeletonTexture: SKTexture) {
        parentNode = parent
        self.lifeTexture = lifeTexture

        setupGameOverLabel()
        setupLevelUpLabels(fireballTexture: fireballTexture, lightningTexture: lightningTexture, orbTexture: orbTexture, beamTexture: beamTexture, meteorTexture: meteorTexture, lifeTexture: lifeTexture, skeletonTexture: skeletonTexture)
        setupChestRewardLabels()
        setupProgressLabels()
        setupBackgroundPanels()
        setupFireballStatus(fireballTexture: fireballTexture)
        setupLightningStatus(lightningTexture: lightningTexture)
        setupOrbStatus(orbTexture: orbTexture)
        setupBeamStatus(beamTexture: beamTexture)
        setupMeteorStatus(meteorTexture: meteorTexture)
        setupSkeletonStatus(skeletonTexture: skeletonTexture)

        parent.addChild(topStatusBackground)
        parent.addChild(combatStatusBackground)
        parent.addChild(levelUpBackground)
        parent.addChild(chestRewardBackground)
        parent.addChild(gameOverBackground)
        parent.addChild(gameOverLabel)
        parent.addChild(restartLabel)
        parent.addChild(exitLabel)
        parent.addChild(levelUpLabel)
        parent.addChild(chestRewardLabel)
        parent.addChild(chestRewardContinueLabel)
        parent.addChild(fireRateLabel)
        parent.addChild(extraFireballLabel)
        parent.addChild(extraLifeLabel)
        parent.addChild(halveSkeletonsLabel)
        parent.addChild(learnLightningLabel)
        parent.addChild(lightningBounceLabel)
        parent.addChild(lightningRateLabel)
        parent.addChild(learnOrbLabel)
        parent.addChild(extraOrbLabel)
        parent.addChild(orbitalSpeedLabel)
        parent.addChild(learnBeamLabel)
        parent.addChild(beamRateLabel)
        parent.addChild(beamKillCountLabel)
        parent.addChild(learnMeteorLabel)
        parent.addChild(extraMeteorLabel)
        parent.addChild(meteorRateLabel)
        parent.addChild(firstLevelUpKeyLabel)
        parent.addChild(secondLevelUpKeyLabel)
        parent.addChild(thirdLevelUpKeyLabel)
        parent.addChild(fireRateIcon)
        parent.addChild(extraFireballIcon)
        parent.addChild(extraLifeIcon)
        parent.addChild(halveSkeletonsIcon)
        parent.addChild(learnLightningIcon)
        parent.addChild(lightningBounceIcon)
        parent.addChild(lightningRateIcon)
        parent.addChild(learnOrbIcon)
        parent.addChild(extraOrbIcon)
        parent.addChild(orbitalSpeedIcon)
        parent.addChild(learnBeamIcon)
        parent.addChild(beamRateIcon)
        parent.addChild(beamKillCountIcon)
        parent.addChild(learnMeteorIcon)
        parent.addChild(extraMeteorIcon)
        parent.addChild(meteorRateIcon)
        parent.addChild(levelLabel)
        parent.addChild(experienceLabel)
        parent.addChild(fireballIcon)
        parent.addChild(fireballCountLabel)
        parent.addChild(fireballIntervalLabel)
        parent.addChild(fireballKillsLabel)
        parent.addChild(lightningIcon)
        parent.addChild(lightningCountLabel)
        parent.addChild(lightningIntervalLabel)
        parent.addChild(lightningKillsLabel)
        parent.addChild(orbIcon)
        parent.addChild(orbCountLabel)
        parent.addChild(orbSpeedLabel)
        parent.addChild(orbKillsLabel)
        parent.addChild(beamIcon)
        parent.addChild(beamKillLabel)
        parent.addChild(beamIntervalLabel)
        parent.addChild(beamKillsLabel)
        parent.addChild(meteorIcon)
        parent.addChild(meteorCountLabel)
        parent.addChild(meteorIntervalLabel)
        parent.addChild(meteorKillsLabel)
        parent.addChild(skeletonIcon)
        parent.addChild(skeletonAliveLabel)
        parent.addChild(skeletonIntervalLabel)
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

    func showLevelUp(level: Int, options: [LevelUpOption], beamKillUpgradeBonus: Int) {
        hideLevelUpOptions()
        activeLevelUpOptions = Array(options.prefix(3))
        levelUpLabel.text = "LEVEL \(level)"
        layoutLevelUpBackground(optionCount: activeLevelUpOptions.count)
        levelUpBackground.run(SKAction.fadeAlpha(to: Self.panelAlpha, duration: 0.2))
        levelUpLabel.setScale(0.75)
        levelUpLabel.run(
            SKAction.group([
                SKAction.fadeIn(withDuration: 0.2),
                SKAction.scale(to: 1, duration: 0.2)
            ])
        )

        for (index, option) in activeLevelUpOptions.enumerated() {
            showLevelUpOption(
                option,
                index: index,
                yPosition: Self.levelUpOptionYPosition(for: index),
                beamKillUpgradeBonus: beamKillUpgradeBonus
            )
        }
    }

    func hideLevelUp() {
        levelUpBackground.removeAllActions()
        levelUpBackground.alpha = 0
        levelUpLabel.removeAllActions()
        levelUpLabel.alpha = 0
        levelUpLabel.setScale(1)
        hideLevelUpOptions()
        activeLevelUpOptions.removeAll()
    }

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

        for keyLabel in levelUpKeyLabels {
            keyLabel.removeAllActions()
            keyLabel.alpha = 0
            keyLabel.setScale(1)
        }
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

    func levelUpOption(at point: CGPoint) -> LevelUpOption? {
        for option in activeLevelUpOptions {
            if hitArea(for: label(for: option)).contains(point) || icon(for: option).frame.insetBy(dx: -14, dy: -14).contains(point) {
                return option
            }
        }

        return nil
    }

    func levelUpOption(atIndex index: Int) -> LevelUpOption? {
        guard activeLevelUpOptions.indices.contains(index) else {
            return nil
        }

        return activeLevelUpOptions[index]
    }

    private func setupGameOverLabel() {
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

    private func setupBackgroundPanels() {
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

    private func setupLevelUpLabels(fireballTexture: SKTexture, lightningTexture: SKTexture, orbTexture: SKTexture, beamTexture: SKTexture, meteorTexture: SKTexture, lifeTexture: SKTexture, skeletonTexture: SKTexture) {
        levelUpLabel.fontSize = 40
        levelUpLabel.fontColor = Self.primaryTextColor
        levelUpLabel.horizontalAlignmentMode = .center
        levelUpLabel.verticalAlignmentMode = .center
        levelUpLabel.position = CGPoint(x: 0, y: 56)
        levelUpLabel.zPosition = 100
        levelUpLabel.alpha = 0

        setupLevelUpOption(fireRateLabel, text: LevelUpOption.fireRate.title)
        setupLevelUpOption(extraFireballLabel, text: LevelUpOption.extraFireball.title)
        setupLevelUpOption(extraLifeLabel, text: LevelUpOption.extraLife.title)
        setupLevelUpOption(halveSkeletonsLabel, text: LevelUpOption.halveSkeletons.title)
        setupLevelUpOption(learnLightningLabel, text: LevelUpOption.learnLightning.title)
        setupLevelUpOption(lightningBounceLabel, text: LevelUpOption.lightningBounce.title)
        setupLevelUpOption(lightningRateLabel, text: LevelUpOption.lightningRate.title)
        setupLevelUpOption(learnOrbLabel, text: LevelUpOption.learnOrb.title)
        setupLevelUpOption(extraOrbLabel, text: LevelUpOption.extraOrb.title)
        setupLevelUpOption(orbitalSpeedLabel, text: LevelUpOption.orbitalSpeed.title)
        setupLevelUpOption(learnBeamLabel, text: LevelUpOption.learnBeam.title)
        setupLevelUpOption(beamRateLabel, text: LevelUpOption.beamRate.title)
        setupLevelUpOption(beamKillCountLabel, text: LevelUpOption.beamKillCount.title)
        setupLevelUpOption(learnMeteorLabel, text: LevelUpOption.learnMeteor.title)
        setupLevelUpOption(extraMeteorLabel, text: LevelUpOption.extraMeteor.title)
        setupLevelUpOption(meteorRateLabel, text: LevelUpOption.meteorRate.title)
        setupLevelUpIcon(fireRateIcon, texture: fireballTexture)
        setupLevelUpIcon(extraFireballIcon, texture: fireballTexture)
        setupLevelUpIcon(extraLifeIcon, texture: lifeTexture)
        setupLevelUpIcon(halveSkeletonsIcon, texture: skeletonTexture)
        setupLevelUpIcon(learnLightningIcon, texture: lightningTexture)
        setupLevelUpIcon(lightningBounceIcon, texture: lightningTexture)
        setupLevelUpIcon(lightningRateIcon, texture: lightningTexture)
        setupLevelUpIcon(learnOrbIcon, texture: orbTexture)
        setupLevelUpIcon(extraOrbIcon, texture: orbTexture)
        setupLevelUpIcon(orbitalSpeedIcon, texture: orbTexture)
        setupLevelUpIcon(learnBeamIcon, texture: beamTexture)
        setupLevelUpIcon(beamRateIcon, texture: beamTexture)
        setupLevelUpIcon(beamKillCountIcon, texture: beamTexture)
        setupLevelUpIcon(learnMeteorIcon, texture: meteorTexture)
        setupLevelUpIcon(extraMeteorIcon, texture: meteorTexture)
        setupLevelUpIcon(meteorRateIcon, texture: meteorTexture)
        setupLevelUpKeyLabel(firstLevelUpKeyLabel, text: "[Q]")
        setupLevelUpKeyLabel(secondLevelUpKeyLabel, text: "[A]")
        setupLevelUpKeyLabel(thirdLevelUpKeyLabel, text: "[Y]")
    }

    private func setupChestRewardLabels() {
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

    private func setupGameOverOption(_ label: SKLabelNode, text: String, yPosition: CGFloat) {
        label.text = text
        label.fontSize = 22
        label.fontColor = Self.deathTextColor
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
        label.position = CGPoint(x: -78, y: 0)
        label.zPosition = 100
        label.alpha = 0
    }

    private func setupLevelUpKeyLabel(_ label: SKLabelNode, text: String) {
        label.text = text
        label.fontSize = 18
        label.fontColor = Self.keyHintTextColor
        label.horizontalAlignmentMode = .center
        label.verticalAlignmentMode = .center
        label.position = CGPoint(x: -150, y: 0)
        label.zPosition = 100
        label.alpha = 0
    }

    private func setupLevelUpIcon(_ icon: SKSpriteNode, texture: SKTexture) {
        icon.texture = texture
        icon.size = CGSize(width: 24, height: 24)
        icon.position = CGPoint(x: -116, y: 0)
        icon.zPosition = 100
        icon.alpha = 0
    }

    private func setupProgressLabels() {
        for label in [levelLabel, experienceLabel] {
            label.fontSize = 16
            label.fontColor = Self.primaryTextColor
            label.horizontalAlignmentMode = .left
            label.verticalAlignmentMode = .center
            label.zPosition = 90
        }
    }

    private func syncLifeIcons() {
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

    private func layoutLifeIcons() {
        for (index, icon) in lifeIcons.enumerated() {
            let column = index % Self.lifeIconsPerRow
            let row = index / Self.lifeIconsPerRow
            icon.position = CGPoint(
                x: lifeIconOrigin.x + CGFloat(column) * Self.lifeIconSpacing,
                y: lifeIconOrigin.y - CGFloat(row) * Self.lifeIconSpacing
            )
        }
    }

    private func layoutTopStatusBackground(left: CGFloat? = nil, top: CGFloat? = nil) {
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

    private func setupFireballStatus(fireballTexture: SKTexture) {
        fireballIcon.texture = fireballTexture
        fireballIcon.size = CGSize(width: 18, height: 18)
        fireballIcon.zPosition = 90

        setupCombatStatusLabels([fireballCountLabel, fireballIntervalLabel, fireballKillsLabel])
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

        setupCombatStatusLabels([lightningCountLabel, lightningIntervalLabel, lightningKillsLabel])
    }

    private func setupOrbStatus(orbTexture: SKTexture) {
        orbIcon.texture = orbTexture
        orbIcon.size = CGSize(width: 18, height: 18)
        orbIcon.zPosition = 90

        setupCombatStatusLabels([orbCountLabel, orbSpeedLabel, orbKillsLabel])
    }

    private func setupBeamStatus(beamTexture: SKTexture) {
        beamIcon.texture = beamTexture
        beamIcon.size = CGSize(width: 18, height: 18)
        beamIcon.zPosition = 90

        setupCombatStatusLabels([beamKillLabel, beamIntervalLabel, beamKillsLabel])
    }

    private func setupMeteorStatus(meteorTexture: SKTexture) {
        meteorIcon.texture = meteorTexture
        meteorIcon.size = CGSize(width: 18, height: 18)
        meteorIcon.zPosition = 90

        setupCombatStatusLabels([meteorCountLabel, meteorIntervalLabel, meteorKillsLabel])
    }

    private func setupCombatStatusLabels(_ labels: [SKLabelNode]) {
        for label in labels {
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

    private func showLevelUpOption(_ option: LevelUpOption, index: Int, yPosition: CGFloat, beamKillUpgradeBonus: Int) {
        let label = label(for: option)
        let icon = icon(for: option)
        let keyLabel = levelUpKeyLabels[index]

        label.text = option.title(beamKillBonus: beamKillUpgradeBonus)
        label.position = CGPoint(x: -78, y: yPosition)
        icon.position = CGPoint(x: -116, y: yPosition)
        keyLabel.position = CGPoint(x: -150, y: yPosition)

        label.run(SKAction.fadeIn(withDuration: 0.2))
        icon.run(SKAction.fadeIn(withDuration: 0.2))
        keyLabel.run(SKAction.fadeIn(withDuration: 0.2))
    }

    private func label(for option: LevelUpOption) -> SKLabelNode {
        switch option {
        case .fireRate:
            return fireRateLabel
        case .extraFireball:
            return extraFireballLabel
        case .extraLife:
            return extraLifeLabel
        case .halveSkeletons:
            return halveSkeletonsLabel
        case .learnLightning:
            return learnLightningLabel
        case .lightningBounce:
            return lightningBounceLabel
        case .lightningRate:
            return lightningRateLabel
        case .learnOrb:
            return learnOrbLabel
        case .extraOrb:
            return extraOrbLabel
        case .orbitalSpeed:
            return orbitalSpeedLabel
        case .learnBeam:
            return learnBeamLabel
        case .beamRate:
            return beamRateLabel
        case .beamKillCount:
            return beamKillCountLabel
        case .learnMeteor:
            return learnMeteorLabel
        case .extraMeteor:
            return extraMeteorLabel
        case .meteorRate:
            return meteorRateLabel
        }
    }

    private func icon(for option: LevelUpOption) -> SKSpriteNode {
        switch option {
        case .fireRate:
            return fireRateIcon
        case .extraFireball:
            return extraFireballIcon
        case .extraLife:
            return extraLifeIcon
        case .halveSkeletons:
            return halveSkeletonsIcon
        case .learnLightning:
            return learnLightningIcon
        case .lightningBounce:
            return lightningBounceIcon
        case .lightningRate:
            return lightningRateIcon
        case .learnOrb:
            return learnOrbIcon
        case .extraOrb:
            return extraOrbIcon
        case .orbitalSpeed:
            return orbitalSpeedIcon
        case .learnBeam:
            return learnBeamIcon
        case .beamRate:
            return beamRateIcon
        case .beamKillCount:
            return beamKillCountIcon
        case .learnMeteor:
            return learnMeteorIcon
        case .extraMeteor:
            return extraMeteorIcon
        case .meteorRate:
            return meteorRateIcon
        }
    }

    private func showChestRewardItems(_ items: [ChestRewardDisplayItem]) {
        for (index, item) in items.enumerated() {
            showChestRewardItem(item, index: index, itemCount: items.count)
        }
    }

    private func showChestRewardItem(_ item: ChestRewardDisplayItem, index: Int, itemCount: Int) {
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

    private func layoutChestRewardItems() {
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

    private var levelUpOptionLabels: [SKLabelNode] {
        [
            fireRateLabel, extraFireballLabel, extraLifeLabel, halveSkeletonsLabel,
            learnLightningLabel, lightningBounceLabel, lightningRateLabel,
            learnOrbLabel, extraOrbLabel, orbitalSpeedLabel,
            learnBeamLabel, beamRateLabel, beamKillCountLabel,
            learnMeteorLabel, extraMeteorLabel, meteorRateLabel
        ]
    }

    private var levelUpOptionIcons: [SKSpriteNode] {
        [
            fireRateIcon, extraFireballIcon, extraLifeIcon, halveSkeletonsIcon,
            learnLightningIcon, lightningBounceIcon, lightningRateIcon,
            learnOrbIcon, extraOrbIcon, orbitalSpeedIcon,
            learnBeamIcon, beamRateIcon, beamKillCountIcon,
            learnMeteorIcon, extraMeteorIcon, meteorRateIcon
        ]
    }

    private var levelUpKeyLabels: [SKLabelNode] {
        [firstLevelUpKeyLabel, secondLevelUpKeyLabel, thirdLevelUpKeyLabel]
    }

    private var centeredPanelWidth: CGFloat {
        min(max(360, currentSceneSize.width - 48), 620)
    }

    private func layoutLevelUpBackground(optionCount: Int) {
        let clampedOptionCount = max(2, optionCount)
        let height = CGFloat(178 + max(0, clampedOptionCount - 2) * 52)

        setPanel(
            levelUpBackground,
            rect: CGRect(x: -centeredPanelWidth / 2, y: 90 - height, width: centeredPanelWidth, height: height),
            cornerRadius: Self.panelCornerRadius
        )
    }

    private func layoutChestRewardBackground(itemCount: Int) {
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

    private static func levelUpOptionYPosition(for index: Int) -> CGFloat {
        -4 - CGFloat(index) * 52
    }

    private static func chestRewardItemYPosition(for index: Int, itemCount: Int) -> CGFloat {
        let spacing: CGFloat = 30
        let totalHeight = CGFloat(max(0, itemCount - 1)) * spacing

        return totalHeight / 2 - CGFloat(index) * spacing - 12
    }

    private static func formattedSeconds(_ value: TimeInterval) -> String {
        if value >= 1 {
            return String(format: "%.1f", value)
        }

        return String(format: "%.2f", value)
    }

    private static let primaryTextColor = SKColor(calibratedRed: 0.96, green: 0.93, blue: 0.83, alpha: 1)
    private static let keyHintTextColor = SKColor(calibratedRed: 1.0, green: 0.86, blue: 0.26, alpha: 1)
    private static let deathTextColor = SKColor(calibratedRed: 0.95, green: 0.05, blue: 0.08, alpha: 1)
    private static let panelColor = SKColor(calibratedWhite: 0.02, alpha: 1)
    private static let panelAlpha: CGFloat = 0.62
    private static let panelCornerRadius: CGFloat = 6
    private static let lifeIconSize: CGFloat = 14
    private static let lifeIconsPerRow = 12
    private static let lifeIconSpacing: CGFloat = 16

    private func setPanel(_ panel: SKShapeNode, rect: CGRect, cornerRadius: CGFloat) {
        panel.path = CGPath(
            roundedRect: rect,
            cornerWidth: cornerRadius,
            cornerHeight: cornerRadius,
            transform: nil
        )
    }
}
