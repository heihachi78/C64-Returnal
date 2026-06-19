# C64-Returnal

C64-Returnal is a platform-independent Go/Ebitengine survival game with code-generated pixel art, auto-casting weapons, level-up choices, chest rewards, coins, and escalating skeleton hordes.

## Features

- Cross-platform Go game using Ebitengine.
- Builds from the same source on macOS and Windows.
- Procedural C64-style pixel rendering, with no external art assets required.
- Infinite scrolling grass field centered on the player.
- Auto-targeting fireballs, chain lightning, orbital orbs, beams, and meteors.
- Level progression with randomized upgrade choices.
- Coins that can be collected and spent to redraw level-up choices.
- Bronze, silver, and gold chest rewards at kill milestones.
- Enemy escalation from regular skeletons to red, purple, black, and late-game blue giant variants.
- Parallel skeleton movement updates for large hordes while combat mutation remains deterministic on the main update thread.
- Go tests for progression, chest policy, and spatial indexing rules.

## Requirements

- Go 1.24 or newer.
- macOS, Windows, or another Ebitengine-supported desktop platform.

## Run

```sh
go run ./cmd/c64-returnal
```

## Build

Build for the current platform:

```sh
go build -o build/c64-returnal ./cmd/c64-returnal
```

Build for macOS:

```sh
GOOS=darwin GOARCH=arm64 go build -o build/c64-returnal-macos-arm64 ./cmd/c64-returnal
GOOS=darwin GOARCH=amd64 go build -o build/c64-returnal-macos-amd64 ./cmd/c64-returnal
```

Build for Windows:

```sh
GOOS=windows GOARCH=amd64 go build -o build/c64-returnal-windows-amd64.exe ./cmd/c64-returnal
```

## Test

```sh
go test ./...
```

## Controls

| Action | Control |
| --- | --- |
| Move | Arrow keys |
| Pick level-up option 1 | `Q` |
| Pick level-up option 2 | `A` |
| Pick level-up option 3 | `C` |
| Pick level-up option 4 | `X` |
| Redraw level-up options | `R` |
| Advance chest reward overlay | `Q` |
| Restart or exit after game over | Click the HUD option |
| Development kill-all shortcut | `1` / keypad `1` |

## Gameplay Notes

You start with three lives and a single auto-cast fireball. Skeletons spawn outside the viewport and chase the player. Defeated skeletons grant experience; enough experience opens a level-up overlay with two choices, sometimes three.

Upgrades can improve existing weapons, add lives, unlock new weapon families, or occasionally halve the current horde. Gold coins spawn once per level well outside the visible screen; collecting one grants 1 to 100 coins. During a level-up choice, redraw spends coins equal to your current level and replaces all visible upgrade options. As levels rise, skeleton spawn timing changes and stronger variants enter the rotation. After level 100, giant blue monsters can enter oversized hordes, thin the enemy count, and accelerate future spawns. Chests appear at kill milestones and grant skill-focused upgrades based on the weapons you have already learned.

## Project Layout

```text
cmd/c64-returnal/      Go application entrypoint
internal/game/         Portable game simulation, rendering, input, and tests
internal/game/assets/  Embedded application assets, including app icons
```

## Architecture

The Go version keeps the original game structure: a central `Tuning` table, a `Session` state, compact entity types, and focused systems for progression, spawning, combat, chests, coins, and rendering.

The most expensive horde operation, skeleton pursuit movement, is split into chunks and processed with goroutines when the enemy count crosses the configured threshold. State-changing operations such as damage, kills, chest spawning, and level-up queuing stay on the main update thread to avoid races and keep gameplay outcomes predictable.
