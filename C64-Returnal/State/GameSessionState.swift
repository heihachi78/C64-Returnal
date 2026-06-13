import CoreGraphics
import Foundation

struct KillCounts {
    var totalSkeletons = 0
    var purpleSkeletons = 0
    var fireball = 0
    var lightning = 0
    var orbitalOrb = 0
    var beam = 0
    var meteor = 0
}

struct AnimationTimers {
    var skeletonTimer: TimeInterval = 0
    var skeletonFrameIndex = 0
    var fireballTimer: TimeInterval = 0
    var fireballFrameIndex = 0
    var orbitalOrbTimer: TimeInterval = 0
    var orbitalOrbFrameIndex = 0
    var meteorTimer: TimeInterval = 0
    var meteorFrameIndex = 0
}

struct CastTimers {
    var skeletonSpawn: TimeInterval = 0
    var fireball: TimeInterval = 0
    var lightning: TimeInterval = 0
    var beam: TimeInterval = 0
    var meteor: TimeInterval = 0
}

struct GameSessionState {
    var progression: Progression
    var pressedKeys = Set<UInt16>()
    var lastUpdateTime: TimeInterval = 0
    var animations = AnimationTimers()
    var casts = CastTimers()
    var playerHitInvulnerabilityTimer: TimeInterval = 0
    var playerLives: Int
    var currentPlayerMovementDirection: CGVector?
    var orbitalOrbAngle: CGFloat = 0
    var pendingLevelUpLevels = [Int]()
    var kills = KillCounts()
    var nextChestMilestone: Int
    var collectedCoins = 0
    var spawnedCoinLevels = Set<Int>()
    var isGameOver = false
    var isLevelUpChoiceActive = false
    var isChestRewardActive = false
    var isSceneConfigured = false

    private let tuning: GameTuning

    init(tuning: GameTuning = GameConfiguration.defaultTuning) {
        self.tuning = tuning
        progression = Progression(tuning: tuning)
        playerLives = tuning.player.initialLives
        nextChestMilestone = tuning.chest.bronzeKillInterval
    }

    mutating func reset() {
        progression.reset()
        pressedKeys.removeAll()
        lastUpdateTime = 0
        animations = AnimationTimers()
        casts = CastTimers()
        playerHitInvulnerabilityTimer = 0
        playerLives = tuning.player.initialLives
        currentPlayerMovementDirection = nil
        orbitalOrbAngle = 0
        pendingLevelUpLevels.removeAll(keepingCapacity: true)
        kills = KillCounts()
        nextChestMilestone = tuning.chest.bronzeKillInterval
        collectedCoins = 0
        spawnedCoinLevels.removeAll(keepingCapacity: true)
        isGameOver = false
        isLevelUpChoiceActive = false
        isChestRewardActive = false
    }

    mutating func registerKill(from attackKind: AttackKind?) {
        kills.totalSkeletons += 1

        switch attackKind {
        case .fireball:
            kills.fireball += 1
        case .lightning:
            kills.lightning += 1
        case .orbitalOrb:
            kills.orbitalOrb += 1
        case .beam:
            kills.beam += 1
        case .meteor:
            kills.meteor += 1
        case .none:
            break
        }
    }
}
