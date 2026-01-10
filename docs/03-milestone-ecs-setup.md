# Milestone 3: Entity Component System (ECS) Setup

## Goal

Set up Donburi ECS and create your first game entities (characters and enemies) that exist on the hex grid.

## Learning Objectives

- Understand the Entity Component System pattern
- Learn Donburi's API (components, entities, queries)
- Create reusable components for game data
- Spawn entities and place them on the hex grid
- Query and iterate over entities
- Separate data (components) from logic (systems)

## Prerequisites

- Completed Milestone 2
- Installed `yohamta/donburi` package

## Tasks

### ECS Foundation
- [ ] Install Donburi: `go get github.com/yohamta/donburi`
- [ ] Create Donburi world in Game struct
- [ ] Understand component definition vs component data

### Define Components
- [ ] Create `internal/components/position.go` - hex grid position (q, r)
- [ ] Create `internal/components/stats.go` - D&D stats (STR, DEX, CON, INT, WIS, CHA)
- [ ] Create `internal/components/health.go` - current HP, max HP
- [ ] Create `internal/components/display.go` - name, color for rendering
- [ ] Create `internal/components/size.go` - how many hexes the entity occupies
- [ ] Write component tests in `internal/components/components_test.go`

### Create Entity Archetypes
- [ ] Create `internal/entities/character.go` - spawn party member
- [ ] Create `internal/entities/boss.go` - spawn boss enemy
- [ ] Define starter stats for different character classes

### Rendering System
- [ ] Create `internal/systems/render.go` - render entities on hex grid
- [ ] Query all entities with Position + Display components
- [ ] Draw entities as colored circles (for now)
- [ ] Display entity name and HP above them

### Integration
- [ ] Spawn 3-4 party members on grid
- [ ] Spawn 1 boss enemy
- [ ] Update game loop to run render system
- [ ] Verify entities appear on correct hexes

### Testing
- [ ] Run tests: `go test ./internal/components`
- [ ] All entities visible on grid
- [ ] Entity stats display correctly
- [ ] Clicking entity shows its info

## ECS Pattern Explained

### What is ECS?

Traditional OOP:
```
Character (class)
  ├─ position
  ├─ health
  ├─ attack()
  └─ takeDamage()

Enemy (class)
  ├─ position
  ├─ health
  ├─ attack()
  └─ patrol()
```

ECS approach:
```
Entity (just an ID)
├─ Position Component (data)
├─ Health Component (data)
└─ Stats Component (data)

Systems (logic):
├─ Render System: draws all entities with Position + Display
├─ Combat System: processes all entities with Stats + Health
└─ AI System: controls all entities with AI component
```

### Why ECS?

**Benefits:**
- **Composition over inheritance** - mix and match components
- **Data-oriented design** - better cache performance
- **Flexibility** - add/remove behaviors by adding/removing components
- **Clarity** - data separate from logic

**Example:**
- Player character = Position + Stats + Health + Display + PlayerControlled
- Boss enemy = Position + Stats + Health + Display + Size + AIControlled
- Prop/obstacle = Position + Display + Blocking

All share Position + Display, so same render system works for all.

### Donburi Specifics

**Components** are defined as types:
```go
// Component definition (type)
var PositionComponent = donburi.NewComponentType[PositionData]()

// Component data (struct)
type PositionData struct {
    Q int
    R int
}
```

**Entities** are created and components added:
```go
entity := world.Create(PositionComponent, HealthComponent)
```

**Queries** find entities with specific components:
```go
query := donburi.NewQuery(filter.Contains(PositionComponent, DisplayComponent))
query.Each(world, func(entry *donburi.Entry) {
    pos := PositionComponent.Get(entry)
    display := DisplayComponent.Get(entry)
    // Use data...
})
```

## Step-by-Step Implementation

### Step 1: Install Donburi

```bash
go get github.com/yohamta/donburi
go get github.com/yohamta/donburi/filter
```

### Step 2: Define Components

Create `internal/components/position.go`:

```go
package components

import "github.com/yohamta/donburi"

// PositionData stores hex grid coordinates
type PositionData struct {
	Q int // Axial coordinate q
	R int // Axial coordinate r
}

// PositionComponent is the component type
var PositionComponent = donburi.NewComponentType[PositionData]()
```

Create `internal/components/stats.go`:

```go
package components

import "github.com/yohamta/donburi"

// StatsData represents D&D-style ability scores
type StatsData struct {
	Strength     int // Physical power, melee damage
	Dexterity    int // Agility, initiative, ranged damage
	Constitution int // Toughness, HP modifier
	Intelligence int // Magical power, spell damage
	Wisdom       int // Perception, willpower
	Charisma     int // Leadership, persuasion
}

// StatModifier calculates the D&D modifier for a stat
// Formula: (stat - 10) / 2, rounded down
func (s *StatsData) Modifier(stat int) int {
	return (stat - 10) / 2
}

// Common modifier getters
func (s *StatsData) StrMod() int { return s.Modifier(s.Strength) }
func (s *StatsData) DexMod() int { return s.Modifier(s.Dexterity) }
func (s *StatsData) ConMod() int { return s.Modifier(s.Constitution) }
func (s *StatsData) IntMod() int { return s.Modifier(s.Intelligence) }
func (s *StatsData) WisMod() int { return s.Modifier(s.Wisdom) }
func (s *StatsData) ChaMod() int { return s.Modifier(s.Charisma) }

var StatsComponent = donburi.NewComponentType[StatsData]()
```

Create `internal/components/health.go`:

```go
package components

import "github.com/yohamta/donburi"

// HealthData tracks hit points
type HealthData struct {
	Current int
	Max     int
}

// IsDead checks if entity is at 0 HP
func (h *HealthData) IsDead() bool {
	return h.Current <= 0
}

// IsFullHealth checks if at max HP
func (h *HealthData) IsFullHealth() bool {
	return h.Current >= h.Max
}

// Damage reduces current HP
func (h *HealthData) Damage(amount int) {
	h.Current -= amount
	if h.Current < 0 {
		h.Current = 0
	}
}

// Heal increases current HP
func (h *HealthData) Heal(amount int) {
	h.Current += amount
	if h.Current > h.Max {
		h.Current = h.Max
	}
}

var HealthComponent = donburi.NewComponentType[HealthData]()
```

Create `internal/components/display.go`:

```go
package components

import (
	"image/color"

	"github.com/yohamta/donburi"
)

// DisplayData holds rendering information
type DisplayData struct {
	Name  string
	Color color.Color
}

var DisplayComponent = donburi.NewComponentType[DisplayData]()
```

Create `internal/components/size.go`:

```go
package components

import "github.com/yohamta/donburi"

// SizeData defines how many hex tiles an entity occupies
type SizeData struct {
	// For now, just a radius (0 = single hex, 1 = 7 hexes, etc.)
	// Later we can make this more complex with specific hex positions
	Radius int
}

// NumHexes returns approximate number of hexes occupied
func (s *SizeData) NumHexes() int {
	if s.Radius == 0 {
		return 1
	}
	// Hex ring formula: 1 + 6 + 12 + 18 + ... = 1 + 3*n*(n+1)
	return 1 + 3*s.Radius*(s.Radius+1)
}

var SizeComponent = donburi.NewComponentType[SizeData]()
```

### Step 3: Create Component Tests

Create `internal/components/components_test.go`:

```go
package components

import (
	"testing"
)

func TestStatModifier(t *testing.T) {
	tests := []struct {
		name     string
		stat     int
		expected int
	}{
		{"very low (3)", 3, -4},
		{"low (8)", 8, -1},
		{"average (10)", 10, 0},
		{"good (14)", 14, 2},
		{"great (18)", 18, 4},
		{"legendary (20)", 20, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := StatsData{Strength: tt.stat}
			got := stats.Modifier(tt.stat)

			if got != tt.expected {
				t.Errorf("Modifier(%d) = %d, want %d", tt.stat, got, tt.expected)
			}
		})
	}
}

func TestHealthDamage(t *testing.T) {
	health := HealthData{Current: 50, Max: 100}

	health.Damage(20)
	if health.Current != 30 {
		t.Errorf("After 20 damage: got %d HP, want 30", health.Current)
	}

	health.Damage(100)
	if health.Current != 0 {
		t.Errorf("After overkill damage: got %d HP, want 0", health.Current)
	}

	if !health.IsDead() {
		t.Error("Entity should be dead at 0 HP")
	}
}

func TestHealthHeal(t *testing.T) {
	health := HealthData{Current: 30, Max: 100}

	health.Heal(20)
	if health.Current != 50 {
		t.Errorf("After healing 20: got %d HP, want 50", health.Current)
	}

	health.Heal(100)
	if health.Current != 100 {
		t.Errorf("After overheal: got %d HP, want 100 (max)", health.Current)
	}

	if !health.IsFullHealth() {
		t.Error("Entity should be at full health")
	}
}

func TestSizeNumHexes(t *testing.T) {
	tests := []struct {
		radius   int
		expected int
	}{
		{0, 1},   // Single hex
		{1, 7},   // Center + 6 surrounding
		{2, 19},  // Center + 6 + 12
		{3, 37},  // Center + 6 + 12 + 18
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			size := SizeData{Radius: tt.radius}
			got := size.NumHexes()

			if got != tt.expected {
				t.Errorf("Radius %d: got %d hexes, want %d", tt.radius, got, tt.expected)
			}
		})
	}
}
```

Run tests:
```bash
go test ./internal/components
```

### Step 4: Create Entity Archetypes

Create `internal/entities/character.go`:

```go
package entities

import (
	"image/color"

	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/components"
)

// CharacterClass defines different character types
type CharacterClass int

const (
	Warrior CharacterClass = iota
	Rogue
	Mage
	Cleric
)

// SpawnCharacter creates a party member entity
func SpawnCharacter(world *donburi.World, class CharacterClass, q, r int) *donburi.Entry {
	entity := world.Create(
		components.PositionComponent,
		components.StatsComponent,
		components.HealthComponent,
		components.DisplayComponent,
		components.SizeComponent,
	)

	// Set position
	components.PositionComponent.Set(entity, &components.PositionData{
		Q: q,
		R: r,
	})

	// Set size (characters are single-hex)
	components.SizeComponent.Set(entity, &components.SizeData{
		Radius: 0,
	})

	// Set stats and display based on class
	switch class {
	case Warrior:
		components.StatsComponent.Set(entity, &components.StatsData{
			Strength:     16, // +3
			Dexterity:    12, // +1
			Constitution: 14, // +2
			Intelligence: 8,  // -1
			Wisdom:       10, // +0
			Charisma:     10, // +0
		})
		components.DisplayComponent.Set(entity, &components.DisplayData{
			Name:  "Warrior",
			Color: color.RGBA{200, 100, 100, 255}, // Reddish
		})
		// HP = 10 base + (level * CON modifier), let's say level 1
		components.HealthComponent.Set(entity, &components.HealthData{
			Max:     10 + 2, // 12 HP
			Current: 10 + 2,
		})

	case Rogue:
		components.StatsComponent.Set(entity, &components.StatsData{
			Strength:     12,
			Dexterity:    16, // +3
			Constitution: 12,
			Intelligence: 10,
			Wisdom:       12,
			Charisma:     14,
		})
		components.DisplayComponent.Set(entity, &components.DisplayData{
			Name:  "Rogue",
			Color: color.RGBA{100, 100, 200, 255}, // Blueish
		})
		components.HealthComponent.Set(entity, &components.HealthData{
			Max:     8 + 1,
			Current: 8 + 1,
		})

	case Mage:
		components.StatsComponent.Set(entity, &components.StatsData{
			Strength:     8,
			Dexterity:    12,
			Constitution: 10,
			Intelligence: 16, // +3
			Wisdom:       14,
			Charisma:     12,
		})
		components.DisplayComponent.Set(entity, &components.DisplayData{
			Name:  "Mage",
			Color: color.RGBA{150, 100, 200, 255}, // Purple
		})
		components.HealthComponent.Set(entity, &components.HealthData{
			Max:     6 + 0,
			Current: 6 + 0,
		})

	case Cleric:
		components.StatsComponent.Set(entity, &components.StatsData{
			Strength:     14,
			Dexterity:    10,
			Constitution: 12,
			Intelligence: 10,
			Wisdom:       16, // +3
			Charisma:     14,
		})
		components.DisplayComponent.Set(entity, &components.DisplayData{
			Name:  "Cleric",
			Color: color.RGBA{200, 200, 100, 255}, // Yellow
		})
		components.HealthComponent.Set(entity, &components.HealthData{
			Max:     8 + 1,
			Current: 8 + 1,
		})
	}

	return entity
}
```

Create `internal/entities/boss.go`:

```go
package entities

import (
	"image/color"

	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/components"
)

// SpawnBoss creates a boss enemy entity
func SpawnBoss(world *donburi.World, name string, q, r int) *donburi.Entry {
	entity := world.Create(
		components.PositionComponent,
		components.StatsComponent,
		components.HealthComponent,
		components.DisplayComponent,
		components.SizeComponent,
	)

	// Set position
	components.PositionComponent.Set(entity, &components.PositionData{
		Q: q,
		R: r,
	})

	// Bosses are big (occupy multiple hexes)
	components.SizeComponent.Set(entity, &components.SizeData{
		Radius: 1, // 7 hexes total
	})

	// Boss stats (powerful!)
	components.StatsComponent.Set(entity, &components.StatsData{
		Strength:     18, // +4
		Dexterity:    12, // +1
		Constitution: 16, // +3
		Intelligence: 14, // +2
		Wisdom:       14, // +2
		Charisma:     10, // +0
	})

	// Boss has lots of HP
	components.HealthComponent.Set(entity, &components.HealthData{
		Max:     50,
		Current: 50,
	})

	// Display
	components.DisplayComponent.Set(entity, &components.DisplayData{
		Name:  name,
		Color: color.RGBA{255, 50, 50, 255}, // Bright red
	})

	return entity
}
```

### Step 5: Create Render System

Create `internal/systems/render.go`:

```go
package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/yourusername/boss-battler/internal/components"
	"github.com/yourusername/boss-battler/internal/hex"
)

// RenderSystem draws all entities on the screen
func RenderSystem(world *donburi.World, screen *ebiten.Image, layout *hex.Layout) {
	// Query all entities that have Position and Display components
	query := donburi.NewQuery(
		filter.Contains(
			components.PositionComponent,
			components.DisplayComponent,
		),
	)

	// Iterate over all matching entities
	query.Each(world, func(entry *donburi.Entry) {
		pos := components.PositionComponent.Get(entry)
		display := components.DisplayComponent.Get(entry)

		// Get size if entity has it
		size := 0
		if entry.HasComponent(components.SizeComponent) {
			sizeData := components.SizeComponent.Get(entry)
			size = sizeData.Radius
		}

		// Get hex screen position
		x, y := layout.HexToPixel(pos.Q, pos.R)

		// Draw entity
		radius := float32(hex.HexSize * 0.6)
		if size > 0 {
			// Larger entities get bigger circles
			radius = float32(hex.HexSize * 0.8 * float64(size+1))
		}

		vector.DrawFilledCircle(screen, float32(x), float32(y), radius, display.Color, false)

		// Draw outline
		vector.StrokeCircle(screen, float32(x), float32(y), radius, 2, color.RGBA{255, 255, 255, 200}, false)

		// Draw name
		nameX := int(x) - len(display.Name)*3
		nameY := int(y) - int(radius) - 15
		ebitenutil.DebugPrintAt(screen, display.Name, nameX, nameY)

		// Draw HP if entity has health component
		if entry.HasComponent(components.HealthComponent) {
			health := components.HealthComponent.Get(entry)
			hpText := fmt.Sprintf("%d/%d HP", health.Current, health.Max)
			hpX := int(x) - len(hpText)*3
			hpY := int(y) + int(radius) + 5
			ebitenutil.DebugPrintAt(screen, hpText, hpX, hpY)
		}
	})
}
```

### Step 6: Integrate into Game

Update `cmd/game/main.go`:

```go
package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/entities"
	"github.com/yourusername/boss-battler/internal/hex"
	"github.com/yourusername/boss-battler/internal/systems"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	gridSize     = 5
)

type Game struct {
	world  *donburi.World
	layout *hex.Layout
}

func NewGame() *Game {
	game := &Game{
		world:  donburi.NewWorld(),
		layout: hex.NewLayout(),
	}

	// Spawn party members
	entities.SpawnCharacter(game.world, entities.Warrior, -2, -1)
	entities.SpawnCharacter(game.world, entities.Rogue, -2, 0)
	entities.SpawnCharacter(game.world, entities.Mage, -2, 1)
	entities.SpawnCharacter(game.world, entities.Cleric, -3, 0)

	// Spawn boss
	entities.SpawnBoss(game.world, "Dragon Boss", 2, 0)

	return game
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw hex grid (lighter, as background)
	for q := -gridSize; q <= gridSize; q++ {
		for r := -gridSize; r <= gridSize; r++ {
			g.drawHex(screen, q, r)
		}
	}

	// Draw entities using render system
	systems.RenderSystem(g.world, screen, g.layout)

	// UI
	ebitenutil.DebugPrint(screen, "Milestone 3: ECS\nParty vs Boss\nPress ESC to quit")
}

func (g *Game) drawHex(screen *ebiten.Image, q, r int) {
	cx, cy := g.layout.HexToPixel(q, r)
	hexColor := color.RGBA{40, 40, 50, 100} // Very dim
	vector.DrawFilledCircle(screen, float32(cx), float32(cy), hex.HexSize*0.7, hexColor, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boss Battler - ECS")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

Don't forget to add missing import:
```go
"github.com/hajimehoshi/ebiten/v2/vector"
```

### Step 7: Run and Verify

```bash
go test ./internal/components
go run cmd/game/main.go
```

You should see:
- Dim hex grid in background
- 4 party members on the left (different colors)
- 1 large red boss on the right
- Names and HP displayed above each entity

## Key Concepts

### Components Are Just Data

Components hold no logic, just data:
```go
type HealthData struct {
    Current int
    Max int
}
```

Helper methods are OK for convenience:
```go
func (h *HealthData) IsDead() bool {
    return h.Current <= 0
}
```

### Systems Are Just Functions

Systems contain logic, operate on components:
```go
func RenderSystem(world, screen, layout) {
    // Query entities with Position + Display
    // For each entity:
    //   Read position
    //   Read display
    //   Draw to screen
}
```

### Queries Find Entities

```go
// Find all entities with these components
query := donburi.NewQuery(
    filter.Contains(ComponentA, ComponentB),
)

// Iterate over results
query.Each(world, func(entry *donburi.Entry) {
    dataA := ComponentA.Get(entry)
    dataB := ComponentB.Get(entry)
    // Use data...
})
```

### Archetypes Are Factories

Instead of manually creating entities everywhere, use factory functions:
```go
func SpawnWarrior(world, q, r) {
    entity := world.Create(Position, Stats, Health, Display)
    // Set component data...
    return entity
}
```

## Common Issues

### "Component not registered" error
- Make sure you defined the component with `donburi.NewComponentType[T]()`
- Check you're using the right component variable

### Can't Get component data
- Verify entity was created with that component
- Use `entry.HasComponent(Component)` to check first

### Entities don't render
- Check RenderSystem is being called in Draw()
- Verify query filters match entity's components
- Make sure entity has Position + Display components

### Stats modifier wrong
- D&D formula: (stat - 10) / 2, **rounded down**
- Go's `/` operator does this automatically for ints

## Next Steps

Milestone 3 complete! You now have:
- A working ECS architecture
- Components for game data
- Entity archetypes (characters, boss)
- Separation of data and logic
- Foundation for combat system

In [Milestone 4](04-milestone-turn-system.md), you'll implement the turn-based system with initiative and action queues!

## Extra Challenges (Optional)

- [ ] Add a tag component to distinguish players from enemies
- [ ] Create more character classes
- [ ] Add experience/level components
- [ ] Create a system that damages entities over time
- [ ] Add equipment component (weapons, armor)
- [ ] Implement component add/remove dynamically
