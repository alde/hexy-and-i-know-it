# Milestone 4: Turn-Based System

## Goal

Implement a turn-based combat system with initiative, turn order, and state management.

## Learning Objectives

- Understand state machines for game flow
- Implement initiative system (D&D-style)
- Manage turn order and active combatant
- Create game states (player turn, enemy turn, action execution)
- Build a turn queue system

## Prerequisites

- Completed Milestone 3
- Understanding of state machines

## Tasks

### Initiative System
- [ ] Create `internal/components/initiative.go` component
- [ ] Add initiative component to all combatants
- [ ] Implement initiative roll (d20 + DEX modifier)
- [ ] Sort entities by initiative (high to low)
- [ ] Create turn order queue

### Battle State Machine
- [ ] Create `internal/states/battle.go` for battle states
- [ ] Define states: BattleStart, PlayerTurn, EnemyTurn, ActionExecute, BattleEnd
- [ ] Implement state transitions
- [ ] Track current combatant
- [ ] Add "next turn" functionality

### Turn Management
- [ ] Create turn component to mark active combatant
- [ ] Highlight active combatant visually
- [ ] Display turn order on screen
- [ ] Show current state
- [ ] Implement end turn action

### Integration
- [ ] Roll initiative when battle starts
- [ ] Advance turns on key press (SPACE for testing)
- [ ] Show whose turn it is
- [ ] Cycle through all combatants
- [ ] Detect when round completes

### Testing
- [ ] Initiative values are correct (d20 + DEX mod)
- [ ] Turn order is correct (high to low)
- [ ] Turns advance properly
- [ ] Active combatant highlighted
- [ ] State transitions work

## Turn-Based Combat Flow

```
BattleStart
    ↓
Roll Initiative
    ↓
Sort Turn Order
    ↓
First Combatant's Turn
    ↓
┌─→ Player Turn? ────────────→ Wait for Player Input
│   │                              ↓
│   │                         Select Action
│   │                              ↓
│   │                         Execute Action
│   │                              ↓
│   Enemy Turn? ────────────→ AI Select Action
│   │                              ↓
│   │                         Execute Action
│   │                              ↓
│   End of Turn
│   │
│   Next Combatant
│   │
│   ↓
└─ More Combatants? ──Yes──┘
    │
    No
    ↓
End of Round
    ↓
Battle Over? ──No──→ Start New Round
    │
    Yes
    ↓
BattleEnd
```

## State Machine Diagram

```
BattleStart ──→ RollInitiative ──→ SortTurnOrder
                                        ↓
BattleEnd ←──────────────────────→ NextTurn
                                        ↓
                              ┌─── PlayerTurn ────┐
                              │         ↓         │
                              │   WaitingInput    │
                              │         ↓         │
                              │   ActionSelected  │
                              │                   │
                              ├─── EnemyTurn ─────┤
                              │         ↓         │
                              │   AI Decision     │
                              │                   │
                              └─→ ExecuteAction ←─┘
                                        ↓
                                  (Back to NextTurn)
```

## Step-by-Step Implementation

### Step 1: Create Initiative Component

Create `internal/components/initiative.go`:

```go
package components

import (
	"math/rand"

	"github.com/yohamta/donburi"
)

// InitiativeData stores initiative for turn order
type InitiativeData struct {
	Roll int  // The d20 + modifier roll
	Went bool // Has this entity taken a turn this round?
}

// RollInitiative rolls d20 + dexterity modifier
func RollInitiative(dexMod int) int {
	d20 := rand.Intn(20) + 1 // 1-20
	return d20 + dexMod
}

var InitiativeComponent = donburi.NewComponentType[InitiativeData]()

// ActiveTurnData marks the entity whose turn it currently is
type ActiveTurnData struct {
	// Empty marker component (just presence matters)
}

var ActiveTurnComponent = donburi.NewComponentType[ActiveTurnData]()
```

### Step 2: Create Battle State

Create `internal/states/battle.go`:

```go
package states

import (
	"sort"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/yourusername/boss-battler/internal/components"
)

// BattleState represents the current state of the battle
type BattleState int

const (
	BattleStart BattleState = iota
	RollInitiative
	PlayerTurn
	EnemyTurn
	ExecuteAction
	BattleEnd
)

// String returns the name of the state
func (s BattleState) String() string {
	return []string{
		"Battle Start",
		"Roll Initiative",
		"Player Turn",
		"Enemy Turn",
		"Execute Action",
		"Battle End",
	}[s]
}

// BattleManager manages turn order and state
type BattleManager struct {
	State      BattleState
	TurnOrder  []*donburi.Entry
	CurrentIdx int
	Round      int
}

// NewBattleManager creates a new battle manager
func NewBattleManager() *BattleManager {
	return &BattleManager{
		State:      BattleStart,
		TurnOrder:  make([]*donburi.Entry, 0),
		CurrentIdx: 0,
		Round:      1,
	}
}

// StartBattle initializes the battle
func (b *BattleManager) StartBattle(world *donburi.World) {
	b.State = RollInitiative
	b.rollInitiative(world)
	b.sortTurnOrder(world)
	b.State = PlayerTurn // Assume first in order is player
	b.markActiveCombatant()
}

// rollInitiative rolls initiative for all combatants
func (b *BattleManager) rollInitiative(world *donburi.World) {
	// Find all entities with Stats and Initiative components
	query := donburi.NewQuery(
		filter.Contains(
			components.StatsComponent,
			components.InitiativeComponent,
		),
	)

	query.Each(world, func(entry *donburi.Entry) {
		stats := components.StatsComponent.Get(entry)
		initiative := components.InitiativeComponent.Get(entry)

		// Roll initiative: d20 + DEX modifier
		initiative.Roll = components.RollInitiative(stats.DexMod())
		initiative.Went = false
	})
}

// sortTurnOrder sorts entities by initiative (high to low)
func (b *BattleManager) sortTurnOrder(world *donburi.World) {
	b.TurnOrder = make([]*donburi.Entry, 0)

	// Collect all combatants
	query := donburi.NewQuery(
		filter.Contains(components.InitiativeComponent),
	)

	query.Each(world, func(entry *donburi.Entry) {
		b.TurnOrder = append(b.TurnOrder, entry)
	})

	// Sort by initiative (high to low)
	sort.Slice(b.TurnOrder, func(i, j int) bool {
		initI := components.InitiativeComponent.Get(b.TurnOrder[i])
		initJ := components.InitiativeComponent.Get(b.TurnOrder[j])
		return initI.Roll > initJ.Roll
	})
}

// NextTurn advances to the next combatant
func (b *BattleManager) NextTurn() {
	// Mark current combatant as having gone
	if b.CurrentIdx < len(b.TurnOrder) {
		current := b.TurnOrder[b.CurrentIdx]
		initiative := components.InitiativeComponent.Get(current)
		initiative.Went = true

		// Remove active turn marker
		if current.HasComponent(components.ActiveTurnComponent) {
			current.RemoveComponent(components.ActiveTurnComponent)
		}
	}

	// Move to next combatant
	b.CurrentIdx++

	// Check if round is over
	if b.CurrentIdx >= len(b.TurnOrder) {
		b.endRound()
		return
	}

	b.markActiveCombatant()
	b.determineState()
}

// endRound handles end of round logic
func (b *BattleManager) endRound() {
	// Reset for new round
	b.CurrentIdx = 0
	b.Round++

	// Reset "went" flags
	for _, entry := range b.TurnOrder {
		initiative := components.InitiativeComponent.Get(entry)
		initiative.Went = false
	}

	b.markActiveCombatant()
	b.determineState()
}

// markActiveCombatant adds ActiveTurnComponent to current combatant
func (b *BattleManager) markActiveCombatant() {
	if b.CurrentIdx >= len(b.TurnOrder) {
		return
	}

	current := b.TurnOrder[b.CurrentIdx]
	if !current.HasComponent(components.ActiveTurnComponent) {
		current.AddComponent(components.ActiveTurnComponent)
	}
}

// determineState figures out if it's player or enemy turn
func (b *BattleManager) determineState() {
	if b.CurrentIdx >= len(b.TurnOrder) {
		return
	}

	current := b.TurnOrder[b.CurrentIdx]

	// For now, check if it has specific display name (we'll improve this later with tags)
	display := components.DisplayComponent.Get(current)

	if display.Name == "Warrior" || display.Name == "Rogue" ||
		display.Name == "Mage" || display.Name == "Cleric" {
		b.State = PlayerTurn
	} else {
		b.State = EnemyTurn
	}
}

// GetActiveCombatant returns the current combatant
func (b *BattleManager) GetActiveCombatant() *donburi.Entry {
	if b.CurrentIdx >= len(b.TurnOrder) {
		return nil
	}
	return b.TurnOrder[b.CurrentIdx]
}
```

### Step 3: Update Entity Creation

Update `internal/entities/character.go` to add initiative:

```go
// In SpawnCharacter function, add InitiativeComponent to entity.Create:
entity := world.Create(
    components.PositionComponent,
    components.StatsComponent,
    components.HealthComponent,
    components.DisplayComponent,
    components.SizeComponent,
    components.InitiativeComponent, // ADD THIS
)

// After setting other components, initialize initiative:
components.InitiativeComponent.Set(entity, &components.InitiativeData{
    Roll: 0,
    Went: false,
})
```

Update `internal/entities/boss.go` similarly:

```go
entity := world.Create(
    components.PositionComponent,
    components.StatsComponent,
    components.HealthComponent,
    components.DisplayComponent,
    components.SizeComponent,
    components.InitiativeComponent, // ADD THIS
)

components.InitiativeComponent.Set(entity, &components.InitiativeData{
    Roll: 0,
    Went: false,
})
```

### Step 4: Update Render System

Update `internal/systems/render.go` to highlight active combatant:

```go
// In RenderSystem function, after getting display:
display := components.DisplayComponent.Get(entry)

// Check if this is the active combatant
isActive := entry.HasComponent(components.ActiveTurnComponent)

// Modify outline color based on active status
var outlineColor color.Color
if isActive {
	outlineColor = color.RGBA{255, 255, 0, 255} // Yellow for active
} else {
	outlineColor = color.RGBA{255, 255, 255, 200} // White for inactive
}

// When drawing outline:
vector.StrokeCircle(screen, float32(x), float32(y), radius, 3, outlineColor, false)
// Note: increased thickness to 3 for active combatant visibility
```

### Step 5: Update Main Game

Update `cmd/game/main.go`:

```go
package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/components"
	"github.com/yourusername/boss-battler/internal/entities"
	"github.com/yourusername/boss-battler/internal/hex"
	"github.com/yourusername/boss-battler/internal/states"
	"github.com/yourusername/boss-battler/internal/systems"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	gridSize     = 5
)

type Game struct {
	world   *donburi.World
	layout  *hex.Layout
	battle  *states.BattleManager
	started bool
}

func NewGame() *Game {
	game := &Game{
		world:   donburi.NewWorld(),
		layout:  hex.NewLayout(),
		battle:  states.NewBattleManager(),
		started: false,
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

	// Start battle on first frame
	if !g.started {
		g.battle.StartBattle(g.world)
		g.started = true
	}

	// Advance turn with SPACE key (for testing)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.battle.NextTurn()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw hex grid
	for q := -gridSize; q <= gridSize; q++ {
		for r := -gridSize; r <= gridSize; r++ {
			g.drawHex(screen, q, r)
		}
	}

	// Draw entities
	systems.RenderSystem(g.world, screen, g.layout)

	// Draw turn order UI
	g.drawTurnOrderUI(screen)

	// Draw controls
	ebitenutil.DebugPrint(screen, "Milestone 4: Turn System\nPress SPACE to advance turn\nPress ESC to quit")
}

func (g *Game) drawTurnOrderUI(screen *ebiten.Image) {
	// Draw state
	stateText := fmt.Sprintf("State: %s\nRound: %d\n\n", g.battle.State.String(), g.battle.Round)
	ebitenutil.DebugPrintAt(screen, stateText, screenWidth-250, 10)

	// Draw turn order
	turnOrderText := "Turn Order:\n"
	for i, entry := range g.battle.TurnOrder {
		display := components.DisplayComponent.Get(entry)
		initiative := components.InitiativeComponent.Get(entry)

		marker := "  "
		if i == g.battle.CurrentIdx {
			marker = "→ " // Arrow for current turn
		} else if initiative.Went {
			marker = "✓ " // Check for went this round
		}

		turnOrderText += fmt.Sprintf("%s%s (Init: %d)\n",
			marker, display.Name, initiative.Roll)
	}

	ebitenutil.DebugPrintAt(screen, turnOrderText, screenWidth-250, 70)
}

func (g *Game) drawHex(screen *ebiten.Image, q, r int) {
	cx, cy := g.layout.HexToPixel(q, r)
	hexColor := color.RGBA{40, 40, 50, 100}
	vector.DrawFilledCircle(screen, float32(cx), float32(cy), hex.HexSize*0.7, hexColor, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boss Battler - Turn System")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

### Step 6: Run and Test

```bash
go run cmd/game/main.go
```

You should see:
- All entities with rolled initiative values
- Turn order list on the right
- Active combatant highlighted with yellow outline
- Arrow pointing to current turn
- Checkmarks for entities that went this round
- Pressing SPACE advances to next turn
- After all entities go, round increments

## Key Concepts

### Initiative (D&D-style)

Initiative determines turn order:
- Roll: d20 + Dexterity modifier
- Higher rolls go first
- Ties can be broken by highest DEX (you can add this later)

### State Machine

Clean way to manage game flow:
```go
type BattleState int

const (
    StateA BattleState = iota
    StateB
    StateC
)

// Transitions
switch currentState {
case StateA:
    currentState = StateB
case StateB:
    currentState = StateC
}
```

### Turn Order Queue

Array of entities sorted by initiative:
```go
[Rogue(19), Warrior(15), Mage(12), Boss(10), Cleric(8)]
 ^          ^           ^         ^         ^
 idx 0      idx 1       idx 2     idx 3     idx 4
```

Increment index each turn, wrap at end.

### Marker Components

Components can be empty (just presence matters):
```go
type ActiveTurnData struct {
    // Empty - we only care if entity has this component
}

if entity.HasComponent(ActiveTurnComponent) {
    // This entity's turn
}
```

## Common Issues

### Initiative values all the same
- Make sure you're using `rand.Seed()` or Go 1.20+ (auto-seeded)
- Check you're actually rolling (not using default 0)

### Turn order doesn't change
- Verify sort is working (check slice after sort)
- Print initiative values to debug

### Active combatant not highlighting
- Check if ActiveTurnComponent is being added/removed correctly
- Verify render system checks for component

### Turns skip entities
- Check CurrentIdx increments correctly
- Verify TurnOrder length matches entity count

### Round never ends
- Make sure CurrentIdx resets after reaching end
- Check if condition (>= len(TurnOrder))

## Next Steps

Milestone 4 complete! You now have:
- Initiative system
- Turn-based state machine
- Turn order management
- Visual feedback for active combatant
- Round tracking

In [Milestone 5](05-milestone-combat-mechanics.md), you'll implement actual combat: attacks, damage calculation, and D&D-style mechanics!

## Extra Challenges (Optional)

- [ ] Add tie-breaker for same initiative (use DEX stat)
- [ ] Implement delay turn action (go later in initiative)
- [ ] Add "surprise round" where some entities act first
- [ ] Show initiative rolls with animation
- [ ] Add sound effect when turn changes
- [ ] Implement persistent turn order (don't re-roll each round)
