import SpriteKit

extension GameScene {
    func chestTier(for milestone: Int) -> ChestTier? {
        ChestSystem(tuning: tuning).tier(
            for: milestone,
            playerLevel: session.progression.level
        )
    }

    func spawnChest(tier: ChestTier) {
        guard let texture = chestTextures[tier] else {
            return
        }

        let node = SKSpriteNode(texture: texture)
        node.size = CGSize(width: 32, height: 28)
        node.position = randomChestPosition()
        node.zPosition = 8.5

        worldNode.addChild(node)
        chests.append(Chest(node: node, tier: tier))
    }

    func randomChestPosition() -> CGPoint {
        let halfWidth = max(48, size.width / 2 - tuning.chest.spawnMargin)
        let halfHeight = max(48, size.height / 2 - tuning.chest.spawnMargin)
        let minimumDistanceSquared = tuning.chest.pickupDistance * tuning.chest.pickupDistance * 4

        for _ in 0..<12 {
            let position = CGPoint(
                x: player.position.x + CGFloat.random(in: -halfWidth...halfWidth),
                y: player.position.y + CGFloat.random(in: -halfHeight...halfHeight)
            )

            if position.distanceSquared(to: player.position) >= minimumDistanceSquared {
                return position
            }
        }

        return CGPoint(x: player.position.x + halfWidth, y: player.position.y)
    }

    func checkChestPickups() {
        let pickupDistanceSquared = tuning.chest.pickupDistance * tuning.chest.pickupDistance

        for index in chests.indices.reversed() {
            guard chests[index].node.position.distanceSquared(to: player.position) <= pickupDistanceSquared else {
                continue
            }

            collectChest(at: index)
            return
        }
    }

    func collectChest(at index: Int) {
        let chest = chests[index]
        chest.node.removeAllActions()
        chest.node.removeFromParent()
        chests.remove(at: index)

        let items = applyChestReward(chest.tier)
        showChestReward(tier: chest.tier, items: items)
    }

    func applyChestReward(_ tier: ChestTier) -> [ChestRewardDisplayItem] {
        let learnedSkills = session.progression.learnedSkills
        let items: [ChestRewardDisplayItem]

        switch tier {
        case .bronze:
            guard let skill = learnedSkills.randomElement(),
                  let option = skill.upgradeOptions.randomElement() else {
                return []
            }

            items = [chestRewardItem(for: option, skill: skill)]
            applyUpgradeEffect(option)
        case .silver:
            guard let skill = learnedSkills.randomElement() else {
                return []
            }

            items = chestRewardItems(for: [skill])
            session.progression.upgradeAllProperties(for: skill)
        case .gold:
            let rewardedSkills = Array(learnedSkills.shuffled().prefix(2))
            items = chestRewardItems(for: rewardedSkills)
            session.progression.upgradeAllProperties(for: rewardedSkills)
        }

        syncOrbitalOrbCount()
        updateHUDProgress()

        return items
    }

    func chestRewardItems(for skills: [LearnedSkill]) -> [ChestRewardDisplayItem] {
        ChestSystem(tuning: tuning).rewardItems(
            for: skills,
            beamKillUpgradeBonus: session.progression.beamKillUpgradeBonus
        )
    }

    func chestRewardItem(for option: LevelUpOption, skill: LearnedSkill) -> ChestRewardDisplayItem {
        ChestRewardDisplayItem(
            option: option,
            title: option.title(beamKillBonus: skill == .beam ? session.progression.beamKillUpgradeBonus : nil)
        )
    }

    func showChestReward(tier: ChestTier, items: [ChestRewardDisplayItem]) {
        guard !items.isEmpty else {
            return
        }

        session.isChestRewardActive = true
        worldNode.isPaused = true
        session.pressedKeys.removeAll()
        stopPlayerAnimation()
        hud.showChestReward(tier: tier, items: items)
    }

    func advanceChestReward(with keyCode: UInt16) {
        guard inputController.isChestRewardAdvance(keyCode) else {
            return
        }

        session.isChestRewardActive = false
        worldNode.isPaused = false
        hud.hideChestReward()
        presentNextLevelUpChoiceIfNeeded()
    }


}
