# Boss-Battler Game - Development Guide

Welcome to your journey building a turn-based boss-battler game in Go! This guide will take you through 7 milestones, each building on the previous one, using Test-Driven Development (TDD) principles.

## What You're Building

A turn-based combat game featuring:
- **3-4 character party** with D&D-inspired stats (STR, DEX, CON, INT, WIS, CHA)
- **Isometric hex grid** positioning system
- **Boss enemies** that can occupy multiple hex tiles
- **Turn-based combat** with initiative, actions, and abilities
- **Graphical rendering** with sprites and UI

Future enhancement potential: dungeon crawler with random encounters.

## Tech Stack

### Ebitengine (Ebiten)
**What**: Production-ready 2D game engine for Go
**Why**:
- Pure Go implementation (mostly)
- Excellent cross-platform support (desktop, web/WASM, mobile)
- Auto-batching for performance
- Used in commercial games with 2M+ downloads
- Clean, simple API

**Installation**:
```bash
go get github.com/hajimehoshi/ebiten/v2
```

**Resources**:
- Official site: https://ebitengine.org/
- Tutorial tour: https://ebitengine.org/en/tour/
- Examples: https://github.com/hajimehoshi/ebiten/tree/main/examples

### Donburi ECS
**What**: Entity Component System library built for Ebiten
**Why**:
- Clean separation of data (components) and logic (systems)
- Perfect for games with many entities and behaviors
- No reflection (better performance)
- Event system built-in
- Made specifically for Go/Ebiten

**Installation**:
```bash
go get github.com/yohamta/donburi
```

**Resources**:
- GitHub: https://github.com/yohamta/donburi
- Examples: https://github.com/yohamta/donburi/tree/main/examples

### pmcxs/hexgrid
**What**: Hex grid coordinate library
**Why**:
- Based on Red Blob Games algorithms (proven, well-documented)
- Supports axial and cube coordinates
- Distance, neighbors, pathfinding functions
- Clean API

**Installation**:
```bash
go get github.com/pmcxs/hexgrid
```

**Resources**:
- GitHub: https://github.com/pmcxs/hexgrid
- Red Blob Games hex guide: https://www.redblobgames.com/grids/hexagons/

## Project Structure

```
game/
├── cmd/
│   └── game/
│       └── main.go              # Entry point
├── internal/
│   ├── components/              # ECS components (data)
│   │   ├── position.go
│   │   ├── stats.go
│   │   ├── health.go
│   │   └── components_test.go
│   ├── systems/                 # ECS systems (logic)
│   │   ├── render.go
│   │   ├── combat.go
│   │   └── input.go
│   ├── entities/                # Entity creation helpers
│   │   └── archetypes.go
│   ├── states/                  # Game state machine
│   │   └── battle.go
│   ├── commands/                # Action commands
│   │   ├── action.go
│   │   └── commands_test.go
│   ├── hex/                     # Hex grid integration
│   │   ├── layout.go
│   │   └── hex_test.go
│   └── combat/                  # Combat mechanics
│       ├── initiative.go
│       ├── damage.go
│       └── combat_test.go
├── assets/
│   ├── sprites/
│   └── ui/
├── docs/                        # This guide
│   ├── 00-overview.md          # You are here
│   ├── 01-milestone-hello-ebiten.md
│   ├── 02-milestone-hex-grid.md
│   ├── 03-milestone-ecs-setup.md
│   ├── 04-milestone-turn-system.md
│   ├── 05-milestone-combat-mechanics.md
│   ├── 06-milestone-input-actions.md
│   └── 07-milestone-demo-battle.md
├── go.mod
├── .gitignore
└── README.md
```

## Initial Setup

### 1. Initialize Go Module

```bash
cd /Users/dybeck/git/private/game
go mod init github.com/yourusername/boss-battler
```

Replace `yourusername` with your GitHub username (or use any module path).

### 2. Install Dependencies

```bash
go get github.com/hajimehoshi/ebiten/v2
go get github.com/yohamta/donburi
go get github.com/pmcxs/hexgrid
```

### 3. Verify Installation

Create a simple test file to verify Ebiten works:

```bash
mkdir -p cmd/game
```

Then create `cmd/game/main.go` with a minimal window (you'll do this in Milestone 1).

## Development Workflow (TDD)

Each milestone follows this pattern:

1. **Read the milestone document** - Understand the goal and tasks
2. **Run the tests** - They will fail (initial test files are already created!)
3. **Implement features** - Make the tests pass
4. **Check off tasks** - Mark completed items with ✅
5. **Test manually** - Run the game and verify visually
6. **Move to next milestone** - Build on what you've learned

**Test files already created for you:**
- `internal/hex/hex_test.go` - Hex coordinate conversion tests
- `internal/components/components_test.go` - Component behavior tests
- `internal/combat/combat_test.go` - Combat mechanics tests
- `internal/commands/commands_test.go` - Action command tests

These tests will initially fail because you haven't implemented the features yet. That's the TDD way!

### Running Tests

```bash
# Run all tests
go test ./...

# Quick test with make
make test

# Run tests for a specific package
go test ./internal/hex

# Run with verbose output
go test -v ./internal/combat

# Run a specific test
go test -run TestDamageCalculation ./internal/combat

# Generate coverage report
make coverage
```

### Running the Game

```bash
go run cmd/game/main.go
```

## Milestones Overview

### Milestone 1: Hello Ebiten
**Goal**: Get a window open and understand the game loop
**Duration**: ~1 hour
**You'll learn**: Ebiten basics, Update/Draw pattern, FPS/TPS

### Milestone 2: Hex Grid
**Goal**: Render a hex grid and detect mouse clicks
**Duration**: ~2-3 hours
**You'll learn**: Coordinate systems, rendering, mouse input

### Milestone 3: ECS Setup
**Goal**: Create entities with components using Donburi
**Duration**: ~2 hours
**You'll learn**: Entity Component System pattern, queries

### Milestone 4: Turn System
**Goal**: Implement turn-based state machine
**Duration**: ~2 hours
**You'll learn**: State machines, initiative, turn order

### Milestone 5: Combat Mechanics
**Goal**: D&D-style combat calculations
**Duration**: ~3 hours
**You'll learn**: Stats, modifiers, damage calculation, effects

### Milestone 6: Input & Actions
**Goal**: Player input and action selection
**Duration**: ~2-3 hours
**You'll learn**: Command pattern, action queuing, UI

### Milestone 7: Demo Battle
**Goal**: Playable battle with party vs boss
**Duration**: ~3-4 hours
**You'll learn**: AI, multi-tile entities, polish

**Total estimated time**: 15-20 hours

## Learning Resources

### Go Specific
- Official Go Tour: https://go.dev/tour/
- Effective Go: https://go.dev/doc/effective_go
- Table-driven tests: https://go.dev/wiki/TableDrivenTests

### Game Development
- Ebiten examples: https://ebitengine.org/en/examples/
- Red Blob Games (hex grids): https://www.redblobgames.com/grids/hexagons/
- ECS pattern: https://github.com/SanderMertens/ecs-faq

### D&D Mechanics (for combat inspiration)
- SRD 5.1: https://dnd.wizards.com/resources/systems-reference-document
- Basic combat rules: https://www.dndbeyond.com/sources/basic-rules

## Tips for Success

1. **Follow the order** - Each milestone builds on the previous
2. **Write tests first** - Even if they fail, they guide your implementation
3. **Commit frequently** - Use conventional commits (see CLAUDE.md)
4. **Ask questions** - Use comments in tests to clarify what's expected
5. **Experiment** - Tweak values, try different approaches
6. **Have fun** - This is a learning project, enjoy the process

## Getting Help

If you get stuck:
1. Check the test files for hints
2. Review the milestone doc carefully
3. Look at Ebiten/Donburi examples
4. Read Red Blob Games articles
5. Check Go documentation

## Ready to Start?

Head to [Milestone 1: Hello Ebiten](01-milestone-hello-ebiten.md) to begin!
