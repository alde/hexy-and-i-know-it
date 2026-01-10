# Milestone 6: Player Input & Actions

## Goal

Implement player-controlled actions: selecting targets, executing attacks, and managing the action queue.

## Learning Objectives

- Implement command pattern for actions
- Handle mouse input for target selection
- Create action selection UI
- Validate targets (range, line of sight)
- Queue and execute player actions
- Distinguish player entities from AI entities

## Prerequisites

- Completed Milestone 5
- Understanding of command pattern

## Tasks

### Tag Components
- [ ] Create `internal/components/tags.go` for PlayerControlled and AIControlled
- [ ] Tag party members as PlayerControlled
- [ ] Tag boss as AIControlled

### Command Pattern
- [ ] Create `internal/commands/action.go` for action interface
- [ ] Create `internal/commands/attack.go` for attack command
- [ ] Create `internal/commands/wait.go` for wait/skip turn command
- [ ] Create tests in `internal/commands/commands_test.go`

### Input System
- [ ] Create `internal/systems/input.go` for player input handling
- [ ] Detect mouse clicks on hexes
- [ ] Show valid targets highlighted
- [ ] Select target with click
- [ ] Show action menu (attack, wait, etc.)

### Action Queue
- [ ] Create action queue in battle manager
- [ ] Queue player actions
- [ ] Execute action when player confirms
- [ ] Auto-advance to next turn after action

### UI Improvements
- [ ] Show "Select Target" prompt
- [ ] Highlight valid targets
- [ ] Show action buttons/menu
- [ ] Display combat log

### Testing
- [ ] Player can select allies and enemies
- [ ] Attack command validates target
- [ ] Invalid targets not selectable
- [ ] Action executes correctly
- [ ] Turn advances after action

## Command Pattern Explained

### What is Command Pattern?

Instead of directly calling functions, wrap them in objects:

**Without Command Pattern**:
```go
if keyPressed {
    attack(target) // Directly call
}
```

**With Command Pattern**:
```go
if keyPressed {
    cmd := AttackCommand{target: target}
    queue.Add(cmd)
}

// Later...
for _, cmd := range queue {
    cmd.Execute()
}
```

### Why Use Commands?

- **Undo/Redo**: Store commands to reverse them
- **Queuing**: Execute actions in sequence
- **Validation**: Check if action is valid before executing
- **Serialization**: Save/load commands for replays
- **Separation**: Decouple input from execution

## Step-by-Step Implementation

### Step 1: Create Tag Components

Create `internal/components/tags.go`:

```go
package components

import "github.com/yohamta/donburi"

// PlayerControlledData marks entities controlled by the player
type PlayerControlledData struct {
	// Empty marker component
}

var PlayerControlledComponent = donburi.NewComponentType[PlayerControlledData]()

// AIControlledData marks entities controlled by AI
type AIControlledData struct {
	// Empty marker component
}

var AIControlledComponent = donburi.NewComponentType[AIControlledData]()
```

Update `internal/entities/character.go`:

```go
// In SpawnCharacter, add PlayerControlledComponent:
entity := world.Create(
    components.PositionComponent,
    components.StatsComponent,
    components.HealthComponent,
    components.DisplayComponent,
    components.SizeComponent,
    components.InitiativeComponent,
    components.ArmorComponent,
    components.WeaponComponent,
    components.PlayerControlledComponent, // ADD THIS
)

// After creating entity:
components.PlayerControlledComponent.Set(entity, &components.PlayerControlledData{})
```

Update `internal/entities/boss.go`:

```go
// Add AIControlledComponent:
entity := world.Create(
    // ... existing components ...
    components.AIControlledComponent, // ADD THIS
)

components.AIControlledComponent.Set(entity, &components.AIControlledData{})
```

### Step 2: Create Command Interface

Create `internal/commands/action.go`:

```go
package commands

import "github.com/yohamta/donburi"

// Action represents a game action that can be executed
type Action interface {
	// Execute performs the action
	Execute(world *donburi.World) *ActionResult

	// Validate checks if the action is legal
	Validate(world *donburi.World) error

	// Description returns a human-readable description
	Description() string
}

// ActionResult contains the outcome of an action
type ActionResult struct {
	Success bool
	Message string
	Logs    []string
}
```

Create `internal/commands/attack.go`:

```go
package commands

import (
	"errors"
	"fmt"

	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/combat"
	"github.com/yourusername/boss-battler/internal/components"
)

// AttackAction represents an attack from one entity to another
type AttackAction struct {
	Attacker *donburi.Entry
	Target   *donburi.Entry
}

// Execute performs the attack
func (a *AttackAction) Execute(world *donburi.World) *ActionResult {
	if err := a.Validate(world); err != nil {
		return &ActionResult{
			Success: false,
			Message: err.Error(),
		}
	}

	// Perform the attack using combat system
	result := combat.PerformAttack(a.Attacker, a.Target)

	return &ActionResult{
		Success: true,
		Message: result.String(),
		Logs:    []string{result.String()},
	}
}

// Validate checks if the attack is valid
func (a *AttackAction) Validate(world *donburi.World) error {
	// Check attacker exists and is alive
	if !a.Attacker.Valid() {
		return errors.New("attacker is not valid")
	}

	attackerHealth := components.HealthComponent.Get(a.Attacker)
	if attackerHealth.IsDead() {
		return errors.New("attacker is dead")
	}

	// Check target exists and is alive
	if !a.Target.Valid() {
		return errors.New("target is not valid")
	}

	targetHealth := components.HealthComponent.Get(a.Target)
	if targetHealth.IsDead() {
		return errors.New("target is already dead")
	}

	// Check they're not the same entity
	if a.Attacker == a.Target {
		return errors.New("cannot attack self")
	}

	// TODO: Check range (for now, all targets are in range)

	return nil
}

// Description returns a human-readable description
func (a *AttackAction) Description() string {
	attackerName := components.DisplayComponent.Get(a.Attacker).Name
	targetName := components.DisplayComponent.Get(a.Target).Name
	return fmt.Sprintf("%s attacks %s", attackerName, targetName)
}
```

Create `internal/commands/wait.go`:

```go
package commands

import (
	"github.com/yohamta/donburi"

	"github.com/yourusername/boss-battler/internal/components"
)

// WaitAction represents skipping a turn
type WaitAction struct {
	Actor *donburi.Entry
}

// Execute performs the wait (does nothing)
func (w *WaitAction) Execute(world *donburi.World) *ActionResult {
	actorName := components.DisplayComponent.Get(w.Actor).Name

	return &ActionResult{
		Success: true,
		Message: actorName + " waits",
		Logs:    []string{actorName + " waits"},
	}
}

// Validate checks if wait is valid (always true)
func (w *WaitAction) Validate(world *donburi.World) error {
	return nil
}

// Description returns a human-readable description
func (w *WaitAction) Description() string {
	actorName := components.DisplayComponent.Get(w.Actor).Name
	return actorName + " waits"
}
```

### Step 3: Update Battle Manager

Update `internal/states/battle.go` to add action queue and player input state:

```go
// Add to imports:
import (
    "github.com/yourusername/boss-battler/internal/commands"
)

// Add to BattleState enum:
const (
    BattleStart BattleState = iota
    RollInitiative
    PlayerTurn
    WaitingForPlayerInput   // NEW
    EnemyTurn
    ExecuteAction          // NEW
    BattleEnd
)

// Update String() method:
func (s BattleState) String() string {
    return []string{
        "Battle Start",
        "Roll Initiative",
        "Player Turn",
        "Waiting For Input",
        "Enemy Turn",
        "Execute Action",
        "Battle End",
    }[s]
}

// Add to BattleManager:
type BattleManager struct {
    State          BattleState
    TurnOrder      []*donburi.Entry
    CurrentIdx     int
    Round          int
    PendingAction  commands.Action       // NEW
    CombatLog      []string              // NEW
    MaxLogEntries  int                   // NEW
}

// Update NewBattleManager:
func NewBattleManager() *BattleManager {
    return &BattleManager{
        State:         BattleStart,
        TurnOrder:     make([]*donburi.Entry, 0),
        CurrentIdx:    0,
        Round:         1,
        CombatLog:     make([]string, 0),
        MaxLogEntries: 10,
    }
}

// Add method to queue action:
func (b *BattleManager) QueueAction(action commands.Action) error {
    if err := action.Validate(nil); err != nil {
        return err
    }

    b.PendingAction = action
    b.State = ExecuteAction
    return nil
}

// Add method to execute pending action:
func (b *BattleManager) ExecutePendingAction(world *donburi.World) {
    if b.PendingAction == nil {
        return
    }

    result := b.PendingAction.Execute(world)

    // Add to combat log
    for _, logEntry := range result.Logs {
        b.AddLog(logEntry)
    }

    b.PendingAction = nil

    // Advance to next turn
    b.NextTurn()
}

// Add method to add log entry:
func (b *BattleManager) AddLog(message string) {
    b.CombatLog = append(b.CombatLog, message)

    // Keep only last N entries
    if len(b.CombatLog) > b.MaxLogEntries {
        b.CombatLog = b.CombatLog[1:]
    }
}

// Update determineState to use WaitingForPlayerInput:
func (b *BattleManager) determineState() {
    if b.CurrentIdx >= len(b.TurnOrder) {
        return
    }

    current := b.TurnOrder[b.CurrentIdx]

    // Check if player-controlled
    if current.HasComponent(components.PlayerControlledComponent) {
        b.State = WaitingForPlayerInput
    } else {
        b.State = EnemyTurn
    }
}
```

### Step 4: Create Input System

Create `internal/systems/input.go`:

```go
package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/yourusername/boss-battler/internal/commands"
	"github.com/yourusername/boss-battler/internal/components"
	"github.com/yourusername/boss-battler/internal/hex"
)

// InputState tracks player input state
type InputState struct {
	SelectedTarget *donburi.Entry
	HoveredHex     *HexCoord
}

type HexCoord struct {
	Q, R int
}

// ProcessPlayerInput handles mouse input for target selection
func ProcessPlayerInput(
	world *donburi.World,
	layout *hex.Layout,
	inputState *InputState,
	activeCombatant *donburi.Entry,
) commands.Action {
	// Get mouse position
	mx, my := ebiten.CursorPosition()
	q, r := layout.PixelToHex(float64(mx), float64(my))

	// Update hovered hex
	inputState.HoveredHex = &HexCoord{Q: q, R: r}

	// Check for click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Find entity at clicked hex
		target := FindEntityAtHex(world, q, r)

		if target != nil {
			inputState.SelectedTarget = target

			// Create attack action
			return &commands.AttackAction{
				Attacker: activeCombatant,
				Target:   target,
			}
		}
	}

	// Check for wait key (W)
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return &commands.WaitAction{
			Actor: activeCombatant,
		}
	}

	return nil
}

// FindEntityAtHex finds an entity at the given hex coordinates
func FindEntityAtHex(world *donburi.World, q, r int) *donburi.Entry {
	query := donburi.NewQuery(
		filter.Contains(components.PositionComponent),
	)

	var found *donburi.Entry

	query.Each(world, func(entry *donburi.Entry) {
		pos := components.PositionComponent.Get(entry)
		if pos.Q == q && pos.R == r {
			found = entry
		}
	})

	return found
}

// GetValidTargets returns entities that can be targeted
func GetValidTargets(world *donburi.World) []*donburi.Entry {
	query := donburi.NewQuery(
		filter.Contains(
			components.HealthComponent,
			components.PositionComponent,
		),
	)

	var targets []*donburi.Entry

	query.Each(world, func(entry *donburi.Entry) {
		health := components.HealthComponent.Get(entry)
		if !health.IsDead() {
			targets = append(targets, entry)
		}
	})

	return targets
}
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
	world      *donburi.World
	layout     *hex.Layout
	battle     *states.BattleManager
	inputState *systems.InputState
	started    bool
}

func NewGame() *Game {
	game := &Game{
		world:      donburi.NewWorld(),
		layout:     hex.NewLayout(),
		battle:     states.NewBattleManager(),
		inputState: &systems.InputState{},
		started:    false,
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

	// Handle different battle states
	switch g.battle.State {
	case states.WaitingForPlayerInput:
		activeCombatant := g.battle.GetActiveCombatant()
		if activeCombatant != nil {
			action := systems.ProcessPlayerInput(g.world, g.layout, g.inputState, activeCombatant)
			if action != nil {
				if err := g.battle.QueueAction(action); err != nil {
					log.Printf("Invalid action: %v", err)
				}
			}
		}

	case states.ExecuteAction:
		g.battle.ExecutePendingAction(g.world)

	case states.EnemyTurn:
		// Simple AI: attack first valid target
		activeCombatant := g.battle.GetActiveCombatant()
		if activeCombatant != nil {
			targets := systems.GetValidTargets(g.world)

			// Find a player-controlled target
			for _, target := range targets {
				if target.HasComponent(components.PlayerControlledComponent) {
					action := &commands.AttackAction{
						Attacker: activeCombatant,
						Target:   target,
					}
					g.battle.QueueAction(action)
					break
				}
			}
		}
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

	// Highlight hovered hex
	if g.inputState.HoveredHex != nil && g.battle.State == states.WaitingForPlayerInput {
		entity := systems.FindEntityAtHex(g.world, g.inputState.HoveredHex.Q, g.inputState.HoveredHex.R)
		if entity != nil {
			pos := components.PositionComponent.Get(entity)
			cx, cy := g.layout.HexToPixel(pos.Q, pos.R)
			vector.StrokeCircle(screen, float32(cx), float32(cy), hex.HexSize*0.9, 3,
				color.RGBA{255, 255, 0, 150}, false)
		}
	}

	// Draw UI
	g.drawUI(screen)
}

func (g *Game) drawUI(screen *ebiten.Image) {
	// Instructions
	instructions := "Milestone 6: Player Input\n"
	if g.battle.State == states.WaitingForPlayerInput {
		instructions += "Click enemy to attack\nPress W to wait\n"
	} else {
		instructions += fmt.Sprintf("State: %s\n", g.battle.State.String())
	}
	instructions += "Press ESC to quit"
	ebitenutil.DebugPrint(screen, instructions)

	// Turn order (top right)
	g.drawTurnOrderUI(screen)

	// Combat log (bottom right)
	g.drawCombatLog(screen)
}

func (g *Game) drawTurnOrderUI(screen *ebiten.Image) {
	turnOrderText := fmt.Sprintf("Round: %d\n\nTurn Order:\n", g.battle.Round)
	for i, entry := range g.battle.TurnOrder {
		display := components.DisplayComponent.Get(entry)
		initiative := components.InitiativeComponent.Get(entry)
		health := components.HealthComponent.Get(entry)

		marker := "  "
		if i == g.battle.CurrentIdx {
			marker = "→ "
		} else if initiative.Went {
			marker = "✓ "
		}

		turnOrderText += fmt.Sprintf("%s%s (%d HP, Init:%d)\n",
			marker, display.Name, health.Current, initiative.Roll)
	}

	ebitenutil.DebugPrintAt(screen, turnOrderText, screenWidth-300, 10)
}

func (g *Game) drawCombatLog(screen *ebiten.Image) {
	logText := "Combat Log:\n"
	for _, entry := range g.battle.CombatLog {
		logText += entry + "\n"
	}

	ebitenutil.DebugPrintAt(screen, logText, 10, screenHeight-200)
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
	ebiten.SetWindowTitle("Boss Battler - Player Input")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

Don't forget to add missing imports at the top:
```go
import (
	"github.com/yourusername/boss-battler/internal/commands"
)
```

### Step 6: Run and Test

```bash
go run cmd/game/main.go
```

You should see:
- Player turn: "Click enemy to attack" prompt
- Hover over enemies: yellow highlight
- Click enemy: attack executes, damage shown in log
- Press W: skip turn
- Enemy turn: auto-attacks player
- Combat log updates with results

## Key Concepts

### Command Pattern Benefits

- **Validation**: Check if action is valid before executing
- **Queueing**: Store actions for later execution
- **Undo**: Could reverse actions (not implemented here)
- **Testing**: Easy to test commands in isolation

### Input State Management

Separate input handling from game logic:
```go
type InputState struct {
    SelectedTarget *Entry
    HoveredHex     *Coord
}
```

Allows UI to show what player is targeting.

### Player vs AI

Use marker components:
```go
if entity.HasComponent(PlayerControlledComponent) {
    // Wait for input
} else {
    // Run AI
}
```

## Common Issues

### Clicks don't register
- Check mouse coordinates are being converted to hex correctly
- Verify FindEntityAtHex is searching correctly
- Print hex coords to debug

### Wrong entity attacked
- Check if multiple entities at same hex
- Verify target selection logic
- Print selected target name

### Action doesn't execute
- Check QueueAction is being called
- Verify ExecutePendingAction runs
- Check action validation passes

### Combat log doesn't update
- Verify AddLog is called in ExecutePendingAction
- Check CombatLog array size
- Ensure result.Logs is populated

## Next Steps

Milestone 6 complete! You now have:
- Full player input system
- Command pattern for actions
- Target selection with mouse
- Action validation
- Combat log
- Player vs AI distinction

In [Milestone 7](07-milestone-demo-battle.md), you'll polish everything into a complete playable demo battle!

## Extra Challenges (Optional)

- [ ] Add action menu UI (buttons instead of just keys)
- [ ] Show attack preview (expected damage range)
- [ ] Highlight valid targets differently from invalid ones
- [ ] Add range checking for attacks
- [ ] Implement movement action
- [ ] Add ability/spell actions
- [ ] Create animation for attacks
- [ ] Add sound effects
