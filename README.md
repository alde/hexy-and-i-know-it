# Hexy and I Know It

[![Test & Build](https://github.com/alde/hexy-and-i-know-it/actions/workflows/test.yml/badge.svg)](https://github.com/alde/hexy-and-i-know-it/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alde/hexy-and-i-know-it)](https://goreportcard.com/report/github.com/alde/hexy-and-i-know-it)

A turn-based boss-battler game built in Go with an isometric hex grid, D&D-inspired combat mechanics, and ECS architecture.

## About

This is a learning project that demonstrates:
- 2D game development with [Ebitengine](https://ebitengine.org/)
- Entity Component System pattern using [Donburi](https://github.com/yohamta/donburi)
- Hex grid coordinate systems with [hexgrid](https://github.com/pmcxs/hexgrid)
- Turn-based combat mechanics (D&D 5E-inspired)
- Command pattern for actions
- State machine for game flow

## Features

- Turn-based tactical combat on a hex grid
- 4 character classes: Warrior, Rogue, Mage, Cleric
- D&D-style stats: STR, DEX, CON, INT, WIS, CHA
- Attack rolls vs Armor Class
- Critical hits and misses
- Boss enemies that occupy multiple hexes
- Initiative-based turn order
- Player-controlled party vs AI-controlled boss
- Victory/defeat conditions

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Basic understanding of Go

### Installation

1. Clone the repository:
```bash
cd /path/to/hexy-and-i-know-it
```

2. Install dependencies:
```bash
go get github.com/hajimehoshi/ebiten/v2
go get github.com/yohamta/donburi
go get github.com/pmcxs/hexgrid
```

3. Run the game:
```bash
go run cmd/game/main.go
```

## Learning Guide

This project includes a comprehensive step-by-step learning guide in the `docs/` folder:

- [00-overview.md](docs/00-overview.md) - Project overview and setup
- [01-milestone-hello-ebiten.md](docs/01-milestone-hello-ebiten.md) - Getting started with Ebiten
- [02-milestone-hex-grid.md](docs/02-milestone-hex-grid.md) - Implementing hex grid
- [03-milestone-ecs-setup.md](docs/03-milestone-ecs-setup.md) - Entity Component System
- [04-milestone-turn-system.md](docs/04-milestone-turn-system.md) - Turn-based mechanics
- [05-milestone-combat-mechanics.md](docs/05-milestone-combat-mechanics.md) - D&D-style combat
- [06-milestone-input-actions.md](docs/06-milestone-input-actions.md) - Player input & commands
- [07-milestone-demo-battle.md](docs/07-milestone-demo-battle.md) - Polish & victory conditions

Each milestone builds on the previous one and includes:
- Clear learning objectives
- Step-by-step implementation tasks
- TDD approach with test files
- Explanations of key concepts
- Common pitfalls and solutions

## Project Structure

```
hexy-and-i-know-it/
├── cmd/
│   └── game/
│       └── main.go              # Entry point
├── internal/
│   ├── components/              # ECS components (data)
│   ├── systems/                 # ECS systems (logic)
│   ├── entities/                # Entity factories
│   ├── states/                  # Game state machine
│   ├── commands/                # Action commands
│   ├── hex/                     # Hex grid utilities
│   └── combat/                  # Combat mechanics
├── assets/                      # Game assets (sprites, audio)
├── docs/                        # Learning guide
├── go.mod
└── README.md
```

## Controls

- **Mouse**: Click to select targets
- **W**: Wait/skip turn
- **SPACE**: Advance turn (for testing)
- **ESC**: Quit game
- **R**: Restart after victory/defeat

## Tech Stack

- **[Ebitengine](https://ebitengine.org/)**: Production-ready 2D game engine for Go
- **[Donburi](https://github.com/yohamta/donburi)**: ECS library designed for Ebiten
- **[hexgrid](https://github.com/pmcxs/hexgrid)**: Hex coordinate system utilities

## Development

### Quick Start with Make

```bash
make help      # Show all available commands
make test      # Run all tests
make build     # Build the game
make run       # Build and run the game
make coverage  # Generate test coverage report
make lint      # Run linter
make fmt       # Format code
```

### Running Tests

```bash
# Simple test run
go test ./...

# With coverage
make coverage

# Verbose with race detection
go test -v -race ./...
```

### Running with Debug Info

```bash
go run -tags=ebitenginedebug cmd/game/main.go
```

### Building

```bash
# Using make
make build

# Using go directly
go build -o hexy cmd/game/main.go
```

### Continuous Integration

This project uses GitHub Actions for CI/CD:
- ✅ Automated testing on Go 1.23 and 1.25
- ✅ Code formatting checks
- ✅ Linting with golangci-lint
- ✅ Cross-platform builds (Linux, macOS, Windows)
- ✅ Test coverage reporting

See [.github/workflows/test.yml](.github/workflows/test.yml) for the full configuration.

## Future Plans

- Dungeon crawler mode with random encounters
- More character classes and abilities
- Equipment and inventory system
- Character progression and leveling
- Multiple boss types
- Procedural dungeon generation
- Save/load functionality

## Contributing

This is a learning project, but improvements to the guide or code are welcome! Please feel free to:
- Report bugs
- Suggest enhancements
- Improve documentation
- Share your own variations

## License

This project is provided as-is for educational purposes. Feel free to learn from it, modify it, and make it your own.

## Acknowledgments

- Inspired by early Final Fantasy games and D&D
- Hex grid algorithms from [Red Blob Games](https://www.redblobgames.com/grids/hexagons/)
- D&D combat mechanics from the 5E SRD
- Built with the excellent Ebiten and Donburi libraries

## Resources

- **Ebiten**: https://ebitengine.org/
- **Donburi**: https://github.com/yohamta/donburi
- **Red Blob Games**: https://www.redblobgames.com/
- **D&D SRD**: https://dnd.wizards.com/resources/systems-reference-document

---

**Hexy and I Know It** - Because every hex knows it's sexy
