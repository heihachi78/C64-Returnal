# C64-Returnal

C64-Returnal is a Go/Ebitengine desktop survival game. The player controls a mage on an infinite C64-style grass field while auto-cast weapons fight escalating skeleton hordes.

The game opens a resizable 1200x900 window and runs the simulation at 120 ticks per second.

## Features

- Go module targeting Go 1.24.0, using Ebitengine v2.9.3.
- In-game pixel art generated in Go; embedded PNG files are used only for application window icons.
- Infinite scrolling grass field centered on the player.
- Auto-cast fireballs, chain lightning, orbital orbs, beams, meteors, and an unlockable death wave.
- Level-up choices, queued level-up overlays, coin-funded redraws, and chest rewards.
- Skeleton variants with current hit points of 1 regular, 3 red, 7 purple, 29 black, and 1000 blue.
- Dynamic spawn pressure based on each level's peak actual damage and current enemy HP budget, with `SkeletonHPPerSecondLevelUpBonus` defaulting to a +1 HP/sec bump per level-up; default tuning caps spawns at 999 active skeletons and one queued spawn per tick.
- A parallel skeleton movement path exists for tuning that allows at least 1024 active skeletons; damage, kills, rewards, and overlays stay on the main update path.
- Tests for the entrypoint, progression, combat, pickups, overlays, rendering fidelity, generated assets, spatial indexing, and branch coverage.

## Requirements

- Go 1.24.0 or newer.
- A desktop platform supported by Ebitengine.

## Run

```sh
go run ./cmd/c64-returnal
```

## Build

Build for the current platform:

```sh
go build -o build/c64-returnal ./cmd/c64-returnal
```

Build a Windows amd64 executable:

```sh
GOOS=windows GOARCH=amd64 go build -o build/c64-returnal-windows-amd64.exe ./cmd/c64-returnal
```

Linux-to-macOS cross-build commands are not listed here because they do not compile in this checkout with Ebitengine v2.9.3. Build on macOS with the current-platform command instead.

## Test

```sh
go test ./...
```

## Controls

| Action | Control |
| --- | --- |
| Move | Arrow keys |
| Choose level-up option 1 | `Q` or left-click the option |
| Choose level-up option 2 | `A` or left-click the option |
| Choose level-up option 3 | `C` or left-click the option |
| Choose level-up option 4 | `X` or left-click the option |
| Redraw level-up options | `R` or left-click the redraw control |
| Advance chest reward overlay | `Q` |
| Restart after game over | Left-click `RESTART` |
| Exit after game over | Left-click `EXIT` |
| Development: jump to level 100, add 5000 coins, and spawn 40 gold chests | `0` or keypad `0` |

## Gameplay

You start at level 1 with three lives, one regular skeleton, one off-screen coin, and one auto-cast fireball. Skeletons spawn outside the viewport and chase the player. Defeated skeletons grant experience; the experience required for a level is `level * level / 2`, with a minimum of 1.

Each level-up presents two random options by default. There is a 25% chance to present three options. If the horde is non-empty, there is a 5% chance that one visible option is replaced with `HALVE HORDE`. The UI can show up to four options because an affordable death scroll may be appended to the random choices.

Level-up options can increase fireball count or rate, add a life, learn lightning, learn orbital orbs, learn beams, learn meteors, improve learned weapons, halve the current horde, or buy a death scroll. A death scroll costs 1000 coins and appears only when affordable; buying five scrolls unlocks the death wave. Once unlocked, the death wave casts every 30 seconds and halves the current HP of non-regular skeletons it touches, leaving at least 1 HP.

One coin is spawned for each level, including level 1. Coins spawn outside the visible viewport, can be collected once, and grant 1 to 100 coins. Redrawing level-up choices costs coins equal to the current player level, with a minimum cost of 1.

Chests spawn outside the visible viewport at kill milestones. A chest milestone is checked every 250 kills. Bronze chests appear on non-silver, non-gold milestones through level 33. Silver chests appear every 1000 kills through level 55. Gold chests appear every 2500 kills at any level. Bronze chests apply one random available upgrade for one learned skill, silver chests apply all currently available upgrades for one learned skill, and gold chests do that for up to two learned skills.

Skeleton spawning uses an HP-per-second budget. The spawn planner spends that budget on the strongest affordable skeleton variants first, then regular skeletons for the remainder. Red, purple, black, and blue variants are slower than regular skeletons and have higher HP; purple, black, and blue variants also grant more experience.

## Project Layout

```text
cmd/c64-returnal/      Application entrypoint and entrypoint tests
internal/game/         Game simulation, rendering, input, generated assets, and tests
internal/game/assets/  Embedded application icon PNGs
```

## Architecture

Runtime state is split between `Game`, `Session`, compact entity structs, and `Progression`. Tuning values live in `DefaultTuning`, while focused files handle progression, spawning, combat, pickups, overlays, rendering, generated assets, and input.

Skeleton movement is the only gameplay update currently parallelized. It runs in chunks once the active skeleton count reaches the configured threshold, though the default active-skeleton cap is lower than that threshold. State-changing work such as damage, kills, level-up queuing, chest spawning, coin collection, and overlay transitions remains serialized.
