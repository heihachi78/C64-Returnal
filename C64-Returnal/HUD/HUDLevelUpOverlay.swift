import SpriteKit

extension GameHUD {
    func showLevelUp(level: Int, options: [LevelUpOption], beamKillUpgradeBonus: Int) {
        hideLevelUpOptions()
        activeLevelUpOptions = Array(options.prefix(4))
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


    func hideLevelUpOptions() {
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


    func setupLevelUpLabels(fireballTexture: SKTexture, lightningTexture: SKTexture, orbTexture: SKTexture, beamTexture: SKTexture, meteorTexture: SKTexture, lifeTexture: SKTexture, skeletonTexture: SKTexture) {
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
        setupLevelUpKeyLabel(thirdLevelUpKeyLabel, text: "[C]")
        setupLevelUpKeyLabel(fourthLevelUpKeyLabel, text: "[X]")
    }


    func setupLevelUpOption(_ label: SKLabelNode, text: String) {
        label.text = text
        label.fontSize = 22
        label.fontColor = Self.primaryTextColor
        label.horizontalAlignmentMode = .left
        label.verticalAlignmentMode = .center
        label.position = CGPoint(x: -78, y: 0)
        label.zPosition = 100
        label.alpha = 0
    }

    func setupLevelUpKeyLabel(_ label: SKLabelNode, text: String) {
        label.text = text
        label.fontSize = 18
        label.fontColor = Self.keyHintTextColor
        label.horizontalAlignmentMode = .center
        label.verticalAlignmentMode = .center
        label.position = CGPoint(x: -150, y: 0)
        label.zPosition = 100
        label.alpha = 0
    }

    func setupLevelUpIcon(_ icon: SKSpriteNode, texture: SKTexture) {
        icon.texture = texture
        icon.size = CGSize(width: 24, height: 24)
        icon.position = CGPoint(x: -116, y: 0)
        icon.zPosition = 100
        icon.alpha = 0
    }


    func showLevelUpOption(_ option: LevelUpOption, index: Int, yPosition: CGFloat, beamKillUpgradeBonus: Int) {
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

    func label(for option: LevelUpOption) -> SKLabelNode {
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

    func icon(for option: LevelUpOption) -> SKSpriteNode {
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


    var levelUpOptionLabels: [SKLabelNode] {
        [
            fireRateLabel, extraFireballLabel, extraLifeLabel, halveSkeletonsLabel,
            learnLightningLabel, lightningBounceLabel, lightningRateLabel,
            learnOrbLabel, extraOrbLabel, orbitalSpeedLabel,
            learnBeamLabel, beamRateLabel, beamKillCountLabel,
            learnMeteorLabel, extraMeteorLabel, meteorRateLabel
        ]
    }

    var levelUpOptionIcons: [SKSpriteNode] {
        [
            fireRateIcon, extraFireballIcon, extraLifeIcon, halveSkeletonsIcon,
            learnLightningIcon, lightningBounceIcon, lightningRateIcon,
            learnOrbIcon, extraOrbIcon, orbitalSpeedIcon,
            learnBeamIcon, beamRateIcon, beamKillCountIcon,
            learnMeteorIcon, extraMeteorIcon, meteorRateIcon
        ]
    }

    var levelUpKeyLabels: [SKLabelNode] {
        [firstLevelUpKeyLabel, secondLevelUpKeyLabel, thirdLevelUpKeyLabel, fourthLevelUpKeyLabel]
    }

    static func levelUpOptionYPosition(for index: Int) -> CGFloat {
        -4 - CGFloat(index) * 52
    }


}
