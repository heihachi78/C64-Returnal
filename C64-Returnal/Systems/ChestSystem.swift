struct ChestSystem {
    let tuning: GameTuning

    func tier(for milestone: Int, playerLevel: Int) -> ChestTier? {
        if milestone.isMultiple(of: tuning.chest.goldKillInterval) {
            return .gold
        }

        if milestone.isMultiple(of: tuning.chest.silverKillInterval) {
            guard playerLevel <= tuning.chest.silverMaximumLevel else {
                return nil
            }

            return .silver
        }

        guard playerLevel <= tuning.chest.bronzeMaximumLevel else {
            return nil
        }

        return .bronze
    }

    func rewardItems(
        for skills: [LearnedSkill],
        beamKillUpgradeBonus: Int
    ) -> [ChestRewardDisplayItem] {
        var items = [ChestRewardDisplayItem]()

        for skill in skills {
            for option in skill.upgradeOptions {
                items.append(
                    ChestRewardDisplayItem(
                        option: option,
                        title: option.title(beamKillBonus: skill == .beam ? beamKillUpgradeBonus : nil)
                    )
                )
            }
        }

        return items
    }
}
