import SpriteKit

extension GameScene {
    func queueLevelUpChoices(_ count: Int) {
        guard count > 0, !session.isGameOver else {
            return
        }

        let firstQueuedLevel = session.progression.level - count + 1
        session.pendingLevelUpLevels.append(contentsOf: firstQueuedLevel...session.progression.level)
        presentNextLevelUpChoiceIfNeeded()
    }

    func presentNextLevelUpChoiceIfNeeded() {
        guard !session.isGameOver, !session.isLevelUpChoiceActive, let level = session.pendingLevelUpLevels.first else {
            return
        }

        session.isLevelUpChoiceActive = true
        worldNode.isPaused = true
        session.pressedKeys.removeAll()
        stopPlayerAnimation()
        hud.showLevelUp(
            level: level,
            options: randomLevelUpOptions(),
            beamKillUpgradeBonus: session.progression.beamKillUpgradeBonus
        )
    }

    func applyLevelUpOption(_ option: LevelUpOption) {
        applyUpgradeEffect(option)

        syncOrbitalOrbCount()
        updateHUDProgress()

        if !session.pendingLevelUpLevels.isEmpty {
            session.pendingLevelUpLevels.removeFirst()
        }

        session.isLevelUpChoiceActive = false
        worldNode.isPaused = false
        hud.hideLevelUp()
        presentNextLevelUpChoiceIfNeeded()
    }

    func applyUpgradeEffect(_ option: LevelUpOption) {
        switch option {
        case .extraLife:
            session.playerLives += 1
        case .halveSkeletons:
            halveSkeletons()
        default:
            session.progression.applyLevelUpOption(option)
        }
    }

    func halveSkeletons() {
        let killCount = skeletons.count / 2

        guard killCount > 0 else {
            return
        }

        let targets = Array(skeletons.shuffled().prefix(killCount))
        var levelUpCount = 0

        for target in targets {
            levelUpCount += destroySkeleton(target, shouldTriggerLevelUpChoice: false, shouldUpdateHUD: false)
        }

        updateHUDProgress()
        queueLevelUpChoices(levelUpCount)
    }

    func selectLevelUpOption(with keyCode: UInt16) {
        guard let index = inputController.levelUpOptionIndex(for: keyCode),
              let option = hud.levelUpOption(atIndex: index) else {
            return
        }

        applyLevelUpOption(option)
    }


    func randomLevelUpOptions() -> [LevelUpOption] {
        ProgressionSystem(tuning: tuning).randomLevelUpOptions(
            from: session.progression.availableLevelUpOptions,
            hasSkeletons: !skeletons.isEmpty
        )
    }

    func shouldShowHalveHordeOption() -> Bool {
        ProgressionSystem(tuning: tuning).shouldShowHalveHordeOption(hasSkeletons: !skeletons.isEmpty)
    }

    func shouldShowExtraLevelUpOption() -> Bool {
        ProgressionSystem(tuning: tuning).shouldShowExtraLevelUpOption()
    }


}
