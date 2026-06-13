struct ProgressionSystem {
    let tuning: GameTuning

    func randomLevelUpOptions(
        from availableOptions: [LevelUpOption],
        hasSkeletons: Bool
    ) -> [LevelUpOption] {
        let optionCount = shouldShowExtraLevelUpOption() ? 3 : 2
        var selectedOptions = Array(
            availableOptions
                .filter { $0 != .halveSkeletons }
                .shuffled()
                .prefix(optionCount)
        )

        if shouldShowHalveHordeOption(hasSkeletons: hasSkeletons), !selectedOptions.isEmpty {
            selectedOptions[Int.random(in: selectedOptions.indices)] = .halveSkeletons
        }

        return selectedOptions
    }

    func shouldShowHalveHordeOption(hasSkeletons: Bool) -> Bool {
        guard hasSkeletons else {
            return false
        }

        return randomChance(
            numerator: tuning.progression.halveHordeChanceNumerator,
            denominator: tuning.progression.halveHordeChanceDenominator
        )
    }

    func shouldShowExtraLevelUpOption() -> Bool {
        randomChance(
            numerator: tuning.progression.extraLevelUpOptionChanceNumerator,
            denominator: tuning.progression.extraLevelUpOptionChanceDenominator
        )
    }

    private func randomChance(numerator: Int, denominator: Int) -> Bool {
        guard numerator > 0, denominator > 0 else {
            return false
        }

        return Int.random(in: 1...denominator) <= numerator
    }
}
