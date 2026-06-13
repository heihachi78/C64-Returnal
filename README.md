# C64-Returnal

C64-Returnal is a native macOS SpriteKit survival game with code-generated pixel art, auto-casting weapons, level-up choices, and escalating skeleton hordes. You move the mage, survive as long as possible, and build a growing arsenal of fireballs, chain lightning, orbital orbs, beams, and meteors.

## Features

- Native Cocoa/SpriteKit macOS app with a resizable game window.
- Procedural pixel-art sprites and grass tiles generated in Swift.
- Infinite scrolling field centered on the player.
- Auto-targeting combat with fireballs, chain lightning, orbital orbs, beams, and meteors.
- Level progression with randomized upgrade choices.
- Gold coins that can be collected and spent to redraw level-up choices.
- Bronze, silver, and gold chest rewards at kill milestones.
- Enemy escalation from regular skeletons to red, purple, and black variants.
- HUD overlays for lives, level, experience, weapon status, chest rewards, level-ups, and game over.
- Unit tests for progression, combat targeting, spatial indexing, input bindings, and reward policies.

## Requirements

- macOS with an SDK/Xcode version that supports the project's configured deployment target.
- Xcode command line tools.
- Swift 5 as configured by the Xcode project.

The project currently sets `MACOSX_DEPLOYMENT_TARGET` to `26.5` in `C64-Returnal.xcodeproj`.

## Build And Run

Use the helper script:

```sh
./script/build_and_run.sh
```

The script builds the `C64-Returnal` scheme into `build/DerivedData` and launches the app.

Useful modes:

```sh
./script/build_and_run.sh run
./script/build_and_run.sh --verify
./script/build_and_run.sh --debug
./script/build_and_run.sh --logs
./script/build_and_run.sh --telemetry
./script/build_and_run.sh --test
```

You can also build directly with Xcode:

```sh
xcodebuild \
  -project C64-Returnal.xcodeproj \
  -scheme C64-Returnal \
  -configuration Debug \
  build
```

## Test

Run the test suite with:

```sh
./script/build_and_run.sh --test
```

Or directly:

```sh
xcodebuild \
  -project C64-Returnal.xcodeproj \
  -scheme C64-Returnal \
  -configuration Debug \
  build test
```

## Controls

| Action | Control |
| --- | --- |
| Move | Arrow keys |
| Pick level-up option 1 | `Q` |
| Pick level-up option 2 | `A` |
| Pick level-up option 3 | `C` |
| Redraw level-up options | `R` |
| Advance chest reward overlay | `Q` |
| Restart or exit after game over | Click the HUD option |

There is also a development shortcut bound to `1` / keypad `1` that kills all current enemies and grants their experience.

## Gameplay Notes

You start with three lives and a single auto-cast fireball. Skeletons spawn outside the viewport and chase the player. Defeated skeletons grant experience; enough experience opens a level-up overlay with two choices, sometimes three.

Upgrades can improve existing weapons, add lives, unlock new weapon families, or occasionally halve the current horde. Gold coins spawn once per level well outside the visible screen; collecting one grants 1 to 100 coins. During a level-up choice, redraw spends coins equal to your current level and replaces all visible upgrade options. As levels rise, skeleton spawn timing changes and stronger variants enter the rotation. Chests appear at kill milestones and grant skill-focused upgrades based on the weapons you have already learned.

## Project Layout

```text
C64-Returnal/
  App/              AppDelegate and macOS window setup
  Configuration/    Game tuning and default input bindings
  Entities/         Lightweight gameplay model types
  HUD/              SpriteKit HUD panels and overlays
  Input/            Keyboard binding and input helpers
  Rendering/        Pixel-art and infinite grass rendering
  Scenes/           Main GameScene
  State/            Session and progression state
  Systems/          Combat, spawning, chests, progression, and scene systems
C64-ReturnalTests/  XCTest coverage for gameplay policies and targeting
script/             Build, run, debug, logging, and test helper script
```

## Architecture

`GameScene` owns the SpriteKit scene graph and the live session state. Most gameplay behavior is split into focused scene extensions under `C64-Returnal/Systems`, while reusable policy logic lives in small structs such as `ProgressionSystem`, `ChestSystem`, and `SkeletonSpatialIndex`.

`GameConfiguration.defaultTuning` is the central place to adjust movement speed, spawn cadence, weapon timing, hit distances, chest milestones, and progression probabilities.

## License

No license file is currently included.
