import Foundation

struct InputBindings {
    let moveLeft: UInt16
    let moveRight: UInt16
    let moveDown: UInt16
    let moveUp: UInt16
    let firstLevelUpOption: UInt16
    let secondLevelUpOption: UInt16
    let thirdLevelUpOption: UInt16
    let fourthLevelUpOption: UInt16
    let advanceChestReward: UInt16
    let killAllAndGrantExperience: Set<UInt16>

    init(
        moveLeft: UInt16 = 123,
        moveRight: UInt16 = 124,
        moveDown: UInt16 = 125,
        moveUp: UInt16 = 126,
        firstLevelUpOption: UInt16 = 12,
        secondLevelUpOption: UInt16 = 0,
        thirdLevelUpOption: UInt16 = 16,
        fourthLevelUpOption: UInt16 = 7,
        advanceChestReward: UInt16 = 12,
        killAllAndGrantExperience: Set<UInt16> = [18, 83]
    ) {
        self.moveLeft = moveLeft
        self.moveRight = moveRight
        self.moveDown = moveDown
        self.moveUp = moveUp
        self.firstLevelUpOption = firstLevelUpOption
        self.secondLevelUpOption = secondLevelUpOption
        self.thirdLevelUpOption = thirdLevelUpOption
        self.fourthLevelUpOption = fourthLevelUpOption
        self.advanceChestReward = advanceChestReward
        self.killAllAndGrantExperience = killAllAndGrantExperience
    }
}

struct InputController {
    let bindings: InputBindings

    func isMovementKey(_ keyCode: UInt16) -> Bool {
        keyCode == bindings.moveLeft
            || keyCode == bindings.moveRight
            || keyCode == bindings.moveDown
            || keyCode == bindings.moveUp
    }

    func levelUpOptionIndex(for keyCode: UInt16) -> Int? {
        switch keyCode {
        case bindings.firstLevelUpOption:
            return 0
        case bindings.secondLevelUpOption:
            return 1
        case bindings.thirdLevelUpOption:
            return 2
        case bindings.fourthLevelUpOption:
            return 3
        default:
            return nil
        }
    }

    func isChestRewardAdvance(_ keyCode: UInt16) -> Bool {
        keyCode == bindings.advanceChestReward
    }

    func isKillAllAndGrantExperience(_ keyCode: UInt16) -> Bool {
        bindings.killAllAndGrantExperience.contains(keyCode)
    }
}
