enum GameOverOption {
    case restart
    case exit
}

enum LevelUpOption: CaseIterable, Hashable {
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
